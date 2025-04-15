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

func (s *UserService) GetUserByUsername(username string) (*UserInternal, error) {
	return s.userRepo.GetUserByUsername(username)
}

func (s *UserService) GetUserById(id uuid.UUID) (*UserInternal, error) {
	return s.userRepo.GetUserById(id)
}