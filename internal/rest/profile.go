package rest

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/llchhh/spektr-account-api/domain"
	"net/http"
	"strings"
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
// @Summary Get user profile
// @Description Retrieve the profile of the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token (Bearer <token>)"
// @Success 200 {object} domain.Profile "User profile"
// @Failure 400 {object} ResponseError "Invalid request"
// @Failure 401 {object} ResponseError "Unauthorized"
// @Router /api/v1/profile [get]
func (h *ProfileHandler) Profile(c echo.Context) error {
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

	profile, err := h.Service.Profile(c.Request().Context(), token)
	if err != nil {
		return handleError(c, err)
	}

	// Return the profile data
	return c.JSON(http.StatusOK, profile)
}

// ChangePassword handles the POST /profile/change-password endpoint.
// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token (Bearer <token>)"
// @Param new_password body string true "New password"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} ResponseError "Invalid request payload"
// @Failure 401 {object} ResponseError "Unauthorized"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /api/v1/profile/change-password [post]
func (h *ProfileHandler) ChangePassword(c echo.Context) error {
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
	if token == "" {
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: "Invalid token",
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
// @Summary Change user email
// @Description Change the email for the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token (Bearer <token>)"
// @Param new_email body string true "New email"
// @Success 200 {object} map[string]string "Email changed successfully"
// @Failure 400 {object} ResponseError "Invalid request payload"
// @Failure 401 {object} ResponseError "Unauthorized"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /api/v1/profile/change-email [post]
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
