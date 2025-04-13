package service

import (
	"fmt"
	"time"

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
		fmt.Println("Sent push successfully!")
	} else {
		fmt.Println("Sent push wrong!", userData.AuthPrivate)
	}

	return result, err
}

func (s *NotifyService) UpdateNotify(id string, userID string, data *domain.NotifyInput) (*domain.Notify, error) {
	existNotify, err := s.repo.FindNotifyPopulate(&domain.NotifyFilter{ID: []*string{&id}})
	if err != nil {
		return nil, err
	}

	if len(existNotify.Data) > 0 && existNotify.Data[0].Status == 1 {
		return &existNotify.Data[0], err
	}

	result, err := s.repo.UpdateNotify(id, userID, data)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *NotifyService) DeleteNotify(id string, userID string, createArchive bool) (*domain.Notify, error) {
	// result, err := s.repo.DeleteNotify(id)
	// if err != nil {
	// 	return result, err
	// }
	result, err := s.repo.DeleteNotify(id)

	result.Status = 1
	result.ReadAt = time.Now()
	result.UpdatedAt = time.Now()

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: result.UserTo.Hex(), Content: result, ID: "room1", Service: "notify"})

	if createArchive {
		_, err = s.Services.ArchiveNotify.CreateArchiveNotify(userID, result)
	}

	return result, err
}

func (s *NotifyService) ClearNotify(userID string) error {
	return s.repo.ClearNotify(userID)
}
