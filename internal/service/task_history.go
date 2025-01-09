package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHistoryService struct {
	repo        repository.TaskHistory
	Hub         *Hub
	userService *UserService
	taskStatus  *TaskStatusService
}

func NewTaskHistoryService(repo repository.TaskHistory, hub *Hub, userService *UserService, TaskStatus *TaskStatusService) *TaskHistoryService {
	return &TaskHistoryService{repo: repo, Hub: hub, userService: userService, taskStatus: TaskStatus}
}

func (s *TaskHistoryService) FindTaskHistory(input domain.TaskHistoryFilter) (domain.Response[domain.TaskHistory], error) {
	return s.repo.FindTaskHistory(input)
}

func (s *TaskHistoryService) FindTaskHistoryPopulate(input domain.TaskHistoryFilter) (domain.Response[domain.TaskHistory], error) {
	return s.repo.FindTaskHistoryPopulate(input)
}

func (s *TaskHistoryService) CreateTaskHistory(userID string, data *domain.TaskHistory) (*domain.TaskHistory, error) {
	var result *domain.TaskHistory

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

	result, err = s.repo.CreateTaskHistory(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	return result, err
}

func (s *TaskHistoryService) UpdateTaskHistory(id string, userID string, data *domain.TaskHistoryInput) (*domain.TaskHistory, error) {
	result, err := s.repo.UpdateTaskHistory(id, userID, data)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *TaskHistoryService) DeleteTaskHistory(id string) (*domain.TaskHistory, error) {
	result, err := s.repo.DeleteTaskHistory(id)

	return result, err
}
