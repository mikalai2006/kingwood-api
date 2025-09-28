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

type UserMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewUserMongo(db *mongo.Database, i18n config.I18nConfig) *UserMongo {
	return &UserMongo{db: db, i18n: i18n}
}

func (r *UserMongo) Iam(userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.User
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, err
	}

	params := domain.RequestParams{}
	params.Filter = bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, params.Filter).Decode(&result)
	if err != nil {
		return domain.User{}, err
	}

	pipe, err := CreatePipeline(params, &r.i18n) // mongo.Pipeline{bson.D{{"_id", userIDPrimitive}}} //
	if err != nil {
		return result, err
	}

	// add populate.
	// images.
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
	// post.
	pipe = append(pipe, bson.D{{
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
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"postObject": bson.M{"$first": "$posts"}}}})
	// role.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": TblRole,
		"as":   "rolea",
		// "localField":   "userId",
		// "foreignField": "_id",
		"let": bson.D{{Key: "roleId", Value: "$roleId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$roleId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$rolea"}}}})
	// // new roles.
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": TblRole,
	// 	"as":   "roles",
	// 	// "localField":   "userId",
	// 	// "foreignField": "_id",
	// 	"let": bson.D{{Key: "rolesId", Value: "$rolesId"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$in": [2]string{"$_id", "$$rolesId"}}}}},
	// 	},
	// }}})

	// add populate auth.
	pipe = append(pipe, bson.D{{
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
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"auth": bson.M{"$first": "$auths"}}}})
	// add populate.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblRole,
	// 		"as":   "roles",
	// 		// "localField":   "_id",
	// 		// "foreignField": "serviceId",
	// 		"let": bson.D{{Key: "userId", Value: bson.D{{"$toString", "$_id"}}}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$userId", "$$userId"}}}}},
	// 		},
	// 	},
	// }})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // .FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if er := cursor.Decode(&result); er != nil {
			return result, er
		}
	}

	return result, nil
}

