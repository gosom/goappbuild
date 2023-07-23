package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/pkg/sqlext"
	"golang.org/x/exp/maps"
)

var _ goappbuild.QueryRepo = (*queryRepo)(nil)

type queryRepo struct {
	conn sqlext.DBTX
}

// NewQueryRepo returns a new instance of a postgres query service
func NewQueryRepo(conn sqlext.DBTX) goappbuild.QueryRepo {
	return &queryRepo{
		conn: conn,
	}
}

// Query executes a query against the database
func (o *queryRepo) Get(ctx context.Context, params goappbuild.Q) (map[string]any, error) {
	builder := NewPostgresQ(params)
	q, args, err := builder.Build()
	if err != nil {
		return nil, err
	}

	sql := o.wrapCte(q)

	var ans map[string]any
	var data []byte
	err = o.conn.QueryRowContext(ctx, sql, args...).Scan(&data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &ans)
	if err != nil {
		return nil, err
	}

	return ans, nil
}

func (o *queryRepo) Create(ctx context.Context, schema, collectionName string, data map[string]any) (map[string]any, error) {
	sb := strings.Builder{}

	sb.WriteString("INSERT INTO ")
	sb.WriteString(escape(schema))
	sb.WriteString(".")
	sb.WriteString(escape(collectionName))
	sb.WriteString(" (")

	keys := maps.Keys(data)
	sort.Strings(keys)

	for i, k := range keys {
		sb.WriteString(escape(k))
		if i < len(keys)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(") VALUES (")
	for i := range keys {
		sb.WriteString("$" + strconv.Itoa(i+1))
		if i < len(keys)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(") RETURNING *")

	args := make([]interface{}, len(data))
	for i, k := range keys {
		args[i] = data[k]
	}

	sql := o.wrapCte(sb.String())

	var result []byte

	err := o.conn.QueryRowContext(ctx, sql, args...).Scan(&result)
	if err != nil {
		return nil, err
	}

	var ans map[string]any

	err = json.Unmarshal(result, &ans)

	return ans, err
}

func (o *queryRepo) Update(ctx context.Context, schema, collectionName string, id uuid.UUID, data map[string]any) (map[string]any, error) {
	sb := strings.Builder{}

	sb.WriteString("UPDATE ")
	sb.WriteString(escape(schema))
	sb.WriteString(".")
	sb.WriteString(escape(collectionName))

	sb.WriteString(" SET ")

	keys := maps.Keys(data)
	sort.Strings(keys)

	args := make([]interface{}, len(data))
	for i, k := range keys {
		sb.WriteString(escape(k))
		sb.WriteString(" = ")
		sb.WriteString("$" + strconv.Itoa(i+1))
		if i < len(keys)-1 {
			sb.WriteString(", ")
		}

		args[i] = data[k]
	}

	sb.WriteString(" WHERE id = $")
	sb.WriteString(strconv.Itoa(len(args) + 1))

	args = append(args, id)

	sb.WriteString(" RETURNING *")

	sql := o.wrapCte(sb.String())

	var result []byte

	err := o.conn.QueryRowContext(ctx, sql, args...).Scan(&result)

	if err != nil {
		return nil, err
	}

	var ans map[string]any

	err = json.Unmarshal(result, &ans)

	return ans, err
}

func (o *queryRepo) Delete(ctx context.Context, schema, collectionName string, id uuid.UUID) error {
	sb := strings.Builder{}

	sb.WriteString("DELETE FROM ")
	sb.WriteString(escape(schema))
	sb.WriteString(".")
	sb.WriteString(escape(collectionName))
	sb.WriteString(" WHERE id = $1")

	res, err := o.conn.ExecContext(ctx, sb.String(), id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (o *queryRepo) wrapCte(q string) string {
	var sb strings.Builder

	sb.WriteString("WITH selection_cte AS (")
	sb.WriteString(q)
	sb.WriteString(") ")
	sb.WriteString("SELECT to_jsonb(selection_cte.*) as keyvals FROM selection_cte")

	return sb.String()
}
