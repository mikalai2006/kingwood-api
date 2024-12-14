package service

import (
	"fmt"
	"sort"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo        repository.Task
	Hub         *Hub
	userService *UserService
	taskStatus  *TaskStatusService
}

func NewTaskService(repo repository.Task, hub *Hub, userService *UserService, TaskStatus *TaskStatusService) *TaskService {
	return &TaskService{repo: repo, Hub: hub, userService: userService, taskStatus: TaskStatus}
}

func (s *TaskService) FindTask(params domain.RequestParams) (domain.Response[domain.Task], error) {
	return s.repo.FindTask(params)
}

func (s *TaskService) FindTaskWithWorkers(params domain.RequestParams) (domain.Response[domain.Task], error) {
	return s.repo.FindTaskWithWorkers(params)
}

func (s *TaskService) CreateTask(userID string, data *domain.Task) (*domain.Task, error) {
	var result *domain.Task

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
	// 	updateReview := &domain.TaskInput{
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

	result, err = s.repo.CreateTask(userID, data)
	if err != nil {
		return nil, err
	}
	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(userID, domain.UserStat{AddReview: 1})
	// }

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
	if result.Status == "finish" {
		// if *taskStatus.Data[0].Finish == 1 {
		allTasksByOrder, err := s.FindTaskWithWorkers(domain.RequestParams{Filter: bson.D{{"orderId", result.OrderId}}})
		if err != nil {
			return result, err
		}
		sort.Slice(allTasksByOrder.Data, func(i, j int) bool {
			return *allTasksByOrder.Data[i].SortOder < *allTasksByOrder.Data[j].SortOder
		})

		nextIndex := *result.SortOder + 1
		if int64(len(allTasksByOrder.Data)) <= nextIndex {
			nextIndex = -1
		}
		// fmt.Println("next task=", nextIndex, " length task=", len(allTasksByOrder.Data))

		if nextIndex >= 0 {
			nextTask := allTasksByOrder.Data[nextIndex]

			// taskWithWorkers, err := s.FindTaskWithWorkers(domain.RequestParams{Filter: bson.D{{"_id", result.ID}}})
			// if err != nil {
			// 	return result, err
			// }
			// fmt.Println("nextTask.Workers=", len(nextTask.Workers))

			if nextTask.SortOder != nil && *nextTask.Active == 1 {
				statusActive := int64(1)

				nextTaskUpdated, err := s.repo.UpdateTask(nextTask.ID.Hex(), userID, &domain.TaskInput{
					Active: &statusActive,
				})

				if err != nil {
					return result, err
				}
				s.Hub.HandleMessage(domain.Message{Type: "message", Method: "ADD", Sender: userID, Recipient: "sobesednikID.Hex()", Content: nextTaskUpdated, ID: "room1", Service: "task"})

				statusDisable := int64(0)

				result, err = s.repo.UpdateTask(id, userID, &domain.TaskInput{Active: &statusDisable})

				if err != nil {
					return result, err
				}
			}
		}
		// }
	}

	// To check the token is valid
	pushToken, err := expo.NewExponentPushToken("ExponentPushToken[uia35pA2ijvbzgRPxnW50M]")
	if err != nil {
		panic(err)
	}

	// Create a new Expo SDK client
	client := expo.NewPushClient(nil)

	// Publish message
	response, err := client.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{pushToken},
			Body:     fmt.Sprintf("Статус задачи %s изменен на %s", result.Name, result.Status),
			Data:     map[string]string{"withSome": "data"},
			Sound:    "default",
			Title:    "Изменение статуса задачи",
			Priority: expo.DefaultPriority,
		},
	)

	// Check errors
	if err != nil {
		panic(err)
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		fmt.Println(response.PushMessage.To, "failed")
	}

	return result, err
}

func (s *TaskService) DeleteTask(id string) (*domain.Task, error) {
	result, err := s.repo.DeleteTask(id)

	return result, err
}
