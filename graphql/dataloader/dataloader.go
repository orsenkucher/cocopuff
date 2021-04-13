//go:generate dataloaden AccountLoader string *github.com/orsenkucher/cocopuff/graphql.Account
//go:generate dataloaden AccountPaginatedLoader *github.com/orsenkucher/cocopuff/graphql/pb.ListAccountsRequest []*github.com/orsenkucher/cocopuff/graphql.Account

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/graphql/env"
	"github.com/orsenkucher/cocopuff/graphql/pb"
)

type Dataloader struct {
	AccountById      *AccountLoader
	AccountPaginated *AccountPaginatedLoader
}

func Middleware(next http.Handler, client *graphql.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), env.Dataloader, &Dataloader{
			AccountById:      NewAccountLoader(NewAccountLoaderConfig(r.Context(), client)),
			AccountPaginated: NewAccountPaginatedLoader(NewAccountPaginatedLoaderConfig(r.Context(), client)),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Dataloader {
	return ctx.Value(env.Dataloader).(*Dataloader)
}

func NewAccountLoaderConfig(ctx context.Context, client *graphql.Client) AccountLoaderConfig {
	return AccountLoaderConfig{
		MaxBatch: 100,
		Wait:     1 * time.Millisecond,
		Fetch: func(ids []string) ([]*graphql.Account, []error) {
			// TODO: GetAccounts()
			res, err := client.ListAccounts(ctx, 0, 0)
			if err != nil {
				// single error for everything
				return nil, []error{err}
			}

			return res, nil
		},
	}
}

func NewAccountPaginatedLoaderConfig(ctx context.Context, client *graphql.Client) AccountPaginatedLoaderConfig {
	return AccountPaginatedLoaderConfig{
		MaxBatch: 10,
		Wait:     1 * time.Millisecond,
		Fetch: func(keys []*pb.ListAccountsRequest) ([][]*graphql.Account, []error) {
			var results [][]*graphql.Account
			var errors []error

			for _, key := range keys {
				res, err := client.ListAccountsPB(ctx, key)
				if err != nil {
					errors = append(errors, err)
					continue
				}

				results = append(results, res)
			}

			return results, errors
		},
	}
}
