package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserNotFound = errors.New("User Not Found")
	ErrUserInvalid  = errors.New("User Invalid")
)

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo,
	}
}

func (r *userService) Find(userId string) (*User, error) {
	return r.userRepo.Find(userId)
}

func (r *userService) Store(user *User) error {
	if err := validate.Validate(user); err != nil {
		return fmt.Errorf( "service.User.Store %w", err)
	}
	user.CreatedAt = time.Now().UTC().Unix()
	return r.userRepo.Store(user)
}
