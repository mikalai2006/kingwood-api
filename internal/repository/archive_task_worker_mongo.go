package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArchiveTaskWorkerMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchiveTaskWorkerMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveTaskWorkerMongo {
	return &ArchiveTaskWorkerMongo{db: db, i18n: i18n}
}

func (r *ArchiveTaskWorkerMongo) CreateArchiveTaskWorker(userID string, data *domain.TaskWorker) (*domain.ArchiveTaskWorker, error) {
	var result *domain.ArchiveTaskWorker

	collection := r.db.Collection(TblArchiveTaskWorker)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newArchiveTaskWorker := domain.ArchiveTaskWorkerInput{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:          data.ID,
		UserID:      data.UserID,
		ObjectId:    data.ObjectId,
		OrderId:     data.OrderId,
		TaskId:      data.TaskId,
		WorkerId:    data.WorkerId,
		OperationId: data.OperationId,
		SortOrder:   data.SortOrder,
		StatusId:    data.StatusId,
		Status:      data.Status,
		From:        *data.From,
		To:          *data.To,
		TypeGo:      data.TypeGo,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	res, err := collection.InsertOne(ctx, newArchiveTaskWorker)
	if err != nil {
		return nil, err
	}

	// idCreatedItem := res.InsertedID.(primitive.ObjectID).Hex();
	// err = r.db.Collection(TblArchiveTaskWorker).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	insertedId := res.InsertedID.(primitive.ObjectID).Hex()
	ArchiveTaskWorkers, err := r.FindArchiveTaskWorker(&domain.ArchiveTaskWorkerFilter{ID: []string{insertedId}})
	if err != nil {
		return nil, err
	}
	if len(ArchiveTaskWorkers.Data) > 0 {
		result = &ArchiveTaskWorkers.Data[0]
	}
	// } else {
	// 	updatedAt := ArchiveTaskWorker.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateArchiveTaskWorker := &domain.ArchiveTaskWorkerInput{
	// 		Rate:      ArchiveTaskWorker.Rate,
	// 		ArchiveTaskWorker:    ArchiveTaskWorker.ArchiveTaskWorker,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateArchiveTaskWorker(existArchiveTaskWorker.ID.Hex(), userID, updateArchiveTaskWorker)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *ArchiveTaskWorkerMongo) FindArchiveTaskWorker(input *domain.ArchiveTaskWorkerFilter) (domain.Response[domain.ArchiveTaskWorker], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.ArchiveTaskWorker
	var response domain.Response[domain.ArchiveTaskWorker]

	// Filters
	q := bson.D{}
	if input.From != nil && !input.From.IsZero() {
		q = append(q, bson.E{"from", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.From)}}})
	}
	if input.To != nil && !input.To.IsZero() {
		q = append(q, bson.E{"to", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.To)}}})
	}
	if input.Date != nil && !input.Date.IsZero() {
		// q = append(q, bson.E{"from", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.Date)}}})
		q = append(q, bson.E{"to", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.Date)}}})
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
	if input.WorkerId != nil && len(input.WorkerId) > 0 {
		workerIds := []primitive.ObjectID{}
		for i, _ := range input.WorkerId {
			workerIDPrimitive, err := primitive.ObjectIDFromHex(input.WorkerId[i])
			if err != nil {
				return response, err
			}

			workerIds = append(workerIds, workerIDPrimitive)
		}

		q = append(q, bson.E{"workerId", bson.D{{"$in", workerIds}}})
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
	if input.TaskId != nil && len(input.TaskId) > 0 {
		taskIds := []primitive.ObjectID{}
		for i, _ := range input.TaskId {
			taskIDPrimitive, err := primitive.ObjectIDFromHex(input.TaskId[i])
			if err != nil {
				return response, err
			}

			taskIds = append(taskIds, taskIDPrimitive)
		}

		q = append(q, bson.E{"taskId", bson.D{{"$in", taskIds}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	// // object.
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": tblObject,
	// 	"as":   "objecta",
	// 	// "localField":   "userId",
	// 	// "foreignField": "_id",
	// 	"let": bson.D{{Key: "objectId", Value: "$objectId"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$objectId"}}}}},
	// 		bson.D{{"$limit", 1}},
	// 	},
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"object": bson.M{"$first": "$objecta"}}}})

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
		"from": TblOrder,
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

	// task.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblTask,
		"as":           "taska",
		"localField":   "taskId",
		"foreignField": "_id",
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$lookup", Value: bson.M{
				"from":         tblOperation,
				"as":           "operationa",
				"localField":   "operationId",
				"foreignField": "_id",
			}}},
			bson.D{{Key: "$set", Value: bson.M{"operation": bson.M{"$first": "$operationa"}}}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"task": bson.M{"$first": "$taska"}}}})

	// worker.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
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
			// add populate auth.
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from":         TblAuth,
					"as":           "auths",
					"localField":   "userId",
					"foreignField": "_id",
					// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
					// "pipeline": mongo.Pipeline{
					// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
					// },
				},
			}},
			bson.D{{Key: "$set", Value: bson.M{"auth": bson.M{"$first": "$auths"}}}},
			bson.D{{Key: "$set", Value: bson.M{"authPrivate": bson.M{"$first": "$auths"}}}},

			// post.
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": TblPost,
					"as":   "posts",
					// "localField":   "_id",
					// "foreignField": "serviceId",
					"let": bson.D{{Key: "postId", Value: "$postId"}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$postId"}}}}},
					},
				},
			}},
			bson.D{{Key: "$set", Value: bson.M{"postObject": bson.M{"$first": "$posts"}}}},
			// role.
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": TblRole,
				"as":   "rolea",
				// "localField":   "userId",
				// "foreignField": "_id",
				"let": bson.D{{Key: "roleId", Value: "$roleId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$roleId"}}}}},
					bson.D{{"$limit", 1}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$rolea"}}}},
		},
	}}},
		bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}},
	)

	// taskStatus.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblTaskStatus,
		"as":           "taskStatusa",
		"localField":   "statusId",
		"foreignField": "_id",
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"taskStatus": bson.M{"$first": "$taskStatusa"}}}})

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

	cursor, err := r.db.Collection(TblArchiveTaskWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(TblArchiveTaskWorker).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.ArchiveTaskWorker]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *ArchiveTaskWorkerMongo) DeleteArchiveTaskWorker(id string) (*domain.ArchiveTaskWorker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.ArchiveTaskWorker{}
	collection := r.db.Collection(TblArchiveTaskWorker)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	ArchiveTaskWorkers, err := r.FindArchiveTaskWorker(&domain.ArchiveTaskWorkerFilter{ID: []string{id}})
	if err != nil {
		return result, err
	}
	if len(ArchiveTaskWorkers.Data) > 0 {
		result = &ArchiveTaskWorkers.Data[0]
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}
