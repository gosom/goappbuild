package goappbuild

import (
	"context"

	"github.com/google/uuid"
)

type QueryRepo interface {
	Get(context.Context, Q) (map[string]any, error)
	Create(ctx context.Context, schema, table string, data map[string]any) (map[string]any, error)
	Update(ctx context.Context, schema, table string, id uuid.UUID, data map[string]any) (map[string]any, error)
	Delete(ctx context.Context, schema, table string, id uuid.UUID) error
}

// QueryService is the interface that provides the Query
type QueryService interface {
	Get(context.Context, Q) (Document, error)
	Create(context.Context, uuid.UUID, string, map[string]any) (Document, error)
	Update(context.Context, uuid.UUID, string, uuid.UUID, map[string]any) (Document, error)
	Delete(context.Context, uuid.UUID, string, uuid.UUID) error
}

type Q struct {
	schema string
	table  string
	cols   []string
	where  []Op
}

func (q Q) GetSchema() string {
	return q.schema
}

func (q Q) Schema(schema string) Q {
	q.schema = schema

	return q
}

func (q Q) GetTable() string {
	return q.table
}

func (q Q) Table(table string) Q {
	q.table = table

	return q
}

func (q Q) Cols() []string {
	return q.cols
}

func (q Q) Where() []Op {
	return q.where
}

func (q Q) Select(cols ...string) Q {
	existings := make(map[string]bool)
	for _, col := range q.cols {
		existings[col] = true
	}

	for _, col := range cols {
		if _, ok := existings[col]; !ok {
			q.cols = append(q.cols, col)
		}
	}

	return q
}

func (q Q) Equal(column string, value any) Q {
	op := Op{
		column: column,
		op:     OpEq,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) NotEqual(column string, value any) Q {
	op := Op{
		column: column,
		op:     OpNeq,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) LessThan(column string, value any) Q {
	op := Op{
		column: column,
		op:     OpLt,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) LessThanOrEqual(column string, value any) Q {
	op := Op{
		column: column,
		op:     OpLte,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) GreaterThan(column string, value any) Q {
	op := Op{
		column: column,
		op:     OpGt,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) GreaterThanOrEqual(column string, value any) Q {
	op := Op{
		column: column,
		op:     OpGte,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) Null(column string) Q {
	op := Op{
		column: column,
		op:     OpNull,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) NotNull(column string) Q {
	op := Op{
		column: column,
		op:     OpNotNull,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) StartsWith(column string, value string) Q {
	op := Op{
		column: column,
		op:     OpStartsWith,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

func (q Q) EndsWith(column string, value string) Q {
	op := Op{
		column: column,
		op:     OpEndsWith,
		value:  value,
	}

	q.where = append(q.where, op)

	return q
}

const (
	OpEq = iota
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte
	OpNull
	OpNotNull
	OpStartsWith
	OpEndsWith
	OpOrderDesc
	OpOrderAsc
	OpLimit
	OpOffset
)

type Op struct {
	column string
	op     int
	value  any
}

func (op Op) Column() string {
	return op.column
}

func (op Op) Op() int {
	return op.op
}

func (op Op) Value() any {
	return op.value
}
