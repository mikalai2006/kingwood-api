package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type AnalyticService struct {
	repo     repository.Analytic
	Hub      *Hub
	Services *Services
}

func NewAnalyticService(repo repository.Analytic, hub *Hub) *AnalyticService {
	return &AnalyticService{repo: repo, Hub: hub}
}

func (s *AnalyticService) GetAnalytic() (domain.Analytic, error) {
	return s.repo.GetAnalytic()
}
