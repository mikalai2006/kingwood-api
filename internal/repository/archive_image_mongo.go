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

type ArchiveImageMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchiveImageMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveImageMongo {
	return &ArchiveImageMongo{db: db, i18n: i18n}
}

func (r *ArchiveImageMongo) CreateArchiveImage(userID string, data *domain.Image) (domain.ArchiveImage, error) {
	var result domain.ArchiveImage

	collection := r.db.Collection(TblArchiveImage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	newArchiveImage := domain.ArchiveImageInput{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:          data.ID,
		UserID:      data.UserID,
		Service:     data.Service,
		ServiceID:   data.ServiceID,
		Path:        data.Path,
		Title:       data.Title,
		Ext:         data.Ext,
		Dir:         data.Dir,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	res, err := collection.InsertOne(ctx, newArchiveImage)
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ArchiveImageMongo) FindArchiveImage(params domain.RequestParams) (domain.Response[domain.ArchiveImage], error) {
	var results []domain.ArchiveImage
	var response domain.Response[domain.ArchiveImage]

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblArchiveImage)

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.ArchiveImage]{}, err
	}

	cursor, err := collection.Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.ArchiveImage, len(results))
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	// options.SetSkip(params.Skip)
	count, err := collection.CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.ArchiveImage]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ArchiveImageMongo) DeleteArchiveImage(id string) (domain.ArchiveImage, error) {
	var result = domain.ArchiveImage{}

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblArchiveImage)

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
