package goappbuild

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a user.
type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate returns an error if the user is invalid.
func (u *User) Validate() error {
	return nil
}

// RegisterUserRequest represents a request to register a user.
type RegisterUserRequest struct {
}

// UserService represents a service for managing users.
type UserService interface {
	Register(context.Context, RegisterUserRequest) (User, error)
}

// UserRepo represents a repository for managing users.
type UserRepo interface {
	Create(context.Context, *User) error
}
