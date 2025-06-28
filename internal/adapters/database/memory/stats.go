package memory

import (
	"context"
	"fmt"
	"main/internal/models"
)

// StatsRepository is an in-memory implementation of the StatsRepository interface
type StatsRepository struct {
	db *InMemoryDB
}

// NewStatsRepository creates a new in-memory StatsRepository
func NewStatsRepository(db *InMemoryDB) *StatsRepository {
	return &StatsRepository{
		db: db,
	}
}

// GetMainStats implements the StatsRepository interface to retrieve main statistics
// from in-memory storage.
func (db *StatsRepository) GetMainStats(ctx context.Context) (models.Stats, error) {
	return models.Stats{}, fmt.Errorf("not implemented")
}
