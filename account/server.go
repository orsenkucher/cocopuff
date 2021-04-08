//go:generate mkdir ./pb -p
//go:generate protoc ../api/proto/account.proto -I ../api/proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative

package account

import (
	"context"
	"fmt"
	"net"

	"github.com/orsenkucher/cocopuff/account/pb"
	"google.golang.org/grpc"
)

type AccountServiceServer struct {
	pb.UnimplementedAccountServiceServer
	service AccountService
}

var _ pb.AccountServiceServer = (*AccountServiceServer)(nil)

func ListenGRPC(service AccountService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)

	pb.RegisterAccountServiceServer(srv, &AccountServiceServer{service: service})

	return srv.Serve(lis)
}

func (s *AccountServiceServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	a, err := s.service.CreateAccount(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.Account{
		Id:   a.ID,
		Name: a.Name,
	}, nil
}

func (s *AccountServiceServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	a, err := s.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Account{
		Id:   a.ID,
		Name: a.Name,
	}, nil
}

func (s *AccountServiceServer) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	res, err := s.service.ListAccounts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}
	for _, p := range res {
		accounts = append(
			accounts,
			&pb.Account{
				Id:   p.ID,
				Name: p.Name,
			},
		)
	}

	return &pb.ListAccountsResponse{Accounts: accounts}, nil
}