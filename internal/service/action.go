package service

import (
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type ActionService struct {
	repo repository.Action
	i18n config.I18nConfig
}

func NewActionService(repo repository.Action, i18n config.I18nConfig) *ActionService {
	return &ActionService{repo: repo, i18n: i18n}
}

func (s *ActionService) FindAction(params domain.RequestParams) (domain.Response[model.Action], error) {
	return s.repo.FindAction(params)
}

func (s *ActionService) GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error) {
	return s.repo.GetAllAction(params)
}

func (s *ActionService) CreateAction(userID string, tag *model.ActionInput) (*model.Action, error) {
	return s.repo.CreateAction(userID, tag)
}

func (s *ActionService) UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error) {
	return s.repo.UpdateAction(id, userID, data)
}

func (s *ActionService) DeleteAction(id string) (model.Action, error) {
	return s.repo.DeleteAction(id)
}
