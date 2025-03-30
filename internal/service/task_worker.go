package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskWorkerService struct {
	repo              repository.TaskWorker
	userService       *UserService
	taskStatusService *TaskStatusService
	taskService       *TaskService
	Hub               *Hub
	Services          *Services
}

func NewTaskWorkerService(repo repository.TaskWorker, userService *UserService, taskStatusService *TaskStatusService, taskService *TaskService, hub *Hub) *TaskWorkerService {
	return &TaskWorkerService{repo: repo, userService: userService, taskStatusService: taskStatusService, taskService: taskService, Hub: hub}
}

func (s *TaskWorkerService) FindTaskWorkerPopulate(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorker], error) {
	var result domain.Response[domain.TaskWorker]

	if input.Query != "" {
		orders, err := s.Services.Order.FindOrder(&domain.OrderFilter{Query: input.Query})
		if err != nil {
			return result, err
		}

		if len(orders.Data) > 0 {
			for i := range orders.Data {
				input.OrderId = append(input.OrderId, orders.Data[i].ID.Hex())
			}
		} else {
			input.OrderId = []string{primitive.NilObjectID.Hex()}
		}
	}

	result, err := s.repo.FindTaskWorkerPopulate(input)

	return result, err
}

// func (s *TaskWorkerService) FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error) {
// 	return s.repo.FindTaskWorker(params)
// }

func (s *TaskWorkerService) CreateTaskWorker(userID string, data *domain.TaskWorker, autoCreate int) (*domain.TaskWorker, error) {
	var result *domain.TaskWorker

	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// получаем пользователя, который добавил задание.
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
	// 	updateReview := &domain.TaskWorkerInput{
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

	result, err = s.repo.CreateTaskWorker(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	_, err = s.CheckStatusTask(userID, result)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "taskWorker"})

	// add notify.
	_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
		UserTo:     result.WorkerId.Hex(),
		Title:      domain.CreateTaskWorkerTitle,
		Message:    fmt.Sprintf(domain.CreateTaskWorker, authorCreate.Name, result.Task.Name, result.Order.Number, result.Order.Name, result.Object.Name),
		Link:       "/(tabs)/order",
		LinkOption: map[string]interface{}{"orderId": result.OrderId.Hex(), "objectId": result.ObjectId.Hex()},
	})

	// находим пользователей для создания уведомлений.
	roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin"}})
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
		if err != nil {
			return nil, err
		}

		users = _users.Data
	}

	// получаем пользователя, для которого изменили рабочую сессию.
	var worker domain.User
	_workers, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{result.WorkerId.Hex()}})
	if err != nil {
		return nil, err
	}
	if len(_workers.Data) > 0 {
		worker = _workers.Data[0]
	}

	// отправляем уведомления админам.
	for i := range users {
		_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo:     users[i].ID.Hex(),
			Title:      domain.CreateTaskWorkerTitle,
			Message:    fmt.Sprintf(domain.CreateTaskWorkerAdmin, authorCreate.Name, worker.Name, result.Task.Name, result.Order.Number, result.Order.Name, result.Object.Name),
			Link:       "/(tabs)/order",
			LinkOption: map[string]interface{}{"orderId": result.OrderId.Hex(), "objectId": result.ObjectId.Hex()},
		})
	}

	// add taskWorker for all task on the object for inserted worker (montaj).
	if autoCreate > 0 {
		allOperation, err := s.Services.Operation.FindOperation(domain.RequestParams{Filter: bson.D{}})
		if err != nil {
			return result, err
		}
		var currentOperation *domain.Operation
		for i := range allOperation.Data {
			if allOperation.Data[i].ID.Hex() == result.OperationId.Hex() {
				currentOperation = &allOperation.Data[i]
			}
		}
		// fmt.Println("allOperation length: ", len(allOperation.Data), currentOperation.Group)
		if currentOperation.Group == "5" {
			allTaskForObject, err := s.Services.Task.FindTaskPopulate(domain.TaskFilter{ObjectId: []string{result.ObjectId.Hex()}, OperationId: []string{result.OperationId.Hex()}})
			if err != nil {
				return result, err
			}
			fmt.Println("allTaskForObject length: ", len(allTaskForObject.Data))
			if len(allTaskForObject.Data) > 0 {
				for i := range allTaskForObject.Data {
					if allTaskForObject.Data[i].Status != "finish" {
						workerIds := []string{}
						for j := range allTaskForObject.Data[i].Workers {
							workerIds = append(workerIds, allTaskForObject.Data[i].Workers[j].WorkerId.Hex())
						}
						if allTaskForObject.Data[i].ID.Hex() != result.TaskId.Hex() && !utils.Contains(workerIds, result.WorkerId.Hex()) {
							newTaskWorker := domain.TaskWorker{
								ObjectId:    allTaskForObject.Data[i].ObjectId,
								OrderId:     allTaskForObject.Data[i].OrderId,
								TaskId:      allTaskForObject.Data[i].ID,
								OperationId: allTaskForObject.Data[i].OperationId,
								WorkerId:    result.WorkerId,
								SortOrder:   result.SortOrder,
								StatusId:    result.StatusId,
								Status:      result.Status,
								From:        result.From,
								To:          result.To,
								TypeGo:      result.TypeGo,
							}
							_, err := s.CreateTaskWorker(userID, &newTaskWorker, 0)
							if err != nil {
								return result, err
							}
						}
					}
				}
			}
		}
	}

	return result, err
}

func (s *TaskWorkerService) UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput, autoUpdate int) (*domain.TaskWorker, error) {
	// инициализация данных.
	var result *domain.TaskWorker

	// запрос на существование задачи для пользователя.
	existTaskWorker, err := s.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{ID: []string{id}})
	if err != nil {
		return result, err
	}

	if len(existTaskWorker.Data) > 0 {
		// если задача есть, обрабатываем в обычном режиме.

		result, err = s.repo.UpdateTaskWorker(id, userID, data)
		if err != nil {
			return result, err
		}

	} else {
		// иначе создаем фейковые данные для поддержки приложения.

		statuses, err := s.Services.TaskStatus.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"name", data.Status}}})
		if err != nil {
			return result, err
		}

		statusId := data.StatusId
		status := data.Status
		if len(statuses.Data) > 0 {
			statusId = statuses.Data[0].ID
			status = statuses.Data[0].Status
		}

		from := time.Now().AddDate(-1, 0, 0)
		to := time.Now().AddDate(1, 0, 0)
		// if !data.From.IsZero() {
		// 	from = data.From
		// }
		// if !data.To.IsZero() {
		// 	to = data.To
		// }

		// получаем пользователя - исполнителя.
		var workerTaskWorker domain.User
		_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{data.WorkerId.Hex()}})
		if err != nil {
			return nil, err
		}
		if len(_users.Data) > 0 {
			workerTaskWorker = _users.Data[0]
		}

		result = &domain.TaskWorker{
			ID:          data.ID,
			UserID:      data.UserID,
			ObjectId:    data.ObjectId,
			OrderId:     data.OrderId,
			TaskId:      data.TaskId,
			WorkerId:    data.WorkerId,
			OperationId: data.OperationId,
			SortOrder:   data.SortOrder,
			StatusId:    statusId,
			Status:      status,
			From:        &from,
			To:          &to,
			Worker:      workerTaskWorker,
			TypeGo:      data.TypeGo,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	// fmt.Println("data: ", data)
	// fmt.Println("result: ", result)

	// if result != nil && !result.ID.IsZero() {
	// получаем пользователя, который добавил задание.
	var author domain.User
	_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
	if err != nil {
		return nil, err
	}
	if len(_users.Data) > 0 {
		author = _users.Data[0]
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "taskWorker"})

	// roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin", "boss"}})
	// if err != nil {
	// 	return nil, err
	// }
	// ids := []string{}
	// var users []domain.User
	// if len(roles.Data) > 0 {
	// 	for i := range roles.Data {
	// 		ids = append(ids, roles.Data[i].ID.Hex())
	// 	}

	// 	_users, err := s.Services.User.FindUser(&domain.UserFilter{RoleId: ids})
	// 	//domain.RequestParams{Filter: bson.M{"roleId": bson.D{{"$in", ids}}}})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	users = _users.Data

	// }

	// for i := range users {
	// 	_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
	// 		UserTo:  users[i].ID.Hex(),
	// 		Title:   domain.NewOrderTitle,
	// 		Message: fmt.Sprintf(domain.NewOrder, result.Name, result.Object.Name),
	// 	})
	// }

	// add notify.
	fmt.Println("workerID:", data.WorkerId)
	fmt.Println("userID:", userID)
	if result.WorkerId.Hex() != userID && !result.OrderId.IsZero() {
		_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo:     result.WorkerId.Hex(),
			Title:      domain.PatchTaskWorkerTitle,
			Message:    fmt.Sprintf(domain.PatchTaskWorker, author.Name, result.Task.Name, result.Order.Number, result.Order.Name, result.Object.Name),
			Link:       "/(tabs)/order",
			LinkOption: map[string]interface{}{"orderId": result.OrderId.Hex(), "objectId": result.ObjectId.Hex()},
		})

		// находим админов.
		roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin"}})
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
			if err != nil {
				return nil, err
			}

			users = _users.Data
		}

		// получаем пользователя, для которого изменили задание.
		var worker domain.User
		_workers, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{result.WorkerId.Hex()}})
		if err != nil {
			return nil, err
		}
		if len(_workers.Data) > 0 {
			worker = _workers.Data[0]
		}

		// отправляем уведомления админам.
		for i := range users {
			_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
				UserTo:     users[i].ID.Hex(),
				Title:      domain.PatchTaskWorkerTitle,
				Message:    fmt.Sprintf(domain.PatchTaskWorkerAdmin, author.Name, worker.Name, result.Task.Name, result.Order.Number, result.Order.Name, result.Object.Name),
				Link:       "/(tabs)/order",
				LinkOption: map[string]interface{}{"orderId": result.OrderId.Hex(), "objectId": result.ObjectId.Hex()},
			})
		}
	}

	// change taskWorker for all task on the object for updated worker (montaj).
	if autoUpdate > 0 && !result.ID.IsZero() {
		_, err = s.CheckStatusTask(userID, result)
		if err != nil {
			return result, err
		}

		allOperation, err := s.Services.Operation.FindOperation(domain.RequestParams{Filter: bson.D{}})
		if err != nil {
			return result, err
		}
		var currentOperation *domain.Operation
		for i := range allOperation.Data {
			if allOperation.Data[i].ID.Hex() == result.OperationId.Hex() {
				currentOperation = &allOperation.Data[i]
			}
		}
		// fmt.Println("allOperation length: ", len(allOperation.Data), currentOperation.Group)
		if currentOperation.Group == "5" {
			objectId := result.ObjectId.Hex()
			operationId := result.OperationId.Hex()
			workerId := result.WorkerId.Hex()
			allTaskWorkerForObject, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{ObjectId: []string{objectId}, WorkerId: []string{workerId}, OperationId: []string{operationId}})
			if err != nil {
				return result, err
			}
			// fmt.Println("allTaskForObject length: ", len(allTaskForObject.Data))
			if len(allTaskWorkerForObject.Data) > 0 {
				stopStatus := []string{"finish", "process", "autofinish"}
				// statusNotChange := []string{"finish",""}
				for i := range allTaskWorkerForObject.Data {
					if allTaskWorkerForObject.Data[i].ID.Hex() != result.ID.Hex() && !utils.Contains(stopStatus, allTaskWorkerForObject.Data[i].Status) {
						newTaskWorker := domain.TaskWorkerInput{
							// ObjectId:    allTaskWorkerForObject.Data[i].ObjectId,
							// OrderId:     allTaskWorkerForObject.Data[i].OrderId,
							// TaskId:      allTaskWorkerForObject.Data[i].ID,
							// OperationId: allTaskWorkerForObject.Data[i].OperationId,
							// WorkerId:    result.WorkerId,
							// SortOrder:   result.SortOrder,
							// StatusId:    result.StatusId,
							// Status:      result.Status,
							From:   *result.From,
							To:     *result.To,
							TypeGo: result.TypeGo,
						}
						_, err := s.UpdateTaskWorker(allTaskWorkerForObject.Data[i].ID.Hex(), userID, &newTaskWorker, 0)
						if err != nil {
							return result, err
						}
					}
				}
			}
		}
	}

	// останавливаем все выполняемые taskHistory для taskWorker.
	status := 0
	existOpenWorkHistory, err := s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{WorkerId: []string{result.WorkerId.Hex()}, TaskWorkerId: []string{id}, Status: &status, Sort: []*domain.FilterSortParams{{Key: "createdAt", Value: -1}}}) //TaskId: []string{result.TaskId.Hex()},
	if err != nil {
		return result, err
	}

	// если сам пользователь меняет, то останавливаем все taskHistory
	if result.WorkerId.Hex() == userID {
		existOpenWorkHistory, err = s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{WorkerId: []string{result.WorkerId.Hex()}, Status: &status, Sort: []*domain.FilterSortParams{{Key: "createdAt", Value: -1}}}) //TaskId: []string{result.TaskId.Hex()},
		if err != nil {
			return result, err
		}
	}

	for i := range existOpenWorkHistory.Data {
		// if len(existOpenWorkHistory.Data) > 0 {
		statusPatch := 1
		s.Services.WorkHistory.UpdateWorkHistory(existOpenWorkHistory.Data[i].ID.Hex(), userID, &domain.WorkHistoryInput{
			Status: &statusPatch,
			To:     time.Now(),
		})
		// }
	}

	// начинаем новый taskHistory, если запрос делает исполнитель задания.
	if result.Status == "process" && result.Worker.ID.Hex() == userID {
		newWorkHistory := domain.WorkHistory{
			ObjectId:     result.ObjectId,
			OrderId:      result.OrderId,
			TaskId:       result.TaskId,
			WorkerId:     result.WorkerId,
			OperationId:  result.OperationId,
			TaskWorkerId: result.ID,
			Status:       0,
			From:         time.Now(),
			Oklad:        result.Worker.Oklad,
		}

		// create workHistory from.
		s.Services.WorkHistory.CreateWorkHistory(userID, &newWorkHistory)
	}
	//  else {
	// 	// close wortHistory to.
	// 	existOpenWorkHistory, err := s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{WorkerId: []string{result.WorkerId.Hex()}, TaskId: []string{result.TaskId.Hex()}, Status: &status, Sort: []*domain.FilterSortParams{{Key: "createdAt", Value: -1}}})
	// 	if err != nil {
	// 		return result, err
	// 	}

	// 	if len(existOpenWorkHistory.Data) > 0 {
	// 		statusPatch := 1
	// 		s.Services.WorkHistory.UpdateWorkHistory(existOpenWorkHistory.Data[0].ID.Hex(), userID, &domain.WorkHistoryInput{
	// 			Status: &statusPatch,
	// 			To:     time.Now(),
	// 		})
	// 	}
	// }
	// }

	return result, err
}

