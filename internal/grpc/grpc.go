package grpc

import (
	"context"
	"database/sql"
	"fmt"
	pb "projectServer/pkg/protos/gen/go"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	Db *sql.DB
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var id int

	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := s.Db.QueryRow(query, req.GetName(), req.GetEmail()).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}

	return &pb.CreateUserResponse{Id: int64(id)}, nil

}

func (s *UserService) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	var name, email string

	query := `SELECT name, email FROM users WHERE id = $1 LIMIT 1`
	err := s.Db.QueryRow(query, req.GetUserId()).Scan(&name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %v not found", req.GetUserId())
		}
		return nil, err
	}

	return &pb.GetUserByIDResponse{Name: name, Email: email}, nil
}

func (s *UserService) GetNewUser(ctx context.Context, req *pb.GetNewUserRequest) (*pb.GetNewUserResponse, error) {
	var users []*pb.User

	query := `SELECT id, name, email, created_at FROM users WHERE created_at >= NOW() - INTERVAL '5 minutes'`
	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user pb.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return &pb.GetNewUserResponse{Users: users}, nil
}
