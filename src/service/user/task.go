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

func (s Service) GetUserWalletBalanceByUserId(ctx context.Context, userId int64) (userrepo.UserBalance, error) {
	if userId == 0 {
		return userrepo.UserBalance{}, errors.New("user id cannot be zero")
	}

	userBalance, err := s.repo.GetUserBalance(ctx, userId)
	if err != nil {
		return userrepo.UserBalance{}, err
	}
	return userBalance, nil
}
