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

type TaskWorkerMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTaskWorkerMongo(db *mongo.Database, i18n config.I18nConfig) *TaskWorkerMongo {
	return &TaskWorkerMongo{db: db, i18n: i18n}
}

func (r *TaskWorkerMongo) FindTaskWorkerFlat(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorkerFlat], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var response domain.Response[domain.TaskWorkerFlat]

	// Filters
	q := bson.D{}
	if input.Status != nil && len(input.Status) > 0 {
		q = append(q, bson.E{"status", bson.D{{"$in", input.Status}}})
	}
	if input.From != nil && !input.From.IsZero() {
		q = append(q, bson.E{"from", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.From)}}})
	}
	if input.To != nil && !input.To.IsZero() {
		q = append(q, bson.E{"to", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.To)}}})
	}
	if input.Date != nil && !input.Date.IsZero() {
		// fmt.Println("from: ", primitive.NewDateTimeFromTime(*input.Date), primitive.NewDateTimeFromTime(*input.Date).Time().String())
		// fmt.Println("to: ", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6)), primitive.NewDateTimeFromTime(*input.Date).Time().AddDate(0, 0, 6).String())
		// q = append(q, bson.E{"from", bson.D{{"$lte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6))}}})
		// q = append(q, bson.E{"to", bson.D{{"$gte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 0))}}})

		queryArr := []bson.M{
			{"$and": []bson.M{
				{"from": bson.D{{"$lte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6))}}},
				{"to": bson.D{{"$gte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 0))}}},
			}},
			{"$and": []bson.M{
				{"from": bson.D{{"$gte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 0))}}},
				{"to": bson.D{{"$lte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6))}}},
			}},
		}
		q = append(q, bson.E{"$or", queryArr})
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
	skip := 0
	limit := 10
	dataOptions := bson.A{}
	if input.Skip != nil {
		skip = *input.Skip
	}
	dataOptions = append(dataOptions, bson.D{{"$skip", skip}})

	if input.Limit != nil {
		limit = *input.Limit
	}
	if limit > 0 {
		dataOptions = append(dataOptions, bson.D{{"$limit", limit}})
	}

	if input.Sort != nil {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		dataOptions = append(dataOptions, bson.D{{"$sort", sortParam}})
	}

	pipe = append(pipe, bson.D{{Key: "$facet", Value: bson.D{
		{"data", dataOptions},
		{Key: "metadata", Value: mongo.Pipeline{
			bson.D{{"$group", bson.D{
				{"_id", nil},
				{"total", bson.D{{"$sum", 1}}}}}},
		}},
	},
	}})

	cursor, err := r.db.Collection(tblTaskWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	resultMap := []bson.M{}
	if er := cursor.All(ctx, &resultMap); er != nil {
		return response, er
	}
	resultFacetOne := domain.ResultFacetTaskWorkerFlat{}
	if len(resultMap) > 0 {
		bsonBytes, errs := bson.Marshal(resultMap[0])
		if errs != nil {
			fmt.Println("rrrrr: errs ", errs)
		}

		bson.Unmarshal(bsonBytes, &resultFacetOne)
	}

	total := 0
	if len(resultFacetOne.Metadata) > 0 {
		total = resultFacetOne.Metadata[0].Total
	}

	response = domain.Response[domain.TaskWorkerFlat]{
		Total: total,
		Skip:  skip,
		Limit: limit,
		Data:  resultFacetOne.Data,
	}
	return response, nil
}

func (r *TaskWorkerMongo) FindTaskWorkerPopulate(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorker], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// var results []domain.TaskWorker
	var response domain.Response[domain.TaskWorker]

	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.TaskWorker]{}, err
	// }
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from":         tblTask,
	// 	"as":           "taska",
	// 	"localField":   "taskId",
	// 	"foreignField": "_id",
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$lookup", Value: bson.M{
	// 			"from":         tblOperation,
	// 			"as":           "operationa",
	// 			"localField":   "operationId",
	// 			"foreignField": "_id",
	// 		}}},
	// 		bson.D{{Key: "$set", Value: bson.M{"operation": bson.M{"$first": "$operationa"}}}},
	// 	},
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"task": bson.M{"$first": "$taska"}}}})
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from":         tblTaskStatus,
	// 	"as":           "taskStatusa",
	// 	"localField":   "statusId",
	// 	"foreignField": "_id",
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"taskStatus": bson.M{"$first": "$taskStatusa"}}}})
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": TblOrder,
	// 	"as":   "ordera",
	// 	"let":  bson.D{{Key: "orderId", Value: "$task.orderId"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$orderId"}}}}},
	// 		bson.D{{"$limit", 1}},
	// 	},
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"order": bson.M{"$first": "$ordera"}}}})

	// Filters
	q := bson.D{}
	if input.Status != nil && len(input.Status) > 0 {
		q = append(q, bson.E{"status", bson.D{{"$in", input.Status}}})
	}
	if input.From != nil && !input.From.IsZero() {
		q = append(q, bson.E{"from", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.From)}}})
	}
	if input.To != nil && !input.To.IsZero() {
		q = append(q, bson.E{"to", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.To)}}})
	}
	if input.Date != nil && !input.Date.IsZero() {
		// fmt.Println("from: ", primitive.NewDateTimeFromTime(*input.Date), primitive.NewDateTimeFromTime(*input.Date).Time().String())
		// fmt.Println("to: ", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6)), primitive.NewDateTimeFromTime(*input.Date).Time().AddDate(0, 0, 6).String())
		// q = append(q, bson.E{"from", bson.D{{"$lte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6))}}})
		// q = append(q, bson.E{"to", bson.D{{"$gte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 0))}}})

		queryArr := []bson.M{
			{"$and": []bson.M{
				{"from": bson.D{{"$lte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6))}}},
				{"to": bson.D{{"$gte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 0))}}},
			}},
			{"$and": []bson.M{
				{"from": bson.D{{"$gte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 0))}}},
				{"to": bson.D{{"$lte", primitive.NewDateTimeFromTime((*input.Date).AddDate(0, 0, 6))}}},
			}},
		}
		q = append(q, bson.E{"$or", queryArr})
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
	skip := 0
	limit := 10
	dataOptions := bson.A{}
	if input.Skip != nil {
		skip = *input.Skip
	}
	dataOptions = append(dataOptions, bson.D{{"$skip", skip}})

	if input.Limit != nil {
		limit = *input.Limit
	}
	if limit > 0 {
		dataOptions = append(dataOptions, bson.D{{"$limit", limit}})
	}

	if input.Sort != nil {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		dataOptions = append(dataOptions, bson.D{{"$sort", sortParam}})
	}

	pipe = append(pipe, bson.D{{Key: "$facet", Value: bson.D{
		{"data", dataOptions},
		{Key: "metadata", Value: mongo.Pipeline{
			bson.D{{"$group", bson.D{
				{"_id", nil},
				{"total", bson.D{{"$sum", 1}}}}}},
		}},
	},
	}})

	cursor, err := r.db.Collection(tblTaskWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	resultMap := []bson.M{}
	if er := cursor.All(ctx, &resultMap); er != nil {
		return response, er
	}
	resultFacetOne := domain.ResultFacetTaskWorker{}
	if len(resultMap) > 0 {
		bsonBytes, errs := bson.Marshal(resultMap[0])
		if errs != nil {
			fmt.Println("rrrrr: errs ", errs)
		}

		bson.Unmarshal(bsonBytes, &resultFacetOne)
	}

	total := 0
	if len(resultFacetOne.Metadata) > 0 {
		total = resultFacetOne.Metadata[0].Total
	}

	response = domain.Response[domain.TaskWorker]{
		Total: total,
		Skip:  skip,
		Limit: limit,
		Data:  resultFacetOne.Data,
	}
	return response, nil
}

// func (r *TaskWorkerMongo) FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
// 	defer cancel()

// 	var results []domain.TaskWorker
// 	var response domain.Response[domain.TaskWorker]
// 	// var response domain.Response[domain.TaskWorker]
// 	// filter, opts, err := CreateFilterAndOptions(params)
// 	// if err != nil {
// 	// 	return domain.Response[domain.TaskWorker]{}, err
// 	// }

// 	// cursor, err := r.db.Collection(TblTaskWorker).Find(ctx, filter, opts)
// 	// if err != nil {
// 	// 	return response, err
// 	// }
// 	// defer cursor.Close(ctx)

// 	// if er := cursor.All(ctx, &results); er != nil {
// 	// 	return response, er
// 	// }

// 	// resultSlice := make([]domain.TaskWorker, len(results))
// 	// // for i, d := range results {
// 	// // 	resultSlice[i] = d
// 	// // }
// 	// copy(resultSlice, results)

// 	pipe, err := CreatePipeline(params, &r.i18n)
// 	if err != nil {
// 		return domain.Response[domain.TaskWorker]{}, err
// 	}
// 	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
// 	// 	"from": "users",
// 	// 	"as":   "usera",
// 	// 	// "localField":   "userId",
// 	// 	// "foreignField": "_id",
// 	// 	"let": bson.D{{Key: "userId", Value: "$userId"}},
// 	// 	"pipeline": mongo.Pipeline{
// 	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
// 	// 		bson.D{{"$limit", 1}},
// 	// 		bson.D{{
// 	// 			Key: "$lookup",
// 	// 			Value: bson.M{
// 	// 				"from": tblImage,
// 	// 				"as":   "images",
// 	// 				// "localField":   "_id",
// 	// 				// "foreignField": "serviceId",
// 	// 				"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
// 	// 				"pipeline": mongo.Pipeline{
// 	// 					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
// 	// 				},
// 	// 			},
// 	// 		}},
// 	// 	},
// 	// }}})
// 	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})

// 	cursor, err := r.db.Collection(tblTaskWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
// 	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
// 	if err != nil {
// 		return response, err
// 	}
// 	defer cursor.Close(ctx)

// 	if er := cursor.All(ctx, &results); er != nil {
// 		return response, er
// 	}

// 	count, err := r.db.Collection(tblTaskWorker).CountDocuments(ctx, params.Filter)
// 	if err != nil {
// 		return response, err
// 	}

// 	response = domain.Response[domain.TaskWorker]{
// 		Total: int(count),
// 		Skip:  int(params.Options.Skip),
// 		Limit: int(params.Options.Limit),
// 		Data:  results,
// 	}
// 	return response, nil
// }

func (r *TaskWorkerMongo) CreateTaskWorker(userID string, data *domain.TaskWorker) (*domain.TaskWorker, error) {
	var result *domain.TaskWorker

	collection := r.db.Collection(tblTaskWorker)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existTaskWorker domain.TaskWorker
	// r.db.Collection(TblTaskWorker).FindOne(ctx, bson.M{"node_id": TaskWorker.NodeID, "userId": userIDPrimitive}).Decode(&existTaskWorker)

	// if existTaskWorker.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	sortOrder := int64(0)

	from := time.Now()
	if data.From != nil && !data.From.IsZero() {
		from = *data.From
	}
	to := time.Now().AddDate(1, 0, 0)
	if data.To != nil && !data.From.IsZero() {
		to = *data.To
	}

	newTaskWorker := domain.TaskWorkerInput{
		UserID:      userIDPrimitive,
		ObjectId:    data.ObjectId,
		OrderId:     data.OrderId,
		TaskId:      data.TaskId,
		WorkerId:    data.WorkerId,
		OperationId: data.OperationId,
		SortOrder:   &sortOrder,
		StatusId:    data.StatusId,
		Status:      data.Status,
		From:        from,
		To:          to,
		TypeGo:      data.TypeGo,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newTaskWorker)
	if err != nil {
		return nil, err
	}

	// idCreatedItem := res.InsertedID.(primitive.ObjectID).Hex();
	// err = r.db.Collection(tblTaskWorker).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	insertedId := res.InsertedID.(primitive.ObjectID).Hex()
	taskWorkers, err := r.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{ID: []string{insertedId}})
	if err != nil {
		return nil, err
	}
	if len(taskWorkers.Data) > 0 {
		result = &taskWorkers.Data[0]
	}
	// } else {
	// 	updatedAt := TaskWorker.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateTaskWorker := &domain.TaskWorkerInput{
	// 		Rate:      TaskWorker.Rate,
	// 		TaskWorker:    TaskWorker.TaskWorker,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateTaskWorker(existTaskWorker.ID.Hex(), userID, updateTaskWorker)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *TaskWorkerMongo) UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput) (*domain.TaskWorker, error) {
	var result *domain.TaskWorker
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblTaskWorker)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if !data.ObjectId.IsZero() {
		newData["objectId"] = data.ObjectId
	}
	if !data.OrderId.IsZero() {
		newData["orderId"] = data.OrderId
	}
	if !data.TaskId.IsZero() {
		newData["taskId"] = data.TaskId
	}
	if !data.WorkerId.IsZero() {
		newData["workerId"] = data.WorkerId
	}
	if !data.OperationId.IsZero() {
		newData["operationId"] = data.OperationId
	}
	if data.SortOrder != nil {
		newData["sortOrder"] = data.SortOrder
	}
	if !data.StatusId.IsZero() {
		newData["statusId"] = data.StatusId
	}
	if data.Status != "" {
		newData["status"] = data.Status
	}
	if !data.From.IsZero() {
		newData["from"] = data.From
	}
	if !data.To.IsZero() {
		newData["to"] = data.To
	}
	if data.TypeGo != "" {
		newData["typeGo"] = data.TypeGo
	}
	newData["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	taskWorkers, err := r.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{ID: []string{id}})
	// taskWorkers, err := r.FindTaskWorkerPopulate(domain.RequestParams{Filter: bson.D{{"_id", idPrimitive}}})
	// collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(taskWorkers.Data) > 0 {
		result = &taskWorkers.Data[0]
	} else {
		fmt.Println("Len taskWorkers.Data = ", len(taskWorkers.Data))
	}

	return result, nil
}

func (r *TaskWorkerMongo) DeleteTaskWorker(id string) (*domain.TaskWorker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.TaskWorker{}
	collection := r.db.Collection(tblTaskWorker)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	taskWorkers, err := r.FindTaskWorkerPopulate(&domain.TaskWorkerFilter{ID: []string{id}})
	if err != nil {
		return result, err
	}
	if len(taskWorkers.Data) > 0 {
		result = &taskWorkers.Data[0]
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}
