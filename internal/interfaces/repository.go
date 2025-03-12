package interfaces

import (
	"context"
	"main/internal/models"
)

type LinksRepository interface {
	Add(ctx context.Context, addedLink models.AddedLink) (string, error)
	AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]models.Result, error)
	Get(ctx context.Context, short string) (string, error)
	Close() error
	Ping() error
}

type FileStorageProducer interface {
	WriteEvent(event *models.Event) error
	Close() error
}

type FileStorageConsumer interface {
	ReadAllEvents() ([]*models.Event, error)
	Close() error
}
