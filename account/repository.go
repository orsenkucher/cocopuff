package account

import (
	"context"

	"github.com/orsenkucher/cocopuff/pub/care"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, ids []string) ([]Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

var _ AccountRepository = (*accountRepository)(nil)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(dsn string) (*accountRepository, error) {
	care := care.With(zap.String("function", "NewAccountRepository"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, care.Of(err, "fail gorm.Open")
	}

	err = db.AutoMigrate(&Account{})
	if err != nil {
		return nil, care.Of(err, "fail db.AutoMigrate")
	}

	return &accountRepository{db}, nil
}

func (r *accountRepository) CreateAccount(ctx context.Context, a Account) error {
	return r.db.Create(&a).Error
}

func (r *accountRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	var a *Account
	tx := r.db.First(a, id)
	return a, tx.Error
}

func (r *accountRepository) GetAccounts(ctx context.Context, ids []string) ([]Account, error) {
	var accounts []Account
	tx := r.db.Find(&accounts, ids)
	return accounts, tx.Error
}

func (r *accountRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	db, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(
		ctx,
		"SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		a := &Account{}
		if err = rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
