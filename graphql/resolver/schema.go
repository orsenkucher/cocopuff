package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/orsenkucher/cocopuff/graphql/dataloader"
	"github.com/orsenkucher/cocopuff/graphql/gql"
	"github.com/orsenkucher/cocopuff/graphql/pb"
	"go.uber.org/zap"
)

func (r *accountResolver) Orders(ctx context.Context, obj *gql.Account) ([]*gql.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

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
	// TODO: how to pass context to grpc call
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()

	sugar := r.sugar.With(
		zap.String("package", "resolver"),
		zap.String("type", fmt.Sprintf("%T", r)),
		zap.String("method", "Accounts"),
	)

	// Get single
	if id != nil {
		a, err := dataloader.For(ctx).AccountById.Load(*id)
		if err != nil {
			sugar.Errorw("failed get by id:", zap.Error(err))
			return nil, err
		}
		return []*gql.Account{{
			ID:   a.ID,
			Name: a.Name,
		}}, nil
	}

	// Else get page
	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.Bounds()
	}

	// TODO: use dataloader?
	// I probably can't
	// accountList, err := r.client.ListAccounts(ctx, skip, take)
	accountList, err := dataloader.For(ctx).AccountPaginated.Load(&pb.ListAccountsRequest{
		Skip: skip, Take: take,
	})
	if err != nil {
		sugar.Errorw("failed to list:", zap.Error(err))
		return nil, err
	}

	var accounts []*gql.Account
	for _, a := range accountList {
		account := &gql.Account{
			ID:   a.ID,
			Name: a.Name,
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *gql.PaginationInput, query *string, id *string) ([]*gql.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

// Account returns gql.AccountResolver implementation.
func (r *Resolver) Account() gql.AccountResolver { return &accountResolver{r} }

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type accountResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
