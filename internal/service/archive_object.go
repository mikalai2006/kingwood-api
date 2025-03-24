package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArchiveObjectService struct {
	repo     repository.ArchiveObject
	Hub      *Hub
	Services *Services
}

func NewArchiveObjectService(repo repository.ArchiveObject, hub *Hub) *ArchiveObjectService {
	return &ArchiveObjectService{repo: repo, Hub: hub}
}

func (s *ArchiveObjectService) FindArchiveObject(input *domain.ArchiveObjectFilter) (domain.Response[domain.ArchiveObject], error) {
	return s.repo.FindArchiveObject(input)
}

func (s *ArchiveObjectService) CreateArchiveObject(userID string, data *domain.Object) (*domain.ArchiveObject, error) {
	var result *domain.ArchiveObject

	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	result, err = s.repo.CreateArchiveObject(userID, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *ArchiveObjectService) DeleteArchiveObject(id string, userID string) (*domain.ArchiveObject, error) {
	var result *domain.ArchiveObject

	// // delete images.
	// allImages, err := s.Services.Image.FindImage(domain.RequestParams{Filter: bson.D{{"serviceId", id}}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allImages.Data {
	// 	_, err = s.Services.Image.DeleteImage(userID, allImages.Data[i].ID.Hex())
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }
	// // delete taskWorkers.
	// allTaskWorkers, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{OrderId: []string{id}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allTaskWorkers.Data {
	// 	_, err = s.Services.TaskWorker.DeleteTaskWorker(allTaskWorkers.Data[i].ID.Hex(), userID, false)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }

	// // delete task.
	// allTasks, err := s.Services.Task.FindTaskPopulate(domain.TaskFilter{OrderId: []string{id}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allTasks.Data {
	// 	_, err = s.Services.Task.DeleteTask(allTasks.Data[i].ID.Hex(), userID, false)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }

	// // delete messages.
	// allMessages, err := s.Services.Message.FindMessage(&domain.MessageFilter{OrderID: []string{id}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allMessages.Data {
	// 	_, err = s.Services.Message.DeleteMessage(allMessages.Data[i].ID.Hex(), userID)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }

	// // delete workHistory.
	// allWorkHistory, err := s.Services.WorkHistory.FindWorkHistory(domain.WorkHistoryFilter{OrderId: []string{id}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allWorkHistory.Data {
	// 	_, err = s.Services.WorkHistory.DeleteWorkHistory(allWorkHistory.Data[i].ID.Hex(), userID)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }

	result, err := s.repo.DeleteArchiveObject(id, userID)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "archiveObject"})

	return result, err
}
