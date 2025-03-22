package interfaces

import (
	"context"
	"main/internal/models"
)

type LinksService interface {
	Add(ctx context.Context, originLink models.OriginLink, host string) (string, error)
	AddBatch(ctx context.Context, originLinks []models.OriginLink, host string) ([]models.Result, error)
	Get(ctx context.Context, shortLink string) (string, error)
	Ping() error
	Close() error
}

type UsersService interface {
	Login(ctx context.Context) (int, error)
	Close() error
}
