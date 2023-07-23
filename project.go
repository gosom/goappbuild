package goappbuild

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Project is a struct that represents a project
type Project struct {
	// ID is the unique identifier of the project
	ID uuid.UUID
	// UserID is the unique identifier of the user that owns the project
	UserID uuid.UUID
	// Name is the name of the project
	Name string
	// CreatedAt is the time the project was created
	CreatedAt time.Time
	// UpdatedAt is the time the project was last updated
	UpdatedAt time.Time
}

// SchemaName returns the name of the schema for the project
func (o *Project) SchemaName() string {
	return `"` + strings.Replace(o.Name, `"`, `""`, -1) + `"`
}

// CreateProjectRequest is the request for the CreateProject method.
type CreateProjectRequest struct {
	UserID uuid.UUID
	Name   string
}

// ProjectRepo is a repository for projects
type ProjectRepo interface {
	Create(context.Context, *Project) error
	Get(context.Context, uuid.UUID) (Project, error)
}

type ProjectService interface {
	Create(context.Context, CreateProjectRequest) (Project, error)
	Get(context.Context, uuid.UUID) (Project, error)
}
