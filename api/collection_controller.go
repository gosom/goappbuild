package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/restapi"
)

// CollectionController is the controller for the collection resource.
type CollectionController struct {
	restapi.Controller

	app *goappbuild.App
}

// NewCollectionController creates a new collection controller.
func NewCollectionController(app *goappbuild.App) CollectionController {
	return CollectionController{
		app: app,
	}
}

// CreateCollectionRequest is the request for the CreateCollection method.
type CreateCollectionRequest struct {
	Name      string
	ProjectID uuid.UUID
}

// Validate ...
func (o *CreateCollectionRequest) Validate() error {
	return nil
}

// CreateCollectionResponse is the response for the CreateCollection method.
type CreateCollectionResponse struct {
	ID uuid.UUID
}

// Create create a collection
//
// @Summary Create a collection
// @Description Create a collection
// @Tags collections
// @Accept json
// @Produce json
// @Param body body CreateCollectionRequest true "The request body"
// @Success 200 {object} CreateCollectionResponse
// @Failure 400 {object} restapi.ErrorResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/collections [post]
func (o CollectionController) Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateCollectionRequest

	if err := o.DecodeBody(r, &payload); err != nil {
		o.Error(w, r, http.StatusBadRequest, err)
		return
	}

	cr := goappbuild.CollectionCreateRequest{
		Name:      payload.Name,
		ProjectID: payload.ProjectID,
	}

	c, err := o.app.Collections.Create(r.Context(), cr)
	if err != nil {
		o.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	ans := CreateCollectionResponse{
		ID: c.ID,
	}

	o.Success(w, r, http.StatusOK, ans)
}