func (s *TaskWorkerService) DeleteTaskWorker(id string, userID string, checkStatus bool) (*domain.TaskWorker, error) {
	result, err := s.repo.DeleteTaskWorker(id)
	if err != nil {
		return result, err
	}

	if checkStatus {
		_, err = s.CheckStatusTask("userID", result)
		if err != nil {
			return result, err
		}
	}

	// получаем инициатора запрос.
	var authorCreate domain.User
	_users, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{userID}})
	if err != nil {
		return nil, err
	}
	if len(_users.Data) > 0 {
		authorCreate = _users.Data[0]
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: "", Content: result, ID: "room1", Service: "taskWorker"})

	// Закрываем work_history если он есть для удаляемого исполнителя
	statusHistory := 0
	existOpenWorkHistory, err := s.Services.WorkHistory.FindWorkHistoryPopulate(domain.WorkHistoryFilter{WorkerId: []string{result.WorkerId.Hex()}, TaskId: []string{result.TaskId.Hex()}, Status: &statusHistory})
	if err != nil {
		return result, err
	}
	if len(existOpenWorkHistory.Data) > 0 {
		statusPatch := 1
		s.Services.WorkHistory.UpdateWorkHistory(existOpenWorkHistory.Data[0].ID.Hex(), userID, &domain.WorkHistoryInput{
			Status: &statusPatch,
			To:     time.Now(),
		})
	}

	// add notify.
	_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
		UserTo:  result.WorkerId.Hex(),
		Title:   domain.DeleteTaskWorkerTitle,
		Message: fmt.Sprintf(domain.DeleteTaskWorker, authorCreate.Name, result.Task.Name, result.Order.Number, result.Order.Name, result.Object.Name),
	})

	// находим пользователей для создания уведомлений.
	roles, err := s.Services.Role.FindRole(&domain.RoleFilter{Code: []string{"admin"}})
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
		if err != nil {
			return nil, err
		}

		users = _users.Data
	}

	// получаем пользователя, для которого удалили задание.
	var worker domain.User
	_workers, err := s.Services.User.FindUser(&domain.UserFilter{ID: []string{result.WorkerId.Hex()}})
	if err != nil {
		return nil, err
	}
	if len(_workers.Data) > 0 {
		worker = _workers.Data[0]
	}

	// отправляем уведомления админам.
	for i := range users {
		_, err = s.Services.Notify.CreateNotify(userID, &domain.NotifyInput{
			UserTo:     users[i].ID.Hex(),
			Title:      domain.CreateTaskWorkerTitle,
			Message:    fmt.Sprintf(domain.DeleteTaskWorkerAdmin, authorCreate.Name, worker.Name, result.Task.Name, result.Order.Number, result.Order.Name, result.Object.Name),
			Link:       "/(tabs)/order",
			LinkOption: map[string]interface{}{"orderId": result.OrderId.Hex(), "objectId": result.ObjectId.Hex()},
		})
	}

	// add to archive.
	_, err = s.Services.CreateArchiveTaskWorker(userID, result)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *TaskWorkerService) CheckStatusTask(userID string, result *domain.TaskWorker) (*domain.TaskWorker, error) {

	// currentTask, err := s.taskService.FindTask(domain.RequestParams{Filter: bson.D{{"_id", result.TaskId}}})
	// if err != nil {
	// 	return result, err
	// }
	if result == nil {
		return result, errors.New("not found task")
	}

	// get all taskWorkers.
	fmt.Println("taskId: ", result.TaskId, result.Task.Operation.Group)
	taskId := result.TaskId.Hex()
	taskWorkers, err := s.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{TaskId: []string{taskId}})
	var taskWorkersStatus []string

	// var stolyarComplete int64
	// stolyarComplete = 1

	// var malyarComplete int64
	// malyarComplete = 1

	// var montajComplete int64
	// montajComplete = 1
	fmt.Println("taskId=", result.TaskId, " len=", len(taskWorkers.Data))

	listStatusTask := map[string]domain.TaskStatus{}

	if len(taskWorkers.Data) > 0 {

		// isProcess := false
		for i := range taskWorkers.Data {
			if !utils.Contains(taskWorkersStatus, taskWorkers.Data[i].Status) && taskWorkers.Data[i].Status != "autofinish" {
				taskWorkersStatus = append(taskWorkersStatus, taskWorkers.Data[i].Status)
				listStatusTask[taskWorkers.Data[i].Status] = taskWorkers.Data[i].TaskStatus
			}
			// if taskWorkers.Data[i].Status == "process" {
			// 	isProcess = true
			// }
			// if taskWorkers.Data[i].Status != "finish" && taskWorkers.Data[i].Task.Operation.Group == "2" {
			// 	stolyarComplete = 0
			// }
			// if taskWorkers.Data[i].Status != "finish" && taskWorkers.Data[i].Task.Operation.Group == "3" {
			// 	malyarComplete = 0
			// }
			// if taskWorkers.Data[i].Status != "finish" && taskWorkers.Data[i].Task.Operation.Group == "5" {
			// 	montajComplete = 0
			// }
		}

		// dataUpdateOrder := &domain.OrderInput{}
		// // if result.Task.Operation.Group == "2" {
		// dataUpdateOrder.StolyarComplete = &stolyarComplete
		// // }
		// // if result.Task.Operation.Group == "3" {
		// dataUpdateOrder.MalyarComplete = &malyarComplete
		// // }
		// // if result.Task.Operation.Group == "5" {
		// dataUpdateOrder.MontajComplete = &montajComplete
		// // }

		// // fmt.Println("update taskWorker dataUpdateOrder: ", *dataUpdateOrder.StolyarComplete, *dataUpdateOrder.MalyarComplete)

		// _, err = s.taskService.orderService.UpdateOrder(currentTask.Data[0].OrderId.Hex(), userID, dataUpdateOrder)
		// if err != nil {
		// 	return result, err
		// }

		// taskStatus, err := s.taskStatusService.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"_id", bson.D{{"$in", taskWorkersStatus}}}}})

		// isProcess := false
		// for i := range taskStatus.Data {
		// 	if *taskStatus.Data[i].Process == 1 && result.StatusId == taskStatus.Data[i].ID {
		// 		isProcess = true
		// 	}
		// }

		// change task status.
		active := int64(1)
		fmt.Println("update taskWorker: ", len(taskWorkersStatus), taskWorkersStatus)
		// if len(taskWorkersStatus) == 1 || (len(taskWorkersStatus) > 1 && isProcess) {
		if len(taskWorkersStatus) > 0 {
			if result.Status == "finish" {
				active = int64(0)
			}

			newStatus := result.Status
			newStatusId := result.TaskStatus.ID
			isFinishCount := 0

			if val, ok := listStatusTask["finish"]; ok {
				newStatus = "finish"
				newStatusId = val.ID
				isFinishCount++
			}

			if val, ok := listStatusTask["wait"]; ok {
				newStatus = "wait"
				newStatusId = val.ID
			}
			if val, ok := listStatusTask["pause"]; ok {
				newStatus = "pause"
				newStatusId = val.ID
			}
			if val, ok := listStatusTask["process"]; ok {
				newStatus = "process"
				newStatusId = val.ID
			}

			if result.Task.Operation.Group == "5" && isFinishCount > 0 {
				if val, ok := listStatusTask["finish"]; ok {
					newStatus = "finish"
					newStatusId = val.ID
				}

				// autofinish all taskWorker if one montaj finish
				var autoFinihStatus domain.TaskStatus
				allStatus, err := s.Services.TaskStatus.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"status", "autofinish"}}})
				if err != nil {
					return result, err
				}
				if len(allStatus.Data) > 0 {
					autoFinihStatus = allStatus.Data[0]
				}

				stopStatus := []string{"finish", "process", "autofinish"}
				if !autoFinihStatus.ID.IsZero() {
					for k := range taskWorkers.Data {
						if !utils.Contains(stopStatus, taskWorkers.Data[k].Status) {
							_, err := s.UpdateTaskWorker(taskWorkers.Data[k].ID.Hex(), userID, &domain.TaskWorkerInput{StatusId: autoFinihStatus.ID, Status: autoFinihStatus.Status}, 0)
							if err != nil {
								return result, err
							}
						}
					}
				}
			}

			// task, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: result.StatusId, Status: result.Status, Active: &active})
			_, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: newStatusId, Status: newStatus, Active: &active})
			if err != nil {
				return result, err
			}

		} else {
			// если есть задания, но не нашлось совпадений(autofinish), выставляем статус wait.
			if len(taskWorkers.Data) > 0 {
				statuses, err := s.Services.TaskStatus.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"status", "wait"}}})
				if err != nil {
					return result, err
				}

				if len(statuses.Data) > 0 {
					_, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: statuses.Data[0].ID, Status: "wait", Active: &active})
					if err != nil {
						return result, err
					}
				}
			}
		}
	} else {
		statuses, err := s.Services.TaskStatus.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"status", "wait"}}})
		if err != nil {
			return result, err
		}
		if len(statuses.Data) > 0 {
			status := int64(0)
			_, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: statuses.Data[0].ID, Status: statuses.Data[0].Status, Active: &status})
			if err != nil {
				return result, err
			}
		}
	}

	// // check all tasks for change status order.
	// tasksForOrder, err := s.taskService.FindTaskPopulate(domain.RequestParams{Filter: bson.D{{"orderId", result.OrderId}}})

	// var stolyarComplete int64
	// stolyarComplete = 1

	// var malyarComplete int64
	// malyarComplete = 1

	// // var montajComplete int64
	// // montajComplete = 1

	// for i := range tasksForOrder.Data {
	// 	fmt.Println("range ", i, ":", tasksForOrder.Data[i].Status, tasksForOrder.Data[i].Operation.Group)
	// 	if tasksForOrder.Data[i].Status != "finish" && tasksForOrder.Data[i].Operation.Group == "2" {
	// 		stolyarComplete = 0
	// 	}
	// 	if tasksForOrder.Data[i].Status != "finish" && tasksForOrder.Data[i].Operation.Group == "3" {
	// 		malyarComplete = 0
	// 	}
	// 	// if tasksForOrder.Data[i].Status != "finish" && tasksForOrder.Data[i].Operation.Group == "5" {
	// 	// 	montajComplete = 0
	// 	// }
	// }

	// dataUpdateOrder := &domain.OrderInput{}
	// // if result.Task.Operation.Group == "2" {
	// dataUpdateOrder.StolyarComplete = &stolyarComplete
	// // }
	// // if result.Task.Operation.Group == "3" {
	// dataUpdateOrder.MalyarComplete = &malyarComplete
	// // }
	// // if result.Task.Operation.Group == "5" {
	// // dataUpdateOrder.MontajComplete = &montajComplete
	// // }

	// // fmt.Println("update taskWorker dataUpdateOrder: ", *dataUpdateOrder.StolyarComplete, *dataUpdateOrder.MalyarComplete)

	// _, err = s.taskService.orderService.UpdateOrder(currentTask.Data[0].OrderId.Hex(), userID, dataUpdateOrder)
	// if err != nil {
	// 	return result, err
	// }

	return result, err
}
