package restapi

import "net/http"

// DefaultController is the default controller.
type DefaultController struct {
	Controller
}

// NewDefaultController creates a new default controller.
func NewDefaultController() DefaultController {
	return DefaultController{}
}

// NotFound is the handler for not found.
func (d DefaultController) NotFound(w http.ResponseWriter, r *http.Request) {
	d.Error(w, r, http.StatusNotFound, ErrNotFound)
}

// NotAllowed is the handler for not allowed.
func (d DefaultController) NotAllowed(w http.ResponseWriter, r *http.Request) {
	d.Error(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
}
