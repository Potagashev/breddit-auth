package users

import "github.com/google/uuid"

type UserCreate struct {
	Username string
	PasswordHash string
	Email string
}

type UserInternal struct {
	Id uuid.UUID
	Username string
	PasswordHash string
	Email string
}