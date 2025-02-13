package service

import (
	"fmt"
	"math"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkTimeService struct {
	repo        repository.WorkTime
	Hub         *Hub
	userService *UserService
	taskStatus  *TaskStatusService
}

func NewWorkTimeService(repo repository.WorkTime, hub *Hub, userService *UserService, TaskStatus *TaskStatusService) *WorkTimeService {
	return &WorkTimeService{repo: repo, Hub: hub, userService: userService, taskStatus: TaskStatus}
}

func (s *WorkTimeService) FindWorkTime(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error) {
	return s.repo.FindWorkTime(input)
}

func (s *WorkTimeService) FindWorkTimePopulate(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error) {
	return s.repo.FindWorkTimePopulate(input)
}

func (s *WorkTimeService) CreateWorkTime(userID string, data *domain.WorkTime) (*domain.WorkTime, error) {
	var result *domain.WorkTime

	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	// existReview, err := s.repo.FindReview(domain.RequestParams{
	// 	Filter:  bson.M{"node_id": review.NodeID, "userId": userIDPrimitive},
	// 	Options: domain.Options{Limit: 1},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// if len(existReview.Data) > 0 {
	// 	updateReview := &domain.TaskInput{
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

	result, err = s.repo.CreateWorkTime(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	return result, err
}

func (s *WorkTimeService) UpdateWorkTime(id string, userID string, data *domain.WorkTimeInput) (*domain.WorkTime, error) {
	// получаем данные из базы.
	existWorkTime, err := s.repo.FindWorkTimePopulate(domain.WorkTimeFilter{ID: []string{id}})
	if err != nil {
		return nil, err
	}

	if len(existWorkTime.Data) > 0 {
		// если данные для патча отличаются от данных из базы
		if existWorkTime.Data[0].From != data.From || existWorkTime.Data[0].To != data.To {
			// заносим старые данные в пропс.
			newProps := map[string]interface{}{}
			if existWorkTime.Data[0].Props != nil {
				newProps = existWorkTime.Data[0].Props
			}
			newItem := make(map[string]interface{})
			newItem["userId"] = userID
			newItem["item"] = domain.WorkTimeInput{
				UserID:    existWorkTime.Data[0].UserID,
				WorkerId:  existWorkTime.Data[0].WorkerId,
				To:        existWorkTime.Data[0].To,
				From:      existWorkTime.Data[0].From,
				Oklad:     existWorkTime.Data[0].Oklad,
				Date:      existWorkTime.Data[0].Date,
				Total:     existWorkTime.Data[0].Total,
				CreatedAt: existWorkTime.Data[0].CreatedAt,
				UpdatedAt: existWorkTime.Data[0].UpdatedAt,
			}
			newItem["time"] = time.Now()
			newProps[time.Now().String()] = newItem

			// дополняем пропс.
			data.Props = newProps
		}
	}

	result, err := s.repo.UpdateWorkTime(id, userID, data)
	if err != nil {
		return result, err
	}

	// update total.
	newRobotUpdateData := &domain.WorkTimeInput{}
	total := int64(0)
	if !result.From.IsZero() && !result.To.IsZero() {
		totalMinutes := result.To.Sub(result.From).Minutes()
		total = int64(math.Ceil(totalMinutes * (float64(*result.Oklad) / 60)))
	}

	if total > 0 {
		newRobotUpdateData.Total = &total
	}

	explodeDate := false
	oldTo := result.To
	var fromNew time.Time
	var toNew time.Time

	eastOfUTC := time.FixedZone("UTC-3", -3*60*60)
	to1 := time.Date(result.To.Year(), result.To.Month(), result.To.Day(), result.To.Hour(), result.To.Minute(), result.To.Second(), 0, eastOfUTC)
	from1 := time.Date(result.From.Year(), result.From.Month(), result.From.Day(), result.From.Hour(), result.From.Minute(), result.From.Second(), 0, eastOfUTC)

	fmt.Println("======================PATCH TIME WORK====================")
	fmt.Println("from: ", from1, "====>", from1.UTC())
	fmt.Println("to: ", to1, "====>", to1.UTC())
	fmt.Println("========================================================")

	// fmt.Println("result.From: ", to1, to1.UTC(), from1, from1.UTC())
	if from1.UTC().Day() != to1.UTC().Day() {
		explodeDate = true

		// prevDay := oldTo.AddDate(0, 0, -1)
		year, month, _ := oldTo.Date()
		yearPrev, monthPrev, dayPrev := from1.Date()
		// fromNew :=  result.From
		// time.Date(year, month, day, 0, 0, 0, 0, prevDay.Location())
		fromNew = time.Date(year, month, dayPrev, 21, 0, 0, 0, time.UTC)
		toNew = time.Date(yearPrev, monthPrev, dayPrev, 20, 59, 59, 0, time.UTC)

		newRobotUpdateData.To = toNew

		totalMinutesPrev := toNew.Sub(result.From).Minutes()
		totalPrev := int64(math.Ceil(totalMinutesPrev * (float64(*result.Oklad) / 60)))
		// update total.
		newRobotUpdateData.Total = &totalPrev
	}

	result, err = s.repo.UpdateWorkTime(id, userID, newRobotUpdateData)
	if err != nil {
		return result, err
	}

	if explodeDate {
		// Переносим часть рабочего времени на другой день
		totalMinutesNext := oldTo.Sub(fromNew).Minutes()
		totalNext := int64(math.Ceil(totalMinutesNext * (float64(*result.Oklad) / 60)))
		// fmt.Println("totalMinutesNext:", totalMinutesNext, " totalNext:", totalNext, " oldTo:", oldTo)
		result, err = s.repo.CreateWorkTime(userID, &domain.WorkTime{
			UserID:   result.UserID,
			WorkerId: result.WorkerId,
			Status:   result.Status,
			Date:     fromNew,
			From:     fromNew,
			To:       oldTo,
			Oklad:    result.Oklad,
			Total:    &totalNext,
		})
		if err != nil {
			return nil, err
		}
	}

	return result, err
}

func (s *WorkTimeService) DeleteWorkTime(id string) (*domain.WorkTime, error) {
	result, err := s.repo.DeleteWorkTime(id)

	return result, err
}
