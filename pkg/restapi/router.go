package restapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RouteEngine is the interface for all routes.
type RouteEngine interface {
	Engine() http.Handler
}

// ChiRouter is a wrapper around chi.Router.
// wrap it in a struct and implement the RouteEngine interface.
type ChiRouter struct {
	// R is the chi.Router.
	R chi.Router
	RouteEngine
}

// NewChiRouter returns a new ChiRouter.
func NewChiRouter() *ChiRouter {
	ans := ChiRouter{
		R: chi.NewRouter(),
	}

	defaultController := NewDefaultController()

	ans.R.NotFound(defaultController.NotFound)
	ans.R.MethodNotAllowed(defaultController.NotAllowed)

	return &ans
}
