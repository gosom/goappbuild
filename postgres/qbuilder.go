package postgres

import (
	"strconv"
	"strings"

	"github.com/gosom/goappbuild"
)

type postgresQ struct {
	goappbuild.Q
	sb   strings.Builder
	args []any
}

func NewPostgresQ(params goappbuild.Q) *postgresQ {
	return &postgresQ{
		Q: params,
	}
}

func (q *postgresQ) Build() (string, []any, error) {
	q.selectColumns()

	q.from()

	if err := q.where(); err != nil {
		return "", nil, err
	}

	return q.sb.String(), q.args, nil
}

func (q *postgresQ) selectColumns() {
	q.sb.WriteString("SELECT ")
	cols := q.Cols()
	if len(cols) == 0 {
		q.sb.WriteString("*")

		return
	}

	escaped := make([]string, len(cols))
	for i := range cols {
		escaped[i] = escape(cols[i])
	}

	columns := strings.Join(escaped, ", ")

	q.sb.WriteString(columns)

	return
}

func (q *postgresQ) from() {
	table := escape(q.GetSchema()) + "." + escape(q.GetTable())

	q.sb.WriteString(" FROM ")
	q.sb.WriteString(table)

	return
}

func (q *postgresQ) where() error {
	where := q.Where()

	if len(where) == 0 {

		return nil
	}

	for i := range where {
		op := postgresOp{op: where[i].Op()}

		if i == 0 {
			q.sb.WriteString(" WHERE ")
		} else {
			q.sb.WriteString(" AND ")
		}

		q.sb.WriteString(escape(where[i].Column()))
		q.sb.WriteString(" ")
		q.sb.WriteString(op.String())
		q.sb.WriteString(" ")

		if where[i].Value() != nil {
			value := getValue(where[i])
			if value != nil {
				q.sb.WriteString("$")
				q.sb.WriteString(strconv.Itoa(len(q.args) + 1))

				q.args = append(q.args, value)
			}
		}
	}

	return nil
}

func getValue(op goappbuild.Op) any {
	if op.Value() == nil {
		return nil
	}

	switch op.Op() {
	case goappbuild.OpStartsWith:
		return op.Value().(string) + "%"
	case goappbuild.OpEndsWith:
		return "%" + op.Value().(string)
	default:
		return op.Value()
	}
}

func escape(val string) string {
	return `"` + strings.Replace(val, `"`, `""`, -1) + `"`
}

type postgresOp struct {
	op int
}

func (o postgresOp) String() string {
	switch o.op {
	case goappbuild.OpEq:
		return "="
	case goappbuild.OpNeq:
		return "!="
	case goappbuild.OpLt:
		return "<"
	case goappbuild.OpLte:
		return "<="
	case goappbuild.OpGt:
		return ">"
	case goappbuild.OpGte:
		return ">="
	case goappbuild.OpStartsWith, goappbuild.OpEndsWith:
		return "LIKE"
	case goappbuild.OpNull:
		return "IS NULL"
	case goappbuild.OpNotNull:
		return "IS NOT NULL"
	default:
		return ""
	}
}
