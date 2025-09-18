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

type TimerMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTimerMongo(db *mongo.Database, i18n config.I18nConfig) *TimerMongo {
	return &TimerMongo{db: db, i18n: i18n}
}

func (r *TimerMongo) FindTimer(params domain.RequestParams) (domain.Response[domain.TimerShedule], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TimerShedule
	var response domain.Response[domain.TimerShedule]

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.TimerShedule]{}, err
	}

	cursor, err := r.db.Collection(tblTimer).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(tblTimer).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.TimerShedule]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *TimerMongo) FindTimerPopulate(input domain.TimerSheduleFilter) (domain.Response[domain.TimerShedule], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TimerShedule
	var response domain.Response[domain.TimerShedule]

	// Filters
	q := bson.D{}
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
	if input.IDTimer != nil && len(input.IDTimer) > 0 {
		ids := []string{}
		for i, _ := range input.IDTimer {
			ids = append(ids, input.IDTimer[i])
		}

		q = append(q, bson.E{"idTimer", bson.D{{"$in", ids}}})
	}
	if input.WorkerId != nil && len(input.WorkerId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.WorkerId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.WorkerId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"workerId", bson.D{{"$in", ids}}})
	}
	if input.WorkHistoryId != nil && len(input.WorkHistoryId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.WorkHistoryId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.WorkHistoryId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"workHistoryId", bson.D{{"$in", ids}}})
	}
	if input.TaskWorkerId != nil && len(input.TaskWorkerId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.TaskWorkerId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.TaskWorkerId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"taskWorkerId", bson.D{{"$in", ids}}})
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
	if input.IsRunning != nil {
		q = append(q, bson.E{"isRunning", input.IsRunning})
	}
	if !input.ExecuteAt.IsZero() {
		t := time.Time(input.ExecuteAt)
		eastOfUTCP := time.FixedZone("UTC+3", 3*60*60)
		from1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, eastOfUTCP)
		to1 := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, eastOfUTCP)
		q = append(q, bson.E{"date", bson.D{{"$gte", primitive.NewDateTimeFromTime(from1.UTC())}}})
		q = append(q, bson.E{"date", bson.D{{"$lte", primitive.NewDateTimeFromTime(to1.UTC())}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
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
		},
	}}},
		bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}})

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

	cursor, err := r.db.Collection(tblTimer).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblTimer).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.TimerShedule]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *TimerMongo) CreateTimer(userID string, data *domain.TimerShedule) (*domain.TimerShedule, error) {
	var result *domain.TimerShedule

	collection := r.db.Collection(tblTimer)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }

	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}
	isRunning := 1
	newTimer := domain.TimerSheduleInput{
		IDTimer:       data.IDTimer,
		WorkerId:      data.WorkerId,
		ExecuteAt:     data.ExecuteAt,
		TaskId:        data.TaskId,
		IsRunning:     &isRunning,
		TaskWorkerId:  data.TaskWorkerId,
		WorkHistoryId: data.WorkHistoryId,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newTimer)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblTimer).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TimerMongo) UpdateTimer(id string, userID string, data *domain.TimerSheduleInput) (*domain.TimerShedule, error) {
	var result *domain.TimerShedule
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblTimer)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if !data.ExecuteAt.IsZero() {
		newData["executeAt"] = data.ExecuteAt
	}

	newData["updatedAt"] = time.Now()

	if data.IsRunning != nil {
		newData["isRunning"] = data.IsRunning
	}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	timers, err := r.FindTimerPopulate(domain.TimerSheduleFilter{ID: []string{id}})
	// taskWorkers, err := r.FindTaskWorkerPopulate(domain.RequestParams{Filter: bson.D{{"_id", idPrimitive}}})
	// collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(timers.Data) > 0 {
		result = &timers.Data[0]
	} else {
		fmt.Println("Len timers.Data = ", len(timers.Data))
	}

	return result, nil
}

func (r *TimerMongo) DeleteTimer(id string) (*domain.TimerShedule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.TimerShedule{}
	collection := r.db.Collection(tblTimer)

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
