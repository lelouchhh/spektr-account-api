package auth

import (
	"context"
	"errors"
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
	token, err := s.authRepo.Login(ctx, user)
	if err != nil {
		// Map repository errors to domain-specific errors
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return "", domain.ErrInvalidCredentials
		}
		if errors.Is(err, domain.ErrAccountLocked) {
			return "", domain.ErrAccountLocked
		}
		// Handle other errors appropriately
		return "", domain.ErrInternalServerError
	}

	return token, nil
}
