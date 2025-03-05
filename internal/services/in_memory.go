package services

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
)

const (
	urlPrefix = "http://"
	delimiter = "/"
)

type dataBase interface {
	Add(id string, link string) (string, error)
	Get(id string) (string, error)
}

type LinksService struct {
	database dataBase
	shortPre string
}

func NewLinksService(db dataBase, shortPre string) *LinksService {
	return &LinksService{
		database: db,
		shortPre: shortPre,
	}
}

func (s *LinksService) Add(origin string, host string) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	id, err := s.database.Add(getKey(u, s.shortPre), origin)
	if err != nil {
		return "", fmt.Errorf("failed to add link: %w", err)
	}
	return getResponseLink(id, s.shortPre, urlPrefix+host), nil
}

func (s *LinksService) Get(id string) (string, error) {
	origin, err := s.database.Get(id)
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
