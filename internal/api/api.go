package api

import (
	"context"
	"projectServer/internal/service"
	pb "projectServer/pkg/protos/gen/go"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	Service *service.UserService
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	id, err := s.Service.CreateUser(ctx, req.GetName(), req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{Id: id}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	id, name, email, err := s.Service.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.GetUserByIDResponse{Id: id, Name: name, Email: email}, nil
}

func (s *UserService) GetNewUser(ctx context.Context, req *pb.GetNewUserRequest) (*pb.GetNewUserResponse, error) {

	users, err := s.Service.GetNewUsers(ctx)
	if err != nil {
		return nil, err
	}

	var pbUsers []*pb.User
	for _, user := range users {
		u := user
		pbUsers = append(pbUsers, &u)
	}

	return &pb.GetNewUserResponse{Users: pbUsers}, nil
}
