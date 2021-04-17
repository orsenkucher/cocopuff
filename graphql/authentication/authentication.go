package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/pub/log"
	"go.uber.org/zap"
)

type ctxKey int

const (
	sugarCtx ctxKey = iota
	accountCtx
)

func Middleware(sugar *zap.SugaredLogger, ja *jwtauth.JWTAuth, client *graphql.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if log, ok := log.For(ctx); ok {
				sugar = log.Sugar()
			}

			sugar := sugar.With(zap.String("package", "authentication"), zap.String("protocol", "http"))
			ctx = context.WithValue(ctx, sugarCtx, sugar)

			token, claims, err := jwtauth.FromContext(ctx)

			if err != nil || token == nil {
				sugar.Info("serving with no token")
				next.ServeHTTP(w, r)
				return
			}

			if err := jwt.Validate(token); err != nil {
				sugar.Infow("invalid token", zap.Error(err))
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}

			id, ok := claims["user_id"].(string)
			if !ok {
				sugar.Infow("claims no user id", zap.Error(err))
				http.Error(w, "claims no user id", http.StatusForbidden)
				return
			}

			account, err := client.GetAccount(ctx, id)
			if err != nil {
				sugar.Warnw("no account found", zap.Error(err))
				http.Error(w, "no account", http.StatusForbidden)
				return
			}

			_, tokenString, _ := ja.Encode(map[string]interface{}{"user_id": id})
			c := &http.Cookie{
				Name:  "jwt",
				Value: tokenString,
			}
			http.SetCookie(w, c)
			sugar.Info("serving authenticated account")
			ctx = context.WithValue(ctx, accountCtx, account)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func WebsocketMiddleware(sugar *zap.SugaredLogger, ja *jwtauth.JWTAuth, client *graphql.Client) transport.WebsocketInitFunc {
	fn := func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
		sugar = sugar.With(zap.String("package", "authentication"), zap.String("protocol", "ws"))
		ctx = context.WithValue(ctx, sugarCtx, sugar)
		token, claims, err := getClaims(jwtauth.VerifyToken(ja, initPayload.GetString("token")))

		if err != nil {
			sugar.Warnw("invalid token payload", zap.Error(err))
			return nil, err
		}

		if err := jwt.Validate(token); err != nil {
			sugar.Warnw("invalid token", zap.Error(err))
			return nil, err
		}

		id, ok := claims["user_id"].(string)
		if !ok {
			err := errors.New("claims no user id")
			sugar.Warn(err)
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

func getClaims(token jwt.Token, err error) (jwt.Token, map[string]interface{}, error) {
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return token, nil, err
		}
	} else {
		claims = map[string]interface{}{}
	}

	return token, claims, err
}

func For(ctx context.Context) (*graphql.Account, bool) {
	a, ok := ctx.Value(accountCtx).(*graphql.Account)
	if !ok {
		if sugar, ok := ctx.Value(sugarCtx).(*zap.SugaredLogger); ok {
			sugar.Infow("fail to retrieve account", zap.String("function", "For"))
		}
	}

	return a, ok
}
