package projects

import (
	"context"

	"github.com/google/uuid"
	"github.com/gosom/goappbuild"
)

var _ goappbuild.ProjectService = (*projectService)(nil)

type projectService struct {
	storage goappbuild.Storage
}

func New(storage goappbuild.Storage) goappbuild.ProjectService {
	return &projectService{
		storage: storage,
	}
}

func (s *projectService) Create(ctx context.Context, req goappbuild.CreateProjectRequest) (goappbuild.Project, error) {
	p := goappbuild.Project{
		UserID: req.UserID,
		Name:   req.Name,
	}

	uw, err := s.storage.New(ctx)
	if err != nil {
		return goappbuild.Project{}, err
	}

	defer uw.Rollback(ctx)

	if err := uw.Projects().Create(ctx, &p); err != nil {
		return goappbuild.Project{}, err
	}

	if err := uw.Databases().CreateSchema(ctx, p.SchemaName()); err != nil {
		return goappbuild.Project{}, err
	}

	if err := uw.Commit(ctx); err != nil {
		return goappbuild.Project{}, err
	}

	return p, nil
}

func (s *projectService) Get(ctx context.Context, id uuid.UUID) (goappbuild.Project, error) {
	return s.storage.Projects().Get(ctx, id)
}
