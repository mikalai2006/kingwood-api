package service

import (
	"math"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkHistoryService struct {
	repo        repository.WorkHistory
	Hub         *Hub
	userService *UserService
	taskStatus  *TaskStatusService
}

func NewWorkHistoryService(repo repository.WorkHistory, hub *Hub, userService *UserService, TaskStatus *TaskStatusService) *WorkHistoryService {
	return &WorkHistoryService{repo: repo, Hub: hub, userService: userService, taskStatus: TaskStatus}
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
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	return result, err
}

func (s *WorkHistoryService) UpdateWorkHistory(id string, userID string, data *domain.WorkHistoryInput) (*domain.WorkHistory, error) {
	result, err := s.repo.UpdateWorkHistory(id, userID, data)
	if err != nil {
		return result, err
	}

	// update total.
	newRobotUpdateData := &domain.WorkHistoryInput{}
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

	// fmt.Println("======================PATCH WORK HISTORY====================")
	// fmt.Println("from: ", from1, "====>", from1.UTC())
	// fmt.Println("to: ", to1, "====>", to1.UTC())
	// fmt.Println("========================================================")

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

	result, err = s.repo.UpdateWorkHistory(id, userID, newRobotUpdateData)
	if err != nil {
		return result, err
	}

	if explodeDate {
		// Переносим часть рабочего времени на другой день
		totalMinutesNext := oldTo.Sub(fromNew).Minutes()
		totalNext := int64(math.Ceil(totalMinutesNext * (float64(*result.Oklad) / 60)))
		// fmt.Println("totalMinutesNext:", totalMinutesNext, " totalNext:", totalNext, " oldTo:", oldTo)
		result, err = s.repo.CreateWorkHistory(userID, &domain.WorkHistory{
			UserID:      result.UserID,
			WorkerId:    result.WorkerId,
			ObjectId:    result.ObjectId,
			OrderId:     result.OrderId,
			TaskId:      result.TaskId,
			OperationId: result.OperationId,
			Status:      result.Status,
			WorkTimeId:  result.WorkTimeId,
			Date:        fromNew,
			From:        fromNew,
			To:          oldTo,
			Oklad:       result.Oklad,
			Total:       &totalNext,
		})
		if err != nil {
			return nil, err
		}
	}

	return result, err
}

func (s *WorkHistoryService) DeleteWorkHistory(id string) (*domain.WorkHistory, error) {
	result, err := s.repo.DeleteWorkHistory(id)

	return result, err
}

func (s *WorkHistoryService) GetStatByOrder(input domain.WorkHistoryFilter) ([]domain.WorkHistoryStatByOrder, error) {
	result, err := s.repo.GetStatByOrder(input)

	return result, err
}
