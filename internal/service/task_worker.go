package service

import (
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
}

func NewTaskWorkerService(repo repository.TaskWorker, userService *UserService, taskStatusService *TaskStatusService, taskService *TaskService, hub *Hub) *TaskWorkerService {
	return &TaskWorkerService{repo: repo, userService: userService, taskStatusService: taskStatusService, taskService: taskService, Hub: hub}
}

func (s *TaskWorkerService) FindTaskWorkerPopulate(params domain.RequestParams) (domain.Response[domain.TaskWorker], error) {
	return s.repo.FindTaskWorkerPopulate(params)
}

func (s *TaskWorkerService) FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error) {
	return s.repo.FindTaskWorker(params)
}

func (s *TaskWorkerService) CreateTaskWorker(userID string, data *domain.TaskWorker) (*domain.TaskWorker, error) {
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

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "ADD", Sender: userID, Recipient: "sobesednikID.Hex()", Content: result, ID: "room1", Service: "taskWorker"})

	return result, err
}

func (s *TaskWorkerService) UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput) (*domain.TaskWorker, error) {
	result, err := s.repo.UpdateTaskWorker(id, userID, data)
	if err != nil {
		return result, err
	}

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: userID, Recipient: "sobesednikID.Hex()", Content: result, ID: "room1", Service: "taskWorker"})

	// get all taskWorkers.
	taskWorkers, err := s.FindTaskWorker(domain.RequestParams{Filter: bson.D{{"taskId", result.TaskId}}})

	var taskWorkersStatus []string

	isProcess := false
	for i := range taskWorkers.Data {
		if !utils.Contains(taskWorkersStatus, taskWorkers.Data[i].Status) {
			taskWorkersStatus = append(taskWorkersStatus, taskWorkers.Data[i].Status)
		}
		if taskWorkers.Data[i].Status == "process" {
			isProcess = true
		}
	}

	// taskStatus, err := s.taskStatusService.FindTaskStatus(domain.RequestParams{Filter: bson.D{{"_id", bson.D{{"$in", taskWorkersStatus}}}}})

	// isProcess := false
	// for i := range taskStatus.Data {
	// 	if *taskStatus.Data[i].Process == 1 && result.StatusId == taskStatus.Data[i].ID {
	// 		isProcess = true
	// 	}
	// }

	// if one worker, change task status.
	fmt.Println("update taskWorker: ", len(taskWorkersStatus), taskWorkersStatus, isProcess)
	if len(taskWorkersStatus) == 1 || (len(taskWorkersStatus) > 1 && isProcess) {
		active := int64(1)
		if result.Status == "finish" {
			active = int64(0)
		}

		task, err := s.taskService.UpdateTask(result.TaskId.Hex(), userID, &domain.TaskInput{StatusId: result.StatusId, Status: result.Status, Active: &active})

		if err != nil {
			return result, err
		}

		s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: userID, Recipient: "sobesednikID.Hex()", Content: task, ID: "room1", Service: "task"})
	}

	return result, err
}

func (s *TaskWorkerService) DeleteTaskWorker(id string) (*domain.TaskWorker, error) {
	result, err := s.repo.DeleteTaskWorker(id)

	return result, err
}
