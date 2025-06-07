package user

import (
	"context"
	"errors"

	userrepo "github.com/cryptonlx/crypto/src/repositories/user"

	"github.com/shopspring/decimal"
)

type Service struct {
	repo *userrepo.Repo
}

func New(repo *userrepo.Repo) *Service {
	return &Service{repo: repo}
}

func (s Service) GetUserWalletBalanceByUserName(ctx context.Context, username string) (userrepo.UserWallets, error) {
	if username == "" {
		return userrepo.UserWallets{}, errors.New("user id cannot be empty")
	}

	walletBalances, err := s.repo.UserWallets(ctx, username)
	if err != nil {
		return userrepo.UserWallets{}, err
	}
	return walletBalances, nil
}

func (s Service) GetUserTransactionsByUserName(ctx context.Context, username string) ([]userrepo.TransactionLedgers, error) {
	if username == "" {
		return []userrepo.TransactionLedgers{}, errors.New("user id cannot be empty")
	}

	transactions, err := s.repo.Transactions(ctx, username)
	if err != nil {
		return []userrepo.TransactionLedgers{}, err
	}
	return transactions, nil
}

func (s Service) CreateUser(ctx context.Context, username string) (userrepo.User, error) {
	if username == "" {
		return userrepo.User{}, errors.New("user name cannot be empty")
	}

	user, err := s.repo.CreateUser(ctx, username)
	if err != nil {
		return userrepo.User{}, err
	}
	return user, nil
}

func (s Service) CreateWallet(ctx context.Context, username string, _currency string) (userrepo.Wallet, error) {
	if username == "" {
		return userrepo.Wallet{}, errors.New("user name cannot be empty")
	}
	currency := userrepo.CurrencyType(_currency)

	wallet, err := s.repo.CreateWallet(ctx, username, currency)
	if err != nil {
		return userrepo.Wallet{}, err
	}
	return wallet, nil
}

func (s Service) Deposit(ctx context.Context, username string, nonce int64, walletId int64, amount decimal.Decimal) (userrepo.Transaction, userrepo.Ledger, error) {
	if !amount.IsPositive() {
		return userrepo.Transaction{}, userrepo.Ledger{}, errors.New("invalid_amount")
	}
	if nonce == 0 {
		return userrepo.Transaction{}, userrepo.Ledger{}, errors.New("invalid_nonce")
	}

	return s.repo.Deposit(username, ctx, nonce, walletId, amount)
}
