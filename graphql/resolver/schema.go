package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/orsenkucher/cocopuff/graphql/gql"
)

func (r *mutationResolver) CreateAccount(ctx context.Context, account gql.AccountInput) (*gql.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product gql.ProductInput) (*gql.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateOrder(ctx context.Context, order gql.OrderInput) (*gql.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *gql.PaginationInput, id *string) ([]*gql.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context, pagination *gql.PaginationInput, query *string, id *string) ([]*gql.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
