package postgres

import (
	"fmt"
	"strings"

	"github.com/gosom/goappbuild"
)

type createTableParams struct {
	schema string
	table  string
}

func createTableStmt(params createTableParams) string {
	tableQ := fmt.Sprintf(`CREATE TABLE %s.%s ()`, params.schema, params.table)

	return tableQ
}

type addAttributeParams struct {
	schema    string
	table     string
	attribute goappbuild.Attribute
}

func addAttributeStmt(params addAttributeParams) (string, error) {
	typeQ, err := attributeType{params.attribute.Type}.postgresType()
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	sb.WriteString("ALTER TABLE ")
	sb.WriteString(params.schema)
	sb.WriteString(".")
	sb.WriteString(params.table)
	sb.WriteString(" ADD COLUMN ")
	sb.WriteString(`"` + params.attribute.Name + `"`)
	sb.WriteString(" ")
	sb.WriteString(typeQ)

	if params.attribute.Required {
		sb.WriteString(" NOT NULL")
	}

	if params.attribute.Primary {
		sb.WriteString(" PRIMARY KEY")
	}

	if params.attribute.Unique {
		sb.WriteString(" UNIQUE")
	}

	return sb.String(), nil
}

func addIndexesStmt(params addAttributeParams) string {
	if !params.attribute.Index {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("CREATE INDEX ")
	sb.WriteString("idx_")
	sb.WriteString(params.schema + "_" + params.table + "_" + params.attribute.Name)
	sb.WriteString(" ON ")
	sb.WriteString(`"` + params.schema + `"."` + params.table + `"`)
	sb.WriteString(" (" + params.attribute.Name + ")")

	return sb.String()
}

type attributeType struct {
	goappbuild.AttributeType
}

func (t attributeType) postgresType() (string, error) {
	var typeQ string

	switch t.AttributeType {
	case goappbuild.AttributeTypeString:
		typeQ = "TEXT"
	case goappbuild.AttributeTypeInteger:
		typeQ = "INT"
	case goappbuild.AttributeTypeNumeric:
		typeQ = "NUMERIC"
	case goappbuild.AttributeTypeFloat:
		typeQ = "FLOAT"
	case goappbuild.AttributeTypeBoolean:
		typeQ = "BOOLEAN"
	case goappbuild.AttributeTypeTime:
		typeQ = "TIMESTAMPTZ"
	case goappbuild.AttributeTypeUUID:
		typeQ = "UUID"
	case goappbuild.AttributeTypeJSON:
		typeQ = "JSONB"
	default:
		return "", fmt.Errorf("invalid attribute type: %s", t.AttributeType)
	}

	return typeQ, nil
}
