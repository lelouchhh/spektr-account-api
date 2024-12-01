package profile

import (
	"context"
	"errors"
	"fmt"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
	"log"
	"strings"

	"github.com/llchhh/spektr-account-api/domain"
)

type ProfileRepository interface {
	Profile(ctx context.Context, token string) (domain.Profile, error)
	ChangePassword(ctx context.Context, token string, password string) error
	ChangeEmail(ctx context.Context, token string, email string) error
	ChangePhone(ctx context.Context, token string, newPhone string) error
}

type Service struct {
	profileRepo ProfileRepository
}

// NewService creates a new Service instance with the provided ProfileRepository.
func NewService(p ProfileRepository) *Service {
	return &Service{
		profileRepo: p,
	}
}

// Profile retrieves the user's profile using the provided token.
func (s *Service) Profile(ctx context.Context, token string) (domain.Profile, error) {
	if token == "" {
		log.Println("Profile request failed: missing authorization token")
		return domain.Profile{}, domain.ErrInvalidToken
	}
	if middleware.IsSuspicious(token) {
		return domain.Profile{}, domain.ErrInvalidToken
	}
	log.Printf("Fetching profile for token: %s", token)

	profile, err := s.profileRepo.Profile(ctx, token)
	if err != nil {
		log.Printf("Error fetching profile for token %s: %v", token, err)
		return domain.Profile{}, err
	}

	log.Printf("Successfully fetched profile for token: %s", token)
	return profile, nil
}

// ChangePassword updates the user's password using the provided token and new password.
func (s *Service) ChangePassword(ctx context.Context, token string, password string) error {

	if token == "" {
		log.Println("ChangePassword request failed: missing authorization token")
		return domain.ErrUnauthorized
	}
	err := middleware.ValidatePassword(password)
	if err != nil {
		log.Println("Invalid phone number format detected")
		return domain.ErrInvalidCredentials
	}
	if middleware.ContainsForbiddenChars(token) {
		log.Println("Invalid token format detected")
		return domain.ErrInvalidToken
	}
	log.Printf("Changing password for token: %s", token)

	err = s.profileRepo.ChangePassword(ctx, token, password)
	if err != nil {
		log.Printf("Error changing password for token %s: %v", token, err)

		// Check if the error indicates an expired token
		if errors.Is(err, domain.ErrSessionExpired) {
			return domain.ErrSessionExpired
		} else {
			return domain.ErrInternalServerError
		}
	}

	log.Printf("Password successfully changed for token: %s", token)
	return nil
}

// ChangeEmail updates the user's email using the provided token and new email.
func (s *Service) ChangeEmail(ctx context.Context, token string, email string) error {
	if token == "" {
		log.Println("ChangeEmail request failed: missing authorization token")
		return fmt.Errorf("authorization token is required")
	}
	if email == "" {
		log.Println("ChangeEmail request failed: missing email")
		return fmt.Errorf("email cannot be empty")
	}
	if middleware.ContainsForbiddenChars(email) {
		log.Println("Invalid phone number format detected")
		return domain.ErrInvalidCredentials
	}
	if middleware.ContainsForbiddenChars(token) {
		log.Println("Invalid token format detected")
		return domain.ErrInvalidToken
	}
	log.Printf("Changing email for token: %s", token)

	err := s.profileRepo.ChangeEmail(ctx, token, email)
	if err != nil {
		log.Printf("Error changing email for token %s: %v", token, err)

		// Check if the error indicates an expired token
		if strings.Contains(err.Error(), "invalid token") {
			log.Println("Token has expired or is invalid.")
			return domain.ErrSessionExpired
		}
		return err
	}

	log.Printf("Email successfully changed for token: %s", token)
	return nil
}

// ChangePhone updates the user's phone using the provided token and new email.
func (s *Service) ChangePhone(ctx context.Context, token string, phone string) error {
	if token == "" {
		log.Println("ChangePhone request failed: missing authorization token")
		return fmt.Errorf("authorization token is required")
	}
	if phone == "" {
		log.Println("ChangePhone request failed: missing email")
		return fmt.Errorf("phone cannot be empty")
	}
	if middleware.ContainsForbiddenChars(phone) {
		log.Println("Invalid phone number format detected")
		return domain.ErrInvalidCredentials
	}
	if middleware.ContainsForbiddenChars(token) {
		log.Println("Invalid token format detected")
		return domain.ErrInvalidToken
	}
	log.Printf("Changing phone for token: %s", token)

	err := s.profileRepo.ChangePhone(ctx, token, phone)
	if err != nil {
		log.Printf("Error changing phone for token %s: %v", token, err)

		// Check if the error indicates an expired token
		if strings.Contains(err.Error(), "invalid token") {
			log.Println("Token has expired or is invalid.")
			return domain.ErrSessionExpired
		}
		return err
	}

	log.Printf("phone successfully changed for token: %s", token)
	return nil
}
