package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/orsenkucher/cocopuff/account"
)

const service = "account"

type Specification struct {
	Port int
	DSN  string
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var spec Specification
	err := envconfig.Process(service, &spec)
	if err != nil {
		return err
	}

	repo, err := account.NewPostgresRepository(spec.DSN)
	if err != nil {
		return err
	}

	serv := account.NewService(repo)

	return account.ListenGRPC(serv, spec.Port)
}
