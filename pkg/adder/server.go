//go:generate protoc -I ../../api/proto --go_out=../api --go_opt=paths=source_relative --go-grpc_out=../api --go-grpc_opt=paths=source_relative ../../api/proto/adder.proto
// problems: api dir has to exist prior to generation
// very long script
package adder

import (
	"context"

	"github.com/orsenkucher/cocopuff/pkg/api"
)

type AdderServer struct {
	api.UnimplementedAdderServer
}

var _ api.AdderServer = (*AdderServer)(nil)

func (s *AdderServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	return &api.AddResponse{Result: req.GetX() + req.GetY()}, nil
}
