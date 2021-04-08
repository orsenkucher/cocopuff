package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/orsenkucher/cocopuff/account"
	"go.uber.org/zap"
)

const service = "account"

type Specification struct {
	Port int
	Dev  bool
	DSN  string
}

func main() {
	var spec Specification
	err := envconfig.Process(service, &spec)
	if err != nil {
		log.Fatalln(err)
	}

	var logger *zap.Logger
	if spec.Dev {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalln(err)
	}

	defer func() { _ = logger.Sync() }()
	sugar := logger.Sugar()

	if err := run(spec); err != nil {
		sugar.Fatal(err)
	}
}

func run(spec Specification) error {
	repo, err := account.NewPostgresRepository(spec.DSN)
	if err != nil {
		return err
	}

	serv := account.NewService(repo)

	return account.ListenGRPC(serv, spec.Port)
}
