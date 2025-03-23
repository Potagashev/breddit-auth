package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		return userId, err
	}
	return userId, nil
}