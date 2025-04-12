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

type ArchivePayMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchivePayMongo(db *mongo.Database, i18n config.I18nConfig) *ArchivePayMongo {
	return &ArchivePayMongo{db: db, i18n: i18n}
}

func (r *ArchivePayMongo) FindArchivePay(input *domain.ArchivePayFilter) (domain.Response[domain.ArchivePay], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.ArchivePay
	var response domain.Response[domain.ArchivePay]

	// Filters
	q := bson.D{}
	if input.Month != nil {
		q = append(q, bson.E{"month", &input.Month})
	}
	if input.Year != nil {
		q = append(q, bson.E{"year", &input.Year})
	}
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
	if input.WorkerId != nil && len(input.WorkerId) > 0 {
		workerIds := []primitive.ObjectID{}
		for i, _ := range input.WorkerId {
			workerIDPrimitive, err := primitive.ObjectIDFromHex(input.WorkerId[i])
			if err != nil {
				return response, err
			}

			workerIds = append(workerIds, workerIDPrimitive)
		}

		q = append(q, bson.E{"workerId", bson.D{{"$in", workerIds}}})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	// worker.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": TblArchiveUser,
		"as":   "usera",
		"let":  bson.D{{Key: "workerId", Value: "$workerId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$workerId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": tblImage,
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$serviceId", "$$serviceId"}}}}},
					},
				},
			}},
			// // add populate auth.
			// bson.D{{
			// 	Key: "$lookup",
			// 	Value: bson.M{
			// 		"from":         TblAuth,
			// 		"as":           "auths",
			// 		"localField":   "userId",
			// 		"foreignField": "_id",
			// 		// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
			// 		// "pipeline": mongo.Pipeline{
			// 		// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
			// 		// },
			// 	},
			// }},
			bson.D{{Key: "$set", Value: bson.M{"auth": bson.M{"$first": "$auths"}}}},
			bson.D{{Key: "$set", Value: bson.M{"authPrivate": bson.M{"$first": "$auths"}}}},

			// post.
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": TblPost,
					"as":   "posts",
					// "localField":   "_id",
					// "foreignField": "serviceId",
					"let": bson.D{{Key: "postId", Value: "$postId"}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$postId"}}}}},
					},
				},
			}},
			bson.D{{Key: "$set", Value: bson.M{"postObject": bson.M{"$first": "$posts"}}}},
			// role.
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": TblRole,
				"as":   "rolea",
				// "localField":   "userId",
				// "foreignField": "_id",
				"let": bson.D{{Key: "roleId", Value: "$roleId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$roleId"}}}}},
					bson.D{{"$limit", 1}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$rolea"}}}},
		},
	}}},
		bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}},
	)

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

	cursor, err := r.db.Collection(TblArchivePay).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	response = domain.Response[domain.ArchivePay]{
		Total: 0,
		Skip:  skip,
		Limit: limit,
		Data:  results,
	}
	return response, nil
}

func (r *ArchivePayMongo) CreateArchivePay(userID string, data *domain.Pay) (*domain.ArchivePay, error) {
	var result *domain.ArchivePay

	collection := r.db.Collection(TblArchivePay)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existArchivePay domain.ArchivePay
	// r.db.Collection(TblArchivePay).FindOne(ctx, bson.M{"node_id": ArchivePay.NodeID, "userId": userIDPrimitive}).Decode(&existArchivePay)

	// if existArchivePay.NodeID.IsZero() {
	updatedAt := data.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	newArchivePay := domain.ArchivePayInput{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:        data.ID,
		Props:     data.Props,
		UserID:    data.UserID,
		Name:      data.Name,
		WorkerId:  data.WorkerId,
		Year:      &data.Year,
		Month:     &data.Month,
		Total:     data.Total,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	res, err := collection.InsertOne(ctx, newArchivePay)
	if err != nil {
		return nil, err
	}

	// err = r.db.Collection(tblArchivePay).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	// if err != nil {
	// 	return nil, err
	// }

	insertedID := res.InsertedID.(primitive.ObjectID)
	ArchivePays, err := r.FindArchivePay(&domain.ArchivePayFilter{ID: []string{insertedID.Hex()}})
	if err != nil {
		return nil, err
	}
	if len(ArchivePays.Data) > 0 {
		result = &ArchivePays.Data[0]
	}

	return result, nil
}

func (r *ArchivePayMongo) DeleteArchivePay(id string, userID string) (*domain.ArchivePay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.ArchivePay{}
	collection := r.db.Collection(TblArchivePay)

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
