package auth

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/orsenkucher/cocopuff/graphql"
	"go.uber.org/zap"
)

type ctxKey int

const (
	sugarCtx ctxKey = iota
	accountCtx
)

func Middleware(sugar *zap.SugaredLogger, client *graphql.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sugar = sugar.With(zap.String("package", "auth"))
			ctx := context.WithValue(r.Context(), sugarCtx, sugar)

			c, err := r.Cookie("auth-cookie")
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			id, err := validateAndGetAccountID(c)
			if err != nil {
				http.Error(w, "invalid cookie", http.StatusForbidden)
				return
			}

			account, err := client.GetAccount(ctx, id)
			if err != nil {
				http.Error(w, "no account", http.StatusForbidden)
				return
			}

			ctx = context.WithValue(ctx, accountCtx, account)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func validateAndGetAccountID(c *http.Cookie) (string, error) {
	return "orsen-id", nil
}

func validateAndGetAccountIDString(token string) (string, error) {
	return "orsen-id", nil
}

func For(ctx context.Context) *graphql.Account {
	if a, ok := ctx.Value(accountCtx).(*graphql.Account); ok {
		return a
	}

	if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
		sugar.DPanicw("fail to retrieve account", zap.String("function", "For"))
	}

	return nil
}

func WebsocketMiddleware(sugar *zap.SugaredLogger, client *graphql.Client) transport.WebsocketInitFunc {
	fn := func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
		id, err := validateAndGetAccountIDString(initPayload.GetString("token"))
		if err != nil {
			return nil, err
		}

		account, err := client.GetAccount(ctx, id)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, accountCtx, account)

		return ctx, nil
	}
	return transport.WebsocketInitFunc(fn)
}
