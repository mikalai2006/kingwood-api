package service

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkHistoryService struct {
	repo     repository.WorkHistory
	Hub      *Hub
	Services *Services
}

func NewWorkHistoryService(repo repository.WorkHistory, hub *Hub) *WorkHistoryService {
	return &WorkHistoryService{repo: repo, Hub: hub}
}

func (s *WorkHistoryService) FindWorkHistory(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error) {
	return s.repo.FindWorkHistory(input)
}

func (s *WorkHistoryService) FindWorkHistoryPopulate(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error) {
	return s.repo.FindWorkHistoryPopulate(input)
}

func (s *WorkHistoryService) CreateWorkHistory(userID string, data *domain.WorkHistory) (*domain.WorkHistory, error) {
	var result *domain.WorkHistory

	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// если есть время от и до, оклад, рассчитываем total
	if !data.From.IsZero() && !data.To.IsZero() && data.Oklad != nil {
		if data.To.Year() != 1 && data.From.Year() != 1 {
			total := int64(0)
			totalMs := int64(0)
			if !data.From.IsZero() && !data.To.IsZero() {
				totalMinutes := data.To.Sub(data.From).Minutes()
				totalMs = data.To.Sub(data.From).Milliseconds()
				total = int64(math.Round(totalMinutes * (float64(*data.Oklad) / 60)))
			}

			if total > 0 {
				data.Total = &total
				data.TotalTime = &totalMs
			}
		}
	}
	// existReview, err := s.repo.FindReview(domain.RequestParams{
	// 	Filter:  bson.M{"node_id": review.NodeID, "userId": userIDPrimitive},
	// 	Options: domain.Options{Limit: 1},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// if len(existReview.Data) > 0 {
	// 	updateReview := &domain.WorkInput{
	// 		Rate:   review.Rate,
	// 		Review: review.Review,
	// 	}
	// 	result, err = s.UpdateReview(existReview.Data[0].ID.Hex(), userID, updateReview)
	// } else {
	// 	result, err = s.repo.CreateReview(userID, review)

	// 	// set user stat
	// 	if err == nil {
	// 		_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// 	}
	// }

	result, err = s.repo.CreateWorkHistory(userID, data)
	if err != nil {
		return nil, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: userID, Content: result, ID: "room1", Service: "workHistory"})

	// находим пользователей(администрацию) для рассылки создания раб.сессии.
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

	// отправляем уведомления администрации.
	for i := range users {

		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: users[i].ID.Hex(), Content: result, ID: "room1", Service: "workHistory"})

	}

	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	return result, err
}

