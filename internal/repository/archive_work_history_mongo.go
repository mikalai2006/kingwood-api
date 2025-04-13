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

type ArchiveWorkHistoryMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchiveWorkHistoryMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveWorkHistoryMongo {
	return &ArchiveWorkHistoryMongo{db: db, i18n: i18n}
}

func (r *ArchiveWorkHistoryMongo) CreateArchiveWorkHistory(userID string, data *domain.WorkHistory) (*domain.ArchiveWorkHistory, error) {
	var result *domain.ArchiveWorkHistory

	collection := r.db.Collection(TblArchiveWorkHistory)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newTask := domain.ArchiveWorkHistoryInput{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:           data.ID,
		OrderId:      &data.OrderId,
		TaskId:       &data.TaskId,
		WorkerId:     data.WorkerId,
		ObjectId:     &data.ObjectId,
		OperationId:  &data.OperationId,
		UserID:       data.UserID,
		TaskWorkerId: &data.TaskWorkerId,
		Status:       &data.Status,
		Date:         data.Date,
		From:         data.From,
		To:           data.To,
		TotalTime:    data.TotalTime,
		Oklad:        data.Oklad,
		Total:        data.Total,
		Props:        data.Props,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}

	res, err := collection.InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblArchiveWorkHistory).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ArchiveWorkHistoryMongo) FindArchiveWorkHistory(input domain.ArchiveWorkHistoryFilter) (domain.Response[domain.ArchiveWorkHistory], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.ArchiveWorkHistory
	var response domain.Response[domain.ArchiveWorkHistory]

	q := bson.D{}

	// Filters
	if !input.From.IsZero() {
		q = append(q, bson.E{"from", bson.D{{"$gte", primitive.NewDateTimeFromTime(input.From)}}})
	}
	if !input.To.IsZero() {
		q = append(q, bson.E{"to", bson.D{{"$lte", primitive.NewDateTimeFromTime(input.To)}}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", input.Status})
	}
	if !input.Date.IsZero() {
		t := time.Time(input.Date)
		// year, month, day := t.Date()
		// from := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
		// to := time.Date(year, month, day, 23, 59, 59, 0, t.Location())

		// eastOfUTC := time.FixedZone("UTC-3", -3*60*60)
		eastOfUTCP := time.FixedZone("UTC+3", 3*60*60)
		from1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, eastOfUTCP)
		to1 := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, eastOfUTCP)

		// fmt.Println("======================FIND TIME WORK====================")
		// fmt.Println("date: ", t, "====>", t.UTC())
		// fmt.Println("from: ", from1, "====>", from1.UTC())
		// fmt.Println("to: ", to1, "====>", to1.UTC())
		// fmt.Println("========================================================")

		q = append(q, bson.E{"date", bson.D{{"$gte", primitive.NewDateTimeFromTime(from1.UTC())}}})
		q = append(q, bson.E{"date", bson.D{{"$lte", primitive.NewDateTimeFromTime(to1.UTC())}}})
	}
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

	cursor, err := r.db.Collection(TblArchiveWorkHistory).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
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

	response = domain.Response[domain.ArchiveWorkHistory]{
		Total: int(0),
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *ArchiveWorkHistoryMongo) DeleteArchiveWorkHistory(id string) (*domain.ArchiveWorkHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.ArchiveWorkHistory{}
	collection := r.db.Collection(TblArchiveWorkHistory)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	removeItems, err := r.FindArchiveWorkHistory(domain.ArchiveWorkHistoryFilter{ID: []string{id}})
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	if len(removeItems.Data) > 0 {
		result = &removeItems.Data[0]
	} else {
		err = collection.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			return result, err
		}
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ArchiveWorkHistoryMongo) ClearArchiveWorkHistory(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblWorkHistory)

	_, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}

	return nil
}
