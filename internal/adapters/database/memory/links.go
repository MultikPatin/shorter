package memory

import (
	"context"
	"fmt"
	"main/internal/models"
)

// LinksRepository manages operations related to adding and retrieving links from the in-memory database.
type LinksRepository struct {
	db *InMemoryDB // Pointer to the in-memory database instance.
}

// NewLinksRepository creates a new instance of LinksRepository bound to a specific InMemoryDB.
func NewLinksRepository(db *InMemoryDB) *LinksRepository {
	return &LinksRepository{
		db: db,
	}
}

// Add inserts a new link into the database and persists the change to file storage.
func (r *LinksRepository) Add(ctx context.Context, addedLink models.AddedLink) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		r.db.links[addedLink.Short] = addedLink.Origin

		event := &models.Event{
			ID:     len(r.db.links),
			Origin: addedLink.Origin,
			Short:  addedLink.Short,
		}
		if err := r.db.producerFS.WriteEvent(event); err != nil {
			return "", err
		}

		return addedLink.Short, nil
	}
}

// AddBatch adds multiple links in batch fashion, persisting changes to file storage.
func (r *LinksRepository) AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]models.Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var results []models.Result

		for _, addedLink := range addedLinks {
			r.db.links[addedLink.Short] = addedLink.Origin

			event := &models.Event{
				ID:     len(r.db.links),
				Origin: addedLink.Origin,
				Short:  addedLink.Short,
			}
			if err := r.db.producerFS.WriteEvent(event); err != nil {
				return nil, err
			}
			result := models.Result{
				CorrelationID: addedLink.CorrelationID,
				Result:        addedLink.Short,
			}
			results = append(results, result)
		}
		return results, nil
	}
}

// Get retrieves the original URL corresponding to a given shortened link.
func (r *LinksRepository) Get(ctx context.Context, short string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if link, ok := r.db.links[short]; ok {
			return link, nil
		}
		return "", fmt.Errorf("link with short code '%s' not found", short)
	}
}
