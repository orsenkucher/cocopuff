//go:generate dataloaden AccountLoader string *github.com/orsenkucher/cocopuff/graphql.Account
//go:generate dataloaden AccountPaginatedLoader *github.com/orsenkucher/cocopuff/graphql/pb.ListAccountsRequest []*github.com/orsenkucher/cocopuff/graphql.Account

package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/graphql/pb"
	"github.com/orsenkucher/cocopuff/pub/care"
	"github.com/orsenkucher/cocopuff/pub/log"
	"go.uber.org/zap"
)

type ctxKey int

const (
	sugarCtx ctxKey = iota
	dataloaderCtx
)

type Dataloader struct {
	AccountById      *AccountLoader
	AccountPaginated *AccountPaginatedLoader
}

func Middleware(sugar *zap.SugaredLogger, client *graphql.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if log, ok := log.For(ctx); ok {
				sugar = log.Sugar()
			}

			sugar := sugar.With(zap.String("package", "dataloader"))
			ctx = context.WithValue(ctx, sugarCtx, sugar)
			ctx = context.WithValue(ctx, dataloaderCtx, &Dataloader{
				AccountById:      NewAccountLoader(NewAccountLoaderConfig(r.Context(), sugar, client)),
				AccountPaginated: NewAccountPaginatedLoader(NewAccountPaginatedLoaderConfig(r.Context(), sugar, client)),
			})
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func For(ctx context.Context) *Dataloader {
	if dl, ok := ctx.Value(dataloaderCtx).(*Dataloader); ok {
		return dl
	}

	if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
		sugar.DPanicw("fail to retrieve dataloader", zap.String("function", "For"))
	}

	return nil
}

func NewAccountLoaderConfig(ctx context.Context, sugar *zap.SugaredLogger, client *graphql.Client) AccountLoaderConfig {
	return AccountLoaderConfig{
		MaxBatch: 100,
		Wait:     1 * time.Millisecond,
		Fetch: func(ids []string) ([]*graphql.Account, []error) {
			w := care.With(zap.String("dataloader", fmt.Sprintf("%T", AccountLoaderConfig{})))
			res, err := client.GetAccounts(ctx, ids)
			if err != nil {
				// single error for everything
				return nil, []error{w.Of(err, "fail to get accounts", zap.String("by", "ids"))}
			}

			sugar.With(w.Fields).Infof("got %v accounts", len(res))
			return res, nil
		},
	}
}

func NewAccountPaginatedLoaderConfig(
	ctx context.Context,
	sugar *zap.SugaredLogger,
	client *graphql.Client,
) AccountPaginatedLoaderConfig {
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
