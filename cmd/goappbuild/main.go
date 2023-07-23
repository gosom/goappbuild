package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/api"
	"github.com/gosom/goappbuild/collections"
	"github.com/gosom/goappbuild/pkg/cfgreader"
	"github.com/gosom/goappbuild/pkg/httpext"
	"github.com/gosom/goappbuild/pkg/restapi"
	"github.com/gosom/goappbuild/pkg/sqlext"
	"github.com/gosom/goappbuild/postgres"
	"github.com/gosom/goappbuild/projects"
	"github.com/gosom/goappbuild/queries"
	"github.com/gosom/goappbuild/users"
)

//go:generate swag i -g main.go --pd

//go:embed docs/swagger.json
var specFs embed.FS

// @title GoAppBuild API
// @version 0.0.1
// @description This is the API for the GoAppBuild application.

// @contact.name Giorgos Komninos
// @contact.url http://blog.gkomninos.com

// @host localhost:8080
// @BasePath /
// @accept json
// @produce json
// @query.collection.format multi
func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		panic(err)
	}

	os.Exit(0)
}
func run(ctx context.Context) error {
	cfg, err := cfgreader.NewConfig[Config]("")
	if err != nil {
		return err
	}

	db, err := sqlext.OpenPsqlConn(cfg.getDBConn())
	if err != nil {
		return err
	}

	storage := postgres.NewUnitOfWork(db)

	app := goappbuild.App{
		Users:       users.New(storage),
		Projects:    projects.New(storage),
		Collections: collections.New(storage),
		Queries:     queries.New(storage),
	}

	router, err := api.NewRouter(&app, specFs)
	if err != nil {
		return err
	}

	apiHn := restapi.NewAPI(router)

	serverCfg := httpext.ServerConfig{
		Host: cfg.ServerHost,
		Port: cfg.ServerPort,
	}

	httpServer := httpext.NewHTTPServer(serverCfg, apiHn)

	defer httpServer.Close(ctx)

	return httpServer.ListenAndServe(ctx)
}

// Config is the configuration for the application.
type Config struct {
	// PostgresHost is the host of the postgres database.
	PostgresHost string `envconfig:"POSTGRES_HOST" default:"localhost"`
	// PostgresPort is the port of the postgres database.
	PostgresPort int `envconfig:"POSTGRES_PORT" default:"5432"`
	// PostgresUser is the user of the postgres database.
	PostgresUser string `envconfig:"POSTGRES_USER" default:"postgres"`
	// PostgresPassword is the password of the postgres database.
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	// PostgresDatabase is the database of the postgres database.
	PostgresDB string `envconfig:"POSTGRES_DB" default:"postgres"`

	// ServerHost is the host of the server.
	ServerHost string `envconfig:"SERVER_HOST" default:"127.0.0.1"`
	// ServerPort is the port of the server.
	ServerPort int `envconfig:"SERVER_PORT" default:"8080"`
}

func (o *Config) getDBConn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		o.PostgresUser,
		o.PostgresPassword,
		o.PostgresHost,
		o.PostgresPort,
		o.PostgresDB,
	)
}
