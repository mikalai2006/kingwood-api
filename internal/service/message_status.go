package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type MessageStatusService struct {
	repo repository.MessageStatus
	Hub  *Hub
}

func NewMessageStatusService(repo repository.MessageStatus, Hub *Hub) *MessageStatusService {
	return &MessageStatusService{repo: repo, Hub: Hub}
}

func (s *MessageStatusService) FindMessageStatus(params *domain.MessageStatusFilter) (domain.Response[domain.MessageStatus], error) {
	return s.repo.FindMessageStatus(params)
}

func (s *MessageStatusService) CreateMessageStatus(userID string, data *domain.MessageStatus) (*domain.MessageStatus, error) {
	result, err := s.repo.CreateMessageStatus(userID, data)

	// s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "ADD", Sender: userID, Recipient: result.TakeUserID.Hex(), Content: result, ID: "room1", Service: "MessageStatus"})

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

func (s *MessageStatusService) GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error) {
	return s.repo.GetGroupForUser(userID)
}
