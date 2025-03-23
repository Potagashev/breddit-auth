package users

type UserCreate struct {
	Username string
	PasswordHash string
	Email string
}