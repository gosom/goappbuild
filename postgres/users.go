package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/sqlext"
)

var _ goappbuild.UserRepo = (*userRepo)(nil)

type userRepo struct {
	conn sqlext.DBTX
}

// NewUserRepo returns a new instance of a postgres user repository
func NewUserRepo(conn sqlext.DBTX) goappbuild.UserRepo {
	return &userRepo{
		conn: conn,
	}
}

// Create creates a new user in the database
func (o *userRepo) Create(ctx context.Context, u *goappbuild.User) error {
	const q = `INSERT INTO users
		(created_at, updated_at)
		VALUES ((NOW() at time zone 'utc'), (NOW() at time zone 'utc'))
		RETURNING id, created_at, updated_at`

	dbu, err := sqlext.QueryRow[dbUser](ctx, o.conn, q)
	if err != nil {
		return err
	}

	u.ID = dbu.ID
	u.CreatedAt = dbu.CreatedAt
	u.UpdatedAt = dbu.UpdatedAt

	return nil
}

type dbUser struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o *dbUser) Bind() []any {
	return []any{
		&o.ID,
		&o.CreatedAt,
		&o.UpdatedAt,
	}
}

func (o *dbUser) toModel() goappbuild.User {
	return goappbuild.User{
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
