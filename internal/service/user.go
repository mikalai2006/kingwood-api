package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type UserService struct {
	repo repository.User
	Hub  *Hub
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

func (s *UserService) DeleteUser(id string) (domain.User, error) {
	return s.repo.DeleteUser(id)
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
