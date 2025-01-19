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

type PayTemplateMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewPayTemplateMongo(db *mongo.Database, i18n config.I18nConfig) *PayTemplateMongo {
	return &PayTemplateMongo{db: db, i18n: i18n}
}

func (r *PayTemplateMongo) FindPayTemplate(params domain.RequestParams) (domain.Response[domain.PayTemplate], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.PayTemplate
	var response domain.Response[domain.PayTemplate]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.PayTemplate]{}, err
	// }
	// cursor, err := r.db.Collection(TblPayTemplate).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(tblPayTemplate).Aggregate(ctx, pipe)
	// fmt.Println("filter PayTemplate:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.PayTemplate, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblPayTemplate).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.PayTemplate]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *PayTemplateMongo) CreatePayTemplate(userID string, data *domain.PayTemplate) (*domain.PayTemplate, error) {
	var result *domain.PayTemplate
	collection := r.db.Collection(tblPayTemplate)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIdPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	enabled := int64(0)
	if data.Enabled != nil {
		enabled = *data.Enabled
	}

	newPayTemplate := domain.PayTemplate{
		UserID:      userIdPrimitive,
		Name:        data.Name,
		Description: data.Description,
		Total:       data.Total,
		Enabled:     &enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newPayTemplate)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblPayTemplate).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PayTemplateMongo) UpdatePayTemplate(id string, userID string, data *domain.PayTemplateInput) (*domain.PayTemplate, error) {
	var result *domain.PayTemplate
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblPayTemplate)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}

	if data.Name != "" {
		newData["name"] = data.Name
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if data.Enabled != nil {
		newData["enabled"] = data.Enabled
	}
	if data.Total != nil {
		newData["total"] = data.Total
	}

	newData["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *PayTemplateMongo) DeletePayTemplate(id string) (domain.PayTemplate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.PayTemplate{}
	collection := r.db.Collection(tblPayTemplate)

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
