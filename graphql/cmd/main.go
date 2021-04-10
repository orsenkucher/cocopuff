package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/orsenkucher/cocopuff/graphql/gql"
	"github.com/orsenkucher/cocopuff/graphql/resolver"
)

const service = "graphql"

type Specification struct {
	Port       int    `default:"9100"`
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
}

func main() {
	var spec Specification

	err := envconfig.Process(service, &spec)
	if err != nil {
		log.Fatalln(err)
	}

	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &resolver.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	port := strconv.Itoa(spec.Port)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