func (r *UserMongo) GetUser(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.User

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.User{}, err
	}

	pipe, err := CreatePipeline(domain.RequestParams{
		Filter: filter,
	}, &r.i18n)
	if err != nil {
		return result, err
	}

	// add populate.
	// post.
	pipe = append(pipe, bson.D{{
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
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"postObject": bson.M{"$first": "$posts"}}}})
	// role.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": TblRole,
		"as":   "rolea",
		// "localField":   "userId",
		// "foreignField": "_id",
		"let": bson.D{{Key: "roleId", Value: "$roleId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$roleId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$rolea"}}}})

	// add populate auth.
	pipe = append(pipe, bson.D{{
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
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"auth": bson.M{"$first": "$auths"}}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"authPrivate": bson.M{"$first": "$auths"}}}})
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": tblImage,
	// 		"as":   "images",
	// 		// "localField":   "_id",
	// 		// "foreignField": "serviceId",
	// 		"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
	// 		},
	// 	},
	// }})

	// // add populate.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": tblImage,
	// 		"as":   "roles",
	// 		// "localField":   "_id",
	// 		// "foreignField": "serviceId",
	// 		"let": bson.D{{Key: "userId", Value: bson.D{{"$toString", "$_id"}}}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$userId", "$$userId"}}}}},
	// 		},
	// 	},
	// }})

	// // add populate.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": tblImage,
	// 		"as":   "post",
	// 		// "localField":   "_id",
	// 		// "foreignField": "serviceId",
	// 		"let": bson.D{{Key: "userId", Value: bson.D{{"$toString", "$_id"}}}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$userId", "$$userId"}}}}},
	// 		},
	// 	},
	// }})

	// // add populate.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblAuth,
	// 		"as":   "authsx",
	// 		// "localField":   "_id",
	// 		// "foreignField": "serviceId",
	// 		"let": bson.D{{Key: "userId", Value: "$userId"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
	// 			bson.D{{"$limit", 1}},
	// 		},
	// 	},
	// }})

	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"md": "$test.max_distance"}}})
	// // pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.M{"$first": "$authsx"}}}})
	// // pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"roles": "$test.roles"}}})

	// // add stat user tag vote.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodedataVote,
	// 		"as":   "tests",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$userId", "$$userId"}}}}},
	// 			bson.D{
	// 				{
	// 					"$group", bson.D{
	// 						{
	// 							"_id", "",
	// 						},
	// 						{"valueTagLike", bson.D{{"$sum", "$value"}}},
	// 						{"countTagLike", bson.D{{"$sum", 1}}},
	// 					},
	// 				},
	// 			},
	// 			bson.D{{Key: "$project", Value: bson.M{"_id": 0, "valueTagLike": "$valueTagLike", "countTagLike": "$countTagLike"}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.M{"$first": "$tests"}}}})

	// // add stat user node votes.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodeVote,
	// 		"as":   "tests2",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$userId", "$$userId"}}}}},
	// 			bson.D{
	// 				{
	// 					"$group", bson.D{
	// 						{
	// 							"_id", "",
	// 						},
	// 						{"valueNodeLike", bson.D{{"$sum", "$value"}}},
	// 						{"countNodeLike", bson.D{{"$sum", 1}}},
	// 					},
	// 				},
	// 			},
	// 			bson.D{{Key: "$project", Value: bson.M{"_id": 0, "valueNodeLike": "$valueNodeLike", "countNodeLike": "$countNodeLike"}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"$first": "$tests2"},
	// 	},
	// }},
	// }}})

	// // add stat user node votes.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNode,
	// 		"as":   "countNodes",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$userId", "$$userId"}}}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"countNodes": bson.M{"$size": "$countNodes"}},
	// 	},
	// }},
	// }}})

	// // add stat user added nodedata.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodedata,
	// 		"as":   "countNodedatas",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$userId", "$$userId"}}}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"countNodedatas": bson.M{"$size": "$countNodedatas"}},
	// 	},
	// }},
	// }}})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if er := cursor.Decode(&result); er != nil {
			return result, er
		}
	}

	return result, nil
}

