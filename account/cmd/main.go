package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/orsenkucher/cocopuff/account"
	"github.com/orsenkucher/cocopuff/account/env"
	"github.com/orsenkucher/cocopuff/account/log"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

const service = "account"

type specification struct {
	Port int `default:"9100"`
	DSN  string

	Release    bool
	Version    string `default:"v0.0.0"`
	Deployment string `default:"local"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

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
	if err := run(ctx, sugar, spec); err != nil {
		sugar.Fatal(err)
	}
}

func ctx(spec specification) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, env.Service, service)
	ctx = context.WithValue(ctx, env.Version, spec.Version)
	ctx = context.WithValue(ctx, env.Release, spec.Release)
	ctx = context.WithValue(ctx, env.Deployment, spec.Deployment)
	ctx = context.WithValue(ctx, env.Tags, []string{spec.Deployment, spec.Version})
	return ctx
}

func run(ctx context.Context, sugar *zap.SugaredLogger, spec specification) error {
	repo, err := account.NewAccountRepository(spec.DSN)
	if err != nil {
		return err
	}

	serv := account.NewService(repo)

	return account.ListenGRPC(serv, spec.Port)
}
