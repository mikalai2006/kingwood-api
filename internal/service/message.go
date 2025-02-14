package service

import (
	"fmt"
	"os"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type MessageService struct {
	repo        repository.Message
	Hub         *Hub
	imageConfig config.IImageConfig
	Services    *Services
}

func NewMessageService(repo repository.Message, Hub *Hub, imageConfig config.IImageConfig) *MessageService {
	return &MessageService{repo: repo, Hub: Hub, imageConfig: imageConfig}
}

func (s *MessageService) FindMessage(params *domain.MessageFilter) (domain.Response[domain.Message], error) {
	return s.repo.FindMessage(params)
}

func (s *MessageService) CreateMessage(userID string, data *domain.MessageInput) (*domain.Message, error) {
	result, err := s.repo.CreateMessage(userID, data)

	// room, err := s.messageRoomService.FindMessageRoom(&domain.MessageRoomFilter{ID: &result.RoomID})

	// if err == nil && len(room.Data) > 0 {
	// 	sobesednikID := room.Data[0].UserID
	// 	if room.Data[0].UserID == result.UserID {
	// 		sobesednikID = room.Data[0].TakeUserID
	// 	}

	// 	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "ADD", Sender: userID, Recipient: sobesednikID.Hex(), Content: result, ID: "room1", Service: "message"})
	// }
	if err != nil {
		return result, err
	}

	// messages, err := s.Services.Message.FindMessage(&domain.MessageFilter{ID: result.MessageID.Hex()})
	// if err != nil {
	// 	return result,err
	// }

	taskWorkers, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{OrderId: []string{result.OrderID.Hex()}})
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "message"})

	for i := range taskWorkers.Data {
		// add notify.
		_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo:     taskWorkers.Data[i].WorkerId.Hex(),
			Title:      domain.NewMessageTitle,
			Message:    fmt.Sprintf(domain.NewMessage, taskWorkers.Data[i].Order.Number, taskWorkers.Data[i].Order.Name, taskWorkers.Data[i].Object.Name),
			Link:       "/[orderId]/message",
			LinkOption: map[string]interface{}{"orderId": result.OrderID.Hex()},
		})

	}

	return result, err
}

func (s *MessageService) UpdateMessage(id string, userID string, data *domain.MessageInput) (*domain.Message, error) {
	return s.repo.UpdateMessage(id, userID, data)
}

func (s *MessageService) DeleteMessage(id string) (domain.Message, error) {
	result, err := s.repo.DeleteMessage(id)

	// Delete images for message.
	for i := range result.Images {
		objImage := result.Images[i]
		pathDir := fmt.Sprintf("public/%s", objImage.Service)

		path := fmt.Sprintf("%s/%s/%s%s", pathDir, objImage.ServiceID, objImage.Path, objImage.Ext)
		os.Remove(path)

		for j := range s.imageConfig.Sizes {
			path := fmt.Sprintf("%s/%s/%s-%s%s", pathDir, objImage.ServiceID, s.imageConfig.Sizes[j].Prefix, objImage.Path, objImage.Ext)
			os.Remove(path)
		}
	}

	return result, err
}

func (s *MessageService) GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error) {
	return s.repo.GetGroupForUser(userID)
}
