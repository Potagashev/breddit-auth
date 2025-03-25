package users

import "github.com/google/uuid"

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService{
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(createData *UserCreate) (uuid.UUID, error) {
	return s.userRepo.CreateUser(createData)
}

func (s *UserService) GetUserByUsername(username string) (*User, error) {
	return s.userRepo.GetUserByUsername(username)
}