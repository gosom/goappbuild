package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/restapi"
)

type ProjectController struct {
	restapi.Controller

	app *goappbuild.App
}

func NewProjectController(app *goappbuild.App) ProjectController {
	return ProjectController{
		app: app,
	}
}

// CreateProjectRequest is the request for the CreateProject method.
type CreateProjectRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

// Validate validates the request.
func (o *CreateProjectRequest) Validate() error {
	return nil
}

// CreateProjectResponse is the response for the CreateProject method.
type CreateProjectResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Create create a project
//
// @Summary Create a project
// @Description Create a project
// @Tags projects
// @Accept json
// @Produce json
// @Param body body CreateProjectRequest true "The request body"
// @Success 200 {object} CreateProjectResponse
// @Failure 400 {object} restapi.ErrorResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/projects [post]
func (o ProjectController) Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateProjectRequest

	if err := o.DecodeBody(r, &payload); err != nil {
		o.Error(w, r, http.StatusBadRequest, err)
		return
	}

	createReq := goappbuild.CreateProjectRequest{
		UserID: payload.UserID,
		Name:   payload.Name,
	}

	project, err := o.app.Projects.Create(r.Context(), createReq)
	if err != nil {
		o.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	ans := CreateProjectResponse{
		ID:   project.ID,
		Name: project.Name,
	}

	o.Success(w, r, http.StatusOK, ans)
}
