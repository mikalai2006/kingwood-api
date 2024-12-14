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

type QuestionMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewQuestionMongo(db *mongo.Database, i18n config.I18nConfig) *QuestionMongo {
	return &QuestionMongo{db: db, i18n: i18n}
}

func (r *QuestionMongo) FindQuestion(params *model.QuestionFilter) (domain.Response[model.Question], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Question
	var response domain.Response[model.Question]

	// filter := params.Filter.(map[string]interface{})
	// if filter["tag_id"] != nil {
	// 	tagIDPrimitive, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", filter["tag_id"]))
	// 	if err != nil {
	// 		return response, err
	// 	}

	// 	filter["tag_id"] = tagIDPrimitive
	// }
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Question]{}, err
	// }

	// cursor, err := r.db.Collection(TblQuestion).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// if params.Filter["tag_id"] {

	// }

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
		q = append(q, bson.E{"userProductId", bson.D{{"$in", params.UserProductID}}})
	}
	if params.UserID != nil && len(params.UserID) > 0 {
		// q = append(q, bson.E{"userId", params.UserID})

		queryArr := []bson.M{}
		queryArr = append(queryArr, bson.M{"userId": bson.D{{"$in", params.UserID}}})
		queryArr = append(queryArr, bson.M{"userProductId": bson.D{{"$in", params.UserID}}})
		q = append(q, bson.E{"$or", queryArr})
	}
	if params.ProductID != nil {
		q = append(q, bson.E{"productId", bson.D{{"$in", params.ProductID}}})
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

	cursor, err := r.db.Collection(TblQuestion).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Question, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblQuestion).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Question]{
		Total: int(count),
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *QuestionMongo) CreateQuestion(userID string, data *model.QuestionInput) (*model.Question, error) {
	var result *model.Question

	collection := r.db.Collection(TblQuestion)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// productIDPrimitive, err := primitive.ObjectIDFromHex(data.ProductID)
	// if err != nil {
	// 	return nil, err
	// }

	newItem := model.QuestionInput{
		UserID:        userIDPrimitive,
		ProductID:     data.ProductID,
		UserProductID: data.UserProductID,
		Question:      data.Question,
		Answer:        data.Answer,
		Status:        data.Status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	defaultStatus := 1
	if data.Status != nil {
		newItem.Status = data.Status
	} else {
		newItem.Status = &defaultStatus
	}

	res, err := collection.InsertOne(ctx, newItem)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblQuestion).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *QuestionMongo) UpdateQuestion(id string, userID string, data *model.QuestionInput) (*model.Question, error) {
	var result *model.Question
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }

	// newTag := model.Tag{
	// 	UserID:      userIDPrimitive,
	// 	Key:         data.Key,
	// 	Title:       data.Title,
	// 	Description: data.Description,
	// 	Props:       data.Props,
	// 	Locale:      data.Locale,
	// 	UpdatedAt:   time.Now(),
	// }
	// obj := data.(map[string]interface{})
	// obj["user_id"] = userIDPrimitive
	// data = obj

	collection := r.db.Collection(TblQuestion)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	// if !data.ProductID.IsZero() {
	// 	// productIDPrimitive, err := primitive.ObjectIDFromHex(data.ProductID)
	// 	// if err != nil {
	// 	// 	return result, err
	// 	// }
	// 	newData["tag_id"] = data.ProductID
	// }

	if data.Question != "" {
		newData["question"] = data.Question
	}
	if data.Answer != "" {
		newData["answer"] = data.Answer
	}
	// if !data.UserProductID.IsZero() {
	// 	newData["userProductId"] = data.UserProductID
	// }
	// if data.Status != nil {
	// 	newData["status"] = data.Status
	// }
	newData["status"] = 1
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

func (r *QuestionMongo) DeleteQuestion(id string) (model.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Question{}
	collection := r.db.Collection(TblQuestion)

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
