package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OfferMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewOfferMongo(db *mongo.Database, i18n config.I18nConfig) *OfferMongo {
	return &OfferMongo{db: db, i18n: i18n}
}

func (r *OfferMongo) GetOffer(id string) (*model.Offer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result *model.Offer
	var pipe mongo.Pipeline

	IDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	pipe = append(pipe, bson.D{{"$match", bson.M{"_id": IDPrimitive}}})
	pipe = append(pipe, bson.D{{"$limit", 1}})

	cursor, err := r.db.Collection(TblOffer).Aggregate(ctx, pipe)
	// fmt.Println("filter Offer:::", pipe)
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

func (r *OfferMongo) FindOffer(params *model.OfferFilter) (domain.Response[model.Offer], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Offer
	var response domain.Response[model.Offer]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Offer]{}, err
	// }
	// cursor, err := r.db.Collection(TblOffer).Find(ctx, filter, opts)
	// fmt.Println(params)
	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return response, err
	// }
	q := bson.D{}
	if params.ID != nil {
		IDPrimitive, err := primitive.ObjectIDFromHex(*params.ID)
		if err != nil {
			return response, err
		}
		q = append(q, bson.E{"_id", IDPrimitive})
	}
	if params.UserProductID != nil {
		q = append(q, bson.E{"_id", bson.D{{"$in", params.UserProductID}}})
	}
	if params.UserID != nil {
		// q = append(q, bson.E{"userId", params.UserID})

		queryArr := []bson.M{}
		queryArr = append(queryArr, bson.M{"userId": params.UserID})
		queryArr = append(queryArr, bson.M{"userProductId": params.UserID})
		q = append(q, bson.E{"$or", queryArr})
	}
	if params.ProductID != nil {
		q = append(q, bson.E{"productId", bson.D{{"$in", params.ProductID}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	limit := 10
	skip := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Skip != nil {
		skip = *params.Skip
	}

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usera",
		"let":  bson.D{{Key: "userId", Value: "$userId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"as":   "images",
					"from": "image",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}},
		bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})

	pipe = append(pipe, bson.D{{"$limit", skip + limit}})
	pipe = append(pipe, bson.D{{"$skip", skip}})

	cursor, err := r.db.Collection(TblOffer).Aggregate(ctx, pipe)
	// fmt.Println("filter Offer:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Offer, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)
	// fmt.Println("results::", len(results))

	count, err := r.db.Collection(TblOffer).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Offer]{
		Total: int(count),
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *OfferMongo) CreateOffer(userID string, data *model.OfferInput) (*model.Offer, error) {
	var result *model.Offer

	collection := r.db.Collection(TblOffer)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	productIDPrimitive, err := primitive.ObjectIDFromHex(data.ProductID)
	if err != nil {
		return nil, err
	}
	userProductIDPrimitive, err := primitive.ObjectIDFromHex(data.UserProductID)
	if err != nil {
		return nil, err
	}

	defaultValue := 0

	newOffer := model.OfferInputMongo{
		UserID:        userIDPrimitive,
		ProductID:     productIDPrimitive,
		UserProductID: userProductIDPrimitive,
		Cost:          data.Cost,
		Status:        1,
		Win:           &defaultValue,
		Take:          &defaultValue,
		Give:          &defaultValue,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	res, err := collection.InsertOne(ctx, newOffer)
	if err != nil {
		return nil, err
	}

	// // change user stat
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": userIDPrimitive}, bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.nodedata", 1},
	// 	}},
	// })

	// err = r.db.Collection(TblOffer).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	// if err != nil {
	// 	return nil, err
	// }
	insertedID := res.InsertedID.(primitive.ObjectID).Hex()

	offers, err := r.FindOffer(&model.OfferFilter{ID: &insertedID})
	if err != nil {
		return result, err
	}

	result = &offers.Data[0]

	return result, nil
}

func (r *OfferMongo) GqlGetOffers(params domain.RequestParams) ([]*model.Offer, error) {
	// fmt.Println("GqlGetOffers", &r.i18n, params.Lang)
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Offer
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)
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

	// // get tag
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": TblTag,
	// 	"as":   "tags",
	// 	"let":  bson.D{{Key: "tagId", Value: "$tag_id"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$tagId"}}}}},
	// 		bson.D{{"$limit", 1}},
	// 	},
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"tag": bson.M{"$first": "$tags"}}}})

	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": "nodedata_vote",
	// 	"as":   "votes",
	// 	"let":  bson.D{{Key: "id", Value: "$_id"}},
	// 	"pipeline": mongo.Pipeline{
	// 		bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$nodedata_id", "$$id"}}}}},
	// 		bson.D{{"$sort", bson.D{{"updated_at", -1}}}},
	// 		bson.D{{"$limit", 1}},

	// 		bson.D{{Key: "$lookup", Value: bson.M{
	// 			"from": "users",
	// 			"as":   "userx",
	// 			"let":  bson.D{{Key: "userId", Value: "$user_id"}},
	// 			"pipeline": mongo.Pipeline{
	// 				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
	// 				bson.D{{"$limit", 1}},
	// 				bson.D{{
	// 					Key: "$lookup",
	// 					Value: bson.M{
	// 						"from": "image",
	// 						"as":   "images",
	// 						"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
	// 						"pipeline": mongo.Pipeline{
	// 							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
	// 						},
	// 					},
	// 				}},
	// 			},
	// 		}}},
	// 		bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userx"}}}},
	// 	},
	// }}})
	pipe = append(pipe, bson.D{{"$sort", bson.D{{"createdAt", -1}}}})

	cursor, err := r.db.Collection(TblOffer).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Offer, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *OfferMongo) UpdateOffer(id string, userID string, data *model.Offer) (*model.Offer, error) {
	var result *model.Offer
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblOffer)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if !data.ProductID.IsZero() {
		newData["productId"] = data.ProductID
	}
	if !data.UserProductID.IsZero() {
		newData["userProductId"] = data.UserProductID
	}
	if !data.RejectUserId.IsZero() {
		newData["rejectUserId"] = data.RejectUserId
	}
	if !data.RoomId.IsZero() {
		newData["roomId"] = data.RoomId
	}
	if data.Status != 0 {
		newData["status"] = data.Status
	}
	if data.Cost != 0 {
		newData["cost"] = data.Cost
	}
	if data.Give != nil {
		newData["give"] = data.Give
	}
	if data.Message != "" {
		newData["message"] = data.Message
	}
	if data.Take != nil {
		newData["take"] = data.Take
	}
	if data.Win != nil {
		newData["win"] = data.Win
	}
	newData["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	offers, err := r.FindOffer(&model.OfferFilter{ID: &id})
	if err != nil {
		return result, err
	}

	result = &offers.Data[0]

	return result, nil
}

func (r *OfferMongo) DeleteOffer(id string) (model.Offer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Offer{}
	collection := r.db.Collection(TblOffer)

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

	// // change user stat
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.nodedata", -1},
	// 	}},
	// }) //, options.Update().SetUpsert(true)

	// // remove likes.
	// _, err = r.db.Collection(TblNodedataVote).DeleteMany(ctx, bson.M{"nodedata_id": idPrimitive})
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}
