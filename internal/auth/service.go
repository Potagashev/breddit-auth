package auth

import (
	"time"

	"github.com/Potagashev/breddit_auth/internal/config"
	"github.com/Potagashev/breddit_auth/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService users.UserService
	cfg config.Config
}

func NewAuthService(userService *users.UserService, cfg *config.Config) *AuthService {
	return &AuthService{userService: *userService, cfg: *cfg}
}

func (s *AuthService) SignUp(signUpData *SignUpData) (SignUpResponse, error) {
	passwordHash := s.hashPassword(signUpData.Password)
	userCreateData := users.UserCreate{
		Username: signUpData.Username,
		PasswordHash: passwordHash,
		Email: signUpData.Email,
	}
	userId, err := s.userService.CreateUser(&userCreateData)
	if err != nil {
		return SignUpResponse{Token: ""}, err
	}
	
	token, err := s.generateAuthToken(userId)
	if err != nil {
		return SignUpResponse{Token: ""}, err
	}
	return SignUpResponse{Token: token}, nil
}

func (s *AuthService) hashPassword(rawPassword string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (s *AuthService) generateAuthToken(userId uuid.UUID) (string, error) {
    expirationTime := time.Now().Add(time.Duration(s.cfg.JWTTokenExpMinutes) * time.Minute)

    claims := &Claims{
        UserId: userId,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "breddit",
            Subject:   "authToken",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