func (s *WorkHistoryService) UpdateWorkHistory(id string, userID string, data *domain.WorkHistoryInput) (*domain.WorkHistory, error) {
	// получаем данные из базы.
	existWorkHistory, err := s.repo.FindWorkHistoryPopulate(domain.WorkHistoryFilter{ID: []string{id}})
	if err != nil {
		return nil, err
	}

	// fmt.Println("data: ", data)
	// fmt.Println("id: ", id)
	// fmt.Println("WorkerId: ", data.WorkerId)
	// fmt.Println("From: ", existWorkHistory.Data[0].From, data.From)
	// fmt.Println("To: ", existWorkHistory.Data[0].To, data.To)

	// статус изменения времени.
	isWorkHistoryChange := false
	if len(existWorkHistory.Data) > 0 {
		// если данные для патча отличаются от данных из базы
		if (existWorkHistory.Data[0].From != data.From || existWorkHistory.Data[0].To != data.To) && existWorkHistory.Data[0].To.Year() != 1 {
			isWorkHistoryChange = true

			// заносим старые данные в пропс.
			newProps := map[string]interface{}{}
			if existWorkHistory.Data[0].Props != nil {
				newProps = existWorkHistory.Data[0].Props
			}
			newItem := make(map[string]interface{})
			newItem["userId"] = userID
			newItem["item"] = domain.WorkHistoryInput{
				UserID:       existWorkHistory.Data[0].UserID,
				WorkerId:     existWorkHistory.Data[0].WorkerId,
				ObjectId:     &existWorkHistory.Data[0].ObjectId,
				OrderId:      &existWorkHistory.Data[0].OrderId,
				TaskId:       &existWorkHistory.Data[0].TaskId,
				OperationId:  &existWorkHistory.Data[0].OperationId,
				TaskWorkerId: &existWorkHistory.Data[0].TaskWorkerId,
				TotalTime:    existWorkHistory.Data[0].TotalTime,
				To:           existWorkHistory.Data[0].To,
				From:         existWorkHistory.Data[0].From,
				Oklad:        existWorkHistory.Data[0].Oklad,
				Date:         existWorkHistory.Data[0].Date,
				Total:        existWorkHistory.Data[0].Total,
				CreatedAt:    existWorkHistory.Data[0].CreatedAt,
				UpdatedAt:    existWorkHistory.Data[0].UpdatedAt,
			}
			newItem["time"] = time.Now().UTC()
			newProps[time.Now().String()] = newItem

			// дополняем пропс.
			data.Props = newProps
		}
	}
	// fmt.Println("isWorkHistoryChange: ", isWorkHistoryChange)

	result, err := s.repo.UpdateWorkHistory(id, userID, data)
	if err != nil {
		return result, err
	}

	if result != nil {
		// update total.
		newRobotUpdateData := &domain.WorkHistoryInput{}
		total := int64(0)
		totalMs := int64(0)
		if !result.From.IsZero() && !result.To.IsZero() {
			totalMinutes := result.To.Sub(result.From).Minutes()
			totalMs = result.To.Sub(result.From).Milliseconds()
			total = int64(math.Round(totalMinutes * (float64(*result.Oklad) / 60)))
		}

		// if total > 0 {
		newRobotUpdateData.Total = &total
		newRobotUpdateData.TotalTime = &totalMs
		// }

		// explodeDate := false
		// oldTo := result.To
		// var fromNew time.Time
		var toNew time.Time

		eastOfUTC := time.FixedZone("UTC-3", -3*60*60)
		to1 := time.Date(result.To.Year(), result.To.Month(), result.To.Day(), result.To.Hour(), result.To.Minute(), result.To.Second(), 0, eastOfUTC)
		from1 := time.Date(result.From.Year(), result.From.Month(), result.From.Day(), result.From.Hour(), result.From.Minute(), result.From.Second(), 0, eastOfUTC)

		// fmt.Println("======================PATCH WORK HISTORY====================")
		// fmt.Println("from: ", from1, "====>", from1.UTC())
		// fmt.Println("to: ", to1, "====>", to1.UTC())
		// fmt.Println("========================================================")

		// fmt.Println("result.From: ", to1, to1.UTC(), from1, from1.UTC())
		if from1.UTC().Day() != to1.UTC().Day() {
			// explodeDate = true

			// // prevDay := oldTo.AddDate(0, 0, -1)
			// year, month, _ := oldTo.Date()
			// // fromNew :=  result.From
			// // time.Date(year, month, day, 0, 0, 0, 0, prevDay.Location())
			// fromNew = time.Date(year, month, dayPrev, 21, 0, 0, 0, time.UTC)
			yearPrev, monthPrev, dayPrev := from1.Date()
			toNew = time.Date(yearPrev, monthPrev, dayPrev, 20, 59, 59, 0, time.UTC)

			newRobotUpdateData.To = toNew

			totalMinutesPrev := toNew.Sub(result.From).Minutes()
			totalPrev := int64(math.Round(totalMinutesPrev * (float64(*result.Oklad) / 60)))
			// update total.
			newRobotUpdateData.Total = &totalPrev
		}

		result, err = s.repo.UpdateWorkHistory(id, userID, newRobotUpdateData)
		if err != nil {
			return result, err
		}

		// создаем новую запись для оставшейся части времени
		// if explodeDate {
		// 	// Переносим часть рабочего времени на другой день
		// 	totalMinutesNext := oldTo.Sub(fromNew).Minutes()
		// 	totalNext := int64(math.Round(totalMinutesNext * (float64(*result.Oklad) / 60)))
		// 	// fmt.Println("totalMinutesNext:", totalMinutesNext, " totalNext:", totalNext, " oldTo:", oldTo)
		// 	result, err = s.repo.CreateWorkHistory(userID, &domain.WorkHistory{
		// 		UserID:      result.UserID,
		// 		WorkerId:    result.WorkerId,
		// 		ObjectId:    result.ObjectId,
		// 		OrderId:     result.OrderId,
		// 		TaskId:      result.TaskId,
		// 		OperationId: result.OperationId,
		// 		Status:      result.Status,
		// 		Date:        fromNew,
		// 		From:        fromNew,
		// 		To:          oldTo,
		// 		Oklad:       result.Oklad,
		// 		Total:       &totalNext,
		// 	})
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// }

		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: result.WorkerId.Hex(), Content: result, ID: "room1", Service: "workHistory"})

		// находим пользователей(администрацию) для рассылки создания раб.сессии.
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

		// отправляем уведомления администрации.
		for i := range users {

			s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: users[i].ID.Hex(), Content: result, ID: "room1", Service: "workHistory"})

		}
	}

	if isWorkHistoryChange {
		// получаем пользователя, который вносит изменения.
		var authorUpdate domain.User
		_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
		if err != nil {
			return nil, err
		}
		if len(_users.Data) > 0 {
			authorUpdate = _users.Data[0]
		}

		// находим пользователей для создания уведомлений.
		roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin", "boss"}})
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

		// получаем пользователя, для которого изменили рабочую сессию.
		var worker domain.User
		_workers, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{existWorkHistory.Data[0].WorkerId.Hex()}})
		if err != nil {
			return nil, err
		}
		if len(_workers.Data) > 0 {
			worker = _workers.Data[0]
		}

		// протяженность рабочей сессии
		durationNew := result.To.Sub(result.From)
		_durationNewText, _ := time.ParseDuration(durationNew.String())
		durationNewText := strings.Replace(_durationNewText.String(), "h", "ч.", 1)
		durationNewText = strings.Replace(durationNewText, "m", "мин.", 1)
		durationNewText = strings.Replace(durationNewText, "s", "сек.", 1)
		durationOld := existWorkHistory.Data[0].To.Sub(existWorkHistory.Data[0].From)
		_durationOldText, _ := time.ParseDuration(durationOld.String())
		durationOldText := strings.Replace(_durationOldText.String(), "h", "ч.", 1)
		durationOldText = strings.Replace(durationOldText, "m", "мин.", 1)
		durationOldText = strings.Replace(durationOldText, "s", "сек.", 1)
		// durationOldText := fmt.Sprintf("%d:%d:%d", int64(_durationOldText.Hours()), int64(_durationOldText.Minutes()), int64(_durationOldText.Seconds()))

		// информация о заказе.
		oldOrderInfo := fmt.Sprintf("№%d-%s", existWorkHistory.Data[0].Order.Number, existWorkHistory.Data[0].Order.Name)
		if existWorkHistory.Data[0].Order.Number == 0 {
			oldOrderInfo = "Хоз.работы"
		}
		newOrderInfo := fmt.Sprintf("№%d-%s", result.Order.Number, result.Order.Name)
		if result.Order.Number == 0 {
			newOrderInfo = "Хоз.работы"
		}
		// смещение времени.
		westOfUTC := time.FixedZone("UTC+3", 3*60*60)

		for i := range users {
			// отправка уведомления администраторам и нач. цеху
			s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
				UserTo: users[i].ID.Hex(),
				Title:  domain.PatchWorkHistoryTitle,
				Message: fmt.Sprintf(
					domain.PatchWorkHistoryAdmin,
					authorUpdate.Name,
					worker.Name,
					existWorkHistory.Data[0].Date.In(westOfUTC).Format("02.01.2006"),
					existWorkHistory.Data[0].From.In(westOfUTC).Format("15:04:05"),
					existWorkHistory.Data[0].To.In(westOfUTC).Format("15:04:05"),
					durationOldText,
					oldOrderInfo,
					result.From.In(westOfUTC).Format("15:04:05"),
					result.To.In(westOfUTC).Format("15:04:05"),
					durationNewText,
					newOrderInfo,
				),
				// Link:       "/(tabs)/finance",
				// LinkOption: map[string]interface{}{},
			})
		}
		// отправка уведомления сотруднику, для кого меняются данные
		s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo: existWorkHistory.Data[0].WorkerId.Hex(),
			Title:  domain.PatchWorkHistoryTitle,
			Message: fmt.Sprintf(
				domain.PatchWorkHistory,
				authorUpdate.Name,
				existWorkHistory.Data[0].Date.In(westOfUTC).Format("02.01.2006"),
				existWorkHistory.Data[0].From.In(westOfUTC).Format("15:04:05"),
				existWorkHistory.Data[0].To.In(westOfUTC).Format("15:04:05"),
				durationOldText,
				oldOrderInfo,
				result.From.In(westOfUTC).Format("15:04:05"),
				result.To.In(westOfUTC).Format("15:04:05"),
				durationNewText,
				newOrderInfo,
			),
			// Link:       "/(tabs)/finance",
			// LinkOption: map[string]interface{}{},
		})

	}

	return result, err
}

