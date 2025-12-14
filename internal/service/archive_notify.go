package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArchiveNotifyService struct {
	repo     repository.ArchiveNotify
	Hub      *Hub
	Services *Services
}

func NewArchiveNotifyService(repo repository.ArchiveNotify, hub *Hub) *ArchiveNotifyService {
	return &ArchiveNotifyService{repo: repo, Hub: hub}
}

func (s *ArchiveNotifyService) FindArchiveNotifyPopulate(input *domain.ArchiveNotifyFilter) (domain.Response[domain.ArchiveNotify], error) {

	// return s.repo.FindArchiveNotifyPopulate(input)

	var result domain.Response[domain.ArchiveNotify]

	result, err := s.repo.FindArchiveNotifyPopulate(input)
	if err != nil {
		return result, err
	}
	// достаем пользователей.
	users := []string{}
	for i, _ := range result.Data {
		users = append(users, result.Data[i].UserID.Hex())
		users = append(users, result.Data[i].UserTo.Hex())
	}
	resultUsers, err := s.Services.User.FindUser(&domain.UserFilter{
		ID: users,
	})
	// result.Users = resultUsers.Data

	usersMap := map[primitive.ObjectID]domain.User{}
	for i, _ := range resultUsers.Data {
		usersMap[resultUsers.Data[i].ID] = resultUsers.Data[i]
	}

	for i, _ := range result.Data {
		result.Data[i].User = usersMap[result.Data[i].UserID]
		result.Data[i].Recepient = usersMap[result.Data[i].UserTo]
	}

	return result, err
}

func (s *ArchiveNotifyService) CreateArchiveNotify(userID string, data *domain.Notify) (*domain.ArchiveNotify, error) {
	var result *domain.ArchiveNotify

	result, err := s.repo.CreateArchiveNotify(userID, data)
	if err != nil {
		return nil, err
	}

	// send by socket.
	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: result.UserTo.Hex(), Content: result, ID: "room1", Service: "archiveNotify"})

	// // send by push notification.
	// userData, err := s.Services.User.GetUser(result.UserTo.Hex())
	// if err != nil {
	// 	return nil, err
	// }

	// // fmt.Println("userData.AuthPrivate.PushToken=", userData.AuthPrivate.PushToken)
	// if userData.AuthPrivate.PushToken != "" {
	// 	// To check the token is valid
	// 	pushToken, err := expo.NewExponentPushToken(userData.AuthPrivate.PushToken)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Create a new Expo SDK client
	// 	client := expo.NewPushClient(nil)

	// 	// Publish message
	// 	_, err = client.Publish(
	// 		&expo.PushMessage{
	// 			To:       []expo.ExponentPushToken{pushToken},
	// 			Body:     result.Message,
	// 			Data:     map[string]string{"withSome": "data"},
	// 			Sound:    "default",
	// 			Title:    result.Title,
	// 			Priority: expo.DefaultPriority,
	// 		},
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println("Sent push successfully!")
	// } else {
	// 	fmt.Println("Sent push wrong!", userData.AuthPrivate)
	// }

	return result, err
}

func (s *ArchiveNotifyService) DeleteArchiveNotify(id string, userID string) (*domain.ArchiveNotify, error) {
	result, err := s.repo.DeleteArchiveNotify(id, userID)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *ArchiveNotifyService) DeleteArchiveNotifyList(data domain.NotifyListQuery) (*[]domain.ArchiveNotify, error) {
	return s.repo.DeleteArchiveNotifyList(data)
}
