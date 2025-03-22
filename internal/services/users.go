package services

import (
	"context"
	"main/internal/interfaces"
)

type UsersService struct {
	usersRepository interfaces.UsersRepository
}

func NewUserService(usersRepository interfaces.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (s *UsersService) Login(ctx context.Context) (int, error) {
	//ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	//defer cancel()

	return 4, nil
}
