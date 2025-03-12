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

type WorkHistoryMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewWorkHistoryMongo(db *mongo.Database, i18n config.I18nConfig) *WorkHistoryMongo {
	return &WorkHistoryMongo{db: db, i18n: i18n}
}

func (r *WorkHistoryMongo) FindWorkHistory(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.WorkHistory
	var response domain.Response[domain.WorkHistory]

	q := bson.D{}

	// Filters
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
	if input.WorkTimeId != nil && len(input.WorkTimeId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.WorkTimeId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.WorkTimeId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"workTimeId", bson.D{{"$in", ids}}})
	}
	if input.TaskId != nil && len(input.TaskId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.TaskId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.TaskId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"taskId", bson.D{{"$in", ids}}})
	}
	if input.OrderId != nil && len(input.OrderId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.OrderId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.OrderId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"orderId", bson.D{{"$in", ids}}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", input.Status})
	}
	// if input.Status != nil {
	// 	q = append(q, bson.E{"status", input.Status})
	// }

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
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

	cursor, err := r.db.Collection(tblWorkHistory).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTaskHistory).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.WorkHistory]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *WorkHistoryMongo) FindWorkHistoryPopulate(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.WorkHistory
	var response domain.Response[domain.WorkHistory]

	q := bson.D{}

	// Filters
	if input.ID != nil && len(input.ID) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.ID {
			idPrimitive, err := primitive.ObjectIDFromHex(input.ID[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, idPrimitive)
		}

		q = append(q, bson.E{"_id", bson.D{{"$in", ids}}})
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
	if input.TaskWorkerId != nil && len(input.TaskWorkerId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.TaskWorkerId {
			pid, err := primitive.ObjectIDFromHex(input.TaskWorkerId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, pid)
		}

		q = append(q, bson.E{"taskWorkerId", bson.D{{"$in", ids}}})
	}
	if input.WorkTimeId != nil && len(input.WorkTimeId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.WorkTimeId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.WorkTimeId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"workTimeId", bson.D{{"$in", ids}}})
	}
	if input.TaskId != nil && len(input.TaskId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.TaskId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.TaskId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"taskId", bson.D{{"$in", ids}}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", input.Status})
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
		"from": TblOrder,
		"as":   "ordera",
		// "localField":   "userId",
		// "foreignField": "_id",
		"let": bson.D{{Key: "orderId", Value: "$orderId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$orderId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": tblObject,
				"as":   "objecta",
				// "localField":   "userId",
				// "foreignField": "_id",
				"let": bson.D{{Key: "objectId", Value: "$objectId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$objectId"}}}}},
					bson.D{{"$limit", 1}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"object": bson.M{"$first": "$objecta"}}}},
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

	if input.Sort != nil && len(input.Sort) > 0 {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
		// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	}

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

	cursor, err := r.db.Collection(tblWorkHistory).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTaskHistory).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.WorkHistory]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *WorkHistoryMongo) GetStatByOrder(input domain.WorkHistoryFilter) ([]domain.WorkHistoryStatByOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var response []domain.WorkHistoryStatByOrder

	q := bson.D{}

	// Filters
	if input.OrderId != nil && len(input.OrderId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.OrderId {
			idPrimitive, err := primitive.ObjectIDFromHex(input.OrderId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, idPrimitive)
		}

		q = append(q, bson.E{"orderId", bson.D{{"$in", ids}}})
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
	if input.WorkTimeId != nil && len(input.WorkTimeId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.WorkTimeId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.WorkTimeId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"workTimeId", bson.D{{"$in", ids}}})
	}
	if input.TaskId != nil && len(input.TaskId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.TaskId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.TaskId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"taskId", bson.D{{"$in", ids}}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", input.Status})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	pipe = append(pipe,
		bson.D{
			{"$group", bson.D{
				{"_id", bson.D{
					{"workerId", "$workerId"},
					{"operationId", "$operationId"},
				},
				},
				// {"orderId", bson.D{{"$first", "$orderId"}}},
				// {"workerId", bson.D{{"$first", "$workerId"}}},
				// // {"average_price", bson.D{{"$avg", "$price"}}},
				{"count", bson.D{{"$sum", 1}}},
				{"total", bson.D{{"$sum", "$total"}}},
			}},
		})
	pipe = append(pipe,
		bson.D{
			{"$group", bson.D{
				{"_id", "$_id.workerId"},
				{"operations", bson.D{
					{"$push", bson.D{
						{"operationId", "$_id.operationId"},
						{"count", "$count"},
						{"total", "$total"},
					}}},
				},
				{"count", bson.D{{"$sum", "$count"}}},
				{"total", bson.D{{"$sum", "$total"}}},
			}},
		})
	// pipe = append(pipe,
	// 	bson.D{
	// 		{"$group", bson.D{
	// 			{"_id", "$_id.orderId"},
	// 			{"orderId", bson.D{{"$first", "$orderId"}}},
	// 			{"workerId", bson.D{{"$first", "$workerId"}}},
	// 			// {"average_price", bson.D{{"$avg", "$price"}}},
	// 			{"count", bson.D{{"$sum", 1}}},
	// 		}}})

	// // operation.
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": tblOperation,
	// 	"as":   "operationa",
	// 	"let":  bson.D{{Key: "operationId", Value: "$operationId"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$operationId"}}}}},
	// 		bson.D{{"$limit", 1}},
	// 	}}}},
	// 	bson.D{{Key: "$set", Value: bson.M{"operation": bson.M{"$first": "$operationa"}}}})

	// worker.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblUsers,
		"as":   "usera",
		"let":  bson.D{{Key: "workerId", Value: "$_id"}},
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

	if input.Sort != nil && len(input.Sort) > 0 {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
		// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	}

	// skip := 0
	// limit := 10
	// if input.Skip != nil {
	// 	pipe = append(pipe, bson.D{{"$skip", input.Skip}})
	// 	skip = *input.Skip
	// }
	// if input.Limit != nil {
	// 	pipe = append(pipe, bson.D{{"$limit", input.Limit}})
	// 	limit = *input.Limit
	// }

	cursor, err := r.db.Collection(tblWorkHistory).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	var test []interface{}
	if er := cursor.All(ctx, &response); er != nil {
		return response, er
	}

	fmt.Println("test: ", test)
	// count, err := r.db.Collection(tblTaskHistory).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	return response, nil
}

