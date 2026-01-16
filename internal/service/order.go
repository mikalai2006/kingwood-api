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

func (s *OrderService) FindOrder(input *domain.OrderFilter) (domain.ResponseOrderFlatData, error) {
	result, err := s.repo.FindOrder(input)
	if err != nil {
		return result, err
	}

	idsObjects := []string{}
	idsOrders := []string{}
	limit := 1000
	for i := range result.Data {
		idsObjects = append(idsObjects, result.Data[i].ObjectId.Hex())
		idsOrders = append(idsOrders, result.Data[i].ID.Hex())
	}

	tasks, err := s.Services.Task.FindTaskFlat(domain.TaskFilter{OrderId: idsOrders})
	if err != nil {
		return result, err
	}
	result.Tasks = tasks.Data

	objects, err := s.Services.Object.FindObject(&domain.ObjectFilter{ID: idsObjects, Limit: &limit})
	if err != nil {
		return result, err
	}
	result.Objects = objects.Data

	taskWorkers, err := s.Services.TaskWorker.FindTaskWorkerFlat(&domain.TaskWorkerFilter{OrderId: idsOrders, Limit: &limit})
	if err != nil {
		return result, err
	}
	result.TaskWorkers = taskWorkers.Data

	idsWorkers := []string{}
	for i := range taskWorkers.Data {
		idsWorkers = append(idsWorkers, taskWorkers.Data[i].WorkerId.Hex())
	}

	users, err := s.Services.User.FindUserFlat(&domain.UserFilter{ID: idsWorkers, Limit: &limit})
	if err != nil {
		return result, err
	}
	result.Users = users.Data

	// проходим по всем заданиям и выбираем workHistory для задания если у него указано maxHours.
	idsTasksForQueryWorkHistorys := []string{}
	for i := range result.Tasks {
		if result.Tasks[i].MaxHours > 0 {
			idsTasksForQueryWorkHistorys = append(idsTasksForQueryWorkHistorys, result.Tasks[i].ID.Hex())
		}
	}

	workHistorys, err := s.Services.WorkHistory.FindWorkHistory(domain.WorkHistoryFilter{TaskId: idsTasksForQueryWorkHistorys, Limit: &limit})
	if err != nil {
		return result, err
	}
	outputWorkHistory := domain.Response[domain.WorkHistoryFlat]{}
	for i := range workHistorys.Data {
		outputWorkHistory.Data = append(outputWorkHistory.Data, domain.WorkHistoryFlat{
			ID:           workHistorys.Data[i].ID,
			UserID:       workHistorys.Data[i].UserID,
			ObjectId:     workHistorys.Data[i].ObjectId,
			OrderId:      workHistorys.Data[i].OrderId,
			TaskId:       workHistorys.Data[i].TaskId,
			WorkerId:     workHistorys.Data[i].WorkerId,
			OperationId:  workHistorys.Data[i].OperationId,
			TaskWorkerId: workHistorys.Data[i].TaskWorkerId,
			Status:       workHistorys.Data[i].Status,
			Date:         workHistorys.Data[i].Date,
			From:         workHistorys.Data[i].From,
			To:           workHistorys.Data[i].To,
			Oklad:        workHistorys.Data[i].Oklad,
			Total:        workHistorys.Data[i].Total,
			TotalTime:    workHistorys.Data[i].TotalTime,
			Props:        workHistorys.Data[i].Props,
			CreatedAt:    workHistorys.Data[i].CreatedAt,
			UpdatedAt:    workHistorys.Data[i].UpdatedAt,
		})
	}
	result.WorkHistorys = outputWorkHistory.Data

	return result, err
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
		} else {
			data.Number = 1
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
			Message: fmt.Sprintf(domain.NewOrder, authorCreate.Name, result.Number, result.Name, ""), // result.Object.Name
		})
	}

	// fmt.Println("===============ORDER CREATE===================")
	// fmt.Println("ids=", ids)
	// fmt.Println("users=", len(users))
	// fmt.Println("==============================================")

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
	unLim := 0

	// delete images.
	allImages, err := s.Services.Image.FindImage(&domain.ImageFilter{ServiceId: []string{id}})
	// domain.RequestParams{Filter: bson.D{{"serviceId", id}}}
	if err != nil {
		return result, err
	}
	for i := range allImages.Data {
		_, err = s.Services.Image.DeleteImage(userID, allImages.Data[i].ID.Hex(), true)
		if err != nil {
			return result, err
		}
	}
	// delete taskWorkers.
	allTaskWorkers, err := s.Services.TaskWorker.FindTaskWorkerFlat(&domain.TaskWorkerFilter{OrderId: []string{id}, Limit: &unLim})
	if err != nil {
		return result, err
	}
	for i := range allTaskWorkers.Data {
		_, err = s.Services.TaskWorker.DeleteTaskWorker(allTaskWorkers.Data[i].ID.Hex(), userID, false, false)
		if err != nil {
			return result, err
		}
	}

	// delete task.
	allTasks, err := s.Services.Task.FindTaskFlat(domain.TaskFilter{OrderId: []string{id}, Limit: &unLim})
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
	allMessages, err := s.Services.Message.FindMessage(&domain.MessageFilter{OrderID: []string{id}, Limit: &unLim})
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
	allWorkHistory, err := s.Services.WorkHistory.FindWorkHistory(domain.WorkHistoryFilter{OrderId: []string{id}, Limit: &unLim})
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
