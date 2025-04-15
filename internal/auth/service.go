package auth

import (
	"context"
	"time"

	"github.com/Potagashev/breddit_auth/internal/config"
	"github.com/Potagashev/breddit_auth/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	pb "github.com/Potagashev/breddit_auth/internal/auth/proto"
)


type AuthService struct {
    pb.UnimplementedAuthServiceServer
	userService users.UserService
	cfg config.Config
}

func NewAuthService(userService *users.UserService, cfg *config.Config) *AuthService {
	return &AuthService{userService: *userService, cfg: *cfg}
}

func (s *AuthService) SignUp(signUpData *SignUpData) (*SignUpResponse, error) {
	passwordHash := s.hashPassword(signUpData.Password)
	userCreateData := users.UserCreate{
		Username: signUpData.Username,
		PasswordHash: passwordHash,
		Email: signUpData.Email,
	}
	userId, err := s.userService.CreateUser(&userCreateData)
	if err != nil {
		return nil, err
	}
	
	token, err := s.generateAuthToken(userId)
	if err != nil {
		return nil, err
	}
	return &SignUpResponse{Token: token}, nil
}

func (s *AuthService) SignIn(signInData *SignInData) (*SignInResponse, error) {
	user, err := s.userService.GetUserByUsername(signInData.Username)
	if user == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signInData.Password))
	if err != nil {
		return nil, err
	}

	token, err := s.generateAuthToken(user.Id)
	if err != nil {
		return nil, err
	}
	return &SignInResponse{Token: token}, nil
}

func (s *AuthService) VerifyAuthToken(token string) (*VerifyAuthTokenHTTPResponse, error) {
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, nil
	}

	user, err := s.userService.GetUserById(claims.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &VerifyAuthTokenHTTPResponse{Id: user.Id, Username: user.Username, Email: user.Email}, nil
}

func (s *AuthService) VerifyAuthTokenRPC(ctx context.Context, req *pb.VerifyAuthTokenRequest) (*pb.VerifyAuthTokenResponse, error) {
    // Call the existing VerifyAuthToken method
    response, err := s.VerifyAuthToken(req.Token)
    if err != nil {
        return nil, err
    }

    return &pb.VerifyAuthTokenResponse{
        Id:       response.Id.String(),
        Username: response.Username,
        Email:    response.Email,
    }, nil
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
