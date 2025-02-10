package service

import (
	"fmt"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo         repository.Task
	Hub          *Hub
	userService  *UserService
	taskStatus   *TaskStatusService
	orderService *OrderService
	Services     *Services
}

func NewTaskService(repo repository.Task, hub *Hub, userService *UserService, TaskStatus *TaskStatusService, OrderService *OrderService) *TaskService {
	return &TaskService{repo: repo, Hub: hub, userService: userService, taskStatus: TaskStatus, orderService: OrderService}
}

func (s *TaskService) FindTask(params domain.RequestParams) (domain.Response[domain.Task], error) {
	return s.repo.FindTask(params)
}

func (s *TaskService) FindTaskPopulate(filter domain.TaskFilter) (domain.Response[domain.Task], error) {
	return s.repo.FindTaskPopulate(filter)
}

func (s *TaskService) CreateTask(userID string, data *domain.Task) (*domain.Task, error) {
	var result *domain.Task

	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// set default status.
	if data.StatusId.IsZero() {
		allStatus, err := s.Services.TaskStatus.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"status", "wait"}}})
		if err != nil {
			return result, err
		}
		if len(allStatus.Data) > 0 {
			data.StatusId = allStatus.Data[0].ID
			data.Status = allStatus.Data[0].Status
		}
	}

	result, err = s.repo.CreateTask(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	_, err = s.CheckStatusOrder(userID, result)
	if err != nil {
		return result, err
	}

	// add taskWorker for all users by work on the object.
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
		date := time.Now().Local()
		allTaskWorkerForObject, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{
			ObjectId:    []string{objectId},
			OperationId: []string{operationId},
			Date:        &date,
		})
		if err != nil {
			return result, err
		}
		// fmt.Println("allTaskWorkerForObject length: ", len(allTaskWorkerForObject.Data))
		createdWorkers := []string{}
		if len(allTaskWorkerForObject.Data) > 0 {
			for i := range allTaskWorkerForObject.Data {
				if !utils.Contains(createdWorkers, allTaskWorkerForObject.Data[i].WorkerId.Hex()) {
					newTaskWorker := domain.TaskWorker{
						ObjectId:    allTaskWorkerForObject.Data[i].ObjectId,
						OrderId:     result.OrderId,
						TaskId:      result.ID,
						OperationId: result.OperationId,
						WorkerId:    allTaskWorkerForObject.Data[i].WorkerId,
						SortOrder:   allTaskWorkerForObject.Data[i].SortOrder,
						StatusId:    result.StatusId,
						Status:      result.Status,
						From:        allTaskWorkerForObject.Data[i].From,
						To:          allTaskWorkerForObject.Data[i].To,
						TypeGo:      allTaskWorkerForObject.Data[i].TypeGo,
					}
					insertTaskWorker, err := s.Services.TaskWorker.CreateTaskWorker(userID, &newTaskWorker, 0)
					if err != nil {
						return result, err
					}
					createdWorkers = append(createdWorkers, insertTaskWorker.WorkerId.Hex())
				}
			}
		}
	}

	return result, err
}

