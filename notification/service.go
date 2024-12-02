package notification

import (
	"context"
	"github.com/llchhh/spektr-account-api/domain"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
	"log"
)

type NotificationRepository interface {
	GetNotifications(ctx context.Context, token string) ([]domain.Notification, error)
}

type Service struct {
	notificationRepo NotificationRepository
}

// NewService creates a new Service instance with the provided NotificationRepository.
func NewService(n NotificationRepository) *Service {
	return &Service{
		notificationRepo: n,
	}
}

// GetNotifications retrieves a list of notifications for the user using the provided token.
func (s *Service) GetNotifications(ctx context.Context, token string) ([]domain.Notification, error) {
	if token == "" {
		log.Println("GetNotifications request failed: missing authorization token")
		return nil, domain.ErrInvalidToken
	}
	if middleware.ContainsForbiddenChars(token) {
		return nil, domain.ErrInvalidToken
	}
	log.Printf("Fetching notifications for token: %s", token)

	notifications, err := s.notificationRepo.GetNotifications(ctx, token)
	if err != nil {
		log.Printf("Error fetching notifications for token %s: %v", token, err)
		return nil, err
	}

	log.Printf("Successfully fetched notifications for token: %s", token)
	return notifications, nil
}
