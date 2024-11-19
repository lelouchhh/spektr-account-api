package auth

import (
	"context"
	"github.com/llchhh/spektr-account-api/domain"
)

type AuthRepository interface {
	Login(ctx context.Context, user domain.Auth) (string, error)
}

type Service struct {
	authRepo AuthRepository
}

func NewService(a AuthRepository) *Service {
	return &Service{
		authRepo: a,
	}
}

func (s *Service) Login(ctx context.Context, user domain.Auth) (string, error) {
	return s.authRepo.Login(ctx, user)
}
