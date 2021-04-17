package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	super "github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/graphql/authentication"
	"github.com/orsenkucher/cocopuff/graphql/gql"
	"github.com/orsenkucher/cocopuff/pub/care"
	"go.uber.org/zap"
)

func New(sugar *zap.SugaredLogger, client *super.Client) gql.DirectiveRoot {
	return gql.DirectiveRoot{
		HasRole: func(ctx context.Context, obj interface{}, next graphql.Resolver, role gql.Role) (interface{}, error) {
			w := care.With(zap.String("package", "directive"), zap.String("directive", "HasRole"))
			sugar.Desugar().Info("", w.Fields...)

			if a, ok := authentication.For(ctx); !ok || !a.HasRole(role) {
				return nil, w.New("unauthorized role")
			}

			return next(ctx)
		},
	}
}
