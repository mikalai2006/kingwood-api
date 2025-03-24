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

type ArchiveObjectMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchiveObjectMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveObjectMongo {
	return &ArchiveObjectMongo{db: db, i18n: i18n}
}

func (r *ArchiveObjectMongo) CreateArchiveObject(userID string, data *domain.Object) (*domain.ArchiveObject, error) {
	var result *domain.ArchiveObject

	collection := r.db.Collection(TblArchiveObject)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	newArchiveObject := domain.ArchiveObjectInput{
		ID:        data.ID,
		UserID:    data.UserID,
		Name:      data.Name,
		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
		Meta:      domain.ArchiveMeta{Author: userIDPrimitive, CreatedAt: time.Now()},
	}

	res, err := collection.InsertOne(ctx, newArchiveObject)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblArchiveObject).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ArchiveObjectMongo) FindArchiveObject(input *domain.ArchiveObjectFilter) (domain.Response[domain.ArchiveObject], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	results := []domain.ArchiveObject{}
	var response domain.Response[domain.ArchiveObject]

	q := bson.D{}

	// Filter by substring name
	if input.Name != nil && *input.Name != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", *input.Name), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

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
	if input.Skip != nil {
		pipe = append(pipe, bson.D{{"$skip", input.Skip}})
		skip = *input.Skip
	}
	if input.Limit != nil {
		pipe = append(pipe, bson.D{{"$limit", input.Limit}})
		limit = *input.Limit
	}

	cursor, err := r.db.Collection(TblArchiveObject).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// count, err := r.db.Collection(TblArchiveObject).CountDocuments(ctx, pipe)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.ArchiveObject]{
		Total: int(0),
		Skip:  int(skip),
		Limit: int(limit),
		Data:  results,
	}
	return response, nil
}

func (r *ArchiveObjectMongo) DeleteArchiveObject(id string, userID string) (*domain.ArchiveObject, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.ArchiveObject{}
	collection := r.db.Collection(TblArchiveObject)

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
