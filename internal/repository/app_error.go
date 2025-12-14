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

type AppErrorMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewAppErrorMongo(db *mongo.Database, i18n config.I18nConfig) *AppErrorMongo {
	return &AppErrorMongo{db: db, i18n: i18n}
}

func (r *AppErrorMongo) FindAppError(input *domain.AppErrorFilter) (domain.Response[domain.AppError], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.AppError
	var response domain.Response[domain.AppError]

	// Filters
	q := bson.D{}
	if input.Code != "" {
		q = append(q, bson.E{"code", input.Code})
	}
	if input.From != nil && !input.From.IsZero() {
		q = append(q, bson.E{"createdAt", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.From)}}})
	}
	if input.To != nil && !input.To.IsZero() {
		q = append(q, bson.E{"createdAt", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.To)}}})
	}
	if input.Error != "" {
		q = append(q, bson.E{"error", bson.E{"$regex", input.Error}})
	}
	if input.Stack != "" {
		q = append(q, bson.E{"stack", bson.E{"$regex", input.Stack}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", *input.Status})
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
	if input.UserId != nil && len(input.UserId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.UserId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.UserId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"_id", bson.D{{"$in", ids}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	// user.
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
	// 		// add populate auth.
	// 		bson.D{{
	// 			Key: "$lookup",
	// 			Value: bson.M{
	// 				"from":         TblAuth,
	// 				"as":           "auths",
	// 				"localField":   "userId",
	// 				"foreignField": "_id",
	// 				// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
	// 				// "pipeline": mongo.Pipeline{
	// 				// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
	// 				// },
	// 			},
	// 		}},
	// 		bson.D{{Key: "$set", Value: bson.M{"auth": bson.M{"$first": "$auths"}}}},
	// 		bson.D{{Key: "$set", Value: bson.M{"authPrivate": bson.M{"$first": "$auths"}}}},

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

	cursor, err := r.db.Collection(tblAppError).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	response = domain.Response[domain.AppError]{
		Total: 0,
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *AppErrorMongo) CreateAppError(userID string, data *domain.AppError) (*domain.AppError, error) {
	var result *domain.AppError

	collection := r.db.Collection(tblAppError)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existAppError domain.AppError
	// r.db.Collection(TblAppError).FindOne(ctx, bson.M{"node_id": AppError.NodeID, "userId": userIDPrimitive}).Decode(&existAppError)

	// if existAppError.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	newAppError := domain.AppErrorInput{
		UserID: userIDPrimitive,
		Error:  data.Error,
		Status: &data.Status,
		Code:   data.Code,
		Stack:  data.Stack,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newAppError)
	if err != nil {
		return nil, err
	}

	// err = r.db.Collection(tblAppError).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	// if err != nil {
	// 	return nil, err
	// }

	insertedID := res.InsertedID.(primitive.ObjectID)
	AppErrors, err := r.FindAppError(&domain.AppErrorFilter{ID: []string{insertedID.Hex()}})
	if err != nil {
		return nil, err
	}
	if len(AppErrors.Data) > 0 {
		result = &AppErrors.Data[0]
	}

	return result, nil
}

func (r *AppErrorMongo) UpdateAppError(id string, userID string, data *domain.AppErrorInput) (*domain.AppError, error) {
	var result *domain.AppError
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblAppError)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Error != "" {
		newData["error"] = data.Error
	}

	if data.Status != nil {
		newData["status"] = &data.Status
	}
	newData["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	AppErrors, err := r.FindAppError(&domain.AppErrorFilter{ID: []string{id}})
	if err != nil {
		return nil, err
	}
	if len(AppErrors.Data) > 0 {
		result = &AppErrors.Data[0]
	}

	return result, nil
}

func (r *AppErrorMongo) DeleteAppError(id string, userID string) (*domain.AppError, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.AppError{}
	collection := r.db.Collection(tblAppError)

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

func (r *AppErrorMongo) DeleteAppErrorList(query domain.AppErrorListQuery) (*[]domain.AppError, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &[]domain.AppError{}
	collection := r.db.Collection(tblAppError)

	ids := bson.A{}

	if len(query.ID) > 0 {
		for i, _ := range query.ID {
			idPrimitive, err := primitive.ObjectIDFromHex(*query.ID[i])
			if err != nil {
				return result, err
			}
			ids = append(ids, idPrimitive)
		}
	}

	filter := bson.M{"_id": bson.D{{"$in", ids}}}

	// cursor, err := collection.Find(ctx, filter)
	// if err != nil {
	// 	return result, err
	// }
	// defer cursor.Close(ctx)

	// if err = cursor.All(ctx, &result); err != nil {
	// 	return result, err
	// }

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *AppErrorMongo) ClearAppError(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblAppError)

	_, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}

	return nil
}
