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

type NotifyMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewNotifyMongo(db *mongo.Database, i18n config.I18nConfig) *NotifyMongo {
	return &NotifyMongo{db: db, i18n: i18n}
}

func (r *NotifyMongo) FindNotifyPopulate(input *domain.NotifyFilter) (domain.Response[domain.Notify], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// var results []domain.Notify
	var response domain.Response[domain.Notify]

	// Filters
	q := bson.D{}
	if input.ID != nil && len(input.ID) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.ID {
			iDPrimitive, err := primitive.ObjectIDFromHex(*input.ID[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"_id", bson.D{{"$in", ids}}})
	}
	if input.UserID != nil && len(input.UserID) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.UserID {
			iDPrimitive, err := primitive.ObjectIDFromHex(*input.UserID[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"userId", bson.D{{"$in", ids}}})
	}
	if input.UserTo != nil && len(input.UserTo) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.UserTo {
			iDPrimitive, err := primitive.ObjectIDFromHex(*input.UserTo[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"userTo", bson.D{{"$in", ids}}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", input.Status})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	// user.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblUsers,
		"as":   "usera",
		"let":  bson.D{{Key: "userId", Value: "$userId"}},
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
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})

	// recepient.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblUsers,
		"as":   "recepienta",
		"let":  bson.D{{Key: "userTo", Value: "$userTo"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userTo"}}}}},
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
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"recepient": bson.M{"$first": "$recepienta"}}}})

	// // order.
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": TblOrder,
	// 	"as":   "ordera",
	// 	// "localField":   "userId",
	// 	// "foreignField": "_id",
	// 	"let": bson.D{{Key: "orderId", Value: "$orderId"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$orderId"}}}}},
	// 		bson.D{{"$limit", 1}},
	// 	},
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"order": bson.M{"$first": "$ordera"}}}})

	// // task.
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

	// // worker.
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": tblUsers,
	// 	"as":   "usera",
	// 	"let":  bson.D{{Key: "workerId", Value: "$workerId"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$workerId"}}}}},
	// 		bson.D{{"$limit", 1}},
	// 		bson.D{{
	// 			Key: "$lookup",
	// 			Value: bson.M{
	// 				"from": tblImage,
	// 				"as":   "images",
	// 				"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 				"pipeline": mongo.Pipeline{
	// 					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
	// 				},
	// 			},
	// 		}},

	// 		// post.
	// 		bson.D{{
	// 			Key: "$lookup",
	// 			Value: bson.M{
	// 				"from": TblPost,
	// 				"as":   "posts",
	// 				// "localField":   "_id",
	// 				// "foreignField": "serviceId",
	// 				"let": bson.D{{Key: "postId", Value: "$postId"}},
	// 				"pipeline": mongo.Pipeline{
	// 					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$postId"}}}}},
	// 				},
	// 			},
	// 		}},
	// 		bson.D{{Key: "$set", Value: bson.M{"postObject": bson.M{"$first": "$posts"}}}},
	// 		// role.
	// 		bson.D{{Key: "$lookup", Value: bson.M{
	// 			"from": TblRole,
	// 			"as":   "rolea",
	// 			// "localField":   "userId",
	// 			// "foreignField": "_id",
	// 			"let": bson.D{{Key: "roleId", Value: "$roleId"}},
	// 			"pipeline": mongo.Pipeline{
	// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$roleId"}}}}},
	// 				bson.D{{"$limit", 1}},
	// 			},
	// 		}}},
	// 		bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$rolea"}}}},
	// 	},
	// }}},
	// 	bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}},
	// )

	// // taskStatus.
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from":         tblTaskStatus,
	// 	"as":           "taskStatusa",
	// 	"localField":   "statusId",
	// 	"foreignField": "_id",
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"taskStatus": bson.M{"$first": "$taskStatusa"}}}})

	if input.Sort != nil && len(input.Sort) > 0 {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
		// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	} else {
		pipe = append(pipe, bson.D{{"$sort", bson.D{{"createdAt", -1}}}})
	}

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
	dataOptions = append(dataOptions, bson.D{{"$limit", limit}})

	if input.Sort != nil {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		dataOptions = append(dataOptions, bson.D{{"$sort", sortParam}})
	} else {
		dataOptions = append(dataOptions, bson.D{{"$sort", bson.D{{"createdAt", -1}}}})
		//pipe = append(pipe, bson.D{{"$sort", bson.D{{"createdAt", -1}}}})
	}

	// pipe = append(pipe, bson.D{{"$skip", skip}})
	// pipe = append(pipe, bson.D{{"$limit", limit}})

	pipe = append(pipe, bson.D{{Key: "$facet", Value: bson.D{
		{"data", dataOptions},
		{Key: "metadata", Value: mongo.Pipeline{
			bson.D{{"$group", bson.D{
				{"_id", nil},
				{"total", bson.D{{"$sum", 1}}}}}},
		}},
	},
	}})

	cursor, err := r.db.Collection(tblNotify).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	resultMap := []bson.M{}
	if er := cursor.All(ctx, &resultMap); er != nil {
		return response, er
	}
	resultFacetOne := domain.ResultFacetNotify{}
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

	// count, err := r.db.Collection(tblNotify).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }
	// fmt.Println("skip=", skip, " limit=", limit, " s+t=", skip+limit, resultFacetOne.Metadata)
	response = domain.Response[domain.Notify]{
		Total: total,
		Skip:  skip,
		Limit: limit,
		Data:  resultFacetOne.Data,
	}
	return response, nil
}

func (r *NotifyMongo) CreateNotify(userID string, data *domain.NotifyInput) (*domain.Notify, error) {
	var result *domain.Notify

	collection := r.db.Collection(tblNotify)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	userToPrimitive, err := primitive.ObjectIDFromHex(data.UserTo)
	if err != nil {
		return nil, err
	}
	newNotify := domain.NotifyInputMongo{
		UserID:     userIDPrimitive,
		UserTo:     userToPrimitive,
		Status:     0,
		Title:      data.Title,
		Message:    data.Message,
		Props:      data.Props,
		Images:     data.Images,
		Link:       data.Link,
		LinkOption: data.LinkOption,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newNotify)
	if err != nil {
		return nil, err
	}

	insertedId := res.InsertedID.(primitive.ObjectID).Hex()
	Notifys, err := r.FindNotifyPopulate(&domain.NotifyFilter{ID: []*string{&insertedId}})
	if err != nil {
		return nil, err
	}
	if len(Notifys.Data) > 0 {
		result = &Notifys.Data[0]
	}

	return result, nil
}

func (r *NotifyMongo) UpdateNotify(id string, userID string, data *domain.NotifyInput) (*domain.Notify, error) {
	var result *domain.Notify
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblNotify)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Status != nil {
		newData["status"] = &data.Status
		newData["readAt"] = time.Now()
	}
	newData["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	notifys, err := r.FindNotifyPopulate(&domain.NotifyFilter{ID: []*string{&id}})
	// Notifys, err := r.FindNotifyPopulate(domain.RequestParams{Filter: bson.D{{"_id", idPrimitive}}})
	// collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(notifys.Data) > 0 {
		result = &notifys.Data[0]
	} else {
		fmt.Println("Len notifys.Data = ", len(notifys.Data))
	}

	return result, nil
}

func (r *NotifyMongo) DeleteNotify(id string) (*domain.Notify, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.Notify{}
	collection := r.db.Collection(tblNotify)

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
