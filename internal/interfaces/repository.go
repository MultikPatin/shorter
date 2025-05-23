package interfaces

import (
	"context"
	"main/internal/models"
)

// HealthRepository outlines methods for checking system health and readiness.
type HealthRepository interface {
	Ping() error // Performs a basic health check.
}

// LinksRepository governs CRUD operations related to link management.
type LinksRepository interface {
	Add(ctx context.Context, addedLink models.AddedLink) (string, error)                  // Adds a single link.
	AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]models.Result, error) // Adds multiple links in batch.
	Get(ctx context.Context, short string) (string, error)                                // Retrieves the original URL for a given short link.
}

// FileStorageProducer abstracts the process of writing events to a persistent storage medium.
type FileStorageProducer interface {
	WriteEvent(event *models.Event) error // Writes an event to storage.
	Close() error                         // Closes the storage producer gracefully.
}

// FileStorageConsumer encapsulates reading historical events from a persistent storage source.
type FileStorageConsumer interface {
	ReadAllEvents() ([]*models.Event, error) // Reads all available events from storage.
	Close() error                            // Closes the storage consumer properly.
}

// UsersRepository handles user-specific operations such as logging in, fetching links, and deleting links.
type UsersRepository interface {
	Login(ctx context.Context) (int64, error)                   // Logs in a user and assigns a unique identifier.
	GetLinks(ctx context.Context) ([]models.UserLinks, error)   // Retrieves all links created by the user.
	DeleteLinks(ctx context.Context, shortLinks []string) error // Removes specified links created by the user.
}
