package user

import (
	"context"
	"tender/internal/repository"
	svc "tender/internal/service"
	"tender/internal/service/model"
)

type userService struct {
	repo repository.UsersRepos
}

func NewUserService(repo repository.UsersRepos) svc.UserService {
	return &userService{repo: repo}
}

func (s *userService) Get(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.Get(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
