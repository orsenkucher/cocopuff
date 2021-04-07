package service

import (
	"context"

	"github.com/orsenkucher/cocopuff/account/pb"
)

type AccountingServer struct {
	pb.UnimplementedAccountingServer
}

var _ pb.AccountingServer = (*AccountingServer)(nil)

func (s *AccountingServer) PostAccount(ctx context.Context, req *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	return &pb.PostAccountResponse{}, nil
}

// PostAccount(context.Context, *PostAccountRequest) (*PostAccountResponse, error)
// GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
// GetAccounts(context.Context, *GetAccountsRequest) (*GetAccountsResponse, error)
