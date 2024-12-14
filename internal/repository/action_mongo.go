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

type ActionMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewActionMongo(db *mongo.Database, i18n config.I18nConfig) *ActionMongo {
	return &ActionMongo{db: db, i18n: i18n}
}

func (r *ActionMongo) FindAction(params domain.RequestParams) (domain.Response[model.Action], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Action
	var response domain.Response[model.Action]

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return response, err
	}
	fmt.Println(pipe)
	cursor, err := r.db.Collection(TblAction).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Action, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblAction).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Action]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ActionMongo) GqlGetActions(params domain.RequestParams) ([]*model.Action, error) {
	fmt.Println("GqlGetActions: ")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Action
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	cursor, err := r.db.Collection(TblAction).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Action, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *ActionMongo) GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Action
	var response domain.Response[model.Action]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Action]{}, err
	}

	cursor, err := r.db.Collection(TblAction).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Action, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblAction).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Action]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ActionMongo) CreateAction(userID string, data *model.ActionInput) (*model.Action, error) {
	var result *model.Action

	collection := r.db.Collection(TblAction)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	serviceIDPrimitive, err := primitive.ObjectIDFromHex(data.ServiceID)
	if err != nil {
		return nil, err
	}

	newData := model.Action{
		UserID:      userIDPrimitive,
		ServiceID:   serviceIDPrimitive,
		Service:     data.Service,
		Type:        data.Type,
		Description: data.Description,
		Props:       data.Props,
		Status:      data.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblAction).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ActionMongo) UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error) {
	var result *model.Action
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblAction)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.ServiceID != "" {
		serviceIDPrimitive, err := primitive.ObjectIDFromHex(data.ServiceID)
		if err != nil {
			return result, err
		}
		newData["service_id"] = serviceIDPrimitive
	}
	if data.Type != 0 {
		newData["type"] = data.Type
	}
	if data.Service != "" {
		newData["service"] = data.Service
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if data.Props != nil {
		newData["props"] = data.Props
	}
	if data.Status != 0 {
		newData["status"] = data.Status
	}
	newData["updated_at"] = time.Now()

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

func (r *ActionMongo) DeleteAction(id string) (model.Action, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Action{}
	collection := r.db.Collection(TblAction)

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
