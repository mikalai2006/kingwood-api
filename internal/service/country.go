package service

import (
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type CountryService struct {
	repo repository.Country
	i18n config.I18nConfig
}

func NewCountryService(repo repository.Country, i18n config.I18nConfig) *CountryService {
	return &CountryService{repo: repo, i18n: i18n}
}

func (s *CountryService) CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error) {
	return s.repo.CreateCountry(userID, data)
}

func (s *CountryService) GetCountry(id string) (domain.Country, error) {
	return s.repo.GetCountry(id)
}

func (s *CountryService) FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error) {
	return s.repo.FindCountry(params)
}

func (s *CountryService) UpdateCountry(id string, data interface{}) (domain.Country, error) {
	return s.repo.UpdateCountry(id, data)
}
func (s *CountryService) DeleteCountry(id string) (domain.Country, error) {
	return s.repo.DeleteCountry(id)
}
