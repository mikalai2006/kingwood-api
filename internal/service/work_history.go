package service

import (
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

	return result, err
}

func (s *WorkHistoryService) DeleteWorkHistory(id string) (*domain.WorkHistory, error) {
	result, err := s.repo.DeleteWorkHistory(id)

	return result, err
}
