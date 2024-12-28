package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewOrderMongo(db *mongo.Database, i18n config.I18nConfig) *OrderMongo {
	return &OrderMongo{db: db, i18n: i18n}
}

func (r *OrderMongo) FindOrder(input *domain.OrderFilter) (domain.Response[domain.Order], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Order
	var response domain.Response[domain.Order]
	// var response domain.Response[domain.Order]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.Order]{}, err
	// }

	// cursor, err := r.db.Collection(TblOrder).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.Order, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.Order]{}, err
	// }
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": "users",
	// 	"as":   "usera",
	// 	// "localField":   "user_id",
	// 	// "foreignField": "_id",
	// 	"let": bson.D{{Key: "userId", Value: "$userId"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
	// 		bson.D{{"$limit", 1}},
	// 		bson.D{{
	// 			Key: "$lookup",
	// 			Value: bson.M{
	// 				"from": tblImage,
	// 				"as":   "images",
	// 				// "localField":   "_id",
	// 				// "foreignField": "service_id",
	// 				"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 				"pipeline": mongo.Pipeline{
	// 					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
	// 				},
	// 			},
	// 		}},
	// 	},
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})
	q := bson.D{}

	// Filters
	if input.Name != nil && *input.Name != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", *input.Name), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
	}
	if input.Group != nil && len(input.Group) > 0 {
		q = append(q, bson.E{"group", bson.M{"$elemMatch": bson.D{{"$in", input.Group}}}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", input.Status})
	}
	if input.ObjectIds != nil {
		objectIds := []primitive.ObjectID{}
		for key, _ := range input.ObjectIds {
			idObjectPrimitive, err := primitive.ObjectIDFromHex(input.ObjectIds[key])
			if err != nil {
				return response, err
			}
			objectIds = append(objectIds, idObjectPrimitive)
		}
		q = append(q, bson.E{"objectId", bson.D{{"$in", objectIds}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblObject,
		"as":   "objecta",
		// "localField":   "user_id",
		// "foreignField": "_id",
		"let": bson.D{{Key: "objectId", Value: "$objectId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$objectId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"object": bson.M{"$first": "$objecta"}}}})

	if input.Sort != nil && len(input.Sort) > 0 {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{*input.Sort[i].Key, *input.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
		// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	}

	skip := 0
	limit := 10
	if input.Skip != nil {
		pipe = append(pipe, bson.D{{"$skip", input.Skip}})
		skip = *input.Skip
	}
	if input.Limit != nil {
		pipe = append(pipe, bson.D{{"$limit", input.Limit}})
		limit = *input.Limit
	}

	cursor, err := r.db.Collection(TblOrder).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(TblOrder).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.Order]{
		Total: 0,
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *OrderMongo) GetAllOrder(params domain.RequestParams) (domain.Response[domain.Order], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Order
	var response domain.Response[domain.Order]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Order]{}, err
	}

	cursor, err := r.db.Collection(TblOrder).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Order, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblOrder).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Order]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *OrderMongo) GqlGetOrders(params domain.RequestParams) ([]*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*domain.Order
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usera",
		"let":  bson.D{{Key: "userId", Value: "$user_id"}},
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

	cursor, err := r.db.Collection(TblOrder).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*domain.Order, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// count, err := r.db.Collection(TblOrder).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return results, err
	// }

	// results = []*domain.Order{
	// 	Total: int(count),
	// 	Skip:  int(params.Options.Skip),
	// 	Limit: int(params.Options.Limit),
	// 	Data:  resultSlice,
	// }
	return results, nil
}

func (r *OrderMongo) CreateOrder(userID string, data *domain.Order) (*domain.Order, error) {
	var result *domain.Order

	collection := r.db.Collection(TblOrder)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existOrder domain.Order
	// r.db.Collection(TblOrder).FindOne(ctx, bson.M{"node_id": Order.NodeID, "user_id": userIDPrimitive}).Decode(&existOrder)

	// if existOrder.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	itemCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	newOrder := domain.OrderInput{
		UserID:        userIDPrimitive,
		Name:          data.Name,
		Description:   data.Description,
		ObjectId:      data.ObjectId,
		Number:        itemCount + 1,
		ConstructorId: data.ConstructorId,
		Priority:      data.Priority,
		Term:          data.Term,
		TermMontaj:    data.TermMontaj,
		Status:        data.Status,
		Group:         data.Group,
		// NeedMontaj:    data.NeedMontaj,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblOrder).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	updatedAt := Order.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateOrder := &domain.OrderInput{
	// 		Rate:      Order.Rate,
	// 		Order:    Order.Order,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateOrder(existOrder.ID.Hex(), userID, updateOrder)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *OrderMongo) UpdateOrder(id string, userID string, data *domain.OrderInput) (*domain.Order, error) {
	var result *domain.Order
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblOrder)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Name != "" {
		newData["name"] = data.Name
	}
	if !data.ObjectId.IsZero() {
		newData["objectId"] = data.ObjectId
	}
	if !data.Term.IsZero() {
		newData["term"] = data.Term
	}
	if !data.TermMontaj.IsZero() {
		newData["termMontaj"] = data.TermMontaj
	}
	if data.Priority != nil {
		newData["priority"] = data.Priority
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	// if data.NeedMontaj != nil {
	// 	newData["needMontaj"] = data.NeedMontaj
	// }
	if !data.ConstructorId.IsZero() {
		newData["constructorId"] = data.ConstructorId
	}
	if len(data.Group) > 0 {
		newData["group"] = data.Group
	}
	if data.Status != nil {
		newData["status"] = data.Status
	}
	newData["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *OrderMongo) DeleteOrder(id string) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.Order{}
	collection := r.db.Collection(TblOrder)

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

	return result, nil
}
