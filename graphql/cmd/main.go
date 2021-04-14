package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/graphql/dataloader"
	"github.com/orsenkucher/cocopuff/graphql/env"
	"github.com/orsenkucher/cocopuff/graphql/gql"
	"github.com/orsenkucher/cocopuff/graphql/log"
	"github.com/orsenkucher/cocopuff/graphql/resolver"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

const service = "graphql"

type specification struct {
	Port       int    `default:"9100"`
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`

	Release    bool
	Version    string `default:"v0.0.0"`
	Deployment string `default:"local"`
}

func main() {
	var spec specification
	err := envconfig.Process(service, &spec)
	if err != nil {
		log.Abortf("can't populate specification: %v", zap.Error(err))
	}

	sugar, err := log.New(service, spec.Deployment, spec.Version, spec.Release)
	if err != nil {
		log.Abortf("can't initialize zap logger: %v", zap.Error(err))
	}

	defer func() { _ = sugar.Sync() }()

	defer func() {
		if r := recover(); r != nil {
			sugar.Error("recovered error", zap.Any("description", r))
		}
	}()

	ctx := ctx(spec)
	_ = ctx // TODO:

	// TODO: move to run()
	// how about error logging in run?
	router := router(spec)

	client, err := graphql.NewClient(sugar, spec.AccountURL)
	if err != nil {
		sugar.Fatal("fail to dial:", zap.Error(err))
	}

	defer client.Close()

	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{
		Resolvers: resolver.NewResolver(sugar, client),
	}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Use(dataloader.Middleware(sugar, client))

	router.Handle("/graphql", srv)
	router.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	port := strconv.Itoa(spec.Port)
	sugar.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	sugar.Fatal(http.ListenAndServe(":"+port, router))
	// TODO: http.GracefulStop
	// TODO: move to server.go
}

func ctx(spec specification) context.Context {
	ctx := context.Background()
	ctx = env.With(ctx, service, spec.Deployment, spec.Version, spec.Release)
	return ctx
}

func router(spec specification) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		Debug:            !spec.Release,
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	return router
}
