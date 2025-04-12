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

type ArchiveUserMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchiveUserMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveUserMongo {
	return &ArchiveUserMongo{db: db, i18n: i18n}
}

func (r *ArchiveUserMongo) FindArchiveUser(input *domain.ArchiveUserFilter) (domain.Response[domain.ArchiveUser], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.ArchiveUser
	var response domain.Response[domain.ArchiveUser]
	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.ArchiveUser]{}, err
	// }
	// fmt.Println("params:::", params)

	q := bson.D{}

	// Filters
	if input.Archive != nil {
		q = append(q, bson.E{"archive", input.Archive})
	}
	if input.Blocked != nil {
		q = append(q, bson.E{"blocked", input.Blocked})
	}
	if input.Hidden != nil {
		q = append(q, bson.E{"hidden", input.Hidden})
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

		q = append(q, bson.E{"ArchiveUserId", bson.D{{"$in", ids}}})
	}
	if input.RoleId != nil && len(input.RoleId) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.RoleId {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.RoleId[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"roleId", bson.D{{"$in", ids}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": tblImage,
			"as":   "images",
			// "localField":   "_id",
			// "foreignField": "serviceId",
			"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
			},
		},
	}})

	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         TblRole,
			"as":           "roles",
			"localField":   "roleId",
			"foreignField": "_id",
			// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
			// "pipeline": mongo.Pipeline{
			// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
			// },
		},
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$roles"}}}})

	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         TblPost,
			"as":           "posts",
			"localField":   "postId",
			"foreignField": "_id",
			// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
			// "pipeline": mongo.Pipeline{
			// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
			// },
		},
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"postObject": bson.M{"$first": "$posts"}}}})

	// add populate auth.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         TblAuth,
			"as":           "auths",
			"localField":   "ArchiveUserId",
			"foreignField": "_id",
			// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
			// "pipeline": mongo.Pipeline{
			// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
			// },
		},
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"auth": bson.M{"$first": "$auths"}}}})

	// workHistorys.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblWorkHistory,
		"as":   "workHistorys",
		"let":  bson.D{{Key: "id", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{
				{
					Key: "$match",
					Value: bson.M{
						"$expr": bson.D{
							{"$and", bson.A{
								bson.M{"$eq": bson.A{"$status", 0}},
								bson.M{"$eq": [2]string{"$workerId", "$$id"}},
								// bson.D{{"$lte", bson.A{"$to", "$$to"}}},
								// bson.D{{"$gte", bson.A{"$from", "$$from"}}},
							}},
						}},
					// Value: bson.M{
					// 	"$expr": bson.M{"$eq": [2]string{"$workerId", "$$id"}},
					// },
				},
			},
		},
	}}})

	// tasks.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblTaskWorker,
		"as":   "taskWorkers",
		"let":  bson.D{{Key: "id", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{
				"$expr":  bson.M{"$eq": [2]string{"$workerId", "$$id"}},
				"status": bson.D{{"$nin", [2]string{"finish", "autofinish"}}},
			}}},

			// object.
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": tblObject,
				"as":   "objecta",
				"let":  bson.D{{Key: "objectId", Value: "$objectId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$objectId"}}}}},
					bson.D{{"$limit", 1}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"object": bson.M{"$first": "$objecta"}}}},

			// order.
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": TblOrder,
				"as":   "ordera",
				"let":  bson.D{{Key: "orderId", Value: "$orderId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$orderId"}}}}},
					bson.D{{"$limit", 1}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"order": bson.M{"$first": "$ordera"}}}},

			// task.
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": TblArchiveTask,
				"as":   "taska",
				"let":  bson.D{{Key: "taskId", Value: "$taskId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$taskId"}}}}},
					bson.D{{"$limit", 1}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"task": bson.M{"$first": "$taska"}}}},

			// worker.
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": TblArchiveUser,
				"as":   "ArchiveUsera",
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
							"localField":   "ArchiveUserId",
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
						// "localField":   "ArchiveUserId",
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
			bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$ArchiveUsera"}}}},
		},
	}}})

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

	cursor, err := r.db.Collection(TblArchiveUser).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.ArchiveUser, len(results))
	copy(resultSlice, results)

	// устанавливаем работает ли пользователь.
	for i := range resultSlice {
		if len(resultSlice[i].WorkHistorys) > 0 {
			resultSlice[i].IsWork = 1
		}
	}

	count, err := r.db.Collection(TblArchiveUser).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.ArchiveUser]{
		Total: int(count),
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ArchiveUserMongo) CreateArchiveUser(userID string, data *domain.User) (*domain.ArchiveUser, error) {
	var result *domain.ArchiveUser

	collection := r.db.Collection(TblArchiveUser)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var auth domain.Auth
	err = r.db.Collection(TblAuth).FindOne(ctx, bson.D{{"_id", data.UserID}}).Decode(&auth)
	if err != nil {
		return result, err
	}

	// remove auth.
	_, err = r.db.Collection(TblAuth).DeleteOne(ctx, bson.D{{"_id", data.UserID}})
	if err != nil {
		return result, err
	}

	newArchiveUser := domain.ArchiveUserInputMongo{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:        data.ID,
		Name:      data.Name,
		UserID:    data.UserID,
		Phone:     data.Phone,
		Hidden:    data.Hidden,
		RoleId:    data.RoleId,
		PostId:    data.PostId,
		Archive:   data.Archive,
		TypeWork:  data.TypeWork,
		Birthday:  data.Birthday,
		TypePay:   data.TypePay,
		Oklad:     data.Oklad,
		Blocked:   data.Blocked,
		LastTime:  data.LastTime,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Props:     data.Props,
		Auth:      auth,
	}

	res, err := collection.InsertOne(ctx, newArchiveUser)
	if err != nil {
		return nil, err
	}

	// err = r.db.Collection(tblArchiveUsers).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	insertedID := res.InsertedID.(primitive.ObjectID).Hex()
	archiveUser, err := r.FindArchiveUser(&domain.ArchiveUserFilter{ID: []string{insertedID}})
	if err != nil {
		return nil, err
	}
	if len(archiveUser.Data) > 0 {
		result = &archiveUser.Data[0]
	}

	return result, nil
}

func (r *ArchiveUserMongo) DeleteArchiveUser(id string) (domain.ArchiveUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.ArchiveUser{}
	collection := r.db.Collection(TblArchiveUser)

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

	// // remove auth.
	// _, err = r.db.Collection(TblAuth).DeleteOne(ctx, bson.D{{"_id", result.UserID}})
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}
