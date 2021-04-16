package authentication

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/pub/log"
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
			ctx := r.Context()
			reqID := middleware.GetReqID(ctx)
			if log, ok := log.For(ctx); ok {
				sugar = log.Sugar()
			}

			sugar := sugar.With(
				zap.String("reqID", reqID), zap.String("package", "authentication"), zap.String("protocol", "http"),
			)
			ctx = context.WithValue(ctx, sugarCtx, sugar)

			c, err := r.Cookie("auth-cookie")
			if err != nil || c == nil {
				sugar.Info("serving with no token")
				next.ServeHTTP(w, r)
				return
			}

			id, err := validateAndGetAccountID(c)
			if err != nil {
				sugar.Infow("invalid token", zap.Error(err))
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}

			account, err := client.GetAccount(ctx, id)
			if err != nil {
				sugar.Warnw("no account found", zap.Error(err))
				http.Error(w, "no account", http.StatusForbidden)
				return
			}

			sugar.Info("serving authenticated account")
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
		sugar = sugar.With(zap.String("package", "authentication"), zap.String("protocol", "ws"))
		ctx = context.WithValue(ctx, sugarCtx, sugar)
		id, err := validateAndGetAccountIDString(initPayload.GetString("token"))
		if err != nil {
			sugar.Warnw("invalid token payload", zap.Error(err))
			return nil, err
		}

		account, err := client.GetAccount(ctx, id)
		if err != nil {
			sugar.Warnw("no account found", zap.Error(err))
			return nil, err
		}

		sugar.Info("serving authenticated account")
		ctx = context.WithValue(ctx, accountCtx, account)
		return ctx, nil
	}
	return transport.WebsocketInitFunc(fn)
}
