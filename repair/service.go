package repair

import (
	"context"
	"github.com/llchhh/spektr-account-api/domain"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
	"log"
	"strings"
)

type RepairRepository interface {
	CreateRepair(ctx context.Context, token string, repair domain.Repair) error
}

type Service struct {
	repairRepo RepairRepository
}

// NewService creates a new RepairService instance with the provided RepairRepository.
func NewService(r RepairRepository) *Service {
	return &Service{
		repairRepo: r,
	}
}

// CreateRepair creates a new repair request using the provided token and repair details.
func (s *Service) CreateRepair(ctx context.Context, token string, repair domain.Repair) error {
	if middleware.ContainsForbiddenChars(repair.Text) {
		return domain.ErrInvalidToken
	}
	if middleware.ContainsForbiddenChars(repair.Subject) {
		return domain.ErrInvalidToken
	}
	if middleware.ContainsForbiddenChars(token) {
		return domain.ErrUnauthorized
	}

	log.Printf("Creating repair request for token: %s", token)

	err := s.repairRepo.CreateRepair(ctx, token, repair)
	if err != nil {
		log.Printf("Error changing password for token %s: %v", token, err)
		if err != nil {
			log.Printf("Error changing phone for token %s: %v", token, err)

			// Check if the error indicates an expired token
			if strings.Contains(err.Error(), "Необходимо авторизоваться") {
				log.Println("Token has expired or is invalid.")
				return domain.ErrSessionExpired
			}
			return err
		}
	}

	log.Printf("Repair created successfully for token: %s", token)
	return nil
}