func (r *UserMongo) FindUser(input *domain.UserFilter) (domain.Response[domain.User], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.User
	var response domain.Response[domain.User]
	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.User]{}, err
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

		q = append(q, bson.E{"userId", bson.D{{"$in", ids}}})
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
			"localField":   "userId",
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
				"from": tblTask,
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

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.User, len(results))
	copy(resultSlice, results)

	// устанавливаем работает ли пользователь.
	for i := range resultSlice {
		if len(resultSlice[i].WorkHistorys) > 0 {
			resultSlice[i].IsWork = 1
		}
	}

	count, err := r.db.Collection(tblUsers).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.User]{
		Total: int(count),
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *UserMongo) CreateUser(userID string, data *domain.User) (*domain.User, error) {
	var result *domain.User

	collection := r.db.Collection(tblUsers)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	typeWork := []string{}
	if len(data.TypeWork) > 0 {
		typeWork = data.TypeWork
	}

	blockedValue := 0
	newUser := domain.UserInputMongo{
		// Avatar: user.Avatar,
		Name:   data.Name,
		UserID: userIDPrimitive,
		// Online:    user.Online,
		Phone:  data.Phone,
		Hidden: 0,
		RoleId: data.RoleId,
		// RolesId:   data.RolesId,
		PostId:    data.PostId,
		Archive:   data.Archive,
		TypeWork:  typeWork,
		Birthday:  data.Birthday,
		TypePay:   data.TypePay,
		Oklad:     data.Oklad,
		MaxTime:   data.MaxTime,
		Blocked:   &blockedValue,
		LastTime:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// err = r.db.Collection(tblUsers).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	insertedID := res.InsertedID.(primitive.ObjectID).Hex()
	user, err := r.GetUser(insertedID)
	if err != nil {
		return nil, err
	}
	result = &user

	return result, nil
}

func (r *UserMongo) DeleteUser(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.User{}
	collection := r.db.Collection(tblUsers)

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

func (r *UserMongo) UpdateUser(id string, data *domain.UserInput) (domain.User, error) {
	var result domain.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblUsers)

	// data, err := utils.GetBodyToData(user)
	// if err != nil {
	// 	return result, err
	// }

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}

	if data.Name != "" {
		newData["name"] = data.Name
	}
	if data.Phone != "" {
		newData["phone"] = data.Phone
	}
	if data.Birthday != nil {
		newData["birthday"] = data.Birthday
	}
	if data.Oklad != nil {
		newData["oklad"] = data.Oklad
	}
	if data.MaxTime != nil {
		newData["maxTime"] = data.MaxTime
	}
	if data.TypePay != nil {
		newData["typePay"] = data.TypePay
	}
	if data.Archive != nil {
		newData["archive"] = data.Archive
	}
	if data.Hidden != nil {
		newData["hidden"] = data.Hidden
	}
	if data.Blocked != nil {
		newData["blocked"] = data.Blocked
	}
	if data.Dops != nil {
		newData["dops"] = data.Dops
	}
	if data.RoleId != "" {
		IDPrimitive, err := primitive.ObjectIDFromHex(data.RoleId)
		if err != nil {
			return result, err
		}
		newData["roleId"] = IDPrimitive
	}
	// if data.RolesId != nil {
	// 	IDsPrimitive := []primitive.ObjectID{}
	// 	for i, _ := range data.RolesId {
	// 		IDPrimitive, err := primitive.ObjectIDFromHex(data.RolesId[i])
	// 		if err != nil {
	// 			return result, err
	// 		}
	// 		IDsPrimitive = append(IDsPrimitive, IDPrimitive)
	// 	}
	// 	newData["rolesId"] = IDsPrimitive
	// }
	// fmt.Println("user.TypeWork=", user.TypeWork)
	if data.TypeWork != nil {
		newData["typeWork"] = data.TypeWork
	}
	//  else {
	// 	newData["typeWork"] = []string{}
	// }

	if !data.LastTime.IsZero() {
		newData["lastTime"] = data.LastTime
	}
	if data.Online != nil {
		newData["online"] = data.Online
	}

	if data.PostId != "" {
		IDPrimitive, err := primitive.ObjectIDFromHex(data.PostId)
		if err != nil {
			return result, err
		}
		newData["postId"] = IDPrimitive
	}

	newData["updatedAt"] = time.Now()

	// fmt.Println("data=", user)
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	limit := 1
	results, err := r.FindUser(&domain.UserFilter{ID: []string{id}, Limit: &limit})
	// domain.RequestParams{Filter: bson.M{"_id": idPrimitive}, Options: domain.Options{Limit: 1}})
	if err != nil {
		return result, err
	}

	if len(results.Data) > 0 {
		result = results.Data[0]
	}

	return result, nil
}

func (r *UserMongo) GqlGetUsers(params domain.RequestParams) ([]*domain.User, error) {
	fmt.Println("GqlGetUsers")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*domain.User
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

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

	// fmt.Println(pipe)

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*domain.User, len(results))

	copy(resultSlice, results)
	return results, nil
}

// func (r *UserMongo) SetStat(userID string, inputData domain.UserStat) (domain.User, error) {
// 	var result domain.User
// 	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
// 	defer cancel()

// 	collection := r.db.Collection(tblUsers)

// 	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		return result, err
// 	}

// 	filter := bson.M{"_id": userIDPrimitive}

// 	err = collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		return result, err
// 	}

