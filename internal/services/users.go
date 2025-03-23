package services

import (
	"context"
	"errors"
	"fmt"
	"main/internal/interfaces"
	"main/internal/models"
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

func (s *UsersService) Login() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userID, err := s.usersRepository.Login(ctx)
	if err != nil {
		return userID, ErrAddUser
	}
	return userID, nil
}

func (s *UsersService) GetLinks(ctx context.Context) ([]models.UserLinks, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	links, err := s.usersRepository.GetLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("links not found: %w", err)
	}
	return links, nil
}
