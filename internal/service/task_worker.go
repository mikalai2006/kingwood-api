package service

import (
	"errors"
	"fmt"

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
	return s.repo.FindTaskWorkerPopulate(input)
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
	// existReview, err := s.repo.FindReview(domain.RequestParams{
	// 	Filter:  bson.M{"node_id": review.NodeID, "user_id": userIDPrimitive},
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

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "ADD", Sender: userID, Recipient: "sobesednikID.Hex()", Content: result, ID: "room1", Service: "taskWorker"})

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
			objectId := result.ObjectId.Hex()
			operationId := result.OperationId.Hex()
			allTaskForObject, err := s.Services.Task.FindTaskPopulate(domain.TaskFilter{ObjectId: []*string{&objectId}, OperationId: []*string{&operationId}})
			if err != nil {
				return result, err
			}
			fmt.Println("allTaskForObject length: ", len(allTaskForObject.Data))
			if len(allTaskForObject.Data) > 0 {
				for i := range allTaskForObject.Data {
					workerIds := []string{}
					for i := range allTaskForObject.Data[i].Workers {
						workerIds = append(workerIds, allTaskForObject.Data[i].Workers[i].WorkerId.Hex())
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

	return result, err
}

func (s *TaskWorkerService) UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput, autoUpdate int) (*domain.TaskWorker, error) {
	result, err := s.repo.UpdateTaskWorker(id, userID, data)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: userID, Recipient: "sobesednikID.Hex()", Content: result, ID: "room1", Service: "taskWorker"})

	_, err = s.CheckStatusTask(userID, result)
	if err != nil {
		return result, err
	}

	// change taskWorker for all task on the object for updated worker (montaj).
	if autoUpdate > 0 {
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
			allTaskWorkerForObject, err := s.Services.TaskWorker.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{ObjectId: []*string{&objectId}, WorkerId: []*string{&workerId}, OperationId: []*string{&operationId}})
			if err != nil {
				return result, err
			}
			// fmt.Println("allTaskForObject length: ", len(allTaskForObject.Data))
			if len(allTaskWorkerForObject.Data) > 0 {
				// statusNotChange := []string{"finish",""}
				for i := range allTaskWorkerForObject.Data {
					if allTaskWorkerForObject.Data[i].ID.Hex() != result.ID.Hex() {
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

	return result, err
}

func (s *TaskWorkerService) DeleteTaskWorker(id string) (*domain.TaskWorker, error) {
	result, err := s.repo.DeleteTaskWorker(id)
	if err != nil {
		return result, err
	}

	_, err = s.CheckStatusTask("userID", result)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "DELETE", Sender: "userID", Recipient: "sobesednikID.Hex()", Content: result, ID: "room1", Service: "taskWorker"})

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
	fmt.Println("taskId: ", result.TaskId)
	taskId := result.TaskId.Hex()
	taskWorkers, err := s.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{TaskId: []*string{&taskId}})
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
		fmt.Println("update taskWorker: ", len(taskWorkersStatus), taskWorkersStatus)
		// if len(taskWorkersStatus) == 1 || (len(taskWorkersStatus) > 1 && isProcess) {
		if len(taskWorkersStatus) > 0 {
			active := int64(1)
			if result.Status == "finish" {
				active = int64(0)
			}

			newStatus := result.Status
			newStatusId := result.TaskStatus.ID

			if val, ok := listStatusTask["finish"]; ok {
				newStatus = "finish"
				newStatusId = val.ID
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

			// task, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: result.StatusId, Status: result.Status, Active: &active})
			_, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: newStatusId, Status: newStatus, Active: &active})
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
