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
	var result *domain.ArchiveOrder

	// delete images.
	allImages, err := s.Services.ArchiveImage.FindArchiveImage(&domain.ArchiveImageFilter{ServiceId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allImages.Data {
		_, err = s.Services.ArchiveImage.DeleteArchiveImage(allImages.Data[i].ID.Hex())
		if err != nil {
			return result, err
		}
	}
	// delete taskWorkers.
	allTaskWorkers, err := s.Services.ArchiveTaskWorker.FindArchiveTaskWorker(&domain.ArchiveTaskWorkerFilter{OrderId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allTaskWorkers.Data {
		_, err = s.Services.ArchiveTaskWorker.DeleteArchiveTaskWorker(allTaskWorkers.Data[i].ID.Hex(), userID)
		if err != nil {
			return result, err
		}
	}

	// delete task.
	allTasks, err := s.Services.ArchiveTask.FindArchiveTask(domain.ArchiveTaskFilter{OrderId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allTasks.Data {
		_, err = s.Services.ArchiveTask.DeleteArchiveTask(allTasks.Data[i].ID.Hex(), userID)
		if err != nil {
			return result, err
		}
	}

	// delete messages.
	allMessages, err := s.Services.ArchiveMessage.FindArchiveMessage(&domain.ArchiveMessageFilter{OrderID: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allMessages.Data {
		_, err = s.Services.ArchiveMessage.DeleteArchiveMessage(allMessages.Data[i].ID.Hex())
		if err != nil {
			return result, err
		}
	}

	// delete workHistory.
	allWorkHistory, err := s.Services.ArchiveWorkHistory.FindArchiveWorkHistory(domain.ArchiveWorkHistoryFilter{OrderId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allWorkHistory.Data {
		_, err = s.Services.ArchiveWorkHistory.DeleteArchiveWorkHistory(allWorkHistory.Data[i].ID.Hex(), userID)
		if err != nil {
			return result, err
		}
	}

	return s.repo.DeleteArchiveOrder(id)
}
