package service

import (
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type LangService struct {
	repo repository.Lang
	i18n config.I18nConfig
}

func NewLangService(repo repository.Lang, i18n config.I18nConfig) *LangService {
	return &LangService{repo: repo, i18n: i18n}
}

func (s *LangService) CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error) {
	return s.repo.CreateLanguage(userID, data)
}

func (s *LangService) GetLanguage(id string) (domain.Language, error) {
	return s.repo.GetLanguage(id)
}

func (s *LangService) FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error) {
	return s.repo.FindLanguage(params)
}

func (s *LangService) UpdateLanguage(id string, data interface{}) (domain.Language, error) {
	return s.repo.UpdateLanguage(id, data)
}
func (s *LangService) DeleteLanguage(id string) (domain.Language, error) {
	return s.repo.DeleteLanguage(id)
}
