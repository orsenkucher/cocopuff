//go:generate mkdir ./pb -p
//go:generate protoc ../api/proto/account.proto -I ../api/proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative

package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/orsenkucher/cocopuff/account/pkg"
)

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
	err := envconfig.Process("account", &spec)
	if err != nil {
		return err
	}

	r, err := pkg.NewPostgresRepository(spec.DSN)
	if err != nil {
		return err
	}

	s := pkg.NewService(r)

	return pkg.ListenGRPC(s, spec.Port)
}
