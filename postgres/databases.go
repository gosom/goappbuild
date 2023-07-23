package postgres

import (
	"context"
	"fmt"

	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/sqlext"
)

type dbRepo struct {
	conn sqlext.DBTX
}

// NewDBRepo returns a new instance of a postgres database repository
func NewDBRepo(conn sqlext.DBTX) goappbuild.DatabaseRepo {
	return &dbRepo{
		conn: conn,
	}
}

// Create creates a new schema in the database
func (o *dbRepo) CreateSchema(ctx context.Context, name string) error {
	schema := fmt.Sprintf("CREATE schema %s", name)

	_, err := o.conn.ExecContext(ctx, schema)
	if err != nil {
		return err
	}

	return nil
}

func (o *dbRepo) CreateTable(ctx context.Context, schema, name string) error {
	tableQ := createTableStmt(createTableParams{
		schema: schema,
		table:  name,
	})

	_, err := o.conn.ExecContext(ctx, tableQ)
	if err != nil {
		return err
	}

	return nil
}

func (o *dbRepo) CreateColumns(ctx context.Context, schema, table string, attributes map[string]goappbuild.Attribute) error {
	attributesQ := make([]string, 0, len(attributes))
	indexesQ := make([]string, 0, len(attributes))
	for _, attr := range attributes {
		attParam := addAttributeParams{
			schema:    schema,
			table:     table,
			attribute: attr,
		}
		alter, err := addAttributeStmt(attParam)
		if err != nil {
			return err
		}

		attributesQ = append(attributesQ, alter)

		indexQ := addIndexesStmt(attParam)
		if indexQ != "" {
			indexesQ = append(indexesQ, indexQ)
		}
	}

	for i := range attributesQ {
		_, err := o.conn.ExecContext(ctx, attributesQ[i])
		if err != nil {
			return err
		}
	}

	for i := range indexesQ {
		_, err := o.conn.ExecContext(ctx, indexesQ[i])
		if err != nil {
			return err
		}
	}

	return nil
}