func (r *WorkHistoryMongo) CreateWorkHistory(userID string, data *domain.WorkHistory) (*domain.WorkHistory, error) {
	var result *domain.WorkHistory

	collection := r.db.Collection(tblWorkHistory)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existTask domain.TaskHistory
	// r.db.Collection(TblTask).FindOne(ctx, bson.M{"node_id": Task.NodeID, "userId": userIDPrimitive}).Decode(&existTask)

	// if existTask.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	// // get sortOrder index.
	// var allTaskByOrder []domain.TaskHistory
	// cursor, err := collection.Find(ctx, bson.D{{"orderId", data.OrderId}})
	// if err != nil {
	// 	return result, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &allTaskByOrder); er != nil {
	// 	return result, er
	// }
	date := time.Now()
	if !data.Date.IsZero() {
		date = data.Date
	}
	defaultTotal := int64(0)
	if data.Total != nil {
		defaultTotal = *data.Total
	}
	defaultTotalTime := int64(0)
	if data.TotalTime != nil {
		defaultTotalTime = *data.TotalTime
	}

	newTask := domain.WorkHistoryInput{
		OrderId:      &data.OrderId,
		TaskId:       &data.TaskId,
		WorkerId:     data.WorkerId,
		ObjectId:     &data.ObjectId,
		OperationId:  &data.OperationId,
		UserID:       userIDPrimitive,
		TaskWorkerId: &data.TaskWorkerId,
		Status:       &data.Status,
		Date:         date,
		From:         data.From,
		To:           data.To,
		WorkTimeId:   data.WorkTimeId,
		TotalTime:    &defaultTotalTime,
		Oklad:        data.Oklad,
		Total:        &defaultTotal,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblWorkHistory).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *WorkHistoryMongo) UpdateWorkHistory(id string, userID string, data *domain.WorkHistoryInput) (*domain.WorkHistory, error) {
	var result *domain.WorkHistory
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblWorkHistory)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	// if !data.WorkerId.IsZero() {
	// 	newData["workerId"] = data.WorkerId
	// }
	// if !data.OperationId.IsZero() {
	// 	newData["operationId"] = data.OperationId
	// }
	if data.Status != nil {
		newData["status"] = &data.Status
	}
	if data.OrderId != nil {
		newData["orderId"] = data.OrderId
	}
	if data.TaskId != nil {
		newData["taskId"] = data.TaskId
	}
	if data.OperationId != nil {
		newData["operationId"] = data.OperationId
	}
	if !data.WorkerId.IsZero() {
		newData["workerId"] = data.WorkerId
	}
	if !data.WorkTimeId.IsZero() {
		newData["workTimeId"] = data.WorkTimeId
	}
	if data.TaskWorkerId != nil {
		newData["taskWorkerId"] = data.TaskWorkerId
	}
	if data.ObjectId != nil {
		newData["objectId"] = data.ObjectId
	}
	if data.TotalTime != nil {
		newData["totalTime"] = data.TotalTime
	}

	if !data.From.IsZero() {
		newData["from"] = data.From
	}
	if !data.To.IsZero() {
		newData["to"] = data.To
	}
	if data.Total != nil {
		newData["total"] = data.Total
	}
	newData["updatedAt"] = time.Now()

	if len(data.Props) > 0 {
		newData["props"] = data.Props
	}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	updatedItems, err := r.FindWorkHistoryPopulate(domain.WorkHistoryFilter{ID: []string{id}})
	if err != nil {
		return result, err
	}
	if len(updatedItems.Data) > 0 {
		result = &updatedItems.Data[0]
	}

	return result, nil
}

func (r *WorkHistoryMongo) DeleteWorkHistory(id string) (*domain.WorkHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.WorkHistory{}
	collection := r.db.Collection(tblWorkHistory)

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
