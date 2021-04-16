//go:generate gqlgen
//go:generate mkdir ./pb -p
//go:generate protoc ../api/proto/account.proto -I ../api/proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative

package graphql

import (
	"context"

	"github.com/orsenkucher/cocopuff/graphql/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TODO: should I remove Account
// and use pb.Account instead?
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	sugar   *zap.SugaredLogger
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(sugar *zap.SugaredLogger, url string) (*Client, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}

	c := pb.NewAccountServiceClient(conn)
	return &Client{sugar, conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.CreateAccount(
		ctx,
		&pb.CreateAccountRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Id,
		Name: r.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(
		ctx,
		&pb.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Id,
		Name: r.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, ids []string) ([]*Account, error) {
	r, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{Ids: ids},
	)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, &Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}

	return accounts, nil
}

func (c *Client) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error) {
	r, err := c.service.ListAccounts(
		ctx,
		&pb.ListAccountsRequest{
			Skip: skip,
			Take: take,
		},
	)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, &Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}

	return accounts, nil
}

func (c *Client) ListAccountsPB(ctx context.Context, req *pb.ListAccountsRequest) ([]*Account, error) {
	r, err := c.service.ListAccounts(ctx, req)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, &Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}

	return accounts, nil
}
