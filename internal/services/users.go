package services

import (
	"context"
	"errors"
	"fmt"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"sync"
	"time"
)

const batchSize = 5

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
			Shorten:  getResponseLink(result.Shorten, shortPre, constants.URLPrefix+host),
			Original: result.Original,
		}
		links = append(links, link)
	}

	return links, nil
}

func (s *UsersService) DeleteLinks(ctx context.Context, shortLinks []string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var batches [][]string
	for i := 0; i < len(shortLinks); i += batchSize {
		end := i + batchSize
		if end > len(shortLinks) {
			end = len(shortLinks)
		}
		batch := shortLinks[i:end]
		batches = append(batches, batch)
	}

	batchChan := shortLinksGenerator(ctx, batches)
	errChan := make(chan error, len(batches))

	var wg sync.WaitGroup
	for batch := range batchChan {
		wg.Add(1)
		data := batch
		go func(batch []string) {
			defer wg.Done()
			err := s.usersRepository.DeleteLinks(ctx, batch)
			if err != nil {
				errChan <- fmt.Errorf("error updating the link patch: %w", err)
			}
		}(data)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("the following errors occurred when updating the links: %v", errs)
	}

	return nil
}

func shortLinksGenerator(ctx context.Context, batches [][]string) chan []string {
	inputCh := make(chan []string, len(batches))

	go func() {
		defer close(inputCh)

		for _, batch := range batches {
			select {
			case <-ctx.Done():
				return
			case inputCh <- batch:
			}
		}
	}()
	return inputCh
}
