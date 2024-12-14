package service

import (
	"fmt"

	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OfferService struct {
	repo               repository.Offer
	userService        *UserService
	messageService     *MessageService
	messageRoomService *MessageRoomService
	Hub                *Hub
}

func NewOfferService(repo repository.Offer, userService *UserService, hub *Hub, messageService *MessageService, messageRoomService *MessageRoomService) *OfferService {
	return &OfferService{repo: repo, userService: userService, Hub: hub, messageService: messageService, messageRoomService: messageRoomService}
}

func (s *OfferService) FindOffer(params *model.OfferFilter) (domain.Response[model.Offer], error) {
	return s.repo.FindOffer(params)
}

func (s *OfferService) CreateOffer(userID string, data *model.OfferInput) (*model.Offer, error) {
	result, err := s.repo.CreateOffer(userID, data)

	s.Hub.HandleMessage(domain.Message{Type: "message", Sender: userID, Recipient: "user2", Content: result, ID: "room1", Service: "offer"})

	return result, err
}

func (s *OfferService) UpdateOffer(id string, userID string, data *model.Offer) (*model.Offer, error) {
	result, err := s.repo.UpdateOffer(id, userID, data)
	if err != nil {
		return result, err
	}

	if data.Win != nil {
		fmt.Println("data: ", data, *data.Win)
		if *data.Win == 1 {
			messageRoom, err := s.messageRoomService.CreateMessageRoom(userID, &model.MessageRoom{ProductID: result.ProductID, OfferID: result.ID, TakeUserID: result.UserID})
			if err != nil {
				return result, err
			}
			s.messageService.CreateMessage(userID, &model.MessageInput{Message: "[start]", RoomID: messageRoom.ID.Hex()})

			data.RoomId = messageRoom.ID

		} else if *data.Win == 0 {

			// _, err := s.messageRoomService.DeleteMessageRoom(result.RoomId.Hex())
			userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				return result, err
			}
			data.RejectUserId = userIDPrimitive

			// if userIDPrimitive == result.UserProductID {
			// 	s.userService.SetBal(result.UserID.Hex(), int(result.Cost))
			// } else if userIDPrimitive == result.UserID {
			// 	s.userService.SetBal(result.UserID.Hex(), int(result.Cost-1))
			// }
			// fmt.Println("userIDPrimitive: ", userIDPrimitive, result.UserProductID, result.UserID)

			status := -1
			data.Status = -1
			_, err = s.messageRoomService.UpdateMessageRoom(result.RoomId.Hex(), userID, &model.MessageRoom{Status: &status})

			if err != nil {
				return result, err
			}
		}
	}

	status := 1
	if *result.Take == status { // && *result.Give == status
		data.Status = -1

		// _, err = s.userService.SetBal(result.UserProductID.Hex(), int(result.Cost))
		// if err != nil {
		// 	return result, err
		// }

		statusX := -1
		// fmt.Println("result.RoomId=", result.RoomId)
		_, err = s.messageRoomService.UpdateMessageRoom(result.RoomId.Hex(), userID, &model.MessageRoom{Status: &statusX})
		if err != nil {
			return result, err
		}

		// find all offers and return bals.
		offers, err := s.FindOffer(&model.OfferFilter{ProductID: []*primitive.ObjectID{&result.ProductID}})
		if err != nil {
			return result, err
		}
		// fmt.Println("offers=", len(offers.Data))
		statusWin := 0
		for i := range offers.Data {
			if offers.Data[i].ID.Hex() != id {
				s.UpdateOffer(offers.Data[i].ID.Hex(), offers.Data[i].UserProductID.Hex(), &model.Offer{Status: -1, Win: &statusWin})
				// s.userService.SetBal(offers.Data[i].UserID.Hex(), int(offers.Data[i].Cost))
			}
		}
	}

	result, err = s.repo.UpdateOffer(id, userID, data)

	s.Hub.HandleMessage(domain.Message{Type: "message", Sender: userID, Recipient: "user2", Content: result, ID: "room1", Service: "offer"})

	return result, err
}

func (s *OfferService) DeleteOffer(id string) (model.Offer, error) {
	result := model.Offer{}

	// removedOffer, err := s.repo.GetOffer(id)
	// if err != nil {
	// 	return result, err
	// }

	// // find all Offer_vote for remove.
	// OfferVotes, err := s.OfferVoteService.FindOfferVote(domain.RequestParams{
	// 	Filter:  bson.M{"Offer_id": removedOffer.ID},
	// 	Options: domain.Options{Limit: 1000},
	// })
	// if err != nil {
	// 	return result, err
	// }

	// for i, _ := range OfferVotes.Data {
	// 	_, err := s.OfferVoteService.DeleteOfferVote(OfferVotes.Data[i].ID.Hex())
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	// fmt.Println("Remove Offer: ", Offer.Data[i].ID.Hex())
	// }

	result, err := s.repo.DeleteOffer(id)

	// // set user stat
	// if err == nil {
	// 	_, _ = s.userService.SetStat(result.UserID.Hex(), domain.UserStat{AddOffer: -1})
	// }

	return result, err
}
