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

type TaskMontajMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTaskMontajMongo(db *mongo.Database, i18n config.I18nConfig) *TaskMontajMongo {
	return &TaskMontajMongo{db: db, i18n: i18n}
}

func (r *TaskMontajMongo) FindTaskMontaj(input domain.TaskMontajFilter) (domain.Response[domain.TaskMontaj], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TaskMontaj
	var response domain.Response[domain.TaskMontaj]
	// var response domain.Response[domain.TaskMontaj]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.TaskMontaj]{}, err
	// }
	// cursor, err := r.db.Collection(TblTask).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)
	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }
	// resultSlice := make([]domain.TaskMontaj, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)
	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.TaskMontaj]{}, err
	// }
	// fmt.Println("FindTaskMontaj params=", params)
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
	// if input.Group != nil && len(input.Group) > 0 {
	// 	q = append(q, bson.E{"group", bson.M{"$elemMatch": bson.D{{"$in", input.Group}}}})
	// }
	// if input.Status != nil {
	// 	q = append(q, bson.E{"status", input.Status})
	// }

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

	cursor, err := r.db.Collection(tblTaskMontaj).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTaskMontaj).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.TaskMontaj]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *TaskMontajMongo) FindTaskPopulate(input domain.TaskMontajFilter) (domain.Response[domain.TaskMontaj], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TaskMontaj
	var response domain.Response[domain.TaskMontaj]

	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.TaskMontaj]{}, err
	// }
	// fmt.Println("FindTaskMontajWithWorkers params=", input)
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
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from":         tblTaskWorker,
	// 	"as":           "workers",
	// 	"localField":   "_id",
	// 	"foreignField": "taskId",
	// }}})

	q := bson.D{}

	// Filters
	if input.Name != nil && *input.Name != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", *input.Name), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
	}
	// if input.Group != nil && len(input.Group) > 0 {
	// 	q = append(q, bson.E{"group", bson.M{"$elemMatch": bson.D{{"$in", input.Group}}}})
	// }
	// if input.Status != nil {
	// 	q = append(q, bson.E{"status", input.Status})
	// }

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

	cursor, err := r.db.Collection(tblTaskMontaj).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTaskMontaj).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.TaskMontaj]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *TaskMontajMongo) CreateTaskMontaj(userID string, data *domain.TaskMontaj) (*domain.TaskMontaj, error) {
	var result *domain.TaskMontaj

	collection := r.db.Collection(tblTaskMontaj)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existTask domain.TaskMontaj
	// r.db.Collection(TblTask).FindOne(ctx, bson.M{"node_id": Task.NodeID, "user_id": userIDPrimitive}).Decode(&existTask)

	// if existTask.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	// // get sortOrder index.
	// var allTaskByOrder []domain.TaskMontaj
	// cursor, err := collection.Find(ctx, bson.D{{"orderId", data.OrderId}})
	// if err != nil {
	// 	return result, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &allTaskByOrder); er != nil {
	// 	return result, er
	// }

	nextSortOrder := int64(0) //int64(len(allTaskByOrder))

	from := *data.From
	if !data.From.IsZero() {
		from = time.Now()
	}

	newTask := domain.TaskMontajInput{
		UserID:      userIDPrimitive,
		ObjectId:    data.ObjectId,
		OperationId: data.OperationId,
		Name:        data.Name,
		SortOder:    &nextSortOrder,
		StatusId:    data.StatusId,
		Status:      data.Status,
		From:        from,
		To:          data.To,
		TypeGo:      data.TypeGo,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblTaskMontaj).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	updatedAt := Task.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateTask := &domain.TaskMontajInput{
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

func (r *TaskMontajMongo) UpdateTaskMontaj(id string, userID string, data *domain.TaskMontajInput) (*domain.TaskMontaj, error) {
	var result *domain.TaskMontaj
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblTaskMontaj)

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
	if !data.StatusId.IsZero() {
		newData["statusId"] = data.StatusId
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

func (r *TaskMontajMongo) DeleteTaskMontaj(id string) (*domain.TaskMontaj, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.TaskMontaj{}
	collection := r.db.Collection(tblTaskMontaj)

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
