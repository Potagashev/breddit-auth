package users

import "github.com/google/uuid"

type UserCreate struct {
	Username string
	PasswordHash string
	Email string
}

type User struct {
	Id uuid.UUID
	Username string
	PasswordHash string
	Email string
}