package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SignUpData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
}

type SignUpResponse struct {
	Token string `json:"token"`
}

type Claims struct {
    UserId uuid.UUID `json:"userId"`
    jwt.RegisteredClaims
}

type SignInData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}