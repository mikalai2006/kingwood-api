package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PayService struct {
	repo        repository.Pay
	userService *UserService
	Hub         *Hub
}

func NewPayService(repo repository.Pay, userService *UserService, hub *Hub) *PayService {
	return &PayService{repo: repo, userService: userService, Hub: hub}
}

func (s *PayService) FindPay(params domain.RequestParams) (domain.Response[domain.Pay], error) {
	return s.repo.FindPay(params)
}

func (s *PayService) CreatePay(userID string, data *domain.Pay) (*domain.Pay, error) {
	var result *domain.Pay

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
	// 	updateReview := &domain.PayInput{
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

	result, err = s.repo.CreatePay(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	return result, err
}

func (s *PayService) UpdatePay(id string, userID string, data *domain.PayInput) (*domain.Pay, error) {
	result, err := s.repo.UpdatePay(id, userID, data)

	return result, err
}

func (s *PayService) DeletePay(id string) (*domain.Pay, error) {
	result, err := s.repo.DeletePay(id)

	return result, err
}
