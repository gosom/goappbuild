package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/sqlext"
)

type projectRepo struct {
	conn sqlext.DBTX
}

// NewProjectRepo returns a new instance of a postgres project repository
func NewProjectRepo(conn sqlext.DBTX) goappbuild.ProjectRepo {
	return &projectRepo{
		conn: conn,
	}
}

// Create creates a new project in the database
func (o *projectRepo) Create(ctx context.Context, p *goappbuild.Project) error {
	const q = `INSERT INTO projects
		(created_at, updated_at, name, user_id)
		VALUES
		((NOW() at time zone 'utc'), (NOW() at time zone 'utc'), $1, $2)
		RETURNING id, created_at, updated_at, name, user_id`

	dbp, err := sqlext.QueryRow[dbProject](ctx, o.conn, q, p.Name, p.UserID)
	if err != nil {
		return err
	}

	p.ID = dbp.ID
	p.CreatedAt = dbp.CreatedAt
	p.UpdatedAt = dbp.UpdatedAt

	return nil
}

func (o *projectRepo) Get(ctx context.Context, id uuid.UUID) (goappbuild.Project, error) {
	const q = `SELECT 
			id, created_at, updated_at, name, user_id
		FROM projects
		WHERE id = $1`

	dbp, err := sqlext.QueryRow[dbProject](ctx, o.conn, q, id)
	if err != nil {
		return goappbuild.Project{}, err
	}

	return dbp.toModel(), nil
}

type dbProject struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	UserID    uuid.UUID
}

func (o *dbProject) Bind() []any {
	return []any{
		&o.ID,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.Name,
		&o.UserID,
	}
}

func (o *dbProject) toModel() goappbuild.Project {
	return goappbuild.Project{
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Name:      o.Name,
		UserID:    o.UserID,
	}
}
