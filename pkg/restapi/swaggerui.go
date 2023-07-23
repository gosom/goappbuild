package restapi

import (
	"embed"
	"net/http"

	"github.com/ismurov/swaggerui"
)

// SwaggerUIConfig is the configuration for Swagger UI.
type SwaggerUIConfig struct {
	// SpecName is the name of the spec.
	SpecName string
	// SpecFile is the path to the spec file (relative to SpecFS).
	SpecFile string
	// SpecFS is the file system containing the spec file.
	SpecFS embed.FS
	// Path is the path to the Swagger UI.
	Path string
}

// NewSwaggerUI creates a new Swagger UI.
func NewSwaggerUI(cfg *SwaggerUIConfig) (http.Handler, error) {
	h, err := swaggerui.New(
		[]swaggerui.SpecFile{{
			Name: cfg.SpecName,
			Path: cfg.SpecFile,
		}},
		cfg.SpecFS,
	)

	return h, err
}
