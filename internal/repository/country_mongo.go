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

type CountryMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewCountryMongo(db *mongo.Database, i18n config.I18nConfig) *CountryMongo {
	return &CountryMongo{db: db, i18n: i18n}
}

func (r *CountryMongo) CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error) {
	var result domain.Country

	collection := r.db.Collection(TblCountry)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }
	// count, err := r.db.Collection(tblCountry).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return response, err
	// }
	// newId := count + 1

	newPage := domain.Country{
		Publish:   data.Publish,
		Flag:      data.Flag,
		Name:      data.Name,
		Code:      data.Code,
		SortOrder: data.SortOrder,
		Image:     data.Image,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newPage)
	if err != nil {
		return result, err
	}

	err = r.db.Collection(TblCountry).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *CountryMongo) GetCountry(id string) (domain.Country, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Country

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Country{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(TblCountry).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Country{}, err
	}

	return result, nil
}

func (r *CountryMongo) FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Country
	var response domain.Response[domain.Country]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Country]{}, err
	}
	cursor, err := r.db.Collection(TblCountry).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Country, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// var options options.CountOptions
	// // options.SetLimit(params.Limit)
	// options.SetSkip(params.Skip)
	count, err := r.db.Collection(TblCountry).CountDocuments(ctx, params.Filter) // , &options
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Country]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *CountryMongo) UpdateCountry(id string, data interface{}) (domain.Country, error) {
	var result domain.Country
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblCountry)

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

func (r *CountryMongo) DeleteCountry(id string) (domain.Country, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Country{}
	collection := r.db.Collection(TblCountry)

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
