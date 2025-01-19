package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type AppErrorService struct {
	repo     repository.AppError
	Hub      *Hub
	Services *Services
}

func NewAppErrorService(repo repository.AppError, hub *Hub) *AppErrorService {
	return &AppErrorService{repo: repo, Hub: hub}
}

func (s *AppErrorService) FindAppError(input *domain.AppErrorFilter) (domain.Response[domain.AppError], error) {
	return s.repo.FindAppError(input)
}

func (s *AppErrorService) CreateAppError(userID string, data *domain.AppError) (*domain.AppError, error) {
	var result *domain.AppError

	result, err := s.repo.CreateAppError(userID, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *AppErrorService) UpdateAppError(id string, userID string, data *domain.AppErrorInput) (*domain.AppError, error) {
	var result *domain.AppError

	result, err := s.repo.UpdateAppError(id, userID, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *AppErrorService) DeleteAppError(id string, userID string) (*domain.AppError, error) {
	var result *domain.AppError

	result, err := s.repo.DeleteAppError(id, userID)
	if err != nil {
		return nil, err
	}

	return result, err
}
