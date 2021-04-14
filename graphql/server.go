package graphql

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/orsenkucher/cocopuff/graphql/gql"
	"github.com/orsenkucher/cocopuff/pub/ec"
	"go.uber.org/zap"
)

type GraphQLServer struct {
	sugar  *zap.SugaredLogger
	router *chi.Mux
}

func NewServer(
	sugar *zap.SugaredLogger,
	config gql.Config,
	upgrader websocket.Upgrader,
	cors func(http.Handler) http.Handler,
	middleware ...func(http.Handler) http.Handler,
) *GraphQLServer {
	server := handler.NewDefaultServer(gql.NewExecutableSchema(config))
	server.AddTransport(&transport.Websocket{Upgrader: upgrader})

	router := chi.NewRouter()
	router.Use(cors)
	router.Use(middleware...)

	router.Handle("/graphql", server)
	router.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	return &GraphQLServer{sugar: sugar, router: router}
}

func (s *GraphQLServer) ListenGraphQL(ctx context.Context, port int) <-chan error {
	return ec.Go(func() error {
		p := strconv.Itoa(port)

		server := &http.Server{
			Addr:    ":" + p,
			Handler: s.router,
		}

		s.sugar.Info("start graphql server",
			zap.String("address", fmt.Sprintf(":%s/graphql", p)),
			zap.String("playground", fmt.Sprintf(":%s/playground", p)),
		)
		select {
		case err := <-ec.Go(listen(server)):
			s.sugar.Info("shutdown graphql server", zap.String("by", "error"), zap.Error(err))
			return server.Shutdown(ctx)
		case <-ctx.Done():
			s.sugar.Info("shutdown graphql server", zap.String("by", "context.Done"))
			return server.Shutdown(ctx)
		}
	})
}

func listen(server *http.Server) func() error {
	return func() error {
		switch err := server.ListenAndServe(); err {
		case http.ErrServerClosed:
			return nil
		default:
			return err
		}
	}
}
