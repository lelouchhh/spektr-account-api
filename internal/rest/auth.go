package rest

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/llchhh/spektr-account-api/domain"
	"net/http"
)

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	Service AuthService
}

// AuthService defines the interface for authentication services.
type AuthService interface {
	Login(ctx context.Context, user domain.Auth) (string, error)
	RequestPasswordResetToken(ctx context.Context, login string) error
	UpdatePassword(ctx context.Context, token, password string) error
}

func NewAuthHandler(e *echo.Echo, svc AuthService) {
	handler := &AuthHandler{
		Service: svc, // Initialize the handler with the service
	}
	authGroup := e.Group("/api/v1/auth")
	authGroup.POST("/sign-in", handler.Login)
	authGroup.POST("/request-password-reset-token", handler.RequestPasswordResetToken)
	authGroup.POST("/update-password", handler.UpdatePassword)
}

// Login handles the /sign-in endpoint.
// @Summary Login a user
// @Description Logs the user in using their credentials (username and password)
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body domain.Auth true "Login credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} ResponseError "Invalid request payload"
// @Failure 401 {object} ResponseError "Invalid credentials"
// @Failure 403 {object} ResponseError "Account is locked"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /api/v1/auth/sign-in [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var auth domain.Auth
	// Bind the incoming JSON payload to the auth struct
	if err := c.Bind(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Message: "Invalid request payload", // Handle binding errors
		})
	}

	// Attempt to log in with the provided credentials
	token, err := h.Service.Login(c.Request().Context(), auth)
	if err != nil {
		// Handle errors from the service layer
		return handleError(c, err)
	}

	// Return the token on successful login
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// RequestPasswordResetToken handles the /request-password-reset-token endpoint.
// @Summary Request a password reset token
// @Description Requests a token to reset the user's password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body domain.Auth true "Login credentials"
// @Success 200 {object} map[string]string "Message"
// @Failure 400 {object} ResponseError "Invalid request payload"
// @Failure 401 {object} ResponseError "Invalid credentials"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /api/v1/auth/request-password-reset-token [post]
func (h *AuthHandler) RequestPasswordResetToken(c echo.Context) error {
	var auth domain.Auth
	// Bind the incoming JSON payload to the auth struct
	if err := c.Bind(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Message: "Invalid request payload", // Handle binding errors
		})
	}

	// Request the password reset token
	err := h.Service.RequestPasswordResetToken(c.Request().Context(), auth.Login)
	if err != nil {
		// Handle errors from the service layer
		return handleError(c, err)
	}

	// Return a success message after requesting the reset token
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password reset token has been sent to your email",
	})
}

// UpdatePassword handles the /update-password endpoint.
// @Summary Update the user's password
// @Description Update the password for a user using the reset token and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.Auth true "Password reset request"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} ResponseError "Invalid request payload"
// @Failure 401 {object} ResponseError "Invalid credentials"
// @Failure 403 {object} ResponseError "Account is locked"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /api/v1/auth/update-password [post]
func (h *AuthHandler) UpdatePassword(c echo.Context) error {
	var request domain.Auth
	// Bind the incoming JSON payload to the request struct
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Message: "Invalid request payload", // Handle binding errors
		})
	}

	// Attempt to update the password with the provided token and password
	err := h.Service.UpdatePassword(c.Request().Context(), request.Token, request.Password)
	if err != nil {
		// Handle errors from the service layer
		return handleError(c, err)
	}

	// Return a success message on successful password update
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password updated successfully",
	})
}

// handleError handles different error cases and returns the appropriate HTTP status code and response
func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: domain.ErrInvalidCredentials.Error(),
		})
	case errors.Is(err, domain.ErrAccountLocked):
		return c.JSON(http.StatusForbidden, ResponseError{
			Message: domain.ErrAccountLocked.Error(),
		})
	case errors.Is(err, domain.ErrSessionExpired):
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: domain.ErrSessionExpired.Error(),
		})
	case errors.Is(err, domain.ErrTooManyRequests):
		return c.JSON(http.StatusTooManyRequests, ResponseError{
			Message: domain.ErrTooManyRequests.Error(),
		})
	case errors.Is(err, domain.ErrUnauthorized):
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: domain.ErrUnauthorized.Error(),
		})
	default:
		// For any other unhandled errors, return 500 Internal Server Error
		return c.JSON(http.StatusInternalServerError, ResponseError{
			Message: domain.ErrInternalServerError.Error(),
		})
	}
}

// ResponseError is used to send error messages to the client
type ResponseError struct {
	Message string `json:"message"`
}
