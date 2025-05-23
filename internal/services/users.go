package services // Package services implements business logic for user-related operations.

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

// Constants for batch-processing of link deletions.
const batchSize = 5

// Custom error types for user-related failures.
var (
	ErrAddUser       = errors.New("failed to insert user")
	ErrNoLinksByUser = errors.New("links by userID %d not found")
)

// UsersService encapsulates the business logic for user management.
type UsersService struct {
	usersRepository interfaces.UsersRepository // Dependency for accessing user-related repository methods.
}

// NewUserService constructs a new UsersService instance bound to a specific users repository.
func NewUserService(usersRepository interfaces.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

// Login initiates a new user session and returns a unique user identifier.
func (s *UsersService) Login() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userID, err := s.usersRepository.Login(ctx)
	if err != nil {
		return userID, ErrAddUser
	}
	return userID, nil
}

// GetLinks retrieves all links created by the current user.
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

// DeleteLinks deletes a collection of short links in parallelized batches.
func (s *UsersService) DeleteLinks(ctx context.Context, shortLinks []string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	batches := setBatches(shortLinks)
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

// setBatches splits a large collection of links into smaller chunks for batch processing.
func setBatches(shortLinks []string) [][]string {
	var batches [][]string
	for i := 0; i < len(shortLinks); i += batchSize {
		end := i + batchSize
		if end > len(shortLinks) {
			end = len(shortLinks)
		}
		batch := shortLinks[i:end]
		batches = append(batches, batch)
	}
	return batches
}

// shortLinksGenerator feeds batches of links into a channel for consumption by worker goroutines.
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
