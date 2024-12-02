package auth

import (
	"context"
	"errors"
	"github.com/llchhh/spektr-account-api/domain"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
)

type AuthRepository interface {
	Login(ctx context.Context, user domain.Auth) (string, error)
	RequestPasswordResetToken(ctx context.Context, login string) error
	UpdatePassword(ctx context.Context, token, password string) error
}

type Service struct {
	authRepo AuthRepository
}

// RequestPasswordResetToken requests a password reset token for the user
func (s *Service) RequestPasswordResetToken(ctx context.Context, login string) error {
	// You might want to validate user input first, e.g., check if email or username exists
	if middleware.ContainsForbiddenChars(login) {
		return domain.ErrInvalidCredentials
	}
	err := s.authRepo.RequestPasswordResetToken(ctx, login)
	if err != nil {
		// Map repository errors to domain-specific errors
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return domain.ErrInvalidCredentials
		}
		if errors.Is(err, domain.ErrTooManyRequests) {
			return domain.ErrTooManyRequests
		}
		// Handle other errors appropriately
		return domain.ErrInternalServerError
	}
	return nil
}

// UpdatePassword updates the user's password using a reset token
func (s *Service) UpdatePassword(ctx context.Context, token, password string) error {
	// Validate the token and password here (e.g., check the token's expiration)
	if middleware.ContainsForbiddenChars(password) {
		return domain.ErrInvalidCredentials
	}
	if middleware.ContainsForbiddenChars(token) {
		return domain.ErrInvalidCredentials
	}
	err := s.authRepo.UpdatePassword(ctx, token, password)
	if err != nil {
		// Map repository errors to domain-specific errors
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return domain.ErrInvalidCredentials
		}
		if errors.Is(err, domain.ErrTooManyRequests) {
			return domain.ErrTooManyRequests
		}
		return domain.ErrInternalServerError
	}
	return nil
}

func NewService(a AuthRepository) *Service {
	return &Service{
		authRepo: a,
	}
}

func (s *Service) Login(ctx context.Context, user domain.Auth) (string, error) {
	// Marshal the user data into a JSON object for arg1
	if middleware.ContainsForbiddenChars(user.Login) {
		return "", domain.ErrInvalidCredentials
	}

	if middleware.ContainsForbiddenChars(user.Password) {
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
