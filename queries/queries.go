package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
)

type queryService struct {
	storage goappbuild.Storage
}

// NewQueryService returns a new instance of a query service
func New(storage goappbuild.Storage) goappbuild.QueryService {
	return &queryService{
		storage: storage,
	}
}

func (q *queryService) Get(ctx context.Context, param goappbuild.Q) (goappbuild.Document, error) {
	m, err := q.storage.Queries().Get(ctx, param)
	if err != nil {
		return goappbuild.Document{}, err
	}

	ans := goappbuild.Document{
		Values: m,
	}

	return ans, nil
}

func (q *queryService) Create(
	ctx context.Context,
	projectID uuid.UUID,
	collectionName string,
	data map[string]any,
) (goappbuild.Document, error) {

	uw, err := q.storage.New(ctx)
	if err != nil {
		return goappbuild.Document{}, err
	}

	defer uw.Rollback(ctx)

	project, err := q.storage.Projects().Get(ctx, projectID)
	if err != nil {
		return goappbuild.Document{}, err
	}

	result, err := q.storage.Queries().Create(ctx, project.Name, collectionName, data)
	if err != nil {
		return goappbuild.Document{}, err
	}

	if err := uw.Commit(ctx); err != nil {
		return goappbuild.Document{}, err
	}

	ans := goappbuild.Document{
		Values: result,
	}

	return ans, nil
}

func (q *queryService) Update(
	ctx context.Context,
	projectID uuid.UUID,
	collectionName string,
	id uuid.UUID,
	data map[string]any,
) (goappbuild.Document, error) {
	uw, err := q.storage.New(ctx)
	if err != nil {
		return goappbuild.Document{}, err
	}

	defer uw.Rollback(ctx)

	project, err := q.storage.Projects().Get(ctx, projectID)
	if err != nil {
		return goappbuild.Document{}, err
	}

	result, err := q.storage.Queries().Update(ctx, project.Name, collectionName, id, data)
	if err != nil {
		return goappbuild.Document{}, err
	}

	if err := uw.Commit(ctx); err != nil {
		return goappbuild.Document{}, err
	}

	ans := goappbuild.Document{
		Values: result,
	}

	return ans, nil
}

func (q *queryService) Delete(
	ctx context.Context,
	projectID uuid.UUID,
	collectionName string,
	id uuid.UUID,
) error {
	uw, err := q.storage.New(ctx)
	if err != nil {
		return err
	}

	defer uw.Rollback(ctx)

	project, err := q.storage.Projects().Get(ctx, projectID)
	if err != nil {
		return err
	}

	if err := q.storage.Queries().Delete(ctx, project.Name, collectionName, id); err != nil {
		return err
	}

	if err := uw.Commit(ctx); err != nil {
		return err
	}

	return nil
}
