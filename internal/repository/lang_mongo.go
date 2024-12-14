package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LangMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewLangMongo(db *mongo.Database, i18n config.I18nConfig) *LangMongo {
	return &LangMongo{db: db, i18n: i18n}
}

func (r *LangMongo) CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error) {
	var result domain.Language

	collection := r.db.Collection(TblLanguage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }
	// count, err := r.db.Collection(tblLanguage).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return response, err
	// }
	// newId := count + 1

	newPage := domain.Language{
		Publish:      data.Publish,
		Flag:         data.Flag,
		Name:         data.Name,
		Code:         data.Code,
		Locale:       data.Locale,
		SortOrder:    data.SortOrder,
		Localization: data.Localization,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	res, err := collection.InsertOne(ctx, newPage)
	if err != nil {
		return result, err
	}

	err = r.db.Collection(TblLanguage).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *LangMongo) GetLanguage(id string) (domain.Language, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Language

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Language{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(TblLanguage).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Language{}, err
	}

	return result, nil
}

func (r *LangMongo) FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Language
	var response domain.Response[domain.Language]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Language]{}, err
	}
	cursor, err := r.db.Collection(TblLanguage).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Language, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	options.SetSkip(params.Skip)
	count, err := r.db.Collection(TblLanguage).CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Language]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *LangMongo) UpdateLanguage(id string, data interface{}) (domain.Language, error) {
	var result domain.Language
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblLanguage)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *LangMongo) DeleteLanguage(id string) (domain.Language, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Language{}
	collection := r.db.Collection(TblLanguage)

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
