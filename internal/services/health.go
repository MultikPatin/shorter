package services // Package services implements business logic for health-related operations.

import (
	"main/internal/interfaces"
)

// HealthService encapsulates the business logic for health checks.
type HealthService struct {
	healthRepository interfaces.HealthRepository // Dependency for accessing health-related repository methods.
}

// NewHealthService constructs a new HealthService instance tied to a specific health repository.
func NewHealthService(healthRepository interfaces.HealthRepository) *HealthService {
	return &HealthService{
		healthRepository: healthRepository,
	}
}

// Ping delegates a health check request to the underlying repository.
func (s *HealthService) Ping() error {
	err := s.healthRepository.Ping()
	return err
}
