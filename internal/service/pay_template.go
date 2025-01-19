package service

import (
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type PayTemplateService struct {
	repo repository.PayTemplate
	i18n config.I18nConfig
}

func NewPayTemplateService(repo repository.PayTemplate, i18n config.I18nConfig) *PayTemplateService {
	return &PayTemplateService{repo: repo, i18n: i18n}
}

func (s *PayTemplateService) FindPayTemplate(params domain.RequestParams) (domain.Response[domain.PayTemplate], error) {
	return s.repo.FindPayTemplate(params)
}

func (s *PayTemplateService) CreatePayTemplate(userID string, data *domain.PayTemplate) (*domain.PayTemplate, error) {
	return s.repo.CreatePayTemplate(userID, data)
}

func (s *PayTemplateService) UpdatePayTemplate(id string, userID string, data *domain.PayTemplateInput) (*domain.PayTemplate, error) {
	return s.repo.UpdatePayTemplate(id, userID, data)
}

func (s *PayTemplateService) DeletePayTemplate(id string) (domain.PayTemplate, error) {
	return s.repo.DeletePayTemplate(id)
}
