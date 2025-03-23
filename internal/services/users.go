package services

import (
	"context"
	"errors"
	"fmt"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"time"
)

var ErrAddUser = errors.New("failed to insert user")
var ErrNoLinksByUser = errors.New("links by userID %d not found")

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

func (s *UsersService) GetLinks(ctx context.Context, host string) ([]models.UserLinks, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var links []models.UserLinks

	results, err := s.usersRepository.GetLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("links not found: %w", err)
	}
	for _, result := range results {
		link := models.UserLinks{
			Shorten:  getResponseLink(result.Shorten, shortPre, constants.UrlPrefix+host),
			Original: result.Original,
		}
		links = append(links, link)
	}

	return links, nil
}
