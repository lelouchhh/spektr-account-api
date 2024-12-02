package rest

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/llchhh/spektr-account-api/domain"
	"github.com/llchhh/spektr-account-api/repair"
	"net/http"
	"strings"
)

// RepairHandler handles repair-related requests.
type RepairHandler struct {
	Service repair.Service
}

// RepairService defines the interface for repair services.
type RepairService interface {
	CreateRepair(ctx context.Context, token string, repair domain.Repair) error
}

// NewRepairHandler initializes the repair handler with the given service and routes.
func NewRepairHandler(e *echo.Echo, svc repair.Service) {
	handler := &RepairHandler{
		Service: svc, // Initialize the handler with the service
	}
	repairGroup := e.Group("/api/v1/repairs")
	repairGroup.POST("", handler.CreateRepair) // Create a new repair request
}

// CreateRepair handles the request to create a new repair request.
// @Summary Create a new repair request
// @Description Submit a new repair request for the authenticated user
// @Tags Repairs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"  // Define the Authorization header with Bearer token
// @Param repair body domain.Repair true "Repair Request"  // Repair details
// @Success 201 {string} string "Repair request created successfully"
// @Failure 401 {object} ResponseError "Unauthorized"
// @Failure 400 {object} ResponseError "Bad request"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /api/v1/repairs [post]
func (h *RepairHandler) CreateRepair(c echo.Context) error {
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

	// Parse the repair request from the body
	var repair domain.Repair
	if err := c.Bind(&repair); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Message: "Invalid repair request format",
		})
	}

	// Call the service to create a new repair request
	err := h.Service.CreateRepair(c.Request().Context(), token, repair)
	if err != nil {
		return handleError(c, err)
	}

	// Return a success message
	return c.JSON(201, map[string]string{
		"message": "Repair created",
	})
}
