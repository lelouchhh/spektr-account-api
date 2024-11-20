package rest

import (
	"context"
	"github.com/labstack/echo"
	"github.com/llchhh/spektr-account-api/domain"
	"net/http"
)

// ProfileHandler handles profile-related requests.
type ProfileHandler struct {
	Service ProfileService
}

// ProfileService defines the interface for profile services.
type ProfileService interface {
	Profile(ctx context.Context, token string) (domain.Profile, error)
	ChangePassword(ctx context.Context, token string, newPassword string) error
	ChangeEmail(ctx context.Context, token string, newEmail string) error
}

// NewProfileHandler initializes the profile handler with the given service and routes.
func NewProfileHandler(e *echo.Echo, svc ProfileService) {
	handler := &ProfileHandler{
		Service: svc, // Initialize the handler with the service
	}
	profileGroup := e.Group("/api/v1/profile")
	profileGroup.GET("", handler.Profile)
	profileGroup.POST("/change-password", handler.ChangePassword)
	profileGroup.POST("/change-email", handler.ChangeEmail)
}

// GetProfile handles the GET /profile endpoint.
func (h *ProfileHandler) Profile(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: "Authorization token is required",
		})
	}

	// Fetch the profile using the token
	profile, err := h.Service.Profile(c.Request().Context(), token)
	if err != nil {
		return handleError(c, err)
	}

	// Return the profile data
	return c.JSON(http.StatusOK, profile)
}

// ChangePassword handles the POST /profile/change-password endpoint.
func (h *ProfileHandler) ChangePassword(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: "Authorization token is required",
		})
	}

	var payload struct {
		NewPassword string `json:"new_password"`
	}

	// Bind the request payload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Message: "Invalid request payload",
		})
	}

	// Attempt to change the password
	err := h.Service.ChangePassword(c.Request().Context(), token, payload.NewPassword)
	if err != nil {
		return handleError(c, err)
	}

	// Respond with success
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}

// ChangeEmail handles the POST /profile/change-email endpoint.
func (h *ProfileHandler) ChangeEmail(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: "Authorization token is required",
		})
	}

	var payload struct {
		NewEmail string `json:"new_email"`
	}

	// Bind the request payload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Message: "Invalid request payload",
		})
	}

	// Attempt to change the email
	err := h.Service.ChangeEmail(c.Request().Context(), token, payload.NewEmail)
	if err != nil {
		return handleError(c, err)
	}

	// Respond with success
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Email changed successfully",
	})
}
