package restapi

import "net/http"

// API is the main struct for the API.
type API struct {
	handler http.Handler
}

// NewAPI creates a new API.
func NewAPI(routes RouteEngine) *API {
	h := &API{
		handler: routes.Engine(),
	}

	return h
}

// ServeHTTP implements the http.Handler interface.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
