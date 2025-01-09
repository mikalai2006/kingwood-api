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

type TaskMontajWorkerMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTaskMontajWorkerMongo(db *mongo.Database, i18n config.I18nConfig) *TaskMontajWorkerMongo {
	return &TaskMontajWorkerMongo{db: db, i18n: i18n}
}

func (r *TaskMontajWorkerMongo) FindTaskMontajWorkerPopulate(input *domain.TaskMontajWorkerFilter) (domain.Response[domain.TaskMontajWorker], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TaskMontajWorker
	var response domain.Response[domain.TaskMontajWorker]

	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.TaskMontajWorker]{}, err
	// }

	// Filters
	q := bson.D{}

	if input.From != nil && !input.From.IsZero() {
		q = append(q, bson.E{"from", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.From)}}})
	}
	if input.To != nil && !input.To.IsZero() {
		q = append(q, bson.E{"to", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.To)}}})
	}
	if input.Date != nil && !input.Date.IsZero() {
		q = append(q, bson.E{"from", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.Date)}}})
		q = append(q, bson.E{"to", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.Date)}}})
	}
	if input.ObjectId != nil && len(input.ObjectId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.ObjectId {
			iDPrimitive, err := primitive.ObjectIDFromHex(*input.ObjectId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"objectId", bson.D{{"$in", ids}}})
	}
	if input.TaskId != nil && len(input.TaskId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.TaskId {
			iDPrimitive, err := primitive.ObjectIDFromHex(*input.TaskId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"taskId", bson.D{{"$in", ids}}})
	}
	if input.WorkerId != nil && len(input.WorkerId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.WorkerId {
			iDPrimitive, err := primitive.ObjectIDFromHex(*input.WorkerId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"workerId", bson.D{{"$in", ids}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

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
	// object.
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

	// task.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblTaskMontaj,
		"as":           "taska",
		"localField":   "taskId",
		"foreignField": "_id",
		// "pipeline": mongo.Pipeline{
		// 	bson.D{{Key: "$lookup", Value: bson.M{
		// 		"from":         tblOperation,
		// 		"as":           "operationa",
		// 		"localField":   "operationId",
		// 		"foreignField": "_id",
		// 	}}},
		// 	bson.D{{Key: "$set", Value: bson.M{"operation": bson.M{"$first": "$operationa"}}}},
		// },
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"taskMontaj": bson.M{"$first": "$taska"}}}})

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
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
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

	cursor, err := r.db.Collection(tblTaskMontajWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTaskMontajWorker).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.TaskMontajWorker]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *TaskMontajWorkerMongo) FindTaskMontajWorker(params domain.RequestParams) (domain.Response[domain.TaskMontajWorker], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TaskMontajWorker
	var response domain.Response[domain.TaskMontajWorker]
	// var response domain.Response[domain.TaskMontajWorker]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.TaskMontajWorker]{}, err
	// }

	// cursor, err := r.db.Collection(TblTaskMontajWorker).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.TaskMontajWorker, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.TaskMontajWorker]{}, err
	}
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

	cursor, err := r.db.Collection(tblTaskMontajWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(tblTaskMontajWorker).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.TaskMontajWorker]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *TaskMontajWorkerMongo) GqlGetTaskMontajWorkers(params domain.RequestParams) ([]*domain.TaskMontajWorker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*domain.TaskMontajWorker
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

	cursor, err := r.db.Collection(tblTaskMontajWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*domain.TaskMontajWorker, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// count, err := r.db.Collection(TblTaskMontajWorker).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return results, err
	// }

	// results = []*domain.TaskMontajWorker{
	// 	Total: int(count),
	// 	Skip:  int(params.Options.Skip),
	// 	Limit: int(params.Options.Limit),
	// 	Data:  resultSlice,
	// }
	return results, nil
}

func (r *TaskMontajWorkerMongo) CreateTaskMontajWorker(userID string, data *domain.TaskMontajWorker) (*domain.TaskMontajWorker, error) {
	var result *domain.TaskMontajWorker

	collection := r.db.Collection(tblTaskMontajWorker)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// find exists taskWorker.
	allExistMontajWorkers, err := r.FindTaskMontajWorker(domain.RequestParams{Filter: bson.D{{"workerId", data.WorkerId}}})
	if err != nil {
		return nil, err
	}

	if len(allExistMontajWorkers.Data) > 0 {
		existMontajWorkers := allExistMontajWorkers.Data[0]

		fmt.Println("exist montajworkers: ", existMontajWorkers)
		fmt.Println("exist after form: ", existMontajWorkers.From.After(data.From))
		fmt.Println("exist before form: ", existMontajWorkers.To.Before(data.To))
	}

	// var existTaskMontajWorker domain.TaskMontajWorker
	// r.db.Collection(TblTaskMontajWorker).FindOne(ctx, bson.M{"node_id": TaskMontajWorker.NodeID, "user_id": userIDPrimitive}).Decode(&existTaskMontajWorker)

	// if existTaskMontajWorker.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	sortOrder := int64(0)

	from := time.Now()
	if !data.From.IsZero() {
		from = data.From
	}

	newTaskMontajWorker := domain.TaskMontajWorkerInput{
		UserID:    userIDPrimitive,
		ObjectId:  data.ObjectId,
		TaskId:    data.TaskId,
		SortOrder: &sortOrder,
		StatusId:  data.StatusId,
		WorkerId:  data.WorkerId,
		Status:    data.Status,
		From:      from,
		To:        data.To,
		TypeGo:    data.TypeGo,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newTaskMontajWorker)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblTaskMontajWorker).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	updatedAt := TaskMontajWorker.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateTaskMontajWorker := &domain.TaskMontajWorkerInput{
	// 		Rate:      TaskMontajWorker.Rate,
	// 		TaskMontajWorker:    TaskMontajWorker.TaskMontajWorker,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateTaskMontajWorker(existTaskMontajWorker.ID.Hex(), userID, updateTaskMontajWorker)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *TaskMontajWorkerMongo) UpdateTaskMontajWorker(id string, userID string, data *domain.TaskMontajWorkerInput) (*domain.TaskMontajWorker, error) {
	var result *domain.TaskMontajWorker
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblTaskMontajWorker)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if !data.ObjectId.IsZero() {
		newData["objectId"] = data.ObjectId
	}
	if !data.TaskId.IsZero() {
		newData["taskId"] = data.TaskId
	}
	if !data.WorkerId.IsZero() {
		newData["workerId"] = data.WorkerId
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

	TaskMontajWorkers, err := r.FindTaskMontajWorkerPopulate(&domain.TaskMontajWorkerFilter{ID: []*string{&id}})
	// collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(TaskMontajWorkers.Data) > 0 {
		result = &TaskMontajWorkers.Data[0]
	} else {
		fmt.Println("Len TaskMontajWorkers.Data = ", len(TaskMontajWorkers.Data))
	}

	return result, nil
}

func (r *TaskMontajWorkerMongo) DeleteTaskMontajWorker(id string) (*domain.TaskMontajWorker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.TaskMontajWorker{}
	collection := r.db.Collection(tblTaskMontajWorker)

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