func (s *TaskService) UpdateTask(id string, userID string, data *domain.TaskInput) (*domain.Task, error) {
	result, err := s.repo.UpdateTask(id, userID, data)
	if err != nil {
		return result, err
	}

	// taskStatus, err := s.taskStatus.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"_id", result.StatusId}}})
	// if err != nil {
	// 	return result, err
	// }

	// if taskStatus.Data[0].Finish != nil {
	// if result.Status == "finish" {
	// 	// if *taskStatus.Data[0].Finish == 1 {
	// 	allTasksByOrder, err := s.FindTaskWithWorkers(domain.RequestParams{Filter: bson.D{{"orderId", result.OrderId}}})
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	sort.Slice(allTasksByOrder.Data, func(i, j int) bool {
	// 		return *allTasksByOrder.Data[i].SortOrder < *allTasksByOrder.Data[j].SortOrder
	// 	})

	// 	nextIndex := *result.SortOrder + 1
	// 	if int64(len(allTasksByOrder.Data)) <= nextIndex {
	// 		nextIndex = -1
	// 	}
	// 	fmt.Println("next task=", nextIndex, " length task=", len(allTasksByOrder.Data))

	// 	if nextIndex >= 0 {
	// 		nextTask := allTasksByOrder.Data[nextIndex]

	// 		// taskWithWorkers, err := s.FindTaskWithWorkers(domain.RequestParams{Filter: bson.D{{"_id", result.ID}}})
	// 		// if err != nil {
	// 		// 	return result, err
	// 		// }
	// 		// fmt.Println("nextTask.Workers=", len(nextTask.Workers))

	// 		if nextTask.SortOrder != nil && *nextTask.Active == 1 {
	// 			statusActive := int64(1)

	// 			nextTaskUpdated, err := s.repo.UpdateTask(nextTask.ID.Hex(), userID, &domain.TaskInput{
	// 				Active: &statusActive,
	// 			})

	// 			if err != nil {
	// 				return result, err
	// 			}
	// 			s.Hub.HandleMessage(domain.Message{Type: "message", Method: "ADD", Sender: userID, Recipient: "sobesednikID.Hex()", Content: nextTaskUpdated, ID: "room1", Service: "task"})

	// 			statusDisable := int64(0)

	// 			result, err = s.repo.UpdateTask(id, userID, &domain.TaskInput{Active: &statusDisable})

	// 			if err != nil {
	// 				return result, err
	// 			}
	// 		}
	// 	} else {
	// 		// получаем заказ для задания
	// 		idOrder := result.OrderId.Hex()
	// 		ordersByTask, err := s.orderService.FindOrder(&domain.OrderFilter{ID: []*string{&idOrder}})
	// 		if err != nil {
	// 			return result, err
	// 		}
	// 		// если стадия цеховые работы, меняем на монтажные работы
	// 		// status := int64(0)
	// 		if utils.Contains(ordersByTask.Data[0].Group, "create") {
	// 			order, err := s.orderService.UpdateOrder(result.OrderId.Hex(), userID, &domain.OrderInput{
	// 				Group: []string{"create_complete"},
	// 				// Status: &status,
	// 			})
	// 			if err != nil {
	// 				return result, err
	// 			}
	// 			s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: userID, Recipient: "sobesednikID.Hex()", Content: order, ID: "room1", Service: "order"})
	// 		}
	// 	}
	// 	// }
	// }

	// // To check the token is valid
	// pushToken, err := expo.NewExponentPushToken("ExponentPushToken[uia35pA2ijvbzgRPxnW50M]")
	// if err != nil {
	// 	panic(err)
	// }

	// // Create a new Expo SDK client
	// client := expo.NewPushClient(nil)

	// // Publish message
	// response, err := client.Publish(
	// 	&expo.PushMessage{
	// 		To:       []expo.ExponentPushToken{pushToken},
	// 		Body:     fmt.Sprintf("Статус задачи %s изменен на %s", result.Name, result.Status),
	// 		Data:     map[string]string{"withSome": "data"},
	// 		Sound:    "default",
	// 		Title:    "Изменение статуса задачи",
	// 		Priority: expo.DefaultPriority,
	// 	},
	// )

	// // Check errors
	// if err != nil {
	// 	panic(err)
	// }

	// // Validate responses
	// if response.ValidateResponse() != nil {
	// 	fmt.Println(response.PushMessage.To, "failed")
	// }

	_, err = s.CheckStatusOrder(userID, result)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "PATCH", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "task"})

	return result, err
}

func (s *TaskService) DeleteTask(id string, userID string, checkStatus bool) (*domain.Task, error) {
	result, err := s.repo.DeleteTask(id)

	if checkStatus {
		_, err = s.CheckStatusOrder("userID", result)
		if err != nil {
			return result, err
		}
	}

	s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: "userID", Recipient: "", Content: result, ID: "room1", Service: "task"})

	return result, err
}

