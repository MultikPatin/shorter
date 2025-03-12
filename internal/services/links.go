package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"main/internal/config"
	"main/internal/interfaces"
	"main/internal/models"
	"net/url"
	"time"
)

var ErrConflict = errors.New("data conflict")

const (
	urlPrefix = "http://"
	delimiter = "/"
)

type LinksService struct {
	linksRepository interfaces.LinksRepository
	shortPre        string
}

func NewLinksService(c *config.Config, linksRepository interfaces.LinksRepository) *LinksService {
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

func (s *LinksService) Add(ctx context.Context, originLink models.OriginLink, host string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	addedLink := models.AddedLink{
		Short:  getKey(u, s.shortPre),
		Origin: originLink.URL,
	}

	id, err := s.linksRepository.Add(ctx, addedLink)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			return getResponseLink(id, s.shortPre, urlPrefix+host), err
		} else {
			return "", fmt.Errorf("failed to add link: %w", err)
		}
	}
	return getResponseLink(id, s.shortPre, urlPrefix+host), nil
}

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
			Short:         getKey(u, s.shortPre),
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
			Result:        getResponseLink(result.Result, s.shortPre, urlPrefix+host),
		}
		responseLinks = append(responseLinks, response)
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
