package service

import (
	"fmt"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type AppErrorService struct {
	repo     repository.AppError
	Hub      *Hub
	Services *Services
}

func NewAppErrorService(repo repository.AppError, hub *Hub) *AppErrorService {
	return &AppErrorService{repo: repo, Hub: hub}
}

func (s *AppErrorService) FindAppError(input *domain.AppErrorFilter) (domain.Response[domain.AppError], error) {
	return s.repo.FindAppError(input)
}

func (s *AppErrorService) CreateAppError(userID string, data *domain.AppError) (*domain.AppError, error) {
	var result *domain.AppError

	result, err := s.repo.CreateAppError(userID, data)
	if err != nil {
		return nil, err
	}

	// находим пользователей(суперадмин) для создания уведомлений.
	roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"systemrole"}})
	if err != nil {
		return nil, err
	}
	ids := []string{}
	var users []domain.User

	if len(roles.Data) > 0 {
		for i := range roles.Data {
			ids = append(ids, roles.Data[i].ID.Hex())
		}

		_users, err := s.Services.User.FindUser(&domain.UserFilter{RoleId: ids})
		if err != nil {
			return nil, err
		}

		users = _users.Data
	}

	// получаем инициатора запроса.
	var authorRequest domain.User
	_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
	if err != nil {
		return nil, err
	}
	if len(_users.Data) > 0 {
		authorRequest = _users.Data[0]
	}

	// отправляем уведомления суперадминам.
	for i := range users {
		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: "userID", Recipient: users[i].ID.Hex(), Content: result, ID: "room1", Service: "appError"})

		_, _ = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo:  users[i].ID.Hex(),
			Title:   domain.AddAppErrorTitle,
			Message: fmt.Sprintf(domain.AddAppError, authorRequest.Name),
		})
	}

	return result, err
}

func (s *AppErrorService) UpdateAppError(id string, userID string, data *domain.AppErrorInput) (*domain.AppError, error) {
	var result *domain.AppError

	result, err := s.repo.UpdateAppError(id, userID, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *AppErrorService) DeleteAppError(id string, userID string) (*domain.AppError, error) {
	var result *domain.AppError

	result, err := s.repo.DeleteAppError(id, userID)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *AppErrorService) ClearAppError(userID string) error {
	return s.repo.ClearAppError(userID)
}
