package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewProductMongo(db *mongo.Database, i18n config.I18nConfig) *ProductMongo {
	return &ProductMongo{db: db, i18n: i18n}
}

func (r *ProductMongo) FindProduct(params *model.ProductFilter) (domain.Response[model.Product], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Product
	var response domain.Response[model.Product]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Product]{}, err
	// }
	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[model.Product]{}, err
	// }
	// fmt.Println(params.Filter)
	q := bson.D{}
	// Filter by substring name
	if params.Query != nil && *params.Query != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", *params.Query), Options: "i"}
		q = append(q, bson.E{"title", bson.D{{"$regex", strName}}})
		// fmt.Println("q:", q)
	}
	if params.UserID != nil && *params.UserID != "" {
		userIDPrimitive, err := primitive.ObjectIDFromHex(*params.UserID)
		if err != nil {
			return response, err
		}
		q = append(q, bson.E{"user_id", userIDPrimitive})
	}
	if params.ID != nil && len(params.ID) > 0 {
		userIDs := make([]primitive.ObjectID, len(params.ID))
		for i := range params.ID {
			IDPrimitive, err := primitive.ObjectIDFromHex(*params.ID[i])
			if err != nil {
				return response, err
			}
			userIDs[i] = IDPrimitive
		}
		q = append(q, bson.E{"_id", bson.D{{"$in", userIDs}}})
	}

	// Filter by category
	if params.CategoryID != nil && len(params.CategoryID) > 0 {
		categories := make([]primitive.ObjectID, len(params.CategoryID))
		for i := range params.CategoryID {
			categoryIDPrimitive, err := primitive.ObjectIDFromHex(*params.CategoryID[i])
			if err != nil {
				return response, err
			}
			categories[i] = categoryIDPrimitive
		}
		// fmt.Println(len(input.CategoryID), categories)
		q = append(q, bson.E{"category_id", bson.D{{"$in", categories}}})
	}

	// Filter by products id.
	if params.ProductID != nil {
		q = append(q, bson.E{"_id", bson.D{{"$in", params.ProductID}}})
	}

	if params.AddressId != nil {
		addressIDPrimitive, err := primitive.ObjectIDFromHex(*params.AddressId)
		if err != nil {
			return response, err
		}
		q = append(q, bson.E{"addressId", addressIDPrimitive})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	if params.Sort != nil && len(params.Sort) > 0 {
		sortParam := bson.D{}
		for i := range params.Sort {
			sortParam = append(sortParam, bson.E{*params.Sort[i].Key, *params.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
		// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	}

	// pipe = append(pipe,
	// 	bson.D{{Key: "$lookup", Value: bson.M{
	// 		"from": "Productdata",
	// 		// "let":  bson.D{{Key: "ProductId", Value: bson.D{{"$toString", "$_id"}}}},
	// 		// "pipeline": mongo.Pipeline{
	// 		// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$Product_id", "$$ProductId"}}}}},
	// 		// },
	// 		"localField":   "_id",
	// 		"foreignField": "Product_id",
	// 		"as":           "data",
	// 	}}})

	// pipe = append(pipe,
	// 	bson.D{{
	// 		Key: "$lookup",
	// 		Value: bson.M{
	// 			"from": "image",
	// 			"as":   "images",
	// 			"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 			"pipeline": mongo.Pipeline{
	// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
	// 			},
	// 		},
	// 	}})
	pipe = append(pipe,
		bson.D{{
			Key: "$lookup",
			Value: bson.M{
				"from": "image",
				"as":   "images",
				"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
				},
			},
		}})

	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"as":   "offers",
	// 	"from": "offer",
	// 	"let":  bson.D{{Key: "productId", Value: "$_id"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$productId", "$$productId"}}}}},
	// 		bson.D{{Key: "$lookup", Value: bson.M{
	// 			"from": "users",
	// 			"as":   "usera",
	// 			"let":  bson.D{{Key: "userId", Value: "$user_id"}},
	// 			"pipeline": mongo.Pipeline{
	// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
	// 				bson.D{{"$limit", 1}},
	// 				bson.D{{
	// 					Key: "$lookup",
	// 					Value: bson.M{
	// 						"as":   "images",
	// 						"from": "image",
	// 						"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 						"pipeline": mongo.Pipeline{
	// 							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
	// 						},
	// 					},
	// 				}},
	// 			},
	// 		}}},
	// 		bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}},

	// 		// // tagopt
	// 		// bson.D{{Key: "$lookup", Value: bson.M{
	// 		// 	"from": "tagopt",
	// 		// 	"as":   "tagopts",
	// 		// 	"let":  bson.D{{Key: "tagoptId", Value: "$tagopt_id"}},
	// 		// 	"pipeline": mongo.Pipeline{
	// 		// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$tagoptId"}}}}},
	// 		// 		bson.D{{"$limit", 1}},
	// 		// 		// bson.D{{
	// 		// 		// 	Key: "$lookup",
	// 		// 		// 	Value: bson.M{
	// 		// 		// 		"as":   "images",
	// 		// 		// 		"from": "image",
	// 		// 		// 		"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 		// 		// 		"pipeline": mongo.Pipeline{
	// 		// 		// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
	// 		// 		// 		},
	// 		// 		// 	},
	// 		// 		// }},
	// 		// 	},
	// 		// }}},
	// 		// bson.D{{Key: "$set", Value: bson.M{"tagopt": bson.M{"$first": "$tagopts"}}}},

	// 		// // audit section
	// 		// bson.D{{Key: "$lookup", Value: bson.M{
	// 		// 	"as":   "audit",
	// 		// 	"from": "nodedata_audit",
	// 		// 	"let":  bson.D{{Key: "nodedataId", Value: "$_id"}},
	// 		// 	"pipeline": mongo.Pipeline{
	// 		// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$nodedata_id", "$$nodedataId"}}}}},
	// 		// 		bson.D{{Key: "$lookup", Value: bson.M{
	// 		// 			"from": "users",
	// 		// 			"as":   "userc",
	// 		// 			"let":  bson.D{{Key: "userId", Value: "$user_id"}},
	// 		// 			"pipeline": mongo.Pipeline{
	// 		// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
	// 		// 				bson.D{{"$limit", 1}},
	// 		// 				bson.D{{
	// 		// 					Key: "$lookup",
	// 		// 					Value: bson.M{
	// 		// 						"as":   "images",
	// 		// 						"from": "image",
	// 		// 						"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 		// 						"pipeline": mongo.Pipeline{
	// 		// 							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
	// 		// 						},
	// 		// 					},
	// 		// 				}},
	// 		// 			},
	// 		// 		}}},
	// 		// 		bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userc"}}}},
	// 		// 	},
	// 		// 	// "localField":   "_id",
	// 		// 	// "foreignField": "node_id",
	// 		// }}},
	// 	},
	// 	// "localField":   "_id",
	// 	// "foreignField": "node_id",
	// }}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "userb",
		"let":  bson.D{{Key: "userId", Value: "$user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": "image",
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userb"}}}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "address",
		"as":   "addressb",
		"let":  bson.D{{Key: "addressId", Value: "$addressId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$addressId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"address": bson.M{"$first": "$addressb"}}}})

	limit := 100
	skip := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Skip != nil {
		skip = *params.Skip
	}

	pipe = append(pipe, bson.D{{"$limit", skip + limit}})
	pipe = append(pipe, bson.D{{"$skip", skip}})

	cursor, err := r.db.Collection(TblProduct).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblProduct).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Product, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// count, err := r.db.Collection(TblProduct).CountDocuments(ctx, params.Filter)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[model.Product]{
		Total: len(resultSlice),
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ProductMongo) CreateProduct(userID string, Product *model.ProductInputData) (*model.Product, error) {
	var result *model.Product

	collection := r.db.Collection(TblProduct)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	categoryIDPrimitive, err := primitive.ObjectIDFromHex(Product.CategoryID)
	if err != nil {
		return nil, err
	}
	addressIDPrimitive, err := primitive.ObjectIDFromHex(Product.AddressId)
	if err != nil {
		return nil, err
	}

	createdAt := Product.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	newProduct := model.ProductInput{
		UserID:      userIDPrimitive,
		CategoryID:  categoryIDPrimitive,
		Props:       Product.Props,
		Title:       Product.Title,
		Description: Product.Description,
		Status:      Product.Status,
		Cost:        Product.Cost,
		Lon:         Product.Lon,
		Lat:         Product.Lat,
		AddressId:   addressIDPrimitive,
		Actions:     Product.Actions,
		CreatedAt:   createdAt,
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	// change user stat
	// var operations []mongo.WriteModel
	// operationA := mongo.NewUpdateOneModel()
	// operationA.SetFilter(bson.M{"_id": userIDPrimitive})
	// operationA.SetUpdate(bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.Product", 1},
	// 	}},
	// })
	// operations = append(operations, operationA)
	// _, err = r.db.Collection(TblProduct).BulkWrite(ctx, operations)
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": userIDPrimitive}, bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.Product", 1},
	// 	}},
	// })

	err = r.db.Collection(TblProduct).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ProductMongo) UpdateProduct(id string, userID string, data *model.Product) (*model.Product, error) {
	var result *model.Product
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblProduct)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// idUser, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return result, err
	// }
	filter := bson.M{"_id": idPrimitive}

	// // Find old data
	// var oldResult *model.Product
	// err = collection.FindOne(ctx, filter).Decode(&oldResult)
	// if err != nil {
	// 	return result, err
	// }
	// oldProductAudit := model.ProductInput{
	// 	UserID:    oldResult.UserID,
	// 	Lon:       oldResult.Lon,
	// 	Lat:       oldResult.Lat,
	// 	Type:      oldResult.Type,
	// 	Name:      oldResult.Name,
	// 	OsmID:     oldResult.OsmID,
	// 	AmenityID: oldResult.AmenityID,
	// 	Props:     oldResult.Props,
	// 	Status:    oldResult.Status,
	// 	Like:      oldResult.Like,
	// 	Dlike:     oldResult.Dlike,
	// 	UpdatedAt: time.Now(),
	// }
	// _, err = r.db.Collection(TblProductHistory).InsertOne(ctx, oldProductAudit)
	// if err != nil {
	// 	return result, err
	// }
	// fmt.Println(data)
	newData := bson.M{}
	if data.Status != 0 {
		newData["status"] = data.Status
	}
	if data.Cost != nil {
		newData["cost"] = data.Cost
	}
	if data.Title != "" {
		newData["title"] = data.Title
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if !data.CategoryID.IsZero() {
		newData["category_id"] = data.CategoryID
	}
	if !data.AddressId.IsZero() {
		newData["addressId"] = data.AddressId
	}
	if data.Lat != 0 {
		newData["lat"] = data.Lat
	}
	if data.Lon != 0 {
		newData["lon"] = data.Lon
	}
	if len(data.Actions) > 0 {
		newData["actions"] = data.Actions
	}
	if data.Props != nil {
		//newProps := make(map[string]interface{})
		newProps := data.Props
		if val, ok := data.Props["status"]; ok {
			if val == -1.0 {
				newDel := make(map[string]interface{})
				newDel["user_id"] = userID
				newDel["del_at"] = time.Now()
				newProps["del"] = newDel
			}
		}
		newData["props"] = newProps
	}
	newData["updated_at"] = time.Now()
	fmt.Println("Edit product: ", id, newData)
	// test := model.ProductLike{}
	// if data.ProductLike != test {
	// 	newData["Product_like"] = data.ProductLike
	// }
	// if data.Status != 0 {
	// 	newData["status"] = data.Status
	// }
	// bson.M{
	// 	"lon":        data.Lon,
	// 	"lat":        data.Lat,
	// 	"type":       data.Type,
	// 	"osm_id":     data.OsmID,
	// 	"amenity_id": data.AmenityID,
	// 	"props":      data.Props,
	// 	"name":       data.Name,
	// 	"updated_at": time.Now(),
	// }
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	// update type for Productdata collection
	// if data.Type != "" {
	// _, err = r.db.Collection(TblProductdata).UpdateMany(
	// 	ctx,
	// 	bson.M{"Product_id": result.ID},
	// 	bson.M{"$set": bson.M{"type": result.Type, "lat": result.Lat, "lon": result.Lon}},
	// )
	// if err != nil {
	// 	return result, err
	// }
	// }

	return result, nil
}

func (r *ProductMongo) DeleteProduct(id string) (model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Product{}
	collection := r.db.Collection(TblProduct)

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

	// change user stat
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.Product", -1},
	// 	}},
	// })

	// // remove Productdata.
	// _, err = r.db.Collection(TblProductdata).DeleteMany(ctx, bson.M{"Product_id": idPrimitive})
	// if err != nil {
	// 	return result, err
	// }

	// // remove reviews.
	// _, err = r.db.Collection(TblReview).DeleteMany(ctx, bson.M{"Product_id": idPrimitive})
	// if err != nil {
	// 	return result, err
	// }

	// // remove audits.
	// _, err = r.db.Collection(TblProductAudit).DeleteMany(ctx, bson.M{"Product_id": idPrimitive})
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}
