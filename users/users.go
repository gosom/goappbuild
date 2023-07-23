package users

import (
	"context"

	"github.com/gosom/goappbuild"
)

var _ goappbuild.UserService = (*service)(nil)

type service struct {
	storage goappbuild.Storage
}

func New(storage goappbuild.Storage) goappbuild.UserService {
	return &service{
		storage: storage,
	}
}

func (s *service) Register(ctx context.Context, _ goappbuild.RegisterUserRequest) (goappbuild.User, error) {
	u := goappbuild.User{}

	if err := u.Validate(); err != nil {
		return goappbuild.User{}, err
	}

	if err := s.storage.Users().Create(ctx, &u); err != nil {
		return goappbuild.User{}, err
	}

	return u, nil
}
