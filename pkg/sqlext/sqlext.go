package sqlext

import (
	"context"
	"database/sql"

	// postgres driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Bindable is an interface that can be used to bind a struct to a sql query.
type Bindable[T any] interface {
	*T
	Bind() []any
}

// RowScanner is an interface that can be used to scan a row into a struct.
type RowScanner interface {
	Scan(dest ...interface{}) error
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// OpenPsqlConn opens a connection to a postgres database.
func OpenPsqlConn(dsn string) (conn *sql.DB, err error) {
	conn, err = sql.Open("pgx", dsn)
	if err != nil {
		return
	}

	err = conn.Ping()

	return
}

// Query returns a slice of items from the database.
func Query[T any, PT Bindable[T]](ctx context.Context, tx DBTX, q string, args ...any) ([]T, error) {
	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []T
	for rows.Next() {
		var item T
		var pt PT = &item
		if err := rows.Scan(pt.Bind()...); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// QueryRow returns a single item from the database.
func QueryRow[T any, PT Bindable[T]](ctx context.Context, tx DBTX, q string, args ...any) (T, error) {
	var item T
	var pt PT = &item
	row := tx.QueryRowContext(ctx, q, args...)
	if err := row.Scan(pt.Bind()...); err != nil {
		return item, err
	}
	return item, nil
}
