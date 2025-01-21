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

	hidden := 0
	if data.Hidden != nil {
		hidden = *data.Hidden
	}

	newRole := domain.Role{
		Name:   data.Name,
		Code:   data.Code,
		Value:  data.Value,
		Hidden: hidden,
		// UserID:    userIDPrimitive,
		SortOrder: *data.SortOrder,
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

func (r *RoleMongo) FindRole(input *domain.RoleFilter) (domain.Response[domain.Role], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Role
	var response domain.Response[domain.Role]

	q := bson.D{}

	// pipe, err := CreatePipeline(params, &r.i18n)
	// if err != nil {
	// 	return domain.Response[domain.Role]{}, err
	// }
	// Filters
	if input.ID != nil && len(input.ID) > 0 {
		ids := []primitive.ObjectID{}
		for i, _ := range input.ID {
			iDPrimitive, err := primitive.ObjectIDFromHex(input.ID[i])
			if err != nil {
				return response, err
			}

			ids = append(ids, iDPrimitive)
		}

		q = append(q, bson.E{"_id", bson.D{{"$in", ids}}})
	}
	if input.Code != nil && len(input.Code) > 0 {
		ids := []string{}
		for i, _ := range input.Code {
			ids = append(ids, input.Code[i])
		}

		q = append(q, bson.E{"code", bson.D{{"$in", ids}}})
	}
	if input.Name != nil && len(input.Name) > 0 {
		ids := []string{}
		for i, _ := range input.Name {
			ids = append(ids, input.Code[i])
		}

		q = append(q, bson.E{"name", bson.D{{"$in", ids}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

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

	// var options options.CountOptions
	// // options.SetLimit(params.Limit)
	// options.SetSkip(params.Skip)
	// count, err := r.db.Collection(TblRole).CountDocuments(ctx, params.Filter, &options)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.Role]{
		Total: 0,
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *RoleMongo) UpdateRole(id string, data *domain.RoleInput) (domain.Role, error) {
	var result domain.Role
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblRole)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Code != "" {
		newData["code"] = data.Code
	}
	if data.Name != "" {
		newData["name"] = data.Name
	}
	if data.Value != nil {
		newData["value"] = data.Value
	}
	if data.Hidden != nil {
		newData["hidden"] = data.Hidden
	}
	if data.SortOrder != nil {
		newData["sortOrder"] = data.SortOrder
	}

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
