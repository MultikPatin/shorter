package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models"
)

// StatsRepository is an in-memory implementation of the StatsRepository interface
type StatsRepository struct {
	db *PostgresDB
}

// NewStatsRepository creates a new in-memory StatsRepository
func NewStatsRepository(db *PostgresDB) *StatsRepository {
	return &StatsRepository{
		db: db,
	}
}

// GetMainStats implements the StatsRepository interface to retrieve main statistics
// from in-memory storage.
func (r *StatsRepository) GetMainStats(ctx context.Context) (models.Stats, error) {
	var urls int
	var users int

	err := r.db.Connection.QueryRowContext(ctx, getMainStatsQuery).Scan(&urls, &users)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Stats{}, fmt.Errorf("no stats found")
	} else if err != nil {
		return models.Stats{}, err
	}

	return models.Stats{
		Urls:  urls,
		Users: users,
	}, nil
}
