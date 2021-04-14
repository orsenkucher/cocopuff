package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func (s *GraphQLServer) routes(graph *handler.Server) {
	s.router.Handle("/graphql", graph)
	s.router.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
}
