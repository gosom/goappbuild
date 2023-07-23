package api

import (
	"embed"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/restapi"
)

var _ restapi.RouteEngine = (*Router)(nil)

type Router struct {
	*restapi.ChiRouter

	swaggerController http.Handler
	swaggerPath       string

	healthController     HealthController
	userController       UserController
	projectController    ProjectController
	collectionController CollectionController
	queryController      QueryController
}

// NewRouter creates a new router.
func NewRouter(l *goappbuild.App, specFS embed.FS) (*Router, error) {
	swagCfg := restapi.SwaggerUIConfig{
		SpecName: "GoAppBuild API",
		SpecFile: "/docs/swagger.json",
		Path:     "/api/v1/docs",
		SpecFS:   specFS,
	}

	sw, err := restapi.NewSwaggerUI(&swagCfg)
	if err != nil {
		return nil, err
	}

	ans := Router{
		ChiRouter:            restapi.NewChiRouter(),
		swaggerController:    sw,
		swaggerPath:          swagCfg.Path,
		healthController:     NewHealthController(),
		userController:       NewUserController(l),
		projectController:    NewProjectController(l),
		collectionController: NewCollectionController(l),
		queryController:      NewQueryController(l),
		//idempotencyMiddleware: idempotencyMiddleware,
		//userContextMiddleware: NewUserContextMiddleware(),
	}

	return &ans, nil
}

// Engine returns the engine (an http.Handler) for the API.
func (router *Router) Engine() http.Handler {
	sp := router.swaggerPath
	if !strings.HasSuffix(sp, "/") {
		sp += "/"
	}

	sp += "*"

	router.R.Handle(sp, http.StripPrefix(router.swaggerPath, router.swaggerController))

	router.R.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", router.healthController.GetHealth)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", router.userController.Register)
		})

		r.Route("/projects", func(r chi.Router) {
			r.Post("/", router.projectController.Create)
		})

		r.Route("/collections", func(r chi.Router) {
			r.Post("/", router.collectionController.Create)
		})

		r.Route("/queries", func(r chi.Router) {
			r.Post("/{collectionName}", router.queryController.Create)
			r.Patch("/{collectionName}/{id}", router.queryController.Update)
			r.Delete("/{collectionName}/{id}", router.queryController.Delete)
			r.Get("/{collectionName}/{id}", router.queryController.Get)
		})
	})

	return router.R
}
