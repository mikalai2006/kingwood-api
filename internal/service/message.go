package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type MessageService struct {
	repo               repository.Message
	Hub                *Hub
	messageRoomService *MessageRoomService
}

func NewMessageService(repo repository.Message, Hub *Hub, messageRoomService *MessageRoomService) *MessageService {
	return &MessageService{repo: repo, Hub: Hub, messageRoomService: messageRoomService}
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

	return result, err
}

func (s *MessageService) UpdateMessage(id string, userID string, data *domain.MessageInput) (*domain.Message, error) {
	return s.repo.UpdateMessage(id, userID, data)
}

func (s *MessageService) DeleteMessage(id string) (domain.Message, error) {
	return s.repo.DeleteMessage(id)
}

func (s *MessageService) GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error) {
	return s.repo.GetGroupForUser(userID)
}
