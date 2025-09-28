package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type TimerService struct {
	repo     repository.Timer
	Hub      *Hub
	Services *Services
}

func NewTimerService(repo repository.Timer, hub *Hub) *TimerService {
	return &TimerService{repo: repo, Hub: hub}
}

func (s *TimerService) FindTimer(params domain.RequestParams) (domain.Response[domain.TimerShedule], error) {
	return s.repo.FindTimer(params)
}

func (s *TimerService) FindTimerPopulate(filter domain.TimerSheduleFilter) (domain.Response[domain.TimerShedule], error) {
	return s.repo.FindTimerPopulate(filter)
}

func (s *TimerService) CreateTimer(userID string, data *domain.TimerShedule) (*domain.TimerShedule, error) {
	var result *domain.TimerShedule

	// _, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }

	result, err := s.repo.CreateTimer(userID, data)
	if err != nil {
		return nil, err
	}

	// создаем таймер
	timer1 := StartTimer(*data)

	// добавляем задачу, которая будет выполнена по таймеру
	if timer1 != nil {
		// Отслеживаем таймер в контексте приложения
		go func() {
			<-timer1.C
			// fmt.Println("Задача "+data.IDTimer+" выполнена: workHistory=>", data.WorkHistoryId)
			s.StopTimer(result.ID.Hex(), userID)
		}()
	}

	return result, err
}

func (s *TimerService) UpdateTimer(id string, userID string, data *domain.TimerSheduleInput) (*domain.TimerShedule, error) {
	result, err := s.repo.UpdateTimer(id, userID, data)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *TimerService) StopTimer(id string, userID string) (*domain.TimerShedule, error) {
	var result *domain.TimerShedule

	// получаем данные для таймера из бд.
	timerData, err := s.FindTimerPopulate(domain.TimerSheduleFilter{
		ID: []string{id},
	})
	if err != nil {
		return nil, err
	}

	if len(timerData.Data) > 0 {
		// если таймер еще не остановлен,
		// выполняем задачу
		if timerData.Data[0].IsRunning == 1 {
			isRunining := 0
			// Обновляем статус задачи в БД после выполнения
			// в базе изменяем статус таймера на отработано
			result, err = s.UpdateTimer(timerData.Data[0].ID.Hex(), userID, &domain.TimerSheduleInput{
				IsRunning: &isRunining,
			})
			if err != nil {
				return nil, err
			}

			// находим данные для статуса паузы
			taskStatus, _ := s.Services.TaskStatus.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"status", "pause"}}})
			// запускаем остановку работы для работника, для которого был создан таймер
			if len(taskStatus.Data) > 0 {
				s.Services.TaskWorker.UpdateTaskWorker(result.TaskWorkerId.Hex(), userID, &domain.TaskWorkerInput{
					Status:   taskStatus.Data[0].Status,
					StatusId: taskStatus.Data[0].ID,
					WorkerId: result.WorkerId,
				}, 1)

				// достаем все сессии за текущие сутки и считаем общее время.
				listWorkHistoryByDay, err := s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{
					WorkerId: []string{result.WorkerId.Hex()},
					Date:     result.ExecuteAt,
					Sort: []*domain.FilterSortParams{
						{
							Key:   "createdAt",
							Value: -1,
						},
					},
				})
				if err != nil {
					return nil, err
				}

				// подсчитываем общее время всех сессий за дату текущей рабочей сессии.
				var allWorkTime int64
				for i := range listWorkHistoryByDay.Data {
					allWorkTime = allWorkTime + *listWorkHistoryByDay.Data[i].TotalTime
				}

				// находим пользователей(администрацию) для подписи уведомлений таймера.
				roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"systemrole"}})
				if err != nil {
					return nil, err
				}
				ids := []string{}
				var userSender domain.User

				if len(roles.Data) > 0 {
					users, err := s.Services.User.FindUser(&domain.UserFilter{RoleId: ids})
					if err != nil {
						return nil, err
					}
					if len(users.Data) > 0 {
						userSender = users.Data[0]
					}
				}

				// протяженность рабочей сессии
				durationNew := time.Duration(allWorkTime * int64(time.Millisecond))
				_durationNewText, _ := time.ParseDuration(durationNew.String())
				durationNewText := strings.Replace(_durationNewText.String(), "h", "ч.", 1)
				durationNewText = strings.Replace(durationNewText, "m", "мин.", 1)
				durationNewText = strings.Replace(durationNewText, "s", "сек.", 1)

				_, err = s.Services.Notify.CreateNotify(userSender.ID.Hex(), &domain.NotifyInput{
					UserTo: timerData.Data[0].WorkerId.Hex(),
					Title:  domain.StopTimerTitle,
					Message: fmt.Sprintf(
						domain.StopTimer,
						// userSender.Name,
						listWorkHistoryByDay.Data[0].Order.Name,
						// fmt.Sprintf("%d-%d", result.Year, result.Month+1),
						durationNewText,
					),
				})
				if err != nil {
					return nil, err
				}
			}
		} else {
			fmt.Println("Таймер сработал, но выполнение задач отменено, так как статус - УЖЕ ВЫПОЛНЕНО")
		}
	}
	return result, err
}

func (s *TimerService) DeleteTimer(id string, userID string) (*domain.TimerShedule, error) {
	result, err := s.repo.DeleteTimer(id)

	return result, err
}

func StartTimer(task domain.TimerShedule) *time.Timer {
	duration := time.Until(task.ExecuteAt)
	if duration > 0 {
		timer := time.NewTimer(duration)
		return timer
	}
	return nil
}

func StopTimer(timer *time.Timer, taskID string) {
	if timer != nil && !timer.Stop() {
		// Если таймер уже сработал, необходимо прочитать из канала
		<-timer.C
	}
	// Обновляем статус задачи в базе данных
	// ...
}

func (s *TimerService) RecoveryTimers() (*domain.TimerShedule, error) {
	var result *domain.TimerShedule

	// находим в базе суперадмина
	roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"systemrole"}})
	if err != nil {
		return nil, err
	}
	ids := []string{}
	var sadmin domain.User

	if len(roles.Data) > 0 {
		for i := range roles.Data {
			ids = append(ids, roles.Data[i].ID.Hex())
		}

		_users, err := s.Services.User.FindUser(&domain.UserFilter{RoleId: ids})
		if err != nil {
			return nil, err
		}

		sadmin = _users.Data[0]
	}

	// запрашиваем все таймеры, которые не выполнены или время выполнения прошло.
	isRunining := 1
	timers, err := s.Services.Timer.FindTimerPopulate(domain.TimerSheduleFilter{
		IsRunning: &isRunining,
	})
	if err != nil {
		return result, err
	}
	fmt.Println("Recovery ", len(timers.Data), " timer(s)")

	if len(timers.Data) > 0 {
		for i, _ := range timers.Data {
			currentTime := time.Now()

			// сравниваем дату таймера из бд с текущим временем и датой.
			if currentTime.After(timers.Data[i].ExecuteAt) {
				// если время вышло, запускаем задачу для выполнения.
				s.StopTimer(timers.Data[i].ID.Hex(), sadmin.ID.Hex())
			} else {
				// если время не вышло, создаем и запускаем таймер.

				// создаем таймер с задачей
				timer1 := StartTimer(timers.Data[i])

				// добавляем задачу, которая будет выполнена по таймеру
				if timer1 != nil {
					// Отслеживаем таймер в контексте приложения
					go func() {
						<-timer1.C
						// fmt.Println("Задача "+timers.Data[i].IDTimer+" выполнена: workHistory=>", timers.Data[i].WorkHistoryId)
						s.StopTimer(timers.Data[i].ID.Hex(), sadmin.ID.Hex())
					}()
				}

			}
		}
	}

	return result, err
}
