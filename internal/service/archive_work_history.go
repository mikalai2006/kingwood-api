package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type ArchiveWorkHistoryService struct {
	repo     repository.ArchiveWorkHistory
	Hub      *Hub
	Services *Services
}

func NewArchiveWorkHistoryService(repo repository.ArchiveWorkHistory, hub *Hub) *ArchiveWorkHistoryService {
	return &ArchiveWorkHistoryService{repo: repo, Hub: hub}
}

func (s *ArchiveWorkHistoryService) FindArchiveWorkHistory(input domain.ArchiveWorkHistoryFilter) (domain.Response[domain.ArchiveWorkHistory], error) {
	return s.repo.FindArchiveWorkHistory(input)
}

func (s *ArchiveWorkHistoryService) CreateArchiveWorkHistory(userID string, data *domain.WorkHistory) (*domain.ArchiveWorkHistory, error) {
	var result *domain.ArchiveWorkHistory

	result, err := s.repo.CreateArchiveWorkHistory(userID, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *ArchiveWorkHistoryService) DeleteArchiveWorkHistory(id string, userID string) (*domain.ArchiveWorkHistory, error) {
	var result *domain.ArchiveWorkHistory

	result, err := s.repo.DeleteArchiveWorkHistory(id)

	return result, err
}

func (s *ArchiveWorkHistoryService) ClearArchiveWorkHistory(userID string) error {
	return s.repo.ClearArchiveWorkHistory(userID)
}
