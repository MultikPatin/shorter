package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"main/internal/config"
	"main/internal/models"
	"net/url"
	"time"
)

const (
	urlPrefix = "http://"
	delimiter = "/"
)

type LinksRepository interface {
	Add(ctx context.Context, addedLink models.AddedLink) (string, error)
	AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]string, error)
	Get(ctx context.Context, short string) (string, error)
	Close() error
	Ping() error
}

type LinksService struct {
	linksRepository LinksRepository
	shortPre        string
}

func NewLinksService(c *config.Config, linksRepository LinksRepository) *LinksService {
	return &LinksService{
		linksRepository: linksRepository,
		shortPre:        c.ShortLinkPrefix,
	}
}

func (s *LinksService) Ping() error {
	err := s.linksRepository.Ping()
	return err
}

func (s *LinksService) Close() error {
	err := s.linksRepository.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *LinksService) Add(ctx context.Context, shortenRequest models.ShortenRequest, host string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	addedLink := models.AddedLink{
		Short:  getKey(u, s.shortPre),
		Origin: shortenRequest.URL,
	}

	id, err := s.linksRepository.Add(ctx, addedLink)
	if err != nil {
		return "", fmt.Errorf("failed to add link: %w", err)
	}
	return getResponseLink(id, s.shortPre, urlPrefix+host), nil
}

func (s *LinksService) AddBatch(ctx context.Context, shortenRequests []models.ShortenRequest, host string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	retries := 0
	var addedLinks []models.AddedLink
	var responseLinks []string

	for i := 0; i <= len(shortenRequests); {
		u, err := uuid.NewRandom()
		if err != nil {
			retries += 1
			continue
		}
		if retries >= 5 {
			return nil, fmt.Errorf("failed to generate UUIDs: %w", err)
		}
		addedLink := models.AddedLink{
			Short:  getKey(u, s.shortPre),
			Origin: shortenRequests[i].URL,
		}
		addedLinks = append(addedLinks, addedLink)
		i++
	}

	ids, err := s.linksRepository.AddBatch(ctx, addedLinks)
	if err != nil {
		return nil, fmt.Errorf("failed to add links: %w", err)
	}

	for _, id := range ids {
		responseLinks = append(responseLinks, getResponseLink(id, s.shortPre, urlPrefix+host))
	}
	return responseLinks, nil
}

func (s *LinksService) Get(ctx context.Context, shortLink string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	originLink, err := s.linksRepository.Get(ctx, shortLink)
	if err != nil {
		return "", fmt.Errorf("origin link not found: %w", err)
	}
	return originLink, nil
}

func getKey(u uuid.UUID, p string) string {
	if isURL(p) {
		return u.String()
	}
	return p + u.String()
}

func getResponseLink(k string, p string, h string) string {
	if isURL(p) {
		return p + delimiter + k + delimiter
	}
	return h + delimiter + k + delimiter
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
