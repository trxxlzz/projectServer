package main

import (
	"database/sql"
	_ "database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	_ "github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"projectServer/internal/api"
	"projectServer/internal/repository"
	"projectServer/internal/service"
	pb "projectServer/pkg/protos/gen/go"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:trxxlzz@localhost:5432/test_database?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = goose.Up(db, "./migrations")
	if err != nil {
		panic(err)
	}

	defer func() {
		log.Println("closing database connection")
		db.Close()
	}()

	repo := &repository.UserRepo{DB: db}

	userService := &service.UserService{Repo: repo}

	apiService := &api.UserService{Service: userService}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, apiService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	log.Println("Server is working!")
}
