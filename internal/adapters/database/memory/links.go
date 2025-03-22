package memory

import (
	"context"
	"fmt"
	"main/internal/models"
)

type LinksRepository struct {
	db *InMemoryDB
}

func NewLinksRepository(db *InMemoryDB) *LinksRepository {
	return &LinksRepository{
		db: db,
	}
}

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

func (r *LinksRepository) AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]models.Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var shortLinks []models.Result

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
			shortLink := models.Result{
				CorrelationID: addedLink.CorrelationID,
				Result:        addedLink.Short,
			}
			shortLinks = append(shortLinks, shortLink)
		}
		return shortLinks, nil
	}
}

func (r *LinksRepository) Get(ctx context.Context, short string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if link, ok := r.db.links[short]; ok {
			return link, nil
		}
		return "", fmt.Errorf("user with short %s not found", short)
	}
}
