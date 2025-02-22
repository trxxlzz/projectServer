package service

import (
	"context"
	pb "projectServer/pkg/protos/gen/go"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name string, email string) (int64, error)
	GetUserByID(ctx context.Context, userID int64) (int64, string, string, error)
	GetNewUsers(ctx context.Context) ([]pb.User, error)
}

type UserService struct {
	Repo UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, name string, email string) (int64, error) {

	id, err := s.Repo.CreateUser(ctx, name, email)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (int64, string, string, error) {

	id, name, email, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		return 0, "", "", err
	}

	return id, name, email, nil
}

func (s *UserService) GetNewUsers(ctx context.Context) ([]pb.User, error) {

	users, err := s.Repo.GetNewUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
