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

type TaskMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTaskMongo(db *mongo.Database, i18n config.I18nConfig) *TaskMongo {
	return &TaskMongo{db: db, i18n: i18n}
}

func (r *TaskMongo) FindTask(params domain.RequestParams) (domain.Response[domain.Task], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Task
	var response domain.Response[domain.Task]
	// var response domain.Response[domain.Task]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.Task]{}, err
	// }

	// cursor, err := r.db.Collection(TblTask).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.Task, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Task]{}, err
	}
	fmt.Println("FindTask params=", params)
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

	cursor, err := r.db.Collection(tblTask).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(tblTask).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Task]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *TaskMongo) FindTaskWithWorkers(params domain.RequestParams) (domain.Response[domain.Task], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Task
	var response domain.Response[domain.Task]

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Task]{}, err
	}
	fmt.Println("FindTaskWithWorkers params=", params)
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
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblTaskWorker,
		"as":           "workers",
		"localField":   "_id",
		"foreignField": "taskId",
	}}})

	cursor, err := r.db.Collection(tblTask).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(tblTask).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Task]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *TaskMongo) CreateTask(userID string, data *domain.Task) (*domain.Task, error) {
	var result *domain.Task

	collection := r.db.Collection(tblTask)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existTask domain.Task
	// r.db.Collection(TblTask).FindOne(ctx, bson.M{"node_id": Task.NodeID, "user_id": userIDPrimitive}).Decode(&existTask)

	// if existTask.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	// get sortOrder index.
	var allTaskByOrder []domain.Task
	cursor, err := collection.Find(ctx, bson.D{{"orderId", data.OrderId}})
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &allTaskByOrder); er != nil {
		return result, er
	}

	nextSortOrder := int64(len(allTaskByOrder))

	newTask := domain.TaskInput{
		OrderId: data.OrderId,
		UserID:  userIDPrimitive,
		// OperationId: data.OperationId,
		Name: data.Name,
		// WorkerId: data.WorkerId,
		SortOder:  &nextSortOrder,
		StatusId:  data.StatusId,
		Active:    data.Active,
		StartAt:   data.StartAt,
		AutoCheck: data.AutoCheck,
		Status:    data.Status,
		From:      data.From,
		To:        data.To,
		TypeGo:    data.TypeGo,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblTask).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	updatedAt := Task.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateTask := &domain.TaskInput{
	// 		Rate:      Task.Rate,
	// 		Task:    Task.Task,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateTask(existTask.ID.Hex(), userID, updateTask)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *TaskMongo) UpdateTask(id string, userID string, data *domain.TaskInput) (*domain.Task, error) {
	var result *domain.Task
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblTask)

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
	if data.Name != "" {
		newData["name"] = data.Name
	}
	if !data.OrderId.IsZero() {
		newData["orderId"] = data.OrderId
	}
	if !data.StatusId.IsZero() {
		newData["statusId"] = data.StatusId
	}
	if data.Active != nil {
		newData["active"] = data.Active
	}
	if data.AutoCheck != nil {
		newData["autoCheck"] = data.AutoCheck
	}
	newData["updatedAt"] = time.Now()
	if data.SortOder != nil {
		newData["sortOrder"] = data.SortOder
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

func (r *TaskMongo) DeleteTask(id string) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.Task{}
	collection := r.db.Collection(tblTask)

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
