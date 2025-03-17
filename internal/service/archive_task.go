package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type ArchiveTaskService struct {
	repo     repository.ArchiveTask
	Services *Services
}

func NewArchiveTaskService(repo repository.ArchiveTask) *ArchiveTaskService {
	return &ArchiveTaskService{repo: repo}
}

func (s *ArchiveTaskService) CreateArchiveTask(userID string, data *domain.Task) (*domain.ArchiveTask, error) {
	return s.repo.CreateArchiveTask(userID, data)
}

func (s *ArchiveTaskService) FindArchiveTask(filter domain.ArchiveTaskFilter) (domain.Response[domain.ArchiveTask], error) {
	return s.repo.FindArchiveTask(filter)
}

func (s *ArchiveTaskService) DeleteArchiveTask(id string, userID string) (*domain.ArchiveTask, error) {
	result, err := s.repo.DeleteArchiveTask(id)

	return result, err
}
