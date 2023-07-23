package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/restapi"
)

type QueryController struct {
	restapi.Controller

	app *goappbuild.App
}

func NewQueryController(app *goappbuild.App) QueryController {
	return QueryController{
		app: app,
	}
}

// Get returns a single document by ID
//
// @Summary Get a document
// @Description Get a document
// @Tags Queries
// @Accept json
// @Produce json
// @Param collectionName path string true "Collection Name"
// @Param id path string true "Document ID"
// @Param projectID header string true "Project ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} restapi.ErrorResponse
// @Failure 404 {object} restapi.ErrorResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/queries/{collectionName}/{id} [get]
func (o QueryController) Get(w http.ResponseWriter, r *http.Request) {
	projectID, err := o.getProjectID(r)
	if err != nil {
		o.Error(w, r, http.StatusBadRequest, err)

		return
	}

	project, err := o.app.Projects.Get(r.Context(), projectID)
	if err != nil {
		o.Error(w, r, http.StatusInternalServerError, fmt.Errorf("error getting project: %v", err))

		return
	}

	collectionName := o.StringURLParam(r, "collectionName")
	id := o.StringURLParam(r, "id")

	q := goappbuild.Q{}.
		Schema(project.Name).
		Table(collectionName).
		Equal("id", id)

	doc, err := o.app.Queries.Get(r.Context(), q)
	if err != nil {
		o.Error(w, r, http.StatusInternalServerError, err)

		return
	}

	o.Success(w, r, http.StatusOK, doc)
}

type CreatePayload map[string]any

func (o *CreatePayload) Validate() error {
	reserved := []string{"id", "created_at", "updated_at"}
	for _, r := range reserved {
		if _, ok := (*o)[r]; ok {
			return fmt.Errorf("reserved field: %s", r)
		}
	}

	return nil
}

// Create creates a new document
//
// @Summary Create a document
// @Description Create a document
// @Tags Queries
// @Accept json
// @Produce json
// @Param collectionName path string true "Collection Name"
// @Param projectID header string true "Project ID"
// @Param body body CreatePayload true "Document"
// @Success 201 {object} map[string]any
// @Failure 400 {object} restapi.ErrorResponse
// @Failure 404 {object} restapi.ErrorResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/queries/{collectionName} [post]
func (o QueryController) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := o.getProjectID(r)
	if err != nil {
		o.Error(w, r, http.StatusBadRequest, err)

		return
	}

	var payload CreatePayload

	if err := o.DecodeBody(r, &payload); err != nil {
		o.Error(w, r, http.StatusBadRequest, err)

		return
	}

	collecitonName := o.StringURLParam(r, "collectionName")

	payload["id"] = uuid.New()
	payload["created_at"] = time.Now().UTC()
	payload["updated_at"] = time.Now().UTC()

	ans, err := o.app.Queries.Create(r.Context(), projectID, collecitonName, payload)
	if err != nil {
		o.Error(w, r, http.StatusInternalServerError, err)

		return
	}

	o.Success(w, r, http.StatusCreated, ans)
}

// Update updates a document
//
// @Summary Update a document
// @Description Update a document
// @Tags Queries
// @Accept json
// @Produce json
// @Param collectionName path string true "Collection Name"
// @Param id path string true "Document ID"
// @Param projectID header string true "Project ID"
// @Param body body CreatePayload true "Document"
// @Success 200 {object} map[string]any
// @Failure 400 {object} restapi.ErrorResponse
// @Failure 404 {object} restapi.ErrorResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/queries/{collectionName}/{id} [patch]
func (o QueryController) Update(w http.ResponseWriter, r *http.Request) {
	projectID, err := o.getProjectID(r)
	if err != nil {
		o.Error(w, r, http.StatusBadRequest, err)

		return
	}

	var payload CreatePayload

	if err := o.DecodeBody(r, &payload); err != nil {
		o.Error(w, r, http.StatusBadRequest, err)

		return
	}

	collectionName := o.StringURLParam(r, "collectionName")

	payload["updated_at"] = time.Now().UTC()

	sid := o.StringURLParam(r, "id")
	if sid == "" {
		o.Error(w, r, http.StatusBadRequest, errors.New("id is required"))

		return

	}

	id, err := uuid.Parse(sid)
	if err != nil {
		o.Error(w, r, http.StatusBadRequest, fmt.Errorf("invalid id: %v", err))

		return
	}

	ans, err := o.app.Queries.Update(r.Context(), projectID, collectionName, id, payload)
	if err != nil {
		o.Error(w, r, http.StatusInternalServerError, err)

		return
	}

	o.Success(w, r, http.StatusOK, ans)
}

// Delete deletes a document
//
// @Summary Delete a document
// @Description Delete a document
// @Tags Queries
// @Accept json
// @Produce json
// @Param collectionName path string true "Collection Name"
// @Param id path string true "Document ID"
// @Param projectID header string true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} restapi.ErrorResponse
// @Failure 404 {object} restapi.ErrorResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/queries/{collectionName}/{id} [delete]
func (o QueryController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(o.StringURLParam(r, "id"))
	if err != nil {
		o.Error(w, r, http.StatusBadRequest, fmt.Errorf("invalid id: %v", err))

		return
	}

	projectID, err := o.getProjectID(r)
	if err != nil {
		o.Error(w, r, http.StatusBadRequest, err)

		return
	}

	collectionName := o.StringURLParam(r, "collectionName")

	if err := o.app.Queries.Delete(r.Context(), projectID, collectionName, id); err != nil {
		o.Error(w, r, http.StatusInternalServerError, err)

		return
	}

	o.Success(w, r, http.StatusNoContent, nil)
}

func (o QueryController) getProjectID(r *http.Request) (uuid.UUID, error) {
	sprojectID := o.HeaderKey(r, "projectID")
	if sprojectID == "" {
		return uuid.UUID{}, errors.New("missing projectID header")
	}

	projectID, err := uuid.Parse(sprojectID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return projectID, nil
}