// 	newData := bson.M{}
// 	if inputData.AddProduct != 0 {
// 		newData["user_stat.addProduct"] = utils.Max(result.UserStat.AddProduct+inputData.AddProduct, 0)
// 	}
// 	if inputData.TakeProduct != 0 {
// 		newData["user_stat.takeProduct"] = utils.Max(result.UserStat.TakeProduct+inputData.TakeProduct, 0)
// 	}
// 	if inputData.GiveProduct != 0 {
// 		newData["user_stat.giveProduct"] = utils.Max(result.UserStat.GiveProduct+inputData.GiveProduct, 0)
// 	}
// 	if inputData.AddOffer != 0 {
// 		newData["user_stat.addOffer"] = utils.Max(result.UserStat.AddOffer+inputData.AddOffer, 0)
// 	}
// 	if inputData.TakeOffer != 0 {
// 		newData["user_stat.takeOffer"] = utils.Max(result.UserStat.TakeOffer+inputData.TakeOffer, 0)
// 	}
// 	if inputData.AddMessage != 0 {
// 		newData["user_stat.addMessage"] = utils.Max(result.UserStat.AddMessage+inputData.AddMessage, 0)
// 	}
// 	if inputData.TakeMessage != 0 {
// 		newData["user_stat.takeMessage"] = utils.Max(result.UserStat.TakeMessage+inputData.TakeMessage, 0)
// 	}
// 	if inputData.AddReview != 0 {
// 		newData["user_stat.addReview"] = utils.Max(result.UserStat.AddReview+inputData.AddReview, 0)
// 	}
// 	if inputData.TakeReview != 0 {
// 		newData["user_stat.takeReview"] = utils.Max(result.UserStat.TakeReview+inputData.TakeReview, 0)
// 	}
// 	if inputData.Warning != 0 {
// 		newData["user_stat.warning"] = utils.Max(result.UserStat.Warning+inputData.Warning, 0)
// 	}
// 	if inputData.Request != 0 {
// 		newData["user_stat.request"] = utils.Max(result.UserStat.Request+inputData.Request, 0)
// 	}
// 	if inputData.Subscribe != 0 {
// 		newData["user_stat.subscribe"] = utils.Max(result.UserStat.Subscribe+inputData.Subscribe, 0)
// 	}
// 	if inputData.Subscriber != 0 {
// 		newData["user_stat.subscriber"] = utils.Max(result.UserStat.Subscriber+inputData.Subscriber, 0)
// 	}
// 	if !inputData.LastRequest.IsZero() {
// 		newData["user_stat.lastRequest"] = result.UserStat.LastRequest
// 	}

// 	// fmt.Println("newData=", newData)
// 	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": newData}).Decode(&result)
// 	if err != nil {
// 		return result, err
// 	}

// 	// var operations []mongo.WriteModel
// 	// operationA := mongo.NewUpdateOneModel()
// 	// operationA.SetFilter(bson.M{"_id": userIDPrimitive})
// 	// operationA.SetUpdate(bson.D{
// 	// 	{"$inc", bson.D{
// 	// 		{"user_stat.node", 1},
// 	// 	}},
// 	// })
// 	// operations = append(operations, operationA)
// 	// _, err = r.db.Collection(TblNode).BulkWrite(ctx, operations)

// 	err = collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		return result, err
// 	}

// 	return result, nil
// }

// func (r *UserMongo) SetBal(userID string, value int) (domain.User, error) {
// 	var result domain.User
// 	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
// 	defer cancel()

// 	collection := r.db.Collection(tblUsers)

// 	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		return result, err
// 	}

// 	filter := bson.M{"_id": userIDPrimitive}

// 	err = collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		return result, err
// 	}

// 	newData := bson.M{}
// 	if value != 0 {
// 		newData["bal"] = (int64)(result.Bal + value) //utils.Max((int64)(result.Bal+value), 0)
// 	}

// 	// fmt.Println("newData=", newData)
// 	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": newData}).Decode(&result)
// 	if err != nil {
// 		return result, err
// 	}

// 	err = collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		return result, err
// 	}

// 	return result, nil
// }
