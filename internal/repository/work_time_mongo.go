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

type WorkTimeMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewWorkTimeMongo(db *mongo.Database, i18n config.I18nConfig) *WorkTimeMongo {
	return &WorkTimeMongo{db: db, i18n: i18n}
}

func (r *WorkTimeMongo) FindWorkTime(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.WorkTime
	var response domain.Response[domain.WorkTime]

	q := bson.D{}

	// Filters
	if input.WorkerId != nil && len(input.WorkerId) > 0 {
		workerIds := []primitive.ObjectID{}
		for i, _ := range input.WorkerId {
			workerIDPrimitive, err := primitive.ObjectIDFromHex(*input.WorkerId[i])
			if err != nil {
				return response, err
			}

			workerIds = append(workerIds, workerIDPrimitive)
		}

		q = append(q, bson.E{"workerId", bson.D{{"$in", workerIds}}})
	}
	// if input.Group != nil && len(input.Group) > 0 {
	// 	q = append(q, bson.E{"group", bson.M{"$elemMatch": bson.D{{"$in", input.Group}}}})
	// }
	// if input.Status != nil {
	// 	q = append(q, bson.E{"status", input.Status})
	// }

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": tblObject,
	// 	"as":   "objecta",
	// 	// "localField":   "user_id",
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

	cursor, err := r.db.Collection(tblWorkTime).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblWorkTime).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.WorkTime]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *WorkTimeMongo) FindWorkTimePopulate(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.WorkTime
	var response domain.Response[domain.WorkTime]

	q := bson.D{}

	// Filters
	if input.WorkerId != nil && len(input.WorkerId) > 0 {
		workerIds := []primitive.ObjectID{}
		for i, _ := range input.WorkerId {
			workerIDPrimitive, err := primitive.ObjectIDFromHex(*input.WorkerId[i])
			if err != nil {
				return response, err
			}

			workerIds = append(workerIds, workerIDPrimitive)
		}

		q = append(q, bson.E{"workerId", bson.D{{"$in", workerIds}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": tblObject,
	// 	"as":   "objecta",
	// 	// "localField":   "user_id",
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

	cursor, err := r.db.Collection(tblWorkTime).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblWorkTime).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.WorkTime]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *WorkTimeMongo) CreateWorkTime(userID string, data *domain.WorkTime) (*domain.WorkTime, error) {
	var result *domain.WorkTime

	collection := r.db.Collection(tblWorkTime)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existTask domain.WorkTime
	// r.db.Collection(TblTask).FindOne(ctx, bson.M{"node_id": Task.NodeID, "user_id": userIDPrimitive}).Decode(&existTask)

	// if existTask.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	// // get sortOrder index.
	// var allTaskByOrder []domain.WorkTime
	// cursor, err := collection.Find(ctx, bson.D{{"orderId", data.OrderId}})
	// if err != nil {
	// 	return result, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &allTaskByOrder); er != nil {
	// 	return result, er
	// }

	newTask := domain.WorkTimeInput{
		// OrderId:  data.OrderId,
		// TaskId:   data.TaskId,
		WorkerId: data.WorkerId,
		UserID:   userIDPrimitive,
		Status:   data.Status,
		From:     data.From,
		To:       data.To,
		Date:     time.Now(),

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblWorkTime).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *WorkTimeMongo) UpdateWorkTime(id string, userID string, data *domain.WorkTimeInput) (*domain.WorkTime, error) {
	var result *domain.WorkTime
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblWorkTime)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	// if !data.OrderId.IsZero() {
	// 	newData["orderId"] = data.OrderId
	// }
	// if !data.TaskId.IsZero() {
	// 	newData["taskId"] = data.TaskId
	// }
	if !data.WorkerId.IsZero() {
		newData["workerId"] = data.WorkerId
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

func (r *WorkTimeMongo) DeleteWorkTime(id string) (*domain.WorkTime, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.WorkTime{}
	collection := r.db.Collection(tblWorkTime)

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
