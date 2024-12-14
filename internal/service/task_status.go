package service

import (
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type TaskStatusService struct {
	repo repository.TaskStatus
	i18n config.I18nConfig
}

func NewTaskStatusService(repo repository.TaskStatus, i18n config.I18nConfig) *TaskStatusService {
	return &TaskStatusService{repo: repo, i18n: i18n}
}

func (s *TaskStatusService) FindTaskStatus(params domain.RequestParams) (domain.Response[domain.TaskStatus], error) {
	return s.repo.FindTaskStatus(params)
}

func (s *TaskStatusService) CreateTaskStatus(userID string, data *domain.TaskStatus) (*domain.TaskStatus, error) {
	return s.repo.CreateTaskStatus(userID, data)
}

func (s *TaskStatusService) UpdateTaskStatus(id string, userID string, data *domain.TaskStatusInput) (*domain.TaskStatus, error) {
	return s.repo.UpdateTaskStatus(id, userID, data)
}

func (s *TaskStatusService) DeleteTaskStatus(id string) (domain.TaskStatus, error) {
	return s.repo.DeleteTaskStatus(id)
}
