package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/sqlext"
)

var _ goappbuild.CollectionRepo = (*collectionRepo)(nil)

type collectionRepo struct {
	conn sqlext.DBTX
}

// NewCollectionRepo returns a new instance of a collection repository
func NewCollectionRepo(conn sqlext.DBTX) goappbuild.CollectionRepo {
	return &collectionRepo{
		conn: conn,
	}
}

// Create creates a new collection
func (r *collectionRepo) Create(ctx context.Context, schema string, collection *goappbuild.Collection) error {
	const (
		q = `INSERT INTO collections
	(created_at, updated_at, name, project_id, attributes)
	VALUES
	((now() at time zone 'utc'), (now() at time zone 'utc'), $1, $2, $3)
	RETURNING id, created_at, updated_at, name, project_id`
	)

	attributesJson, err := json.Marshal(collection.Attributes)
	if err != nil {
		return err
	}

	dbCollection, err := sqlext.QueryRow[dbCollection](ctx, r.conn, q, collection.Name, collection.ProjectID, attributesJson)

	if err != nil {
		return err
	}

	collection.ID = dbCollection.ID
	collection.CreatedAt = dbCollection.CreatedAt
	collection.UpdatedAt = dbCollection.UpdatedAt

	return nil
}

type dbCollection struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ProjectID uuid.UUID
}

func (c *dbCollection) Bind() []any {
	return []any{
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.Name,
		&c.ProjectID,
	}
}

func (c *dbCollection) toModel() goappbuild.Collection {
	return goappbuild.Collection{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
		ProjectID: c.ProjectID,
	}
}