func (s *TaskService) CheckStatusOrder(userID string, result *domain.Task) (*domain.Task, error) {
	// check all tasks for change status order.
	tasksForOrder, err := s.FindTaskPopulate(domain.TaskFilter{OrderId: []string{result.OrderId.Hex()}})

	type CheckedStruct struct {
		Status      int64
		CountFinish int
		CountAll    int
	}

	goComplete := CheckedStruct{
		Status:      0,
		CountFinish: 0,
		CountAll:    0,
	}
	stolyarComplete := CheckedStruct{
		Status:      0,
		CountFinish: 0,
		CountAll:    0,
	}
	malyarComplete := CheckedStruct{
		Status:      0,
		CountFinish: 0,
		CountAll:    0,
	}
	shlifComplete := CheckedStruct{
		Status:      0,
		CountFinish: 0,
		CountAll:    0,
	}
	montajComplete := CheckedStruct{
		Status:      0,
		CountFinish: 0,
		CountAll:    0,
	}
	// var stolyarComplete int64
	// stolyarComplete = 1

	// var malyarComplete int64
	// malyarComplete = 1

	// var montajComplete int64
	// montajComplete = 1
	allTasksStatus := []string{}

	for i := range tasksForOrder.Data {
		fmt.Println("range ", i, ":", tasksForOrder.Data[i].Status, tasksForOrder.Data[i].Operation.Group)
		if tasksForOrder.Data[i].Operation.Group == "2" {
			stolyarComplete.CountAll = stolyarComplete.CountAll + 1
			if utils.Contains([]string{"finish", "autofinish"}, tasksForOrder.Data[i].Status) {
				stolyarComplete.CountFinish += 1
			}
			// stolyarComplete.Status = 0
		}
		if tasksForOrder.Data[i].Operation.Group == "3" {
			malyarComplete.CountAll = malyarComplete.CountAll + 1
			if utils.Contains([]string{"finish", "autofinish"}, tasksForOrder.Data[i].Status) {
				malyarComplete.CountFinish += 1
			}
			// malyarComplete.Status = 0
		}
		if tasksForOrder.Data[i].Operation.Group == "6" {
			shlifComplete.CountAll = shlifComplete.CountAll + 1
			if utils.Contains([]string{"finish", "autofinish"}, tasksForOrder.Data[i].Status) {
				shlifComplete.CountFinish += 1
			}
			// shlifComplete.Status = 0
		}
		if tasksForOrder.Data[i].Operation.Group == "4" {
			goComplete.CountAll = goComplete.CountAll + 1
			if utils.Contains([]string{"finish", "autofinish"}, tasksForOrder.Data[i].Status) {
				goComplete.CountFinish += 1
			}
			// goComplete.Status = 0
		}
		if tasksForOrder.Data[i].Operation.Group == "5" {
			montajComplete.CountAll = montajComplete.CountAll + 1
			if utils.Contains([]string{"finish", "autofinish"}, tasksForOrder.Data[i].Status) {
				montajComplete.CountFinish += 1
			}
			// montajComplete.Status = 0
		}

		if !utils.Contains(allTasksStatus, tasksForOrder.Data[i].Status) {
			allTasksStatus = append(allTasksStatus, tasksForOrder.Data[i].Status)
		}
	}

	if stolyarComplete.CountAll == stolyarComplete.CountFinish && stolyarComplete.CountAll > 0 {
		stolyarComplete.Status = 1
	}
	if malyarComplete.CountAll == malyarComplete.CountFinish && malyarComplete.CountAll > 0 { // || malyarComplete.CountAll == 0
		malyarComplete.Status = 1
	}
	if goComplete.CountAll == goComplete.CountFinish && goComplete.CountAll > 0 {
		goComplete.Status = 1
	}
	if montajComplete.CountFinish > 0 && montajComplete.CountAll > 0 { //montajComplete.CountAll == montajComplete.CountFinish
		montajComplete.Status = 1
	}
	if shlifComplete.CountAll == shlifComplete.CountFinish && shlifComplete.CountAll > 0 { // || shlifComplete.CountAll == 0
		shlifComplete.Status = 1
	}
	// если нет задания упаковки.
	if goComplete.CountAll == 0 {
		if stolyarComplete.Status == 1 && (malyarComplete.Status == 1 || malyarComplete.CountAll == 0) && (shlifComplete.Status == 1 || shlifComplete.CountAll == 0) {
			goComplete.Status = 1
		}
	}

	// fmt.Println(stolyarComplete, malyarComplete, montajComplete)

	dataUpdateOrder := &domain.OrderInput{}
	// if result.Task.Operation.Group == "2" {
	dataUpdateOrder.StolyarComplete = &stolyarComplete.Status
	// }
	// if result.Task.Operation.Group == "3" {
	dataUpdateOrder.MalyarComplete = &malyarComplete.Status
	// }
	// if result.Task.Operation.Group == "5" {
	dataUpdateOrder.MontajComplete = &montajComplete.Status
	// }
	dataUpdateOrder.GoComplete = &goComplete.Status

	dataUpdateOrder.ShlifComplete = &shlifComplete.Status

	if len(allTasksStatus) > 0 {
		// если есть задания, меняем статус заказа на 1
		status := int64(1)
		dataUpdateOrder.Status = &status
	} else if len(allTasksStatus) == 0 {
		// если нет заданий, меняем статус заказа на 0
		status := int64(0)
		dataUpdateOrder.Status = &status
	}
	// если у всех заданий статус - finish, помечаем заказ как выполненный
	if len(allTasksStatus) == 1 {
		if allTasksStatus[0] == "finish" {
			status := int64(100)
			dataUpdateOrder.Status = &status
		}
	}
	fmt.Println("update taskWorker dataUpdateOrder: ", dataUpdateOrder.Status, allTasksStatus, len(allTasksStatus))

	_, err = s.orderService.UpdateOrder(result.OrderId.Hex(), userID, dataUpdateOrder)
	if err != nil {
		return result, err
	}

	return result, err
}
