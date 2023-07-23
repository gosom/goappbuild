package restapi

import "net/http"

// Middleware is the middleware interface.
type Middleware interface {
	Handle(next http.Handler) http.Handler
}
