package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArchiveTaskWorkerService struct {
	repo     repository.ArchiveTaskWorker
	Hub      *Hub
	Services *Services
}

func NewArchiveTaskWorkerService(repo repository.ArchiveTaskWorker, hub *Hub) *ArchiveTaskWorkerService {
	return &ArchiveTaskWorkerService{repo: repo, Hub: hub}
}

func (s *ArchiveTaskWorkerService) CreateArchiveTaskWorker(userID string, data *domain.TaskWorker) (*domain.ArchiveTaskWorker, error) {
	var result *domain.ArchiveTaskWorker

	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	result, err = s.repo.CreateArchiveTaskWorker(userID, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *ArchiveTaskWorkerService) FindArchiveTaskWorker(input *domain.ArchiveTaskWorkerFilter) (domain.Response[domain.ArchiveTaskWorker], error) {
	var result domain.Response[domain.ArchiveTaskWorker]

	if input.Query != "" {
		orders, err := s.Services.Order.FindOrder(&domain.OrderFilter{Query: input.Query})
		if err != nil {
			return result, err
		}

		if len(orders.Data) > 0 {
			for i := range orders.Data {
				input.OrderId = append(input.OrderId, orders.Data[i].ID.Hex())
			}
		} else {
			input.OrderId = []string{primitive.NilObjectID.Hex()}
		}
	}

	result, err := s.repo.FindArchiveTaskWorker(input)

	return result, err
}

func (s *ArchiveTaskWorkerService) DeleteArchiveTaskWorker(id string, userID string) (*domain.ArchiveTaskWorker, error) {
	result, err := s.repo.DeleteArchiveTaskWorker(id)
	if err != nil {
		return result, err
	}

	return result, err
}
