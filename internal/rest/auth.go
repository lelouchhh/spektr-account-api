package rest

import (
	"context"
	"github.com/labstack/echo"
	"net/http"

	"github.com/llchhh/spektr-account-api/domain"
)

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	Service AuthService
}

// AuthService defines the interface for authentication services.
type AuthService interface {
	Login(ctx context.Context, user domain.Auth) (string, error)
}

func NewAuthHandler(e *echo.Echo, svc AuthService) {
	handler := &AuthHandler{
		Service: svc, // Initialize the handler with the service
	}
	authGroup := e.Group("/api/v1/auth")
	authGroup.POST("/sign-in", handler.Login)
}

// Login handles the /sign-in endpoint.
func (h *AuthHandler) Login(c echo.Context) error {
	var auth domain.Auth
	if err := c.Bind(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid request payload"})
	}
	token, err := h.Service.Login(c.Request().Context(), auth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

type ResponseError struct {
	Message string `json:"message"`
}
