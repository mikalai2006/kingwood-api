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
	// 	// "localField":   "userId",
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
	// 				// "foreignField": "serviceId",
	// 				"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 				"pipeline": mongo.Pipeline{
	// 					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
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

func (r *TaskMongo) FindTaskPopulate(input domain.TaskFilter) (domain.Response[domain.Task], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Task
	var response domain.Response[domain.Task]

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
	// workers.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         tblTaskWorker,
		"as":           "workers",
		"localField":   "_id",
		"foreignField": "taskId",
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

	cursor, err := r.db.Collection(tblTask).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTask).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.Task]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *TaskMongo) FindTaskFlat(input domain.TaskFilter) (domain.Response[domain.Task], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Task
	var response domain.Response[domain.Task]

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
	if input.Limit != nil && *input.Limit > 0 {
		pipe = append(pipe, bson.D{{"$limit", input.Limit}})
		limit = *input.Limit
	}

	cursor, err := r.db.Collection(tblTask).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTask).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.Task]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
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
	// r.db.Collection(TblTask).FindOne(ctx, bson.M{"node_id": Task.NodeID, "userId": userIDPrimitive}).Decode(&existTask)

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
	from := time.Now()
	if data.From != nil && !data.From.IsZero() {
		from = *data.From
	}
	to := time.Now().AddDate(1, 0, 0)
	if data.To != nil && !data.From.IsZero() {
		to = *data.To
	}

	newTask := domain.TaskInput{
		OrderId:     data.OrderId,
		UserID:      userIDPrimitive,
		OperationId: data.OperationId,
		Name:        data.Name,
		// WorkerId: data.WorkerId,
		SortOrder: &nextSortOrder,
		StatusId:  data.StatusId,
		Active:    data.Active,
		StartAt:   data.StartAt,
		AutoCheck: data.AutoCheck,
		Status:    data.Status,
		From:      from,
		To:        to,
		TypeGo:    data.TypeGo,
		MaxHours:  &data.MaxHours,
		ObjectId:  data.ObjectId,

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
	if !data.OperationId.IsZero() {
		newData["operationId"] = data.OperationId
	}
	if data.Name != "" {
		newData["name"] = data.Name
	}
	if !data.ObjectId.IsZero() {
		newData["objectId"] = data.ObjectId
	}
	if !data.OrderId.IsZero() {
		newData["orderId"] = data.OrderId
	}
	if data.MaxHours != nil {
		newData["maxHours"] = data.MaxHours
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
	if data.SortOrder != nil {
		newData["sortOrder"] = data.SortOrder
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

	// err = collection.FindOne(ctx, filter).Decode(&result)
	// if err != nil {
	// 	return result, err
	// }
	tasks, err := r.FindTaskPopulate(domain.TaskFilter{ID: []string{id}})
	// taskWorkers, err := r.FindTaskWorkerPopulate(domain.RequestParams{Filter: bson.D{{"_id", idPrimitive}}})
	// collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(tasks.Data) > 0 {
		result = &tasks.Data[0]
	} else {
		fmt.Println("Len tasks.Data = ", len(tasks.Data))
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
