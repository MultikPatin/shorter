package app

import (
	"fmt"
	"github.com/google/uuid"
)

type InMemoryDB struct {
	links map[uuid.UUID]string
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		links: make(map[uuid.UUID]string),
	}
}

func (db *InMemoryDB) AddLink(id uuid.UUID, link string) (uuid.UUID, error) {
	db.links[id] = link
	return id, nil
}

func (db *InMemoryDB) GetByID(id uuid.UUID) (string, error) {
	if link, ok := db.links[id]; ok {
		return link, nil
	}
	return "", fmt.Errorf("user with id %d not found", id)
}
