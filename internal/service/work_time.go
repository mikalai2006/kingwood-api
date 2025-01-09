package service

import (
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
	// 	Filter:  bson.M{"node_id": review.NodeID, "user_id": userIDPrimitive},
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
	result, err := s.repo.UpdateWorkTime(id, userID, data)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *WorkTimeService) DeleteWorkTime(id string) (*domain.WorkTime, error) {
	result, err := s.repo.DeleteWorkTime(id)

	return result, err
}
