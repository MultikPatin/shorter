package services

import (
	"main/internal/interfaces"
)

type HealthService struct {
	healthRepository interfaces.HealthRepository
}

func NewHealthService(healthRepository interfaces.HealthRepository) *HealthService {
	return &HealthService{
		healthRepository: healthRepository,
	}
}

func (s *HealthService) Ping() error {
	err := s.healthRepository.Ping()
	return err
}
