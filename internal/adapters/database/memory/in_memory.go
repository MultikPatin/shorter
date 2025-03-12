package memory

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"main/internal/interfaces"
	"main/internal/models"
)

type InMemoryDB struct {
	links      map[string]string
	producerFS interfaces.FileStorageProducer
	consumerFS interfaces.FileStorageConsumer
}

func NewInMemoryRepository(path string, logger *zap.SugaredLogger) (*InMemoryDB, error) {
	producerFS, err := NewFileProducer(path)
	if err != nil {
		logger.Infow(
			"Create producerFS",
			"error", err.Error(),
		)
		return nil, err
	}
	consumerFS, err := NewFileConsumer(path)
	if err != nil {
		logger.Infow(
			"Create consumerFS",
			"error", err.Error(),
		)
		return nil, err
	}
	db := InMemoryDB{
		links:      make(map[string]string),
		producerFS: producerFS,
		consumerFS: consumerFS,
	}
	err = db.loadFromFile()
	if err != nil {
		logger.Infow(
			"Load events from file",
			"error", err.Error(),
		)
		return nil, err
	}
	return &db, err
}

func (db *InMemoryDB) Close() error {
	err := db.producerFS.Close()
	if err != nil {
		return err
	}
	err = db.consumerFS.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *InMemoryDB) Ping() error {
	return nil
}

func (db *InMemoryDB) loadFromFile() error {
	events, err := db.consumerFS.ReadAllEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		db.links[event.Short] = event.Origin
	}
	return nil
}

func (db *InMemoryDB) Add(ctx context.Context, addedLink models.AddedLink) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		db.links[addedLink.Short] = addedLink.Origin

		event := &models.Event{
			ID:     len(db.links),
			Origin: addedLink.Origin,
			Short:  addedLink.Short,
		}
		if err := db.producerFS.WriteEvent(event); err != nil {
			return "", err
		}

		return addedLink.Short, nil
	}
}

func (db *InMemoryDB) AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]models.Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var shortLinks []models.Result

		for _, addedLink := range addedLinks {
			db.links[addedLink.Short] = addedLink.Origin

			event := &models.Event{
				ID:     len(db.links),
				Origin: addedLink.Origin,
				Short:  addedLink.Short,
			}
			if err := db.producerFS.WriteEvent(event); err != nil {
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

func (db *InMemoryDB) Get(ctx context.Context, short string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if link, ok := db.links[short]; ok {
			return link, nil
		}
		return "", fmt.Errorf("user with short %s not found", short)
	}
}
