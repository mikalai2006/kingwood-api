package service

import (
	"fmt"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type PayService struct {
	repo     repository.Pay
	Hub      *Hub
	Services *Services
}

func NewPayService(repo repository.Pay, hub *Hub) *PayService {
	return &PayService{repo: repo, Hub: hub}
}

func (s *PayService) FindPay(input *domain.PayFilter) (domain.Response[domain.Pay], error) {
	return s.repo.FindPay(input)
}

func (s *PayService) CreatePay(userID string, data *domain.Pay) (*domain.Pay, error) {
	var result *domain.Pay

	result, err := s.repo.CreatePay(userID, data)
	if err != nil {
		return nil, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: result.WorkerId.Hex(), Content: result, ID: "room1", Service: "pay"})

	// получаем пользователя, который вносит изменения.
	var authorUpdate domain.User
	_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
	if err != nil {
		return nil, err
	}
	if len(_users.Data) > 0 {
		authorUpdate = _users.Data[0]
	}

	// получаем пользователя, для которого изменили рабочую сессию.
	var worker domain.User
	_workers, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{result.Worker.ID.Hex()}})
	if err != nil {
		return nil, err
	}
	if len(_workers.Data) > 0 {
		worker = _workers.Data[0]
	}

	// add notify.
	_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
		UserTo: result.WorkerId.Hex(),
		Title: fmt.Sprintf(
			domain.CreatePayTitle,
			fmt.Sprintf("%d-%d", result.Year, result.Month+1),
		),
		Message: fmt.Sprintf(
			domain.CreatePay,
			authorUpdate.Name,
			result.Name,
			*result.Total,
			fmt.Sprintf("%d-%d", result.Year, result.Month+1),
		),
	})

	// находим пользователей для создания уведомлений.
	roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin", "boss", "systemrole"}})
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

	for i := range users {
		// add notify.
		_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo: users[i].ID.Hex(),
			Title:  fmt.Sprintf(domain.CreatePayTitle, fmt.Sprintf("%d-%d", result.Year, result.Month+1)),
			Message: fmt.Sprintf(
				domain.CreatePayAdmin,
				authorUpdate.Name,
				worker.Name,
				result.Name,
				*result.Total,
				fmt.Sprintf("%d-%d", result.Year, result.Month+1),
			),
		})
	}

	return result, err
}

func (s *PayService) UpdatePay(id string, userID string, data *domain.PayInput) (*domain.Pay, error) {
	var result *domain.Pay

	prevResults, err := s.repo.FindPay(&domain.PayFilter{ID: []string{id}})
	if err != nil {
		return nil, err
	}

	if len(prevResults.Data) > 0 {
		prevResult := prevResults.Data[0]

		newProps := map[string]interface{}{}
		if prevResult.Props != nil {
			newProps = prevResult.Props
		}
		newItem := make(map[string]interface{})
		newItem["userId"] = userID
		newItem["item"] = domain.PayInput{
			UserID:    prevResult.UserID,
			WorkerId:  prevResult.WorkerId,
			Month:     &prevResult.Month,
			Year:      &prevResult.Year,
			Name:      prevResult.Name,
			Total:     prevResult.Total,
			CreatedAt: prevResult.CreatedAt,
			UpdatedAt: prevResult.UpdatedAt,
		}
		newItem["time"] = time.Now()
		newProps[time.Now().String()] = newItem

		data.Props = newProps

		result, err = s.repo.UpdatePay(id, userID, data)
		if err != nil {
			return nil, err
		}

		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: result.WorkerId.Hex(), Content: result, ID: "room1", Service: "pay"})

		// получаем пользователя, который вносит изменения.
		var authorUpdate domain.User
		_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
		if err != nil {
			return nil, err
		}
		if len(_users.Data) > 0 {
			authorUpdate = _users.Data[0]
		}

		// получаем пользователя, для которого изменили рабочую сессию.
		var worker domain.User
		_workers, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{result.Worker.ID.Hex()}})
		if err != nil {
			return nil, err
		}
		if len(_workers.Data) > 0 {
			worker = _workers.Data[0]
		}

		// добавляем уведомление пользователю, которому меняем счет.
		_, _ = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo: result.WorkerId.Hex(),
			Title:  fmt.Sprintf(domain.PatchPayTitle, fmt.Sprintf("%d-%d", result.Year, result.Month+1)),
			Message: fmt.Sprintf(
				domain.PatchPay,
				authorUpdate.Name,
				prevResult.Name,
				*prevResult.Total,
				result.Name,
				*result.Total,
				fmt.Sprintf("%d-%d", result.Year, result.Month+1)),
		})

		// находим пользователей для создания уведомлений.
		roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin", "boss", "systemrole"}})
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

		for i := range users {
			// добавляем уведомление пользователю, которому меняем счет.
			_, _ = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
				UserTo: users[i].ID.Hex(),
				Title:  fmt.Sprintf(domain.PatchPayTitle, fmt.Sprintf("%d-%d", result.Year, result.Month+1)),
				Message: fmt.Sprintf(
					domain.PatchPayAdmin,
					authorUpdate.Name,
					worker.Name,
					prevResult.Name,
					*prevResult.Total,
					result.Name,
					*result.Total,
					fmt.Sprintf("%d-%d", result.Year, result.Month+1)),
			})
		}
	}

	return result, err
}

func (s *PayService) DeletePay(id string, userID string) (*domain.Pay, error) {
	var result *domain.Pay
	prevResults, err := s.repo.FindPay(&domain.PayFilter{ID: []string{id}})
	if err != nil {
		return nil, err
	}
	if len(prevResults.Data) > 0 {
		// prevResult := prevResults.Data[0]

		result, err = s.repo.DeletePay(id, userID)
		if err != nil {
			return nil, err
		}
		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "pay"})
	}

	return result, err
}
