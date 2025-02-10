package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

type NotifyService struct {
	repo     repository.Notify
	Hub      *Hub
	Services *Services
}

func NewNotifyService(repo repository.Notify, hub *Hub) *NotifyService {
	return &NotifyService{repo: repo, Hub: hub}
}

func (s *NotifyService) FindNotifyPopulate(input *domain.NotifyFilter) (domain.Response[domain.Notify], error) {
	return s.repo.FindNotifyPopulate(input)
}

func (s *NotifyService) CreateNotify(userID string, data *domain.NotifyInput) (*domain.Notify, error) {
	var result *domain.Notify

	result, err := s.repo.CreateNotify(userID, data)
	if err != nil {
		return nil, err
	}

	// send by socket.
	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: result.UserTo.Hex(), Content: result, ID: "room1", Service: "notify"})

	// send by push notification.
	userData, err := s.Services.User.GetUser(result.UserTo.Hex())
	if err != nil {
		return nil, err
	}

	// fmt.Println("userData.AuthPrivate.PushToken=", userData.AuthPrivate.PushToken)
	if userData.AuthPrivate.PushToken != "" {
		// To check the token is valid
		pushToken, err := expo.NewExponentPushToken(userData.AuthPrivate.PushToken)
		if err != nil {
			return nil, err
		}

		// Create a new Expo SDK client
		client := expo.NewPushClient(nil)

		// Publish message
		_, err = client.Publish(
			&expo.PushMessage{
				To:       []expo.ExponentPushToken{pushToken},
				Body:     result.Message,
				Data:     map[string]string{"withSome": "data"},
				Sound:    "default",
				Title:    result.Title,
				Priority: expo.DefaultPriority,
			},
		)
		if err != nil {
			return nil, err
		}

	}

	return result, err
}

func (s *NotifyService) UpdateNotify(id string, userID string, data *domain.NotifyInput) (*domain.Notify, error) {
	result, err := s.repo.UpdateNotify(id, userID, data)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *NotifyService) DeleteNotify(id string) (*domain.Notify, error) {
	result, err := s.repo.DeleteNotify(id)
	if err != nil {
		return result, err
	}

	return result, err
}
