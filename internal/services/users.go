package services

import (
	"context"
	"main/internal/config"
	"main/internal/interfaces"
	"time"
)

type UsersService struct {
	usersRepository interfaces.UsersRepository
}

func NewUserService(c *config.Config, usersRepository interfaces.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: UsersRepository,
	}
}

func (s *UsersService) Close() error {
	err := s.usersRepository.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersService) Login(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return 4, nil
}
