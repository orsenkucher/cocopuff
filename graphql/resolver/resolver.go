package resolver

import (
	"github.com/orsenkucher/cocopuff/graphql"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your service, add any dependencies you require here.

type Resolver struct {
	sugar  *zap.SugaredLogger
	client *graphql.Client
}

func New(sugar *zap.SugaredLogger, client *graphql.Client) *Resolver {
	return &Resolver{sugar: sugar, client: client}
}
