package rest

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/llchhh/spektr-account-api/domain"
	"net/http"
	"strings"
)

// NotificationHandler handles notification-related requests.
type NotificationHandler struct {
	Service NotificationService
}

// NotificationService defines the interface for notification services.
type NotificationService interface {
	GetNotifications(ctx context.Context, token string) ([]domain.Notification, error)
}

// NewNotificationHandler initializes the notification handler with the given service and routes.
func NewNotificationHandler(e *echo.Echo, svc NotificationService) {
	handler := &NotificationHandler{
		Service: svc, // Initialize the handler with the service
	}
	notificationGroup := e.Group("/api/v1/notifications")
	notificationGroup.GET("", handler.GetNotifications) // Retrieve notifications
}

// GetNotifications handles the request to get notifications for a user.
// @Summary Get user notifications
// @Description Retrieve notifications for the authenticated user
// @Tags Notifications
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"  // Define the Authorization header with Bearer token
// @Success 200 {array} domain.Notification "List of notifications"
// @Failure 401 {object} ResponseError "Unauthorized"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) GetNotifications(c echo.Context) error {
	// Extract the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: "Authorization token is required",
		})
	}

	// Check if the header starts with "Bearer "
	const bearerPrefix = "Bearer "

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))

	notifications, err := h.Service.GetNotifications(c.Request().Context(), token)
	if err != nil {
		return handleError(c, err)
	}

	// Return the notifications data
	return c.JSON(http.StatusOK, notifications)
}
