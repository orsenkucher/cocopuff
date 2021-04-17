package main

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/orsenkucher/cocopuff/pub/care"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

func errorPresenter(sugar *zap.SugaredLogger) func(ctx context.Context, e error) *gqlerror.Error {
	return func(ctx context.Context, e error) *gqlerror.Error {
		gqlerr := graphql.DefaultErrorPresenter(ctx, e)
		err := gqlerr.Unwrap()
		sugar.Desugar().Warn("graphql error", care.ToZap(err))
		return gqlerr
	}
}

func recoverFn(sugar *zap.SugaredLogger) func(ctx context.Context, err interface{}) error {
	return func(ctx context.Context, err interface{}) error {
		sugar.Errorw("graphql panic", zap.Any("error", err))
		return errors.New("Internal server error!")
	}
}
