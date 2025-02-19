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
	grpcService "projectServer/internal/grpc"
	pb "projectServer/pkg/protos/gen/go"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:trxxlzz@localhost:5432/test_database?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = goose.Up(db, "C:/Users/Nikita/GolandProjects/proj1/internal/migrations")
	if err != nil {
		panic(err)
	}

	defer func() {
		log.Println("closing database connection")
		db.Close()
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	userService := &grpcService.UserService{Db: db}
	pb.RegisterUserServiceServer(grpcServer, userService)

	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
