package service

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
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

	simpleCreateItem := false
	// если есть время от и до, оклад, рассчитываем total
	if !data.From.IsZero() && !data.To.IsZero() && data.Oklad != nil {
		if data.To.Year() != 1 && data.From.Year() != 1 {
			simpleCreateItem = true
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

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: result.WorkerId.Hex(), Content: result, ID: "room1", Service: "workHistory"})

	// находим пользователей(администрацию) для рассылки создания раб.сессии.
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

	// отправляем уведомления администрации.
	for i := range users {

		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: users[i].ID.Hex(), Content: result, ID: "room1", Service: "workHistory"})

	}

	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	// если у работника задано максимальное время работы,
	// нужно запустить таймер
	if result.Worker.MaxTime != nil && *result.Worker.MaxTime != 0 {
		// если это не просто создание сессии (вставка модератором или при переходе сессии между днями)
		if !simpleCreateItem {
			// узнаем сколько он отработал, чтобы понять через какое время запустить задачу по таймеру
			listWorkHistoryByDay, err := s.repo.FindWorkHistoryPopulate(domain.WorkHistoryFilter{
				WorkerId: []string{result.WorkerId.Hex()},
				Date:     result.From,
			})
			if err != nil {
				return nil, err
			}

			// подсчитываем общее время всех сессий за дату текущей рабочей сессии.
			var allWorkTime int64
			for i := range listWorkHistoryByDay.Data {
				allWorkTime = allWorkTime + *listWorkHistoryByDay.Data[i].TotalTime
			}

			// вычисляем сколько осталось поработать до максимально возможного времени
			maxTimeDuration := (time.Duration(*result.Worker.MaxTime) * time.Hour).Milliseconds()
			durationForTimerMs := maxTimeDuration - allWorkTime
			durationForTimer := (time.Duration(durationForTimerMs) * time.Millisecond)

			// проверяем, чтобы время таймера не вышло за пределы суток (если время до полуночи меньше положенной смены, таймер ставим на полночь).
			// eastOfUTC := time.FixedZone("UTC-3", -3*60*60)
			// from1 := time.Date(result.From.Year(), result.From.Month(), result.From.Day(), result.From.Hour(), result.From.Minute(), result.From.Second(), 0, eastOfUTC)
			yearPrev, monthPrev, dayPrev := result.From.Date()
			timePolnoc := time.Date(yearPrev, monthPrev, dayPrev, 20, 59, 59, 0, time.UTC)

			raznicaFromToNew := timePolnoc.Sub(result.From)
			if durationForTimer > raznicaFromToNew {
				durationForTimer = raznicaFromToNew
			}

			// создаем таймер с задачей
			_, err = s.Services.Timer.CreateTimer(userID, &domain.TimerShedule{
				IDTimer:       fmt.Sprintf("timer_%v", durationForTimerMs),
				ExecuteAt:     time.Now().Add(durationForTimer), //.Add(5 * time.Duration(time.Second)),
				IsRunning:     1,
				WorkerId:      result.WorkerId,
				TaskWorkerId:  result.TaskWorkerId,
				TaskId:        result.TaskId,
				WorkHistoryId: result.ID,
			})
			if err != nil {
				return nil, err
			}
		} else {
			// если создание произведено модератором или автоматически для переноса части сессии на другой день
			// и у работника стоит максимальное время, проверяем максимальное количество часов для дня текущей сессии
			// если в аккаунте задано макс. время работы, делаем округление времени,
			// нужно достать все его рабочие сессии за день и
			// в изменяемой сессии округлить время, если общее время за все сессии больше положенного
			// положенное время устанавливается в аккаунте работника.
			if result.Worker.MaxTime != nil {
				result, err = s.ValidWorkHistory(userID, result.ID.Hex(), result)
				if err != nil {
					return nil, err
				}
			}
		}

	}

	return result, err
}

