package database

import (
	"fmt"
)

type InMemoryDB struct {
	links map[string]string
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		links: make(map[string]string),
	}
}

func (db *InMemoryDB) AddLink(id string, link string) (string, error) {
	db.links[id] = link
	return id, nil
}

func (db *InMemoryDB) GetByID(id string) (string, error) {
	if link, ok := db.links[id]; ok {
		return link, nil
	}
	return "", fmt.Errorf("user with id %s not found", id)
}
