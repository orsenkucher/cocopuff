package main

import (
	"context"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/orsenkucher/cocopuff/graphql"
	"github.com/orsenkucher/cocopuff/graphql/dataloader"
	"github.com/orsenkucher/cocopuff/graphql/env"
	"github.com/orsenkucher/cocopuff/graphql/gql"
	"github.com/orsenkucher/cocopuff/graphql/log"
	"github.com/orsenkucher/cocopuff/graphql/resolver"
	"github.com/orsenkucher/cocopuff/pub/care"
	"github.com/orsenkucher/cocopuff/pub/gs"
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
			sugar.Errorw("recovered error", zap.Any("description", r))
		}
	}()

	ctx := ctx(sugar, spec)
	if err := run(ctx, sugar, spec); err != nil {
		sugar.Fatalw("run failed", care.ToZap(err))
	}
}

func ctx(sugar *zap.SugaredLogger, spec specification) context.Context {
	ctx := gs.With(context.Background())
	ctx = env.With(ctx, sugar, service, spec.Deployment, spec.Version, spec.Release)
	return ctx
}

func run(ctx context.Context, sugar *zap.SugaredLogger, spec specification) error {
	cors := corsMiddleware(spec)
	upgrader := websocketUpgrader()

	client, err := graphql.NewClient(sugar, spec.AccountURL)
	if err != nil {
		return care.Of(err, "fail to dial grpc client", zap.String("function", "run"))
	}

	defer client.Close()

	config := gql.Config{
		Resolvers: resolver.NewResolver(sugar, client),
	}

	dataloader := dataloader.Middleware(sugar, client)

	middleware := []func(http.Handler) http.Handler{dataloader}
	server := graphql.NewServer(sugar, config, upgrader, cors, middleware...)
	return <-server.ListenGraphQL(ctx, spec.Port)
}

func corsMiddleware(spec specification) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		Debug:            !spec.Release,
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func websocketUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
}
