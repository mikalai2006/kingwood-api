package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRoomMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewMessageRoomMongo(db *mongo.Database, i18n config.I18nConfig) *MessageRoomMongo {
	return &MessageRoomMongo{db: db, i18n: i18n}
}

func (r *MessageRoomMongo) FindMessageRoom(params *model.MessageRoomFilter) (domain.Response[model.MessageRoom], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.MessageRoom
	var response domain.Response[model.MessageRoom]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Node]{}, err
	// }
	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[model.Message]{}, err
	// }
	// fmt.Println(params)
	q := bson.D{}
	if params.UserID != nil && !params.UserID.IsZero() {
		// userIDPrimitive, err := primitive.ObjectIDFromHex(*params.UserID)
		// if err != nil {
		// 	return response, err
		// }
		// q = append(q, bson.E{"userId", params.UserID})
		queryArr := []bson.M{}
		queryArr = append(queryArr, bson.M{"userId": params.UserID})
		queryArr = append(queryArr, bson.M{"takeUserId": params.UserID})
		q = append(q, bson.E{"$or", queryArr})
	}
	if params.ID != nil && !params.ID.IsZero() {
		// userIDPrimitive, err := primitive.ObjectIDFromHex(*params.ID)
		// if err != nil {
		// 	return response, err
		// }
		q = append(q, bson.E{"_id", params.ID})
	}
	if params.TakeUserID != nil && !params.TakeUserID.IsZero() {
		// userProductIDPrimitive, err := primitive.ObjectIDFromHex(*params.UserProductID)
		// if err != nil {
		// 	return response, err
		// }
		q = append(q, bson.E{"takeUserId", params.TakeUserID})
	}

	// Filter by products id.
	if params.ProductID != nil && !params.ProductID.IsZero() {
		q = append(q, bson.E{"productId", params.ProductID})
	}

	// q = append(q, bson.E{"status", bson.M{"$gte": 0}})

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	if params.Sort != nil && len(params.Sort) > 0 {
		sortParam := bson.D{}
		for i := range params.Sort {
			sortParam = append(sortParam, bson.E{*params.Sort[i].Key, *params.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
		// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	}

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usera",
		"let":  bson.D{{Key: "userId", Value: "$userId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": tblImage,
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})

	limit := 100
	skip := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Skip != nil {
		skip = *params.Skip
	}

	pipe = append(pipe, bson.D{{"$limit", skip + limit}})
	pipe = append(pipe, bson.D{{"$skip", skip}})

	cursor, err := r.db.Collection(TblMessageRoom).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.MessageRoom, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count := len(resultSlice)
	// count, err := r.db.Collection(TblNode).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[model.MessageRoom]{
		Total: count,
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *MessageRoomMongo) CreateMessageRoom(userID string, messageRoom *model.MessageRoom) (*model.MessageRoom, error) {
	var result *model.MessageRoom

	collection := r.db.Collection(TblMessageRoom)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	createdAt := messageRoom.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	statusDefault := 1

	newMessageRoom := model.MessageRoomMongo{
		UserID:     userIDPrimitive,
		ProductID:  messageRoom.ProductID,
		TakeUserID: messageRoom.TakeUserID,
		OfferID:    messageRoom.OfferID,
		Status:     &statusDefault,
		Props:      messageRoom.Props,
		CreatedAt:  createdAt,
		UpdatedAt:  time.Now(),
	}

	res, err := collection.InsertOne(ctx, newMessageRoom)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblMessageRoom).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *MessageRoomMongo) UpdateMessageRoom(id string, userID string, data *model.MessageRoom) (*model.MessageRoom, error) {
	var result *model.MessageRoom
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblMessageRoom)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// idUser, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return result, err
	// }
	filter := bson.M{"_id": idPrimitive}

	// // Find old data
	// var oldResult *model.Message
	// err = collection.FindOne(ctx, filter).Decode(&oldResult)
	// if err != nil {
	// 	return result, err
	// }
	// oldMessage := model.Message{
	// 	UserID:  oldResult.UserID,
	// 	NodeID:  oldResult.NodeID,
	// 	Message: oldResult.Message,
	// 	Status:  oldResult.Status,
	// 	Props:   oldResult.Props,
	// }
	// _, err = r.db.Collection(TblMessage).UpdateOne(ctx, filter, bson.M{"$set": oldMessage})
	// if err != nil {
	// 	return result, err
	// }

	newData := bson.M{}
	if data.Status != nil {
		newData["status"] = data.Status

		status := *data.Status
		// fmt.Println("*data.Status=", *data.Status)
		if *data.Status < 0 {
			_, err = r.db.Collection(TblMessage).UpdateMany(ctx, bson.M{"roomId": idPrimitive}, bson.M{"$set": bson.M{"status": status}})
			if err != nil {
				return result, err
			}
		}
	}

	if data.Props != nil {
		newData["props"] = data.Props
	}

	// if data.Props != nil {
	// 	//newProps := make(map[string]interface{})
	// 	newProps := data.Props
	// 	if val, ok := data.Props["status"]; ok {
	// 		if val == -1.0 {
	// 			newDel := make(map[string]interface{})
	// 			newDel["user_id"] = userID
	// 			newDel["del_at"] = time.Now()
	// 			newProps["del"] = newDel
	// 		}
	// 	}
	// 	newData["props"] = newProps
	// }
	newData["updatedAt"] = time.Now()
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	// if err != nil {
	// 	return result, err
	// }
	resultResponse, err := r.FindMessageRoom(&model.MessageRoomFilter{ID: &idPrimitive})
	if err != nil {
		return result, err
	}

	if len(resultResponse.Data) > 0 {
		result = &resultResponse.Data[0]
	}

	return result, nil
}

func (r *MessageRoomMongo) DeleteMessageRoom(id string) (model.MessageRoom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.MessageRoom{}
	collection := r.db.Collection(TblMessageRoom)
	collectionMessage := r.db.Collection(TblMessage)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	// Delete messages for room.
	_, err = collectionMessage.DeleteMany(ctx, bson.M{"roomId": idPrimitive})
	if err != nil {
		return result, err
	}

	// Delete images for room.
	_, err = r.db.Collection(TblMessageImage).DeleteMany(ctx, bson.M{"roomId": idPrimitive})
	if err != nil {
		return result, err
	}

	return result, nil
}

// func (r *MessageRoomMongo) GetGroupForUser(userID string) ([]model.MessageGroupForUser, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
// 	defer cancel()

// 	var results []model.MessageGroupForUser

// 	q := bson.D{}

// 	if userID != "" {
// 		userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
// 		if err != nil {
// 			return results, err
// 		}
// 		queryArr := []bson.M{}
// 		queryArr = append(queryArr, bson.M{"userId": userIDPrimitive})
// 		queryArr = append(queryArr, bson.M{"userProductId": userIDPrimitive})
// 		q = append(q, bson.E{"$or", queryArr})
// 		// q = append(q, bson.E{"status", 1})
// 	}

// 	pipe := mongo.Pipeline{}
// 	pipe = append(pipe, bson.D{{"$match", q}})
// 	pipe = append(pipe,
// 		bson.D{
// 			{"$group", bson.D{
// 				// {"_id", "$productId"},
// 				{"_id", bson.D{
// 					{"productId", "$productId"},
// 					{"userId", "$userId"},
// 				}},
// 				{"productId", bson.D{{"$first", "$productId"}}},
// 				{"userId", bson.D{{"$first", "$userId"}}},
// 				// {"average_price", bson.D{{"$avg", "$price"}}},
// 				{"count", bson.D{{"$sum", 1}}},
// 			}}})
// 	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
// 		"from": "product",
// 		"as":   "products",
// 		"let":  bson.D{{Key: "productId", Value: "$productId"}},
// 		"pipeline": mongo.Pipeline{
// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$productId"}}}}},
// 			bson.D{{
// 				Key: "$lookup",
// 				Value: bson.M{
// 					"from": "image",
// 					"as":   "images",
// 					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
// 					"pipeline": mongo.Pipeline{
// 						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
// 					},
// 				},
// 			}},

// 			bson.D{{Key: "$lookup", Value: bson.M{
// 				"from": "users",
// 				"as":   "userb",
// 				"let":  bson.D{{Key: "userId", Value: "$user_id"}},
// 				"pipeline": mongo.Pipeline{
// 					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
// 					bson.D{{"$limit", 1}},
// 					bson.D{{
// 						Key: "$lookup",
// 						Value: bson.M{
// 							"from": "image",
// 							"as":   "images",
// 							"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
// 							"pipeline": mongo.Pipeline{
// 								bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
// 							},
// 						},
// 					}},
// 				},
// 			}}},
// 			bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userb"}}}},
// 		},
// 	}}})
// 	// pipe = append(pipe, bson.D{{"$unwind", "$product"}})
// 	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"product": bson.M{"$first": "$products"}}}})

// 	cursorGroup, err := r.db.Collection(TblMessage).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
// 	if err != nil {
// 		return results, err
// 	}
// 	defer cursorGroup.Close(ctx)

// 	if er := cursorGroup.All(ctx, &results); er != nil {
// 		return results, er
// 	}

// 	return results, nil
// }
