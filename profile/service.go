package profile

import (
	"context"
	"github.com/llchhh/spektr-account-api/domain"
)

type ProfileRepository interface {
	Profile(ctx context.Context, token string) (domain.Profile, error)
	ChangePassword(ctx context.Context, token string, password string) error
	ChangeEmail(ctx context.Context, token string, email string) error
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
	// Call the ProfileRepository's Profile method to fetch the profile data
	profile, err := s.profileRepo.Profile(ctx, token)
	if err != nil {
		return domain.Profile{}, err
	}
	return profile, nil
}

// ChangePassword updates the user's password using the provided token and new password.
func (s *Service) ChangePassword(ctx context.Context, token string, password string) error {
	// Call the ProfileRepository's ChangePassword method to update the password
	if err := s.profileRepo.ChangePassword(ctx, token, password); err != nil {
		return err
	}
	return nil
}

// ChangeEmail updates the user's email using the provided token and new email.
func (s *Service) ChangeEmail(ctx context.Context, token string, email string) error {
	// Call the ProfileRepository's ChangeEmail method to update the email
	if err := s.profileRepo.ChangeEmail(ctx, token, email); err != nil {
		return err
	}
	return nil
}
