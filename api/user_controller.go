package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/restapi"
)

// UserController is a controller for the user API.
type UserController struct {
	restapi.Controller

	app *goappbuild.App
}

// NewUserController creates a new user controller.
func NewUserController(app *goappbuild.App) UserController {
	return UserController{
		app: app,
	}
}

// RegisterUserRequest is the request for the RegisterUser method.
type RegisterUserRequest struct {
}

// Validate validates the request.
func (o *RegisterUserRequest) Validate() error {
	return nil
}

// RegisterUserResponse is the response for the RegisterUser method.
type RegisterUserResponse struct {
	ID uuid.UUID `json:"id"`
}

// Register register a user
//
// @Summary Register a user
// @Description Register a user
// @Tags users
// @Accept json
// @Produce json
// @Param body body RegisterUserRequest true "The request body"
// @Success 200 {object} RegisterUserResponse
// @Failure 400 {object} restapi.ErrorResponse
// @Failure 500 {object} restapi.ErrorResponse
// @Router /api/v1/users [post]
func (o UserController) Register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserRequest

	if err := o.DecodeBody(r, &payload); err != nil {
		o.Error(w, r, http.StatusBadRequest, err)
		return
	}

	registerReq := goappbuild.RegisterUserRequest{}

	u, err := o.app.Users.Register(r.Context(), registerReq)
	if err != nil {
		o.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	ans := RegisterUserResponse{
		ID: u.ID,
	}

	o.Success(w, r, http.StatusOK, ans)
}
