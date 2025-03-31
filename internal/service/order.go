package service

import (
	"fmt"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
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

	// получаем пользователя, который добавил изделие.
	var authorCreate domain.User
	_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
	if err != nil {
		return nil, err
	}
	if len(_users.Data) > 0 {
		authorCreate = _users.Data[0]
	}

	// existReview, err := s.repo.FindReview(domain.RequestParams{
	// 	Filter:  bson.M{"node_id": review.NodeID, "userId": userIDPrimitive},
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

	// запрос на последний номер заказа
	if data.Number == 0 {
		lastOrders, err := s.FindOrder(&domain.OrderFilter{Year: data.Year, Sort: []*domain.FilterSortParams{{Key: "number", Value: -1}}})
		if err != nil {
			return nil, err
		}
		if len(lastOrders.Data) > 0 {
			data.Number = lastOrders.Data[0].Number + 1
		}
	} else {
		// // если в данных есть number, проверяем на существование такого номера
		// existOrders, err := s.FindOrder(&domain.OrderFilter{Year: data.Year, Number: &data.Number})
		// if err != nil {
		// 	return nil, err
		// }
		// if len(existOrders.Data) > 0 {
		// 	return nil, domain.ErrExistNumberOrder
		// }
	}

	result, err = s.repo.CreateOrder(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin", "boss", "systemrole"}})
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
			Message: fmt.Sprintf(domain.NewOrder, authorCreate.Name, result.Number, result.Name, result.Object.Name),
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

	// allOrderTasks, err := s.Services.Task.FindTaskPopulate(domain.TaskFilter{OrderId: []string{id}})
	// if err != nil {
	// 	return result, err
	// }

	// if data.StolyarComplete != nil || data.MalyarComplete != nil || data.MontajComplete != nil {
	// 	statusCompleted := int64(1)

	// 	dataUpdate := domain.OrderInput{}
	// 	// fmt.Println("dataUpdate: ", data.MalyarComplete, data.StolyarComplete, data.MontajComplete)
	// 	// fmt.Println("dataUpdate result: ", *result.MalyarComplete, *result.StolyarComplete, *result.MontajComplete)
	// 	if (result.MalyarComplete != nil && *result.MalyarComplete == statusCompleted) &&
	// 		(result.StolyarComplete != nil && *result.StolyarComplete == statusCompleted) &&
	// 		(result.MontajComplete != nil && *result.MontajComplete == statusCompleted) &&
	// 		(result.ShlifComplete != nil && *result.ShlifComplete == statusCompleted) {
	// 		val := int64(100)
	// 		dataUpdate.Status = &val
	// 	} else {
	// 		statusOrder := int64(1)
	// 		if len(result.Tasks) == 0 {
	// 			statusOrder = 0
	// 		}
	// 		dataUpdate.Status = &statusOrder
	// 	}

	// 	result, err = s.repo.UpdateOrder(id, userID, &dataUpdate)
	// }

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "order"})

	return result, err
}

func (s *OrderService) DeleteOrder(id string, userID string) (*domain.Order, error) {
	var result *domain.Order

	// delete images.
	allImages, err := s.Services.Image.FindImage(domain.RequestParams{Filter: bson.D{{"serviceId", id}}})
	if err != nil {
		return result, err
	}
	for i := range allImages.Data {
		_, err = s.Services.Image.DeleteImage(userID, allImages.Data[i].ID.Hex())
		if err != nil {
			return result, err
		}
	}
	// delete taskWorkers.
	allTaskWorkers, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{OrderId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allTaskWorkers.Data {
		_, err = s.Services.TaskWorker.DeleteTaskWorker(allTaskWorkers.Data[i].ID.Hex(), userID, false)
		if err != nil {
			return result, err
		}
	}

	// delete task.
	allTasks, err := s.Services.Task.FindTaskPopulate(domain.TaskFilter{OrderId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allTasks.Data {
		_, err = s.Services.Task.DeleteTask(allTasks.Data[i].ID.Hex(), userID, false)
		if err != nil {
			return result, err
		}
	}

	// delete messages.
	allMessages, err := s.Services.Message.FindMessage(&domain.MessageFilter{OrderID: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allMessages.Data {
		_, err = s.Services.Message.DeleteMessage(allMessages.Data[i].ID.Hex(), userID)
		if err != nil {
			return result, err
		}
	}

	// delete workHistory.
	allWorkHistory, err := s.Services.WorkHistory.FindWorkHistory(domain.WorkHistoryFilter{OrderId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allWorkHistory.Data {
		_, err = s.Services.WorkHistory.DeleteWorkHistory(allWorkHistory.Data[i].ID.Hex(), userID, false)
		if err != nil {
			return result, err
		}
	}

	// // delete workHistory.
	// allWorkHistory, err := s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{WorkerId: []string{id}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allWorkHistory.Data {
	// 	_, err = s.Services.WorkHistory.DeleteWorkHistory(allWorkHistory.Data[i].ID.Hex())
	// }

	// // delete pay.
	// allPay, err := s.Services.Pay.FindPay(&domain.PayFilter{WorkerId: []string{id}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allPay.Data {
	// 	_, err = s.Services.Pay.DeletePay(allPay.Data[i].ID.Hex(), userID)
	// }

	// // delete notify.
	// allNotify, err := s.Services.Notify.FindNotifyPopulate(&domain.NotifyFilter{UserTo: []*string{&id}})
	// if err != nil {
	// 	return result, err
	// }
	// for i := range allNotify.Data {
	// 	_, err = s.Services.Notify.DeleteNotify(allNotify.Data[i].ID.Hex())
	// }
	result, err = s.repo.DeleteOrder(id)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "order"})

	_, err = s.Services.CreateArchiveOrder(userID, result)

	return result, err
}
