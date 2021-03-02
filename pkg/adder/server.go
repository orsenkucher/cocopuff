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
