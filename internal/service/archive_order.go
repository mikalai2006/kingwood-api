package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type ArchiveOrderService struct {
	repo     repository.ArchiveOrder
	Services *Services
}

func NewArchiveOrderService(repo repository.ArchiveOrder) *ArchiveOrderService {
	return &ArchiveOrderService{repo: repo}
}

func (s *ArchiveOrderService) CreateArchiveOrder(userID string, data *domain.Order) (*domain.ArchiveOrder, error) {
	return s.repo.CreateArchiveOrder(userID, data)
}

func (s *ArchiveOrderService) FindArchiveOrder(input *domain.ArchiveOrderFilter) (domain.Response[domain.ArchiveOrder], error) {
	return s.repo.FindArchiveOrder(input)
}

func (s *ArchiveOrderService) DeleteArchiveOrder(id string, userID string) (*domain.ArchiveOrder, error) {
	return s.repo.DeleteArchiveOrder(id)
}
