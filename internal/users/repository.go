package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	DB *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(createData *UserCreate) (uuid.UUID, error) {
	var userId uuid.UUID
	query := `INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id as userId`
	err := r.DB.QueryRow(
		context.Background(),
		query,
		createData.Username,
		createData.PasswordHash,
		createData.Email,
	).Scan(&userId)
	if err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // Unique violation error code
				if pgErr.ConstraintName == "unique_username" {
					return uuid.Nil, fmt.Errorf("username already exists")
				}
				if pgErr.ConstraintName == "unique_email" {
					return uuid.Nil, fmt.Errorf("email already exists")
				}
			}
		}
		return uuid.Nil, err
	}
	return userId, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	query := `SELECT id, username, password_hash, email FROM users WHERE username = $1`
	err := r.DB.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
            // Return nil for both user and error if no rows are found
            return nil, nil
        }
		return nil, err
	}
	return &user, nil
}