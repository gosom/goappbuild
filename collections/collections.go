package collections

import (
	"context"

	"github.com/gosom/goappbuild"
)

type collectionService struct {
	storage goappbuild.Storage
}

// NewCollectionService returns a new instance of a collection service
func New(storage goappbuild.Storage) goappbuild.CollectionService {
	return &collectionService{
		storage: storage,
	}
}

// Create creates a new collection
func (s *collectionService) Create(ctx context.Context, req goappbuild.CollectionCreateRequest) (goappbuild.Collection, error) {
	uw, err := s.storage.New(ctx)
	if err != nil {
		return goappbuild.Collection{}, err
	}

	defer uw.Rollback(ctx)

	project, err := s.storage.Projects().Get(ctx, req.ProjectID)
	if err != nil {
		return goappbuild.Collection{}, err
	}

	collection := goappbuild.Collection{
		ProjectID:  req.ProjectID,
		Name:       req.Name,
		Attributes: s.getDefaultAttributes(),
	}

	err = s.storage.Collections().Create(ctx, project.SchemaName(), &collection)
	if err != nil {
		return goappbuild.Collection{}, err
	}

	err = s.storage.Databases().CreateTable(ctx, project.SchemaName(), collection.TableName())
	if err != nil {
		return goappbuild.Collection{}, err
	}

	err = s.storage.Databases().CreateColumns(
		ctx,
		project.SchemaName(),
		collection.TableName(), collection.Attributes,
	)

	if err != nil {
		return goappbuild.Collection{}, err
	}

	if err := uw.Commit(ctx); err != nil {
		return goappbuild.Collection{}, err
	}

	return collection, nil
}

func (s *collectionService) getDefaultAttributes() map[string]goappbuild.Attribute {
	attributes := make(map[string]goappbuild.Attribute)

	attributes["id"] = goappbuild.Attribute{
		Name:     "id",
		Type:     goappbuild.AttributeTypeUUID,
		Required: true,
		Primary:  true,
	}
	attributes["created_at"] = goappbuild.Attribute{
		Name:     "created_at",
		Type:     goappbuild.AttributeTypeTime,
		Required: true,
	}
	attributes["updated_at"] = goappbuild.Attribute{
		Name:     "updated_at",
		Type:     goappbuild.AttributeTypeTime,
		Required: true,
	}

	return attributes
}
