package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"main/internal/config"
	"net/url"
	"time"
)

const (
	urlPrefix = "http://"
	delimiter = "/"
)

type LinksRepository interface {
	Add(ctx context.Context, short string, origin string) (string, error)
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

func (s *LinksService) Add(ctx context.Context, origin string, host string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	id, err := s.linksRepository.Add(ctx, getKey(u, s.shortPre), origin)
	if err != nil {
		return "", fmt.Errorf("failed to add link: %w", err)
	}
	return getResponseLink(id, s.shortPre, urlPrefix+host), nil
}

func (s *LinksService) Get(ctx context.Context, id string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	origin, err := s.linksRepository.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("origin not found: %w", err)
	}
	return origin, nil
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
