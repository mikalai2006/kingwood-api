package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewOrderMongo(db *mongo.Database, i18n config.I18nConfig) *OrderMongo {
	return &OrderMongo{db: db, i18n: i18n}
}

func (r *OrderMongo) FindOrder(input *domain.OrderFilter) (domain.Response[domain.Order], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// var results []domain.Order
	var response domain.Response[domain.Order]
	// var response domain.Response[domain.Order]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.Order]{}, err
	// }

	// cursor, err := r.db.Collection(TblOrder).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.Order, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.Order]{}, err
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
	// Filters
	if input.From != nil && !input.From.IsZero() {
		q = append(q, bson.E{"createdAt", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.From)}}})
	}
	if input.To != nil && !input.To.IsZero() {
		q = append(q, bson.E{"createdAt", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.To)}}})
	}
	if input.Date != nil && !input.Date.IsZero() {
		// q = append(q, bson.E{"from", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.Date)}}})
		q = append(q, bson.E{"dateStart", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.Date)}}})
	}
	if input.CountTaskMontaj != nil {
		// fmt.Println("countTaskMontaj=", *input.CountTaskMontaj)
		q = append(q, bson.E{"countTaskMontaj", bson.D{{"$gt", *input.CountTaskMontaj}}})
	}
	if input.Year != nil && *input.Year > 0 {
		q = append(q, bson.E{"year", input.Year})
	}
	if input.ID != nil && len(input.ID) > 0 {
		ids := []primitive.ObjectID{}
		for key, _ := range input.ID {
			idObjectPrimitive, err := primitive.ObjectIDFromHex(input.ID[key])
			if err != nil {
				return response, err
			}
			ids = append(ids, idObjectPrimitive)
		}
		q = append(q, bson.E{"_id", bson.D{{"$in", ids}}})
	}
	if input.ConstructorId != nil && len(input.ConstructorId) > 0 {
		ids := []primitive.ObjectID{}
		for key, _ := range input.ConstructorId {
			idCPrimitive, err := primitive.ObjectIDFromHex(input.ConstructorId[key])
			if err != nil {
				return response, err
			}
			ids = append(ids, idCPrimitive)
		}
		q = append(q, bson.E{"constructorId", bson.D{{"$in", ids}}})
	}
	if input.Name != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", input.Name), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
	}
	if input.Group != nil && len(input.Group) > 0 {
		q = append(q, bson.E{"group", bson.M{"$elemMatch": bson.D{{"$in", input.Group}}}})
	}
	if len(input.Status) > 0 {
		q = append(q, bson.E{"status", bson.D{{"$in", input.Status}}})
	}
	if input.Number != nil {
		q = append(q, bson.E{"number", input.Number})
	}
	if input.Query != "" {
		qu := []interface{}{}
		i, err := strconv.Atoi(input.Query)
		if err == nil {
			qu = append(qu, bson.D{{"number", i}})
		}
		qu = append(qu, bson.D{{"name", bson.D{{"$regex", primitive.Regex(primitive.Regex{Pattern: input.Query, Options: "i"})}}}}) //, {"$options", "ig"}

		q = append(q, bson.E{
			"$or",
			qu,
		})

	}
	if input.ObjectId != nil {
		objectIds := []primitive.ObjectID{}
		for key, _ := range input.ObjectId {
			idObjectPrimitive, err := primitive.ObjectIDFromHex(input.ObjectId[key])
			if err != nil {
				return response, err
			}
			objectIds = append(objectIds, idObjectPrimitive)
		}
		q = append(q, bson.E{"objectId", bson.D{{"$in", objectIds}}})
	}
	if input.StolyarComplete != nil {
		q = append(q, bson.E{"stolyarComplete", input.StolyarComplete})
	}
	if input.ShlifComplete != nil {
		q = append(q, bson.E{"shlifComplete", input.ShlifComplete})
	}
	if input.MalyarComplete != nil {
		q = append(q, bson.E{"malyarComplete", input.MalyarComplete})
	}
	if input.GoComplete != nil {
		q = append(q, bson.E{"goComplete", input.GoComplete})
	}
	if input.MontajComplete != nil {
		q = append(q, bson.E{"montajComplete", input.MontajComplete})
	}

	// fmt.Println("q: ", q)

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

	// tasks.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblTask,
		"as":   "tasks",
		// "localField":   "userId",
		// "foreignField": "_id",
		"let": bson.D{{Key: "id", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{
				"$expr": bson.M{"$eq": [2]string{"$orderId", "$$id"}},
			}}},
			// workers.
			bson.D{{Key: "$lookup", Value: bson.M{
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
							// bson.D{{Key: "$lookup", Value: bson.M{
							// 	"from": tblTaskWorker,
							// 	"as":   "taskWorkers",
							// 	"let":  bson.D{{Key: "id", Value: "$_id"}},
							// 	"pipeline": mongo.Pipeline{
							// 		bson.D{{Key: "$match", Value: bson.M{
							// 			"$expr":  bson.M{"$eq": [2]string{"$workerId", "$$id"}},
							// 			"status": bson.D{{"$nin", [2]string{"finish", "autofinish"}}},
							// 		}}},

							// 		// object.
							// 		bson.D{{Key: "$lookup", Value: bson.M{
							// 			"from": tblObject,
							// 			"as":   "objecta",
							// 			"let":  bson.D{{Key: "objectId", Value: "$objectId"}},
							// 			"pipeline": mongo.Pipeline{
							// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$objectId"}}}}},
							// 				bson.D{{"$limit", 1}},
							// 			},
							// 		}}},
							// 		bson.D{{Key: "$set", Value: bson.M{"object": bson.M{"$first": "$objecta"}}}},

							// 		// order.
							// 		bson.D{{Key: "$lookup", Value: bson.M{
							// 			"from": TblOrder,
							// 			"as":   "ordera",
							// 			"let":  bson.D{{Key: "orderId", Value: "$orderId"}},
							// 			"pipeline": mongo.Pipeline{
							// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$orderId"}}}}},
							// 				bson.D{{"$limit", 1}},
							// 			},
							// 		}}},
							// 		bson.D{{Key: "$set", Value: bson.M{"order": bson.M{"$first": "$ordera"}}}},

							// 		// task.
							// 		bson.D{{Key: "$lookup", Value: bson.M{
							// 			"from": tblTask,
							// 			"as":   "taska",
							// 			"let":  bson.D{{Key: "taskId", Value: "$taskId"}},
							// 			"pipeline": mongo.Pipeline{
							// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$taskId"}}}}},
							// 				bson.D{{"$limit", 1}},
							// 			},
							// 		}}},
							// 		bson.D{{Key: "$set", Value: bson.M{"task": bson.M{"$first": "$taska"}}}},

							// 		// worker.
							// 		bson.D{{Key: "$lookup", Value: bson.M{
							// 			"from": tblUsers,
							// 			"as":   "usera",
							// 			"let":  bson.D{{Key: "workerId", Value: "$workerId"}},
							// 			"pipeline": mongo.Pipeline{
							// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$workerId"}}}}},
							// 				bson.D{{"$limit", 1}},
							// 				bson.D{{
							// 					Key: "$lookup",
							// 					Value: bson.M{
							// 						"from": tblImage,
							// 						"as":   "images",
							// 						"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
							// 						"pipeline": mongo.Pipeline{
							// 							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
							// 						},
							// 					},
							// 				}},
							// 				// add populate auth.
							// 				bson.D{{
							// 					Key: "$lookup",
							// 					Value: bson.M{
							// 						"from":         TblAuth,
							// 						"as":           "auths",
							// 						"localField":   "userId",
							// 						"foreignField": "_id",
							// 						// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
							// 						// "pipeline": mongo.Pipeline{
							// 						// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
							// 						// },
							// 					},
							// 				}},
							// 				bson.D{{Key: "$set", Value: bson.M{"auth": bson.M{"$first": "$auths"}}}},
							// 				bson.D{{Key: "$set", Value: bson.M{"authPrivate": bson.M{"$first": "$auths"}}}},

							// 				// post.
							// 				bson.D{{
							// 					Key: "$lookup",
							// 					Value: bson.M{
							// 						"from": TblPost,
							// 						"as":   "posts",
							// 						// "localField":   "_id",
							// 						// "foreignField": "serviceId",
							// 						"let": bson.D{{Key: "postId", Value: "$postId"}},
							// 						"pipeline": mongo.Pipeline{
							// 							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$postId"}}}}},
							// 						},
							// 					},
							// 				}},
							// 				bson.D{{Key: "$set", Value: bson.M{"postObject": bson.M{"$first": "$posts"}}}},
							// 				// role.
							// 				bson.D{{Key: "$lookup", Value: bson.M{
							// 					"from": TblRole,
							// 					"as":   "rolea",
							// 					// "localField":   "userId",
							// 					// "foreignField": "_id",
							// 					"let": bson.D{{Key: "roleId", Value: "$roleId"}},
							// 					"pipeline": mongo.Pipeline{
							// 						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$roleId"}}}}},
							// 						bson.D{{"$limit", 1}},
							// 					},
							// 				}}},
							// 				bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$rolea"}}}},
							// 			},
							// 		}}},
							// 		bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}},
							// 	},
							// }}},
						},
					}}},
					bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}},
				},
			}}},
		},
	}}})

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
	dataOptions := bson.A{}
	if input.Skip != nil {
		skip = *input.Skip
		dataOptions = append(dataOptions, bson.D{{"$skip", skip}})
	}
	if input.Limit != nil {
		limit = *input.Limit
		dataOptions = append(dataOptions, bson.D{{"$limit", limit}})
	}
	if input.Sort != nil {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		dataOptions = append(dataOptions, bson.D{{"$sort", sortParam}})
	}

	pipe = append(pipe, bson.D{{Key: "$facet", Value: bson.D{
		{"data", dataOptions},
		{Key: "metadata", Value: mongo.Pipeline{
			bson.D{{"$group", bson.D{
				{"_id", nil},
				{"total", bson.D{{"$sum", 1}}}}}},
		}},
	},
	}})

	// $facet: {
	// 	results: [
	// 	  {
	// 		$skip: 1
	// 	  },
	// 	  {
	// 		$limit: 1
	// 	  }
	// 	],
	// 	count: [
	// 	  {
	// 		$group: {
	// 		  _id: null,
	// 		  count: {
	// 			$sum: 1
	// 		  }
	// 		}
	// 	  }
	// 	]
	//   }

	cursor, err := r.db.Collection(TblOrder).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	resultMap := []bson.M{}
	if er := cursor.All(ctx, &resultMap); er != nil {
		return response, er
	}
	resultFacetOne := domain.ResultFacetOrder{}
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
	// fmt.Println("rrrrr: ", resultMap[0])

	// count, err := r.db.Collection(TblOrder).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.Order]{
		Total: total,
		Skip:  skip,
		Limit: limit,
		Data:  resultFacetOne.Data, //results,
	}
	return response, nil
}

