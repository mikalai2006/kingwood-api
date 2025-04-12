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

type ArchiveMessageMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

type ResultMetadataArchiveMessage struct {
	ID    interface{} `json:"_id" bson:"_id"`
	Total int         `json:"total" bson:"total"`
}
type ResultFacetArchiveMessage struct {
	Metadata []ResultMetadataArchiveMessage `json:"metadata" bson:"metadata"`
	Data     []domain.ArchiveMessage        `json:"data" bson:"data"`
}

func NewArchiveMessageMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveMessageMongo {
	return &ArchiveMessageMongo{db: db, i18n: i18n}
}

func (r *ArchiveMessageMongo) CreateArchiveMessage(userID string, data *domain.Message) (*domain.ArchiveMessage, error) {
	var result *domain.ArchiveMessage

	collection := r.db.Collection(TblArchiveMessage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newArchiveMessage := domain.ArchiveMessageInput{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:        data.ID,
		UserID:    data.UserID,
		Message:   data.Message,
		Props:     data.Props,
		Images:    data.Images,
		OrderID:   data.OrderID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	res, err := collection.InsertOne(ctx, newArchiveMessage)
	if err != nil {
		return nil, err
	}

	insertedID := res.InsertedID.(primitive.ObjectID).Hex()
	insertedArchiveMessage, err := r.FindArchiveMessage(&domain.ArchiveMessageFilter{ID: insertedID})
	if err != nil {
		return nil, err
	}

	result = &insertedArchiveMessage.Data[0]

	return result, nil
}

func (r *ArchiveMessageMongo) FindArchiveMessage(params *domain.ArchiveMessageFilter) (domain.Response[domain.ArchiveMessage], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// var results []domain.ArchiveMessage
	var response domain.Response[domain.ArchiveMessage]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Node]{}, err
	// }
	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.ArchiveMessage]{}, err
	// }
	// fmt.Println(params)
	q := bson.D{}
	if params.UserID != "" {
		userIDPrimitive, err := primitive.ObjectIDFromHex(params.UserID)
		if err != nil {
			return response, err
		}
		q = append(q, bson.E{"userId", userIDPrimitive})
	}
	if params.ID != "" {
		iDPrimitive, err := primitive.ObjectIDFromHex(params.ID)
		if err != nil {
			return response, err
		}
		q = append(q, bson.E{"_id", iDPrimitive})
	}
	if len(params.OrderID) > 0 {
		ids := []primitive.ObjectID{}
		for i := range params.OrderID {
			idPrimitive, err := primitive.ObjectIDFromHex(params.OrderID[i])
			if err != nil {
				return response, err
			}
			ids = append(ids, idPrimitive)
		}
		q = append(q, bson.E{"orderId", bson.D{{"$in", ids}}})
	}

	// // Filter by products id.
	// if params.ProductID != nil && !params.ProductID.IsZero() {
	// 	q = append(q, bson.E{"productId", params.ProductID})
	// }

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	// if params.Sort != nil && len(params.Sort) > 0 {
	// 	sortParam := bson.D{}
	// 	for i := range params.Sort {
	// 		sortParam = append(sortParam, bson.E{*params.Sort[i].Key, *params.Sort[i].Value})
	// 	}
	// 	pipe = append(pipe, bson.D{{"$sort", sortParam}})
	// 	// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	// }

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": TblArchiveMessageStatus,
		"as":   "statuses",
		"let":  bson.D{{Key: "messageId", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$messageId", "$$messageId"}}}}},
			bson.D{{"$limit", 100}},
		},
	}}})

	if params.Sort != nil && len(params.Sort) > 0 {
		sortParam := bson.D{}
		for i := range params.Sort {
			sortParam = append(sortParam, bson.E{params.Sort[i].Key, params.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
		// fmt.Println("sortParam: ", len(input.Sort), sortParam, pipe)
	}

	skip := 0
	limit := 10
	dataOptions := bson.A{}
	if params.Skip != nil {
		skip = *params.Skip
		dataOptions = append(dataOptions, bson.D{{"$skip", skip}})
	}
	if params.Limit != nil {
		limit = *params.Limit
		dataOptions = append(dataOptions, bson.D{{"$limit", limit}})
	}
	if params.Sort != nil {
		sortParam := bson.D{}
		for i := range params.Sort {
			sortParam = append(sortParam, bson.E{params.Sort[i].Key, params.Sort[i].Value})
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

	// pipe = append(pipe, bson.D{{"$limit", skip + limit}})
	// pipe = append(pipe, bson.D{{"$skip", skip}})

	cursor, err := r.db.Collection(TblArchiveMessage).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]domain.ArchiveMessage, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	// count := len(resultSlice)
	// // count, err := r.db.Collection(TblNode).CountDocuments(ctx, params.Filter)
	// // if err != nil {
	// // 	return response, err
	// // }

	resultMap := []bson.M{}
	if er := cursor.All(ctx, &resultMap); er != nil {
		return response, er
	}
	resultFacetOne := ResultFacetArchiveMessage{}
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

	response = domain.Response[domain.ArchiveMessage]{
		Total: total,
		Skip:  skip,
		Limit: limit,
		Data:  resultFacetOne.Data,
	}
	return response, nil
}

func (r *ArchiveMessageMongo) DeleteArchiveMessage(id string) (domain.ArchiveMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.ArchiveMessage{}
	collection := r.db.Collection(TblArchiveMessage)

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
