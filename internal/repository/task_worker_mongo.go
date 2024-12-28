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

type TaskWorkerMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTaskWorkerMongo(db *mongo.Database, i18n config.I18nConfig) *TaskWorkerMongo {
	return &TaskWorkerMongo{db: db, i18n: i18n}
}

func (r *TaskWorkerMongo) FindTaskWorkerPopulate(params domain.RequestParams) (domain.Response[domain.TaskWorker], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TaskWorker
	var response domain.Response[domain.TaskWorker]
	// var response domain.Response[domain.TaskWorker]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.TaskWorker]{}, err
	// }

	// cursor, err := r.db.Collection(TblTaskWorker).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.TaskWorker, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.TaskWorker]{}, err
	}
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblTask,
		"as":           "taska",
		"localField":   "taskId",
		"foreignField": "_id",
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"task": bson.M{"$first": "$taska"}}}})
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblTaskStatus,
		"as":           "taskStatusa",
		"localField":   "statusId",
		"foreignField": "_id",
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"taskStatus": bson.M{"$first": "$taskStatusa"}}}})
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": TblOrder,
		"as":   "ordera",
		"let":  bson.D{{Key: "orderId", Value: "$task.orderId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$orderId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"order": bson.M{"$first": "$ordera"}}}})

	cursor, err := r.db.Collection(tblTaskWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(tblTaskWorker).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.TaskWorker]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *TaskWorkerMongo) FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TaskWorker
	var response domain.Response[domain.TaskWorker]
	// var response domain.Response[domain.TaskWorker]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.TaskWorker]{}, err
	// }

	// cursor, err := r.db.Collection(TblTaskWorker).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.TaskWorker, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.TaskWorker]{}, err
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

	cursor, err := r.db.Collection(tblTaskWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(tblTaskWorker).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.TaskWorker]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *TaskWorkerMongo) GqlGetTaskWorkers(params domain.RequestParams) ([]*domain.TaskWorker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*domain.TaskWorker
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

	cursor, err := r.db.Collection(tblTaskWorker).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*domain.TaskWorker, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// count, err := r.db.Collection(TblTaskWorker).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return results, err
	// }

	// results = []*domain.TaskWorker{
	// 	Total: int(count),
	// 	Skip:  int(params.Options.Skip),
	// 	Limit: int(params.Options.Limit),
	// 	Data:  resultSlice,
	// }
	return results, nil
}

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
	// r.db.Collection(TblTaskWorker).FindOne(ctx, bson.M{"node_id": TaskWorker.NodeID, "user_id": userIDPrimitive}).Decode(&existTaskWorker)

	// if existTaskWorker.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	sortOrder := int64(0)

	from := time.Now()
	if !data.From.IsZero() {
		from = data.From
	}

	newTaskWorker := domain.TaskWorkerInput{
		UserID:    userIDPrimitive,
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

	res, err := collection.InsertOne(ctx, newTaskWorker)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblTaskWorker).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
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

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
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
