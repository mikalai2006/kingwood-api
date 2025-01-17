package service

import (
	"fmt"
	"os"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type MessageRoomService struct {
	repo repository.MessageRoom
	Hub  *Hub
}

func NewMessageRoomService(repo repository.MessageRoom, Hub *Hub) *MessageRoomService {
	return &MessageRoomService{repo: repo, Hub: Hub}
}

func (s *MessageRoomService) FindMessageRoom(params *domain.MessageRoomFilter) (domain.Response[domain.MessageRoom], error) {
	return s.repo.FindMessageRoom(params)
}

func (s *MessageRoomService) CreateMessageRoom(userID string, data *domain.MessageRoom) (*domain.MessageRoom, error) {
	result, err := s.repo.CreateMessageRoom(userID, data)

	// s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "ADD", Sender: userID, Recipient: result.TakeUserID.Hex(), Content: result, ID: "room1", Service: "messageRoom"})

	return result, err
}

func (s *MessageRoomService) UpdateMessageRoom(id string, userID string, data *domain.MessageRoom) (*domain.MessageRoom, error) {
	// return s.repo.UpdateMessageRoom(id, userID, data)
	result, err := s.repo.UpdateMessageRoom(id, userID, data)

	// if result != nil && err == nil {
	// 	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "patch", Sender: userID, Recipient: result.TakeUserID.Hex(), Content: result, ID: "room1", Service: "messageRoom"})
	// }

	return result, err
}

func (s *MessageRoomService) DeleteMessageRoom(id string) (domain.MessageRoom, error) {
	result, err := s.repo.DeleteMessageRoom(id)

	// Delete dir with images for room.
	pathOfRemove := fmt.Sprintf("public/%s/%s", "message", result.ID.Hex())
	os.RemoveAll(pathOfRemove)

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

func (s *MessageRoomService) GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error) {
	return s.repo.GetGroupForUser(userID)
}
