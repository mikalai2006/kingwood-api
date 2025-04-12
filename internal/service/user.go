package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type UserService struct {
	repo     repository.User
	Hub      *Hub
	Services *Services
}

func NewUserService(repo repository.User, hub *Hub) *UserService {
	return &UserService{repo: repo, Hub: hub}
}

func (s *UserService) GetUser(id string) (domain.User, error) {
	return s.repo.GetUser(id)
}

func (s *UserService) FindUser(filter *domain.UserFilter) (domain.Response[domain.User], error) {
	return s.repo.FindUser(filter)
}

func (s *UserService) CreateUser(userID string, user *domain.User) (*domain.User, error) {
	return s.repo.CreateUser(userID, user)
}

func (s *UserService) DeleteUser(id string, userID string) (domain.User, error) {
	var result domain.User

	// delete images.
	allImages, err := s.Services.Image.FindImage(&domain.ImageFilter{ServiceId: []string{id}})
	//domain.RequestParams{Filter: bson.D{{"serviceId", id}}}
	if err != nil {
		return result, err
	}
	for i := range allImages.Data {
		// fmt.Println("Remove image: ", allImages.Data[i].ID)
		_, err = s.Services.Image.DeleteImage(userID, allImages.Data[i].ID.Hex(), true)
	}

	// delete taskWorkers.
	allTaskWorkers, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{WorkerId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allTaskWorkers.Data {
		// fmt.Println("Remove TaskWorkers: ", allTaskWorkers.Data[i].ID)
		_, err = s.Services.TaskWorker.DeleteTaskWorker(allTaskWorkers.Data[i].ID.Hex(), userID, false)
	}

	// delete workHistory.
	allWorkHistory, err := s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{WorkerId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allWorkHistory.Data {
		// fmt.Println("Remove WorkHistory: ", allWorkHistory.Data[i].ID)
		_, err = s.Services.WorkHistory.DeleteWorkHistory(allWorkHistory.Data[i].ID.Hex(), userID, false)
	}

	// delete pay.
	allPay, err := s.Services.Pay.FindPay(&domain.PayFilter{WorkerId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allPay.Data {
		// fmt.Println("Remove Pay: ", allPay.Data[i].ID)
		_, err = s.Services.Pay.DeletePay(allPay.Data[i].ID.Hex(), userID)
	}

	// delete notify.
	allNotify, err := s.Services.Notify.FindNotifyPopulate(&domain.NotifyFilter{UserTo: []*string{&id}})
	if err != nil {
		return result, err
	}
	for i := range allNotify.Data {
		// fmt.Println("Remove Notify: ", allNotify.Data[i].ID)
		_, err = s.Services.Notify.DeleteNotify(allNotify.Data[i].ID.Hex(), userID, false)
	}

	result, err = s.repo.DeleteUser(id)

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: "", Content: result, ID: "room1", Service: "user"})

	_, err = s.Services.CreateArchiveUser(userID, &result)

	return result, err
}

func (s *UserService) UpdateUser(id string, user *domain.UserInput) (domain.User, error) {
	result, err := s.repo.UpdateUser(id, user)
	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: id, Recipient: "", Content: result, ID: "room1", Service: "user"})

	return result, err
}

func (s *UserService) Iam(userID string) (domain.User, error) {
	user, err := s.repo.Iam(userID)
	if err != nil {
		return user, err
	}

	// user, err = s.UpdateUser(userID, &domain.User{Online: true})
	// s.Hub.HandleMessage(domain.Message{Type: "message", Sender: "user1", Recipient: "user2", Content: user, ID: "room1", Service: "user"})

	return user, err
}
