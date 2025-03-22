package services

import (
	"context"
	"errors"
	"main/internal/interfaces"
	"time"
)

var ErrAddUser = errors.New("failed to insert user")

type UsersService struct {
	usersRepository interfaces.UsersRepository
}

func NewUserService(usersRepository interfaces.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (s *UsersService) Login(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userID, err := s.usersRepository.Login(ctx)
	if err != nil {
		return userID, ErrAddUser
	}
	return userID, nil
}
