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

func (s Service) GetUserWalletBalanceByUserName(ctx context.Context, username string) (userrepo.UserBalance, error) {
	if username == "" {
		return userrepo.UserBalance{}, errors.New("user id cannot be empty")
	}

	userBalance, err := s.repo.GetUserBalance(ctx, username)
	if err != nil {
		return userrepo.UserBalance{}, err
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
