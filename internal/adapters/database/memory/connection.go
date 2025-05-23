package memory

import (
	"go.uber.org/zap"
	"main/internal/interfaces"
)

// InMemoryDB represents an in-memory database backed by file storage.
type InMemoryDB struct {
	links      map[string]string              // Map holding the short-to-long URL mappings.
	producerFS interfaces.FileStorageProducer // Interface implementation for writing to persistent storage.
	consumerFS interfaces.FileStorageConsumer // Interface implementation for reading from persistent storage.
}

// Close closes both the producer and consumer file storages.
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

// Ping performs a health check operation, always returning success since it's an in-memory DB.
func (db *InMemoryDB) Ping() error {
	return nil
}

// NewInMemoryDB initializes a new in-memory database instance with file-backed persistence.
func NewInMemoryDB(path string, logger *zap.SugaredLogger) (*InMemoryDB, error) {
	producerFS, err := NewFileProducer(path)
	if err != nil {
		logger.Infow("Failed to create producer file storage", "error", err.Error())
		return nil, err
	}
	consumerFS, err := NewFileConsumer(path)
	if err != nil {
		logger.Infow("Failed to create consumer file storage", "error", err.Error())
		return nil, err
	}
	db := InMemoryDB{
		links:      make(map[string]string),
		producerFS: producerFS,
		consumerFS: consumerFS,
	}
	err = db.loadFromFile()
	if err != nil {
		logger.Infow("Failed to load events from file", "error", err.Error())
		return nil, err
	}
	return &db, nil
}

// loadFromFile loads existing URL mapping events from the consumer file storage into memory.
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
