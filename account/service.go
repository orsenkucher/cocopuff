package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type AccountService interface {
	CreateAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, ids []string) ([]Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

var _ AccountService = (*accountService)(nil)

type Account struct {
	ID   string `json:"id" gorm:"primarykey"`
	Name string `json:"name"`
}

type accountService struct {
	repository AccountRepository
}

func NewAccountService(r AccountRepository) *accountService {
	return &accountService{r}
}

func (s *accountService) CreateAccount(ctx context.Context, name string) (*Account, error) {
	a := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}

	if err := s.repository.CreateAccount(ctx, *a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *accountService) GetAccounts(ctx context.Context, ids []string) ([]Account, error) {
	return s.repository.GetAccounts(ctx, ids)
}

func (s *accountService) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.ListAccounts(ctx, skip, take)
}
