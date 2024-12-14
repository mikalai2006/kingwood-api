package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ObjectService struct {
	repo        repository.Object
	Hub         *Hub
	userService *UserService
}

func NewObjectService(repo repository.Object, hub *Hub, userService *UserService) *ObjectService {
	return &ObjectService{repo: repo, Hub: hub, userService: userService}
}

func (s *ObjectService) FindObject(input *domain.ObjectFilter) (domain.Response[domain.Object], error) {
	return s.repo.FindObject(input)
}

func (s *ObjectService) CreateObject(userID string, data *domain.Object) (*domain.Object, error) {
	var result *domain.Object

	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	result, err = s.repo.CreateObject(userID, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *ObjectService) UpdateObject(id string, userID string, data *domain.ObjectInput) (*domain.Object, error) {
	result, err := s.repo.UpdateObject(id, userID, data)
	if err != nil {
		return result, err
	}

	// ObjectStatus, err := s.ObjectStatus.FindObjectStatus(domain.RequestParams{Filter: bson.D{{"_id", result.StatusId}}})
	// if err != nil {
	// 	return result, err
	// }

	// if ObjectStatus.Data[0].Finish != nil {
	// 	if *ObjectStatus.Data[0].Finish == 1 {
	// 		allObjectsByOrder, err := s.FindObjectWithWorkers(domain.RequestParams{Filter: bson.D{{"orderId", result.OrderId}}})
	// 		if err != nil {
	// 			return result, err
	// 		}
	// 		sort.Slice(allObjectsByOrder.Data, func(i, j int) bool {
	// 			return *allObjectsByOrder.Data[i].SortOder < *allObjectsByOrder.Data[j].SortOder
	// 		})

	// 		nextObject := allObjectsByOrder.Data[*result.SortOder+1]

	// 		// ObjectWithWorkers, err := s.FindObjectWithWorkers(domain.RequestParams{Filter: bson.D{{"_id", result.ID}}})
	// 		// if err != nil {
	// 		// 	return result, err
	// 		// }
	// 		// fmt.Println("nextObject.Workers=", len(nextObject.Workers))

	// 		if nextObject.SortOder != nil {
	// 			statusActive := int64(1)

	// 			nextObjectUpdated, err := s.repo.UpdateObject(nextObject.ID.Hex(), userID, &domain.ObjectInput{
	// 				Active: &statusActive,
	// 			})

	// 			if err != nil {
	// 				return result, err
	// 			}
	// 			s.Hub.HandleMessage(domain.Message{Type: "message", Method: "ADD", Sender: userID, Recipient: "sobesednikID.Hex()", Content: nextObjectUpdated, ID: "room1", Service: "Object"})

	// 			statusDisable := int64(0)

	// 			result, err = s.repo.UpdateObject(id, userID, &domain.ObjectInput{Active: &statusDisable})

	// 			if err != nil {
	// 				return result, err
	// 			}
	// 		}
	// 	}
	// }

	return result, err
}

func (s *ObjectService) DeleteObject(id string) (*domain.Object, error) {
	result, err := s.repo.DeleteObject(id)

	return result, err
}
