package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"main/internal/config"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"time"
)

// Custom error types for handling link conflicts and deleted links.
var (
	ErrConflict    = errors.New("data conflict")
	ErrDeletedLink = errors.New("link is deleted")
)

// LinksService encapsulates the business logic for link management.
type LinksService struct {
	linksRepository interfaces.LinksRepository // Dependency for accessing link-related repository methods.
}

// NewLinksService constructs a new LinksService instance wired to a specific links repository.
func NewLinksService(c *config.Config, linksRepository interfaces.LinksRepository) *LinksService {
	shortPre = c.ShortLinkPrefix
	return &LinksService{
		linksRepository: linksRepository,
	}
}

// Add creates a new link record, assigning a unique short identifier.
func (s *LinksService) Add(ctx context.Context, originLink models.OriginLink, host string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	addedLink := models.AddedLink{
		Short:  getKey(u, shortPre),
		Origin: originLink.URL,
	}

	id, err := s.linksRepository.Add(ctx, addedLink)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			return getResponseLink(id, shortPre, constants.URLPrefix+host), err
		} else {
			return "", fmt.Errorf("failed to add link: %w", err)
		}
	}
	return getResponseLink(id, shortPre, constants.URLPrefix+host), nil
}

// AddBatch allows batch-adding multiple links simultaneously.
func (s *LinksService) AddBatch(ctx context.Context, originLinks []models.OriginLink, host string) ([]models.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	retries := 0
	var addedLinks []models.AddedLink

	for i := 0; i < len(originLinks); {
		u, err := uuid.NewRandom()
		if err != nil {
			retries += 1
			continue
		}
		if retries >= 5 {
			return nil, fmt.Errorf("failed to generate UUIDs: %w", err)
		}
		addedLink := models.AddedLink{
			CorrelationID: originLinks[i].CorrelationID,
			Short:         getKey(u, shortPre),
			Origin:        originLinks[i].URL,
		}
		addedLinks = append(addedLinks, addedLink)
		i++
	}

	results, err := s.linksRepository.AddBatch(ctx, addedLinks)
	if err != nil {
		return nil, fmt.Errorf("failed to add links: %w", err)
	}

	var responseLinks []models.Result

	for _, result := range results {
		response := models.Result{
			CorrelationID: result.CorrelationID,
			Result:        getResponseLink(result.Result, shortPre, constants.URLPrefix+host),
		}
		responseLinks = append(responseLinks, response)
	}
	return responseLinks, nil
}

// Get resolves a short link to its original URL.
func (s *LinksService) Get(ctx context.Context, shortLink string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	originLink, err := s.linksRepository.Get(ctx, shortLink)
	if err != nil {
		return "", fmt.Errorf("origin link not found: %w", err)
	}
	return originLink, nil
}
