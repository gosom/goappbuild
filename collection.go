package goappbuild

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Collection is a struct that represents a collection of documents
type Collection struct {
	// ID is the unique identifier of the collection
	ID uuid.UUID
	// ProjectID is the unique identifier of the project that owns the collection
	ProjectID uuid.UUID
	// Name is the name of the collection
	Name string
	// Attributes holds the attributes of the document
	Attributes map[string]Attribute
	// CreatedAt is the time the collection was created
	CreatedAt time.Time
	// UpdatedAt is the time the collection was last updated
	UpdatedAt time.Time
}

// TableName returns the name of the table for the collection
func (o *Collection) TableName() string {
	return `"` + strings.Replace(o.Name, `"`, `""`, -1) + `"`
}

// CollectionRepo is the interface that wraps the basic CRUD operations for a collection
type CollectionRepo interface {
	Create(context.Context, string, *Collection) error
}

// CollectionCreateRequest is a struct that represents a request to create a collection
type CollectionCreateRequest struct {
	Name      string
	ProjectID uuid.UUID
}

// CollectionService is an interface that represents a service for managing collections
type CollectionService interface {
	Create(context.Context, CollectionCreateRequest) (Collection, error)
}
