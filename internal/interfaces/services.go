package interfaces // Package interfaces defines service-layer contracts for business logic components.

import (
	"context"
	"main/internal/models"
)

// HealthService encapsulates high-level health monitoring functionalities.
type HealthService interface {
	Ping() error // Executes a basic health check.
}

// LinksService governs core link-related operations including creation, batch processing, and retrieval.
type LinksService interface {
	Add(ctx context.Context, originLink models.OriginLink, host string) (string, error)                  // Adds a single link.
	AddBatch(ctx context.Context, originLinks []models.OriginLink, host string) ([]models.Result, error) // Batch-adds multiple links.
	Get(ctx context.Context, shortLink string) (string, error)                                           // Retrieves the original URL for a given short link.
}

// UsersService manages user-specific activities such as login, link retrieval, and deletion.
type UsersService interface {
	Login() (int64, error)                                                 // Logs in a user and generates a unique identifier.
	GetLinks(ctx context.Context, host string) ([]models.UserLinks, error) // Retrieves all links created by the logged-in user.
	DeleteLinks(ctx context.Context, shortLinks []string) error            // Deletes specified links created by the user.
}
