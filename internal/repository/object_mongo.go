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

type ObjectMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewObjectMongo(db *mongo.Database, i18n config.I18nConfig) *ObjectMongo {
	return &ObjectMongo{db: db, i18n: i18n}
}

func (r *ObjectMongo) FindObject(input *domain.ObjectFilter) (domain.Response[domain.Object], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	results := []domain.Object{}
	var response domain.Response[domain.Object]
	// var response domain.Response[domain.Object]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.Object]{}, err
	// }

	// cursor, err := r.db.Collection(TblObject).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.Object, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.Object]{}, err
	// }

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

	q := bson.D{}

	// Filter by substring name
	if input.Name != nil && *input.Name != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", *input.Name), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
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
	if input.Limit != nil {
		pipe = append(pipe, bson.D{{"$limit", input.Limit}})
		limit = *input.Limit
	}

	cursor, err := r.db.Collection(tblObject).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(tblObject).CountDocuments(ctx, pipe)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.Object]{
		Total: int(0),
		Skip:  int(skip),
		Limit: int(limit),
		Data:  results,
	}
	return response, nil
}

func (r *ObjectMongo) FindObjectWithWorkers(params domain.RequestParams) (domain.Response[domain.Object], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Object
	var response domain.Response[domain.Object]

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Object]{}, err
	}
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
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from":         tblObjectWorker,
	// 	"as":           "workers",
	// 	"localField":   "_id",
	// 	"foreignField": "ObjectId",
	// }}})

	cursor, err := r.db.Collection(tblObject).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(tblObject).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Object]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *ObjectMongo) CreateObject(userID string, data *domain.Object) (*domain.Object, error) {
	var result *domain.Object

	collection := r.db.Collection(tblObject)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existObject domain.Object
	// r.db.Collection(TblObject).FindOne(ctx, bson.M{"node_id": Object.NodeID, "userId": userIDPrimitive}).Decode(&existObject)

	// if existObject.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	newObject := domain.ObjectInput{
		UserID: userIDPrimitive,
		Name:   data.Name,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newObject)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblObject).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	updatedAt := Object.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateObject := &domain.ObjectInput{
	// 		Rate:      Object.Rate,
	// 		Object:    Object.Object,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateObject(existObject.ID.Hex(), userID, updateObject)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *ObjectMongo) UpdateObject(id string, userID string, data *domain.ObjectInput) (*domain.Object, error) {
	var result *domain.Object
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblObject)

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

func (r *ObjectMongo) DeleteObject(id string) (*domain.Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.Object{}
	collection := r.db.Collection(tblObject)

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
