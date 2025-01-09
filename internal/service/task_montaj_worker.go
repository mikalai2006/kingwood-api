package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskMontajWorkerService struct {
	repo              repository.TaskMontajWorker
	userService       *UserService
	taskStatusService *TaskStatusService
	taskService       *TaskService
	Hub               *Hub
}

func NewTaskMontajWorkerService(repo repository.TaskMontajWorker, userService *UserService, taskStatusService *TaskStatusService, taskService *TaskService, hub *Hub) *TaskMontajWorkerService {
	return &TaskMontajWorkerService{repo: repo, userService: userService, taskStatusService: taskStatusService, taskService: taskService, Hub: hub}
}

func (s *TaskMontajWorkerService) FindTaskMontajWorkerPopulate(input *domain.TaskMontajWorkerFilter) (domain.Response[domain.TaskMontajWorker], error) {
	return s.repo.FindTaskMontajWorkerPopulate(input)
}

func (s *TaskMontajWorkerService) FindTaskMontajWorker(params domain.RequestParams) (domain.Response[domain.TaskMontajWorker], error) {
	return s.repo.FindTaskMontajWorker(params)
}

func (s *TaskMontajWorkerService) CreateTaskMontajWorker(userID string, data *domain.TaskMontajWorker) (*domain.TaskMontajWorker, error) {
	var result *domain.TaskMontajWorker

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
	// 	updateReview := &domain.TaskMontajWorkerInput{
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

	result, err = s.repo.CreateTaskMontajWorker(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "ADD", Sender: userID, Recipient: "sobesednikID.Hex()", Content: result, ID: "room1", Service: "TaskMontajWorker"})

	return result, err
}

func (s *TaskMontajWorkerService) UpdateTaskMontajWorker(id string, userID string, data *domain.TaskMontajWorkerInput) (*domain.TaskMontajWorker, error) {
	result, err := s.repo.UpdateTaskMontajWorker(id, userID, data)
	if err != nil {
		return result, err
	}

	// currentTask, err := s.taskService.FindTask(domain.RequestParams{Filter: bson.D{{"_id", result.TaskId}}})
	// if err != nil {
	// 	return result, err
	// }

	// s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: userID, Recipient: "sobesednikID.Hex()", Content: result, ID: "room1", Service: "TaskMontajWorker"})

	// // get all TaskMontajWorkers.
	// TaskMontajWorkers, err := s.FindTaskMontajWorkerPopulate(domain.RequestParams{Filter: bson.D{{"taskId", result.TaskId}}})

	// var TaskMontajWorkersStatus []string

	// var stolyarComplete int64
	// stolyarComplete = 1

	// var malyarComplete int64
	// malyarComplete = 1

	// var montajComplete int64
	// montajComplete = 1

	// isProcess := false
	// for i := range TaskMontajWorkers.Data {
	// 	if !utils.Contains(TaskMontajWorkersStatus, TaskMontajWorkers.Data[i].Status) {
	// 		TaskMontajWorkersStatus = append(TaskMontajWorkersStatus, TaskMontajWorkers.Data[i].Status)
	// 	}
	// 	if TaskMontajWorkers.Data[i].Status == "process" {
	// 		isProcess = true
	// 	}
	// 	if TaskMontajWorkers.Data[i].Status != "finish" && TaskMontajWorkers.Data[i].Task.Operation.Group == "2" {
	// 		stolyarComplete = 0
	// 	}
	// 	if TaskMontajWorkers.Data[i].Status != "finish" && TaskMontajWorkers.Data[i].Task.Operation.Group == "3" {
	// 		malyarComplete = 0
	// 	}
	// 	if TaskMontajWorkers.Data[i].Status != "finish" && TaskMontajWorkers.Data[i].Task.Operation.Group == "5" {
	// 		montajComplete = 0
	// 	}
	// }

	// dataUpdateOrder := &domain.OrderInput{}
	// if result.Task.Operation.Group == "2" {
	// 	dataUpdateOrder.StolyarComplete = &stolyarComplete
	// }
	// if result.Task.Operation.Group == "3" {
	// 	dataUpdateOrder.MalyarComplete = &malyarComplete
	// }
	// if result.Task.Operation.Group == "5" {
	// 	dataUpdateOrder.MontajComplete = &montajComplete
	// }

	// // fmt.Println("update TaskMontajWorker dataUpdateOrder: ", *dataUpdateOrder.StolyarComplete, result.Task.Operation.Group)

	// _, err = s.taskService.orderService.UpdateOrder(currentTask.Data[0].OrderId.Hex(), userID, dataUpdateOrder)
	// if err != nil {
	// 	return result, err
	// }

	// // taskStatus, err := s.taskStatusService.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"_id", bson.D{{"$in", TaskMontajWorkersStatus}}}}})

	// // isProcess := false
	// // for i := range taskStatus.Data {
	// // 	if *taskStatus.Data[i].Process == 1 && result.StatusId == taskStatus.Data[i].ID {
	// // 		isProcess = true
	// // 	}
	// // }

	// // if one worker, change task status.
	// fmt.Println("update TaskMontajWorker: ", len(TaskMontajWorkersStatus), TaskMontajWorkersStatus, isProcess, stolyarComplete, malyarComplete)
	// if len(TaskMontajWorkersStatus) == 1 || (len(TaskMontajWorkersStatus) > 1 && isProcess) {
	// 	active := int64(1)
	// 	if result.Status == "finish" {
	// 		active = int64(0)
	// 	}

	// 	task, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: result.StatusId, Status: result.Status, Active: &active})

	// 	if err != nil {
	// 		return result, err
	// 	}

	// 	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: userID, Recipient: "sobesednikID.Hex()", Content: task, ID: "room1", Service: "task"})
	// }

	return result, err
}

func (s *TaskMontajWorkerService) DeleteTaskMontajWorker(id string) (*domain.TaskMontajWorker, error) {
	result, err := s.repo.DeleteTaskMontajWorker(id)

	return result, err
}
