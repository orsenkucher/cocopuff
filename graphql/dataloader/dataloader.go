//go:generate dataloaden AccountLoader string *github.com/orsenkucher/cocopuff/graphql.Account

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/graphql/env"
)

type Dataloader struct {
	AccountById *AccountLoader
}

func Middleware(next http.Handler, client *graphql.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), env.Dataloader, &Dataloader{
			AccountById: NewAccountLoader(NewAccountLoaderConfig(r.Context(), client)),
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
