package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OperationService struct {
	repo        repository.Operation
	userService *UserService
}

func NewOperationService(repo repository.Operation, userService *UserService) *OperationService {
	return &OperationService{repo: repo, userService: userService}
}

func (s *OperationService) FindOperation(params domain.RequestParams) (domain.Response[domain.Operation], error) {
	return s.repo.FindOperation(params)
}

func (s *OperationService) CreateOperation(userID string, data *domain.Operation) (*domain.Operation, error) {
	var result *domain.Operation

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
	// 	updateReview := &domain.OperationInput{
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

	result, err = s.repo.CreateOperation(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	return result, err
}

func (s *OperationService) UpdateOperation(id string, userID string, data *domain.OperationInput) (*domain.Operation, error) {
	return s.repo.UpdateOperation(id, userID, data)
}

func (s *OperationService) DeleteOperation(id string) (*domain.Operation, error) {
	result, err := s.repo.DeleteOperation(id)

	return result, err
}
