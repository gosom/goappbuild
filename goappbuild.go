package goappbuild

import "context"

var (
	// Version is the current version of the app.
	Version string
	// Commit is the current git commit of the app.
	Commit string
)

// App is a struct that represents the app
type App struct {
	Users       UserService
	Projects    ProjectService
	Collections CollectionService
	Queries     QueryService
}

// Storage  is a struct that represents the unit of work
type Storage interface {
	New(ctx context.Context) (Storage, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Users() UserRepo
	Projects() ProjectRepo
	Collections() CollectionRepo
	Databases() DatabaseRepo
	Queries() QueryRepo
}
