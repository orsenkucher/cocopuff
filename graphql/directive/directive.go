package directive

import (
	"context"

	graph "github.com/99designs/gqlgen/graphql"
	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/graphql/gql"
	"github.com/orsenkucher/cocopuff/pub/care"
	"go.uber.org/zap"
)

func New(sugar *zap.SugaredLogger, client *graphql.Client) gql.DirectiveRoot {
	return gql.DirectiveRoot{
		HasRole: func(ctx context.Context, obj interface{}, next graph.Resolver, role gql.Role) (interface{}, error) {
			w := care.With(zap.String("package", "directive"), zap.String("directive", "HasRole"))
			sugar.Desugar().Info("", w.Fields...)

			if !auth.UserFor(ctx).HasRole(role) {
				return nil, w.New("unauthorized role")
			}

			return next(ctx)
		},
	}
}

// TODO: Unmock
var auth = struct {
	UserFor func(ctx context.Context) *User
}{UserFor: func(ctx context.Context) *User {
	return &User{Name: "Orsen", IsAdmin: true}
}}

type User struct {
	Name    string
	IsAdmin bool
}

func (u *User) HasRole(role gql.Role) bool {
	if role.IsValid() {
		switch role {
		case gql.RoleAdmin:
			return u.IsAdmin
		case gql.RoleUser:
			return !u.IsAdmin
		}
	}
	return false
}
