package api

import (
	"net/http"

	"github.com/gosom/goappbuild/pkg/restapi"
)

// HealthController is the controller for health.
type HealthController struct {
	restapi.Controller
}

// NewHealthController returns a new HealthController.
func NewHealthController() HealthController {
	return HealthController{}
}

// HealthResponse is the health response.
type HealthResponse struct {
	// Status is the status of the service.
	Status string `json:"status"`
}

// GetHealth returns the health of the service.
//
// @Summary Get the health of the service
// @Description Get the health of the service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/health [get]
func (o HealthController) GetHealth(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status: "ok",
	}

	o.Success(w, r, http.StatusOK, resp)
}
