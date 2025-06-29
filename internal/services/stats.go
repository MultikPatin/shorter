package services

import (
	"context"
	"errors"
	"main/internal/config"
	"main/internal/interfaces"
	"main/internal/models"
	"time"
)

// ErrNotTrustedSubnet is returned when a request comes from an untrusted subnet
var ErrNotTrustedSubnet = errors.New("not a trusted subnet")

// StatsService provides statistics-related operations
type StatsService struct {
	statsRepository interfaces.StatsRepository
}

// NewStatsService creates a new StatsService instance
func NewStatsService(c *config.Config, statsRepository interfaces.StatsRepository) *StatsService {
	trustedSubnet = c.TrustedSubnet
	return &StatsService{
		statsRepository: statsRepository,
	}
}

// GetMainStats retrieves main statistics with a timeout
func (s *StatsService) GetMainStats(ctx context.Context, ip string) (models.Stats, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if !isTrustedSubnet(ip) {
		return models.Stats{}, ErrNotTrustedSubnet
	}

	mainStats, err := s.statsRepository.GetMainStats(ctx)
	if err != nil {
		return models.Stats{}, err
	}
	return mainStats, nil
}
