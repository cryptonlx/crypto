package user

import (
	"context"
	"errors"
	userrepo "github.com/cryptonlx/crypto/src/repositories/user"
)

type Service struct {
	repo *userrepo.Repo
}

func New(repo *userrepo.Repo) *Service {
	return &Service{repo: repo}
}

func (s Service) GetUserWalletBalanceByUserName(ctx context.Context, username string) (userrepo.WalletBalances, error) {
	if username == "" {
		return userrepo.WalletBalances{}, errors.New("user id cannot be empty")
	}

	userBalance, err := s.repo.WalletBalances(ctx, username)
	if err != nil {
		return userrepo.WalletBalances{}, err
	}
	return userBalance, nil
}

func (s Service) GetUserTransactionsByUserName(ctx context.Context, username string) (userrepo.UserTransactions, error) {
	if username == "" {
		return userrepo.UserTransactions{}, errors.New("user id cannot be empty")
	}

	userBalance, err := s.repo.Transactions(ctx, username)
	if err != nil {
		return userrepo.UserTransactions{}, err
	}
	return userBalance, nil
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
