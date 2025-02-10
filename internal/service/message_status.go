package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type MessageStatusService struct {
	repo     repository.MessageStatus
	Hub      *Hub
	Services *Services
}

func NewMessageStatusService(repo repository.MessageStatus, Hub *Hub) *MessageStatusService {
	return &MessageStatusService{repo: repo, Hub: Hub}
}

func (s *MessageStatusService) FindMessageStatus(params *domain.MessageStatusFilter) (domain.Response[domain.MessageStatus], error) {
	return s.repo.FindMessageStatus(params)
}

func (s *MessageStatusService) CreateMessageStatus(userID string, data *domain.MessageStatus) (*domain.MessageStatus, error) {
	result, err := s.repo.CreateMessageStatus(userID, data)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "message"})

	// messages, err := s.Services.Message.FindMessage(&domain.MessageFilter{ID: result.MessageID.Hex()})
	// if err != nil {
	// 	return result, err
	// }

	// if len(messages.Data) > 0 {

	// 	taskWorkers, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{OrderId: []string{messages.Data[0].ID.Hex()}})
	// 	if err != nil {
	// 		return result, err
	// 	}

	// 	for i := range taskWorkers.Data {
	// 		// // add notify.
	// 		// _, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
	// 		// 	UserTo:  taskWorkers.Data[i].WorkerId.Hex(),
	// 		// 	Title:   domain.NewMessageTitle,
	// 		// 	Message: fmt.Sprintf(domain.NewMessage, taskWorkers.Data[i].Order.Number, taskWorkers.Data[i].Order.Name, taskWorkers.Data[i].Object.Name),
	// 		// })

	// 		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "message"})

	// 	}
	// }

	return result, err
}

func (s *MessageStatusService) UpdateMessageStatus(id string, userID string, data *domain.MessageStatus) (*domain.MessageStatus, error) {
	// return s.repo.UpdateMessageStatus(id, userID, data)
	result, err := s.repo.UpdateMessageStatus(id, userID, data)

	// if result != nil && err == nil {
	// 	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "patch", Sender: userID, Recipient: result.TakeUserID.Hex(), Content: result, ID: "room1", Service: "MessageStatus"})
	// }

	return result, err
}

func (s *MessageStatusService) DeleteMessageStatus(id string) (domain.MessageStatus, error) {
	result, err := s.repo.DeleteMessageStatus(id)

	// // Delete dir with images for room.
	// pathOfRemove := fmt.Sprintf("public/%s/%s", "message", result.ID.Hex())
	// os.RemoveAll(pathOfRemove)

	// isEmpty, err := utils.IsEmptyDir(pathOfRemove)
	// if err != nil {
	// 	return result, err
	// }
	// if isEmpty {
	// 	err = os.Remove(pathOfRemove)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }

	return result, err
}

// func (s *MessageStatusService) GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error) {
// 	return s.repo.GetGroupForUser(userID)
// }
