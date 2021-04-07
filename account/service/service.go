package service

import (
	"context"

	"github.com/orsenkucher/cocopuff/account/pb"
)

type AccountServiceServer struct {
	pb.UnimplementedAccountServiceServer
}

var _ pb.AccountServiceServer = (*AccountServiceServer)(nil)

func (s *AccountServiceServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
}

func (s *AccountServiceServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
}

func (s *AccountServiceServer) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
}
