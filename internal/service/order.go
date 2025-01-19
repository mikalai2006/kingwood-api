package service

import (
	"fmt"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	repo             repository.Order
	userService      *UserService
	Hub              *Hub
	operationService *OperationService
	Services         *Services
}

func NewOrderService(repo repository.Order, userService *UserService, hub *Hub, operationService *OperationService) *OrderService {
	return &OrderService{repo: repo, userService: userService, Hub: hub, operationService: operationService}
}

func (s *OrderService) FindOrder(input *domain.OrderFilter) (domain.Response[domain.Order], error) {
	return s.repo.FindOrder(input)
}

// func (s *OrderService) GetAllOrder(params domain.RequestParams) (domain.Response[domain.Order], error) {
// 	return s.repo.GetAllOrder(params)
// }

func (s *OrderService) CreateOrder(userID string, data *domain.Order) (*domain.Order, error) {
	var result *domain.Order

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
	// 	updateReview := &domain.OrderInput{
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

	result, err = s.repo.CreateOrder(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin", "boss"}})
	//domain.RequestParams{Filter: bson.M{"code": bson.D{{"$in", bson.A{"admin", "boss"}}}}})
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
		//domain.RequestParams{Filter: bson.M{"roleId": bson.D{{"$in", ids}}}})
		if err != nil {
			return nil, err
		}
		users = _users.Data

	}

	for i := range users {
		// add notify.
		_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo:  users[i].ID.Hex(),
			Title:   domain.NewOrderTitle,
			Message: fmt.Sprintf(domain.NewOrder, result.Name, result.Object.Name),
		})
	}

	fmt.Println("===============ORDER CREATE===================")
	fmt.Println("ids=", ids)
	fmt.Println("users=", len(users))
	fmt.Println("==============================================")

	return result, err
}

func (s *OrderService) UpdateOrder(id string, userID string, data *domain.OrderInput) (*domain.Order, error) {
	result, err := s.repo.UpdateOrder(id, userID, data)

	if data.StolyarComplete != nil || data.MalyarComplete != nil || data.MontajComplete != nil {
		statusCompleted := int64(1)

		dataUpdate := domain.OrderInput{}
		// fmt.Println("dataUpdate: ", data.MalyarComplete, data.StolyarComplete, data.MontajComplete)
		// fmt.Println("dataUpdate result: ", *result.MalyarComplete, *result.StolyarComplete, *result.MontajComplete)
		if (result.MalyarComplete != nil && *result.MalyarComplete == statusCompleted) &&
			(result.StolyarComplete != nil && *result.StolyarComplete == statusCompleted) &&
			(result.MontajComplete != nil && *result.MontajComplete == statusCompleted) {
			val := int64(100)
			dataUpdate.Status = &val
		} else {
			val := int64(1)
			dataUpdate.Status = &val
		}

		result, err = s.repo.UpdateOrder(id, userID, &dataUpdate)
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "order"})

	return result, err
}

func (s *OrderService) DeleteOrder(id string) (*domain.Order, error) {
	result, err := s.repo.DeleteOrder(id)

	return result, err
}
