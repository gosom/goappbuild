package postgres

import (
	"context"
	"database/sql"

	"github.com/gosom/goappbuild"
)

var _ goappbuild.Storage = (*storage)(nil)

type storage struct {
	db *sql.DB
	tx *sql.Tx

	users       goappbuild.UserRepo
	projects    goappbuild.ProjectRepo
	collections goappbuild.CollectionRepo
	databases   goappbuild.DatabaseRepo
	queries     goappbuild.QueryRepo
}

func NewUnitOfWork(db *sql.DB) goappbuild.Storage {
	return &storage{
		db:          db,
		users:       NewUserRepo(db),
		projects:    NewProjectRepo(db),
		collections: NewCollectionRepo(db),
		databases:   NewDBRepo(db),
		queries:     NewQueryRepo(db),
	}
}

func (uw *storage) New(ctx context.Context) (goappbuild.Storage, error) {
	tx, err := uw.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	ans := storage{
		db:          uw.db,
		tx:          tx,
		users:       NewUserRepo(tx),
		projects:    NewProjectRepo(tx),
		collections: NewCollectionRepo(tx),
		databases:   NewDBRepo(tx),
	}

	return &ans, nil
}

func (uw *storage) Commit(ctx context.Context) error {
	return uw.tx.Commit()
}

func (uw *storage) Rollback(ctx context.Context) error {
	return uw.tx.Rollback()
}

func (uw *storage) Users() goappbuild.UserRepo {
	return uw.users
}

func (uw *storage) Projects() goappbuild.ProjectRepo {
	return uw.projects
}

func (uw *storage) Collections() goappbuild.CollectionRepo {
	return uw.collections
}

func (uw *storage) Databases() goappbuild.DatabaseRepo {
	return uw.databases
}

func (uw *storage) Queries() goappbuild.QueryRepo {
	return uw.queries
}