func (s *WorkHistoryService) DeleteWorkHistory(id string, userID string, createNotify bool) (*domain.WorkHistory, error) {
	var result *domain.WorkHistory

	// получаем инициатора запроса.
	var authorRequest domain.User
	_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
	if err != nil {
		return nil, err
	}
	if len(_users.Data) > 0 {
		authorRequest = _users.Data[0]
	}

	// Находим рабочую сессию для удаления
	existWorkHistory, err := s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{ID: []string{id}})
	if err != nil {
		return result, err
	}

	if len(existWorkHistory.Data) > 0 {
		result = &existWorkHistory.Data[0]
	}

	// add notify.
	for i := range existWorkHistory.Data {
		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: existWorkHistory.Data[i].WorkerId.Hex(), Content: existWorkHistory.Data[i], ID: "room1", Service: "WorkHistory"})

		if createNotify {
			message := fmt.Sprintf(domain.DeleteWorkHistory, authorRequest.Name, existWorkHistory.Data[i].Date.Format("02.01.2006"), result.Order.Number, result.Order.Name, result.Object.Name)
			if result.Object.Name == "" {
				message = fmt.Sprintf(domain.DeleteWorkHistoryNotOrder, authorRequest.Name, existWorkHistory.Data[i].Date.Format("02.01.2006"))
			}

			_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
				UserTo:  result.WorkerId.Hex(),
				Title:   domain.DeleteWorkHistoryTitle,
				Message: message,
			})
		}

		// находим пользователей(администрацию) для создания уведомлений.
		roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin", "boss"}})
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

		// получаем пользователя, для которого удалили задание.
		var worker domain.User
		_workers, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{result.WorkerId.Hex()}})
		if err != nil {
			return nil, err
		}
		if len(_workers.Data) > 0 {
			worker = _workers.Data[0]
		}

		// отправляем уведомления админам.
		for i := range users {
			s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: users[i].ID.Hex(), Content: result, ID: "room1", Service: "WorkHistory"})

			if createNotify {
				message := fmt.Sprintf(domain.DeleteWorkHistoryAdmin, authorRequest.Name, worker.Name, result.Date.Format("02.01.2006"), result.Order.Number, result.Order.Name, result.Object.Name)
				if result.Object.Name == "" {
					message = fmt.Sprintf(domain.DeleteWorkHistoryAdminNotOrder, authorRequest.Name, worker.Name, result.Date.Format("02.01.2006"))
				}
				_, _ = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
					UserTo:  users[i].ID.Hex(),
					Title:   domain.DeleteWorkHistoryTitle,
					Message: message,
				})
			}
		}
	}

	result, err = s.repo.DeleteWorkHistory(id)

	_, err = s.Services.CreateArchiveWorkHistory(userID, result)

	return result, err
}

func (s *WorkHistoryService) GetStatByOrder(input domain.WorkHistoryFilter) ([]domain.WorkHistoryStatByOrder, error) {
	result, err := s.repo.GetStatByOrder(input)

	return result, err
}

func (s *WorkHistoryService) GetStatByMonth(input domain.WorkHistoryFilter) ([]domain.WorkHistoryStatByMonth, error) {
	result, err := s.repo.GetStatByMonth(input)

	return result, err
}

func (s *WorkHistoryService) ClearWorkHistory(userID string) error {
	return s.repo.ClearWorkHistory(userID)
}
