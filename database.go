package goappbuild

import "context"

type DatabaseRepo interface {
	CreateSchema(context.Context, string) error
	CreateTable(context.Context, string, string) error
	CreateColumns(context.Context, string, string, map[string]Attribute) error
}
