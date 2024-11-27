package auth

import (
	"context"
	"errors"
	"github.com/llchhh/spektr-account-api/domain"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
)

type AuthRepository interface {
	Login(ctx context.Context, user domain.Auth) (string, error)
	RequestPasswordResetToken(ctx context.Context, user domain.Auth) (string, error)
	UpdatePassword(ctx context.Context, token, password string) error
}

type Service struct {
	authRepo AuthRepository
}

func (s *Service) RequestPasswordResetToken(ctx context.Context, user domain.Auth) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) UpdatePassword(ctx context.Context, token, password string) error {
	//TODO implement me
	panic("implement me")
}

func NewService(a AuthRepository) *Service {
	return &Service{
		authRepo: a,
	}
}

func (s *Service) Login(ctx context.Context, user domain.Auth) (string, error) {
	// Marshal the user data into a JSON object for arg1
	if middleware.IsSuspicious(user.Login) {
		return "", domain.ErrInvalidCredentials
	}

	if middleware.IsSuspicious(user.Password) {
		return "", domain.ErrInvalidCredentials
	}
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
