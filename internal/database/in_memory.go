package database

import (
	"fmt"
)

type producer interface {
	WriteEvent(event *Event) error
	Close() error
}

type consumer interface {
	ReadEvent() ([]*Event, error)
	Close() error
}

type fileStorage interface {
	producer
	consumer
}

type InMemoryDB struct {
	links      map[string]string
	producerFS fileStorage
	consumerFS fileStorage
}

func NewInMemoryDB(producerFS, consumerFS fileStorage) *InMemoryDB {
	return &InMemoryDB{
		links:      make(map[string]string),
		producerFS: producerFS,
		consumerFS: consumerFS,
	}
}

func (db *InMemoryDB) LoadFromFile() error {
	events, err := db.consumerFS.ReadEvent()
	if err != nil {
		return err
	}
	for _, event := range events {
		db.links[event.Short] = event.Original
	}
	return nil
}

func (db *InMemoryDB) AddLink(id string, link string) (string, error) {
	db.links[id] = link
	l := len(db.links)

	event := &Event{
		ID:       l,
		Original: link,
		Short:    id,
	}
	if err := db.producerFS.WriteEvent(event); err != nil {
		return "", err
	}

	return id, nil
}

func (db *InMemoryDB) GetByID(id string) (string, error) {
	if link, ok := db.links[id]; ok {
		return link, nil
	}
	return "", fmt.Errorf("user with id %s not found", id)
}
