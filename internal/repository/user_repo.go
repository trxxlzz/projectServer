package repository

import (
	"context"
	"database/sql"
	"log"
	pb "projectServer/pkg/protos/gen/go"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) CreateUser(ctx context.Context, name string, email string) (int64, error) {
	var id int64

	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := r.DB.QueryRowContext(ctx, query, name, email).Scan(&id)
	if err != nil {
		log.Printf("Could not create user: %v", err)
		return 0, err
	}

	return id, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID int64) (int64, string, string, error) {
	var id int64
	var name, email string

	query := `SELECT id, name, email FROM users WHERE id = $1 LIMIT 1`
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&id, &name, &email)
	if err != nil {
		log.Printf("Could not get user: %v", err)
		return 0, "", "", err
	}

	return id, name, email, nil
}

func (r *UserRepo) GetNewUsers(ctx context.Context) ([]pb.User, error) {
	var users []pb.User

	query := `SELECT id, name, email, created_at FROM users WHERE created_at >= NOW() - INTERVAL '5 minutes'`
	rows, err := r.DB.QueryContext(ctx, query)
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
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
