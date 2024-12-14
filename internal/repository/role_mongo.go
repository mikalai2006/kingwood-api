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

type RoleMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewRoleMongo(db *mongo.Database, i18n config.I18nConfig) *RoleMongo {
	return &RoleMongo{db: db, i18n: i18n}
}

func (r *RoleMongo) CreateRole(userID string, data *domain.RoleInput) (domain.Role, error) {
	var result domain.Role

	collection := r.db.Collection(TblRole)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// userIDPrimitive, err := primitive.ObjectIDFromHex(data.UserID)
	// if err != nil {
	// 	return result, err
	// }
	// count, err := r.db.Collection(tblLanguage).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return response, err
	// }
	// newId := count + 1

	newRole := domain.Role{
		Name:  data.Name,
		Code:  data.Code,
		Value: data.Value,
		// UserID:    userIDPrimitive,
		SortOrder: data.SortOrder,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newRole)
	if err != nil {
		return result, err
	}

	err = r.db.Collection(TblRole).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *RoleMongo) GetRole(id string) (domain.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Role

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Role{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(TblRole).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Role{}, err
	}

	return result, nil
}

func (r *RoleMongo) FindRole(params domain.RequestParams) (domain.Response[domain.Role], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Role
	var response domain.Response[domain.Role]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Role]{}, err
	}
	cursor, err := r.db.Collection(TblRole).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Role, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	options.SetSkip(params.Skip)
	count, err := r.db.Collection(TblRole).CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Role]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *RoleMongo) UpdateRole(id string, data interface{}) (domain.Role, error) {
	var result domain.Role
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblRole)

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

func (r *RoleMongo) DeleteRole(id string) (domain.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Role{}
	collection := r.db.Collection(TblRole)

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
