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

type ArchiveTaskMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchiveTaskMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveTaskMongo {
	return &ArchiveTaskMongo{db: db, i18n: i18n}
}

func (r *ArchiveTaskMongo) CreateArchiveTask(userID string, data *domain.Task) (*domain.ArchiveTask, error) {
	var result *domain.ArchiveTask

	collection := r.db.Collection(TblArchiveTask)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newArchiveTask := domain.ArchiveTaskInput{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:          data.ID,
		OrderId:     data.OrderId,
		UserID:      data.UserID,
		OperationId: data.OperationId,
		Name:        data.Name,
		SortOrder:   data.SortOrder,
		StatusId:    data.StatusId,
		Active:      data.Active,
		StartAt:     data.StartAt,
		AutoCheck:   data.AutoCheck,
		Status:      data.Status,
		From:        *data.From,
		To:          *data.To,
		TypeGo:      data.TypeGo,
		ObjectId:    data.ObjectId,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	res, err := collection.InsertOne(ctx, newArchiveTask)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblArchiveTask).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ArchiveTaskMongo) FindArchiveTask(input domain.ArchiveTaskFilter) (domain.Response[domain.ArchiveTask], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.ArchiveTask
	var response domain.Response[domain.ArchiveTask]

	// Filters
	q := bson.D{}

	if input.Name != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", input.Name), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
	}
	if input.ID != nil && len(input.ID) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.ID {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.ID[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"_id", bson.D{{"$in", ids}}})
	}
	if input.OrderId != nil && len(input.OrderId) > 0 {
		orderIds := []primitive.ObjectID{}
		for i, _ := range input.OrderId {
			orderIDPrimitive, err := primitive.ObjectIDFromHex(input.OrderId[i])
			if err != nil {
				return response, err
			}

			orderIds = append(orderIds, orderIDPrimitive)
		}

		q = append(q, bson.E{"orderId", bson.D{{"$in", orderIds}}})
	}
	if input.ObjectId != nil && len(input.ObjectId) > 0 {
		objectIds := []primitive.ObjectID{}
		for i, _ := range input.ObjectId {
			objectIDPrimitive, err := primitive.ObjectIDFromHex(input.ObjectId[i])
			if err != nil {
				return response, err
			}

			objectIds = append(objectIds, objectIDPrimitive)
		}

		q = append(q, bson.E{"objectId", bson.D{{"$in", objectIds}}})
	}
	if input.OperationId != nil && len(input.OperationId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.OperationId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.OperationId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"operationId", bson.D{{"$in", ids}}})
	}
	if input.Status != nil && len(input.Status) > 0 {
		arr := []string{}
		for i, _ := range input.Status {
			arr = append(arr, input.Status[i])
		}

		q = append(q, bson.E{"status", bson.D{{"$in", arr}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	// object.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblObject,
		"as":   "objecta",
		// "localField":   "userId",
		// "foreignField": "_id",
		"let": bson.D{{Key: "objectId", Value: "$objectId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$objectId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"object": bson.M{"$first": "$objecta"}}}})
	// order.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": TblArchiveOrder,
		"as":   "ordera",
		// "localField":   "userId",
		// "foreignField": "_id",
		"let": bson.D{{Key: "orderId", Value: "$orderId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$orderId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"order": bson.M{"$first": "$ordera"}}}})
	// workers.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         TblArchiveTaskWorker,
		"as":           "workers",
		"localField":   "_id",
		"foreignField": "ArchiveTaskId",
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": tblUsers,
				"as":   "usera",
				"let":  bson.D{{Key: "workerId", Value: "$workerId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$workerId"}}}}},
					bson.D{{"$limit", 1}},
					bson.D{{
						Key: "$lookup",
						Value: bson.M{
							"from": tblImage,
							"as":   "images",
							"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
							"pipeline": mongo.Pipeline{
								bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
							},
						},
					}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}},
		},
	}}})
	// operation.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblOperation,
		"as":           "operationa",
		"localField":   "operationId",
		"foreignField": "_id",
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"operation": bson.M{"$first": "$operationa"}}}})

	if input.Sort != nil && len(input.Sort) > 0 {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
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

	cursor, err := r.db.Collection(TblArchiveTask).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(TblArchiveTask).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.ArchiveTask]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *ArchiveTaskMongo) DeleteArchiveTask(id string) (*domain.ArchiveTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.ArchiveTask{}
	collection := r.db.Collection(TblArchiveTask)

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
