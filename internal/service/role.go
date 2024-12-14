package service

import (
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type RoleService struct {
	repo repository.Role
	i18n config.I18nConfig
}

func NewRoleService(repo repository.Role, i18n config.I18nConfig) *RoleService {
	return &RoleService{repo: repo, i18n: i18n}
}

func (s *RoleService) CreateRole(userID string, data *domain.RoleInput) (domain.Role, error) {
	return s.repo.CreateRole(userID, data)
}

func (s *RoleService) GetRole(id string) (domain.Role, error) {
	return s.repo.GetRole(id)
}

func (s *RoleService) FindRole(params domain.RequestParams) (domain.Response[domain.Role], error) {
	return s.repo.FindRole(params)
}

func (s *RoleService) UpdateRole(id string, data interface{}) (domain.Role, error) {
	return s.repo.UpdateRole(id, data)
}
func (s *RoleService) DeleteRole(id string) (domain.Role, error) {
	return s.repo.DeleteRole(id)
}
