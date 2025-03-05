package memory

import (
	"fmt"
	"go.uber.org/zap"
	"main/internal/models"
)

type producer interface {
	WriteEvent(event *models.Event) error
	Close() error
}

type consumer interface {
	ReadAllEvents() ([]*models.Event, error)
	Close() error
}

type InMemoryDB struct {
	links      map[string]string
	producerFS producer
	consumerFS consumer
}

func NewInMemoryDB(path string, logger *zap.SugaredLogger) (*InMemoryDB, error) {
	producerFS, err := NewFileProducer(path)
	if err != nil {
		logger.Infow(
			"Create producerFS",
			"error", err.Error(),
		)
	}

	consumerFS, err := NewFileConsumer(path)
	if err != nil {
		logger.Infow(
			"Create consumerFS",
			"error", err.Error(),
		)
	}

	db := InMemoryDB{
		links:      make(map[string]string),
		producerFS: producerFS,
		consumerFS: consumerFS,
	}

	err = db.LoadFromFile()
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
	db.producerFS.Close()
	db.consumerFS.Close()
	return nil
}

func (db *InMemoryDB) LoadFromFile() error {
	events, err := db.consumerFS.ReadAllEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		db.links[event.Short] = event.Original
	}
	return nil
}

func (db *InMemoryDB) Add(id string, link string) (string, error) {
	db.links[id] = link
	l := len(db.links)

	event := &models.Event{
		ID:       l,
		Original: link,
		Short:    id,
	}
	if err := db.producerFS.WriteEvent(event); err != nil {
		return "", err
	}

	return id, nil
}

func (db *InMemoryDB) Get(id string) (string, error) {
	if link, ok := db.links[id]; ok {
		return link, nil
	}
	return "", fmt.Errorf("user with id %s not found", id)
}