func (r *OrderMongo) CreateOrder(userID string, data *domain.Order) (*domain.Order, error) {
	var result *domain.Order

	collection := r.db.Collection(TblOrder)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existOrder domain.Order
	// r.db.Collection(TblOrder).FindOne(ctx, bson.M{"node_id": Order.NodeID, "userId": userIDPrimitive}).Decode(&existOrder)

	// if existOrder.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	// itemCount, err := collection.CountDocuments(ctx, bson.M{"year": time.Now().Year()})
	// if err != nil {
	// 	return nil, err
	// }

	year := time.Now().Year()
	if data.Year != nil {
		year = *data.Year
	}

	defaultStatus := int64(0)
	newOrder := domain.OrderInput{
		UserID:          userIDPrimitive,
		Name:            data.Name,
		Description:     data.Description,
		ObjectId:        data.ObjectId,
		Number:          int64(data.Number), //itemCount + 1,
		ConstructorId:   data.ConstructorId,
		Priority:        data.Priority,
		Term:            data.Term,
		TermMontaj:      data.TermMontaj,
		Status:          &defaultStatus,
		Group:           data.Group,
		StolyarComplete: &defaultStatus,
		MalyarComplete:  &defaultStatus,
		ShlifComplete:   &defaultStatus,
		GoComplete:      &defaultStatus,
		MontajComplete:  &defaultStatus,
		Year:            year,

		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	insertedID := res.InsertedID.(primitive.ObjectID).Hex()
	insertedItem, err := r.FindOrder(&domain.OrderFilter{ID: []string{insertedID}})
	// r.db.Collection(TblOrder).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(insertedItem.Data) > 0 {
		result = &insertedItem.Data[0]
	}
	// } else {
	// 	updatedAt := Order.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateOrder := &domain.OrderInput{
	// 		Rate:      Order.Rate,
	// 		Order:    Order.Order,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateOrder(existOrder.ID.Hex(), userID, updateOrder)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *OrderMongo) UpdateOrder(id string, userID string, data *domain.OrderInput) (*domain.Order, error) {
	var result *domain.Order
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblOrder)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Name != "" {
		newData["name"] = data.Name
	}
	if !data.ObjectId.IsZero() {
		newData["objectId"] = data.ObjectId
	}
	if !data.Term.IsZero() {
		newData["term"] = data.Term
	}
	if !data.DateStart.IsZero() {
		newData["dateStart"] = data.DateStart
	}
	if !data.TermMontaj.IsZero() {
		newData["termMontaj"] = data.TermMontaj
	}
	if data.Priority != nil {
		newData["priority"] = data.Priority
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if data.StolyarComplete != nil {
		newData["stolyarComplete"] = data.StolyarComplete
	}
	if data.MalyarComplete != nil {
		newData["malyarComplete"] = data.MalyarComplete
	}
	if data.ShlifComplete != nil {
		newData["shlifComplete"] = data.ShlifComplete
	}
	if data.GoComplete != nil {
		newData["goComplete"] = data.GoComplete
	}
	if data.MontajComplete != nil {
		newData["montajComplete"] = data.MontajComplete
	}
	if data.CountTaskMontaj != nil {
		newData["countTaskMontaj"] = data.CountTaskMontaj
	}

	if !data.DateOtgruzka.IsZero() {
		newData["dateOtgruzka"] = data.DateOtgruzka
	}
	// if data.NeedMontaj != nil {
	// 	newData["needMontaj"] = data.NeedMontaj
	// }
	if !data.ConstructorId.IsZero() {
		newData["constructorId"] = data.ConstructorId
	}
	if len(data.Group) > 0 {
		newData["group"] = data.Group
	}
	if data.Status != nil {
		newData["status"] = data.Status
	}
	newData["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}
	// collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": newData}).Decode(&result)
	// // if err != nil {
	// // 	return result, err
	// // }
	// if data.StolyarComplete != nil ||
	// 	data.MalyarComplete != nil ||
	// 	data.MontajComplete != nil {
	// 	statusCompleted := int64(1)
	// 	dataUpdate := bson.M{}
	// 	// if result.MalyarComplete == &statusCompleted  && result.StolyarComplete == &statusCompleted {
	// 	// 	dataUpdate["goComplete"] = 1
	// 	// }
	// 	if result.MalyarComplete == &statusCompleted && result.StolyarComplete == &statusCompleted && result.MontajComplete == &statusCompleted {
	// 		dataUpdate["status"] = 100
	// 	} else {
	// 		dataUpdate["status"] = 1
	// 	}

	// 	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }

	orders, err := r.FindOrder(&domain.OrderFilter{ID: []string{id}})
	// collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(orders.Data) > 0 {
		result = &orders.Data[0]
	}

	return result, nil
}

func (r *OrderMongo) DeleteOrder(id string) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.Order{}
	collection := r.db.Collection(TblOrder)

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