func (s *WorkHistoryService) UpdateWorkHistory(id string, userID string, data *domain.WorkHistoryInput) (*domain.WorkHistory, error) {
	// получаем данные из базы.
	existWorkHistory, err := s.repo.FindWorkHistoryPopulate(domain.WorkHistoryFilter{ID: []string{id}})
	if err != nil {
		return nil, err
	}

	if len(existWorkHistory.Data) == 0 {
		return nil, errors.New("not found work session")
	}

	// округляем даты "от" и "до" до секунд.
	if !data.From.IsZero() {
		data.From = data.From.Truncate(time.Millisecond)
	}
	if !data.To.IsZero() {
		data.To = data.To.Truncate(time.Millisecond)
	}

	// fmt.Println("data: ", data)
	// fmt.Println("id: ", id)
	// fmt.Println("WorkerId: ", data.WorkerId)
	// fmt.Println("From: ", existWorkHistory.Data[0].From, data.From)
	// fmt.Println("To: ", existWorkHistory.Data[0].To, data.To)

	// статус изменения времени.
	isWorkHistoryChange := false

	if len(existWorkHistory.Data) > 0 {
		// Блок работает при изменении данных сессии, после закрытия сессии.
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

	// вносим изменения в данные сессии.
	result, err := s.repo.UpdateWorkHistory(id, userID, data)
	if err != nil {
		return result, err
	}

	// пересчитываем общее время и другие данные.
	if result != nil {
		// update total.
		newRobotUpdateData := &domain.WorkHistoryInput{}
		total := int64(0)
		totalMs := int64(0)

		// если время начала и завершения не нулевые,
		// пересчитываем общее время
		if !result.From.IsZero() && !result.To.IsZero() {
			totalMinutes := result.To.Sub(result.From).Minutes()
			totalMs = result.To.Sub(result.From).Milliseconds()
			total = int64(math.Round(totalMinutes * (float64(*result.Oklad) / 60)))
		}

		// if total > 0 {
		newRobotUpdateData.Total = &total
		newRobotUpdateData.TotalTime = &totalMs
		// }

		explodeDate := false
		// oldTo := result.To
		// var fromNew time.Time

		// начало - функционал обрезки рабочего времени до полуночи.
		var toNew time.Time

		eastOfUTC := time.FixedZone("UTC-3", -3*60*60)
		to1 := time.Date(result.To.Year(), result.To.Month(), result.To.Day(), result.To.Hour(), result.To.Minute(), result.To.Second(), 0, eastOfUTC)
		from1 := time.Date(result.From.Year(), result.From.Month(), result.From.Day(), result.From.Hour(), result.From.Minute(), result.From.Second(), 0, eastOfUTC)

		yearPrev, monthPrev, dayPrev := from1.Date()
		// fmt.Println("======================PATCH WORK HISTORY====================")
		// fmt.Println("from: ", from1, "====>", from1.UTC())
		// fmt.Println("to: ", to1, "====>", to1.UTC())
		// fmt.Println("========================================================")
		// fmt.Println("result.From: ", to1, to1.UTC(), from1, from1.UTC())
		if from1.UTC().Day() != to1.UTC().Day() {
			explodeDate = true
			// // prevDay := oldTo.AddDate(0, 0, -1)
			// // fromNew :=  result.From
			// // time.Date(year, month, day, 0, 0, 0, 0, prevDay.Location())
			toNew = time.Date(yearPrev, monthPrev, dayPrev, 20, 59, 59, 0, time.UTC)

			newRobotUpdateData.To = toNew

			totalMinutesPrev := toNew.Sub(result.From).Minutes()
			totalPrev := int64(math.Round(totalMinutesPrev * (float64(*result.Oklad) / 60)))
			// update total.
			newRobotUpdateData.Total = &totalPrev

			// если задано макс. время
			maxTimeDuration := (time.Duration(*result.Worker.MaxTime) * time.Hour)
			maxTime := maxTimeDuration.Milliseconds()
			// делаем запрос всех сессий за прошлый день, чтобы узнать сколько уже отработано времени и сколько еще можно добавить до заданного макс. времени.
			if result.Worker.MaxTime != nil && *result.Worker.MaxTime != 0 {
				listWorkHistoryByDay, err := s.repo.FindWorkHistoryPopulate(domain.WorkHistoryFilter{
					WorkerId: []string{result.WorkerId.Hex()},
					Date:     result.From,
				})
				if err != nil {
					return nil, err
				}
				// подсчитываем общее время всех сессий за дату текущей рабочей сессии.
				var allWorkTime int64
				for i := range listWorkHistoryByDay.Data {
					if listWorkHistoryByDay.Data[i].ID.Hex() != id {
						allWorkTime = allWorkTime + *listWorkHistoryByDay.Data[i].TotalTime
					}
				}
				// fmt.Println("listWorkHistoryByDay=", len(listWorkHistoryByDay.Data), ", allWorkTime=", allWorkTime)

				// if (allWorkTime + *result.TotalTime) > maxTime {
				// устанавливаем общее время текущей сессии, как разницу от положенного времени и других сессий (без текущей).
				cuteTotalTime := time.Duration(maxTime-allWorkTime) * time.Millisecond

				// проверяем, чтобы время только до полуночи считалось, даже и в случае если не выйдет полная положенная смена.
				raznicaFromToNew := toNew.Sub(result.From)
				if cuteTotalTime > raznicaFromToNew {
					cuteTotalTime = raznicaFromToNew
				}

				cuteTotalTimeMinutes := cuteTotalTime.Minutes()
				totalPrev = int64(math.Round(cuteTotalTimeMinutes * (float64(*result.Oklad) / 60)))
				newRobotUpdateData.Total = &totalPrev

				cuteTotalTimeMs := cuteTotalTime.Milliseconds()
				newRobotUpdateData.TotalTime = &cuteTotalTimeMs

				newRobotUpdateData.To = result.From.Add(cuteTotalTime).Truncate(time.Millisecond)
				// fmt.Println("cuteTotalTime=", cuteTotalTime, ", newRobotUpdateData.To=", newRobotUpdateData.To)
				// }
			}

			// остаток рабочей сессии, который перешел на другой день.
			// ostatokFrom := time.Date(to1.Year(), to1.Month(), to1.Minute(), 21, 00, 00, 0, time.UTC)
			// ostatokTo := to1
			// ostatokTotalMinutes := ostatokTo.Sub(ostatokFrom).Minutes()
			// ostatokTotal := int64(math.Round(ostatokTotalMinutes * (float64(*result.Oklad) / 60)))

			// fmt.Println("ostatok Minutes = ", ostatokTotalMinutes, ", ostatokTotal = ", ostatokTotal, ", ostatokFrom=", ostatokFrom, ", ostatokTo=", ostatokTo)

		}
		// конец - функционала обрезки времени до полуночи.

		result, err = s.repo.UpdateWorkHistory(id, userID, newRobotUpdateData)
		if err != nil {
			return result, err
		}

		// // создаем новую запись для оставшейся части времени
		// if explodeDate {
		// 	year, month, day := oldTo.Date()
		// 	eastOfUTCPlus3 := time.FixedZone("UTC+3", 3*60*60)
		// 	fromNew = time.Date(year, month, day, 0, 0, 0, 0, eastOfUTCPlus3)
		// 	// Переносим часть рабочего времени на другой день
		// 	totalMinutesNext := oldTo.Sub(fromNew).Minutes()
		// 	totalNext := int64(math.Round(totalMinutesNext * (float64(*result.Oklad) / 60)))
		// 	// fmt.Println("totalMinutesNext:", totalMinutesNext, " totalNext:", totalNext, " oldTo:", oldTo)

		// 	newWorkHistory := domain.WorkHistory{
		// 		UserID:       result.UserID,
		// 		WorkerId:     result.WorkerId,
		// 		ObjectId:     result.ObjectId,
		// 		OrderId:      result.OrderId,
		// 		TaskId:       result.TaskId,
		// 		OperationId:  result.OperationId,
		// 		TaskWorkerId: result.TaskWorkerId,
		// 		Status:       result.Status,
		// 		Date:         oldTo,
		// 		From:         fromNew,
		// 		To:           oldTo,
		// 		Oklad:        result.Oklad,
		// 		Total:        &totalNext,
		// 	}
		// 	// result, err = s.repo.CreateWorkHistory(userID, &newWorkHistory)

		// 	_, err := s.CreateWorkHistory(userID, &newWorkHistory)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// }

		// если в аккаунте задано макс. время работы, делаем округление времени,
		// нужно достать все его рабочие сессии за день и
		// в изменяемой сессии округлить время, если общее время за все сессии больше положенного
		// положенное время устанавливается в аккаунте работника.
		//
		// также, если было разделение рабочей сессии - это отменяет валидацию первой половины сессии.
		if result.Worker.MaxTime != nil && !explodeDate && !isWorkHistoryChange {
			result, err = s.ValidWorkHistory(userID, id, result)
			if err != nil {
				return nil, err
			}
		}

		// проверяем доплаты
		if result.Worker.Dops != nil && len(result.Worker.Dops) > 0 {
			_, err = s.CheckDoplats(userID, id, result)
			if err != nil {
				return nil, err
			}
		}

		// в сокеты отправляем информацию для работника, что начата рабочая сессия
		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: result.WorkerId.Hex(), Content: result, ID: "room1", Service: "workHistory"})

		// находим пользователей(администрацию) для рассылки создания раб.сессии.
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

		// отправляем уведомления администрации.
		for i := range users {

			s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: users[i].ID.Hex(), Content: result, ID: "room1", Service: "workHistory"})

		}
	}

	// отправляем уведомления, если было произведено изменение данных рабочей сессии
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

		if authorUpdate.RoleObject.Code != "systemrole" {
			// находим пользователей для создания уведомлений.
			roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"systemrole"}}) // "admin", "boss"
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
				// отправка уведомления администрации
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
	}

	// проверяем есть ли таймер для текущей рабочей сессии.
	isRunining := 1
	timers, err := s.Services.Timer.FindTimerPopulate(domain.TimerSheduleFilter{
		WorkHistoryId: []string{id},
		IsRunning:     &isRunining,
	})
	fmt.Println("length workHistory: ", len(timers.Data))
	// если есть таймеры, проходим по всем и отключаем таймеры, путем записи в базу, что таймер выполнен
	if len(timers.Data) > 0 {
		isRuniningStatus := 0
		for i, _ := range timers.Data {
			_, err = s.Services.Timer.UpdateTimer(timers.Data[i].ID.Hex(), userID, &domain.TimerSheduleInput{
				IsRunning: &isRuniningStatus,
			})
			if err != nil {
				return nil, err
			}
		}
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
		s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: existWorkHistory.Data[i].WorkerId.Hex(), Content: existWorkHistory.Data[i], ID: "room1", Service: "workHistory"})

		if createNotify && authorRequest.RoleObject.Code != "systemrole" {
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
			s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: users[i].ID.Hex(), Content: result, ID: "room1", Service: "workHistory"})

			if createNotify && authorRequest.RoleObject.Code != "systemrole" {
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

func (s *WorkHistoryService) CheckDoplats(userID string, id string, item *domain.WorkHistory) (*domain.WorkHistory, error) {
	result := item
	var err error

	if item.Worker.Dops == nil || len(item.Worker.Dops) == 0 {
		return result, err
	}

	// если рабочая сессия завершается или просто содержит "from" и "to"
	// достаем все рабочие сессии за сутки этой сессии
	// проверяем общее рабочее время и если есть доплаты, проходим по ним и создаем доплаты
	if !item.From.IsZero() && !item.To.IsZero() {
		listWorkHistoryByDay, err := s.repo.FindWorkHistoryPopulate(domain.WorkHistoryFilter{
			WorkerId: []string{item.WorkerId.Hex()},
			Date:     item.Date,
		})
		if err != nil {
			return nil, err
		}

		// подсчитываем общее время всех сессий за дату текущей рабочей сессии.
		var allWorkTime int64
		for i := range listWorkHistoryByDay.Data {
			allWorkTime = allWorkTime + *listWorkHistoryByDay.Data[i].TotalTime
		}

		// получаем день недели.
		weekDay := item.Date.Weekday()

		// находим пользователя суперадмина, чтобы от его имени создать платежи.
		// находим пользователей для создания уведомлений.
		superUser, err := s.Services.User.GetSuperAdmin()
		if err != nil {
			return nil, err
		}

		statusAuto := 1

		// проходи по всем доплатам работника.
		for i, _ := range item.Worker.Dops {
			// если день недели есть среди дней доплаты.
			if utils.Contains(item.Worker.Dops[i].Days, int(weekDay)) {
				// получаем время, которое нужно отработать, чтобы получить доплату.
				needTimeMs := time.Duration(item.Worker.Dops[i].MinHours * int(time.Hour)).Milliseconds()
				// размер доплаты.
				doplata := int64(item.Worker.Dops[i].Doplata)
				// если отработанное время больше или равно времени,
				// которое нужно для начисления доплаты.
				if allWorkTime >= int64(needTimeMs) {
					name := fmt.Sprintf("Доплата за выходной день (%s)", item.Date.Format("02-01-2006"))
					// проверяем есть уже автодоплата с таким именем или нет.
					autoPay, err := s.Services.Pay.FindPay(&domain.PayFilter{
						Auto: &statusAuto,
						Name: name,
					})
					if err != nil {
						return nil, err
					}
					if len(autoPay.Data) == 0 {
						// создаем доплату.
						_, err = s.Services.Pay.CreatePay(userID, &domain.Pay{
							UserID:   superUser.ID,
							WorkerId: item.WorkerId,
							Month:    int64(item.Date.Month()) - 1,
							Year:     int64(item.Date.Year()),
							Total:    &doplata,
							Auto:     &statusAuto,
							Name:     name,
						})
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}

	return result, err
}

func (s *WorkHistoryService) ValidWorkHistory(userID string, id string, item *domain.WorkHistory) (*domain.WorkHistory, error) {
	result := item
	var err error
	var total int64
	maxTimeDuration := (time.Duration(*item.Worker.MaxTime) * time.Hour)
	maxTime := maxTimeDuration.Milliseconds()
	if !item.From.IsZero() && !item.To.IsZero() && maxTime != 0 {
		listWorkHistoryByDay, err := s.repo.FindWorkHistoryPopulate(domain.WorkHistoryFilter{
			WorkerId: []string{item.WorkerId.Hex()},
			Date:     item.Date,
		})
		if err != nil {
			return nil, err
		}

		// подсчитываем общее время всех сессий за дату текущей рабочей сессии.
		var allWorkTime int64
		for i := range listWorkHistoryByDay.Data {
			allWorkTime = allWorkTime + *listWorkHistoryByDay.Data[i].TotalTime
		}

		// если до выполнения смены не хватает заданного времени (15 минут)
		// ставим, что выполнено заданное время
		timeOstatok := time.Duration(time.Minute * 1).Milliseconds()
		minNeedTime := maxTime - timeOstatok
		if allWorkTime >= minNeedTime {
			allWorkTime = allWorkTime + timeOstatok
			*item.TotalTime = *item.TotalTime + timeOstatok
		}
		// fmt.Println("minNeedTime=", minNeedTime, ", allWorkTime=", allWorkTime, ", maxTime=", maxTime, ", maxTime-minNeedTime=", maxTime-minNeedTime)

		// если время больше положенного
		if (allWorkTime + 1000) > maxTime {
			// устанавливаем общее время текущей сессии, как разницу от положенного времени и других сессий (без текущей).
			allWorkTimeWithoutCurrent := allWorkTime - *item.TotalTime
			cuteTotalTime := time.Duration(maxTime-allWorkTimeWithoutCurrent) * time.Millisecond
			cuteTotalTimeMinutes := cuteTotalTime.Minutes()
			total = int64(math.Round(cuteTotalTimeMinutes * (float64(*item.Oklad) / 60)))
			cutTo := item.From.Add(cuteTotalTime).Truncate(time.Millisecond)
			// fmt.Println("maxTime=", maxTime, " cuteTotalTime=", cuteTotalTime, " cuteTotalTimeMinutes=", cuteTotalTimeMinutes, " allWorkTimeWithoutCurrent=", allWorkTimeWithoutCurrent)

			// заносим старые данные в пропс.
			newProps := map[string]interface{}{}
			if item.Props != nil {
				newProps = item.Props
			}
			newItem := make(map[string]interface{})
			newItem["userId"] = userID
			newItem["item"] = domain.WorkHistoryInput{
				UserID:       item.UserID,
				WorkerId:     item.WorkerId,
				ObjectId:     &item.ObjectId,
				OrderId:      &item.OrderId,
				TaskId:       &item.TaskId,
				OperationId:  &item.OperationId,
				TaskWorkerId: &item.TaskWorkerId,
				TotalTime:    item.TotalTime,
				To:           item.To,
				From:         item.From,
				Oklad:        item.Oklad,
				Date:         item.Date,
				Total:        item.Total,
				CreatedAt:    item.CreatedAt,
				UpdatedAt:    item.UpdatedAt,
			}
			newItem["time"] = time.Now().UTC()
			newProps[time.Now().String()] = newItem

			totalTime := cuteTotalTime.Milliseconds()
			newUpdateData := &domain.WorkHistoryInput{
				TotalTime: &totalTime,
				Total:     &total,
				To:        cutTo,
				Props:     newProps,
			}
			result, err = s.repo.UpdateWorkHistory(id, userID, newUpdateData)
			if err != nil {
				return result, err
			}
		}
	}
	return result, err
}
