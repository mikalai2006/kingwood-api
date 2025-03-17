package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArchiveOrderMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewArchiveOrderMongo(db *mongo.Database, i18n config.I18nConfig) *ArchiveOrderMongo {
	return &ArchiveOrderMongo{db: db, i18n: i18n}
}

func (r *ArchiveOrderMongo) CreateArchiveOrder(userID string, data *domain.Order) (*domain.ArchiveOrder, error) {
	var result *domain.ArchiveOrder

	collection := r.db.Collection(TblArchiveOrder)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newArchiveOrder := domain.ArchiveOrderInput{
		Meta: domain.ArchiveMeta{
			Author:    userIDPrimitive,
			CreatedAt: time.Now(),
		},
		ID:              data.ID,
		UserID:          data.UserID,
		Name:            data.Name,
		Description:     data.Description,
		ObjectId:        data.ObjectId,
		Number:          int64(data.Number),
		ConstructorId:   data.ConstructorId,
		Priority:        data.Priority,
		Term:            data.Term,
		TermMontaj:      data.TermMontaj,
		Status:          data.Status,
		Group:           data.Group,
		StolyarComplete: data.StolyarComplete,
		MalyarComplete:  data.MalyarComplete,
		ShlifComplete:   data.ShlifComplete,
		GoComplete:      data.GoComplete,
		MontajComplete:  data.MontajComplete,
		Year:            *data.Year,
		DateStart:       data.DateStart,
		DateOtgruzka:    data.DateOtgruzka,

		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	res, err := collection.InsertOne(ctx, newArchiveOrder)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblArchiveOrder).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ArchiveOrderMongo) FindArchiveOrder(input *domain.ArchiveOrderFilter) (domain.Response[domain.ArchiveOrder], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var response domain.Response[domain.ArchiveOrder]

	q := bson.D{}
	// Filters
	if input.From != nil && !input.From.IsZero() {
		q = append(q, bson.E{"createdAt", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.From)}}})
	}
	if input.To != nil && !input.To.IsZero() {
		q = append(q, bson.E{"createdAt", bson.D{{"$lte", primitive.NewDateTimeFromTime(*input.To)}}})
	}
	if input.Date != nil && !input.Date.IsZero() {
		q = append(q, bson.E{"dateStart", bson.D{{"$gte", primitive.NewDateTimeFromTime(*input.Date)}}})
	}
	if input.Year != nil && *input.Year > 0 {
		q = append(q, bson.E{"year", input.Year})
	}
	if input.ID != nil && len(input.ID) > 0 {
		ids := []primitive.ObjectID{}
		for key, _ := range input.ID {
			idObjectPrimitive, err := primitive.ObjectIDFromHex(input.ID[key])
			if err != nil {
				return response, err
			}
			ids = append(ids, idObjectPrimitive)
		}
		q = append(q, bson.E{"_id", bson.D{{"$in", ids}}})
	}
	if input.Name != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", input.Name), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
	}
	if input.Group != nil && len(input.Group) > 0 {
		q = append(q, bson.E{"group", bson.M{"$elemMatch": bson.D{{"$in", input.Group}}}})
	}
	if input.Status != nil {
		q = append(q, bson.E{"status", input.Status})
	}
	if input.Number != nil {
		q = append(q, bson.E{"number", input.Number})
	}
	if input.Query != "" {
		qu := []interface{}{}
		i, err := strconv.Atoi(input.Query)
		if err == nil {
			qu = append(qu, bson.D{{"number", i}})
		}
		qu = append(qu, bson.D{{"name", bson.D{{"$regex", primitive.Regex(primitive.Regex{Pattern: input.Query, Options: "i"})}}}})

		q = append(q, bson.E{
			"$or",
			qu,
		})

	}
	if input.ObjectId != nil {
		objectIds := []primitive.ObjectID{}
		for key, _ := range input.ObjectId {
			idObjectPrimitive, err := primitive.ObjectIDFromHex(input.ObjectId[key])
			if err != nil {
				return response, err
			}
			objectIds = append(objectIds, idObjectPrimitive)
		}
		q = append(q, bson.E{"objectId", bson.D{{"$in", objectIds}}})
	}
	if input.StolyarComplete != nil {
		q = append(q, bson.E{"stolyarComplete", input.StolyarComplete})
	}
	if input.ShlifComplete != nil {
		q = append(q, bson.E{"shlifComplete", input.ShlifComplete})
	}
	if input.MalyarComplete != nil {
		q = append(q, bson.E{"malyarComplete", input.MalyarComplete})
	}
	if input.GoComplete != nil {
		q = append(q, bson.E{"goComplete", input.GoComplete})
	}
	if input.MontajComplete != nil {
		q = append(q, bson.E{"montajComplete", input.MontajComplete})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	// object.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": tblObject,
		"as":   "objecta",
		"let":  bson.D{{Key: "objectId", Value: "$objectId"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$objectId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"object": bson.M{"$first": "$objecta"}}}})

	// tasks.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": TblArchiveTask,
		"as":   "tasks",
		"let":  bson.D{{Key: "id", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{
				"$expr": bson.M{"$eq": [2]string{"$ArchiveOrderId", "$$id"}},
			}}},
			// workers.
			bson.D{{Key: "$lookup", Value: bson.M{
				"from":         tblTaskWorker,
				"as":           "workers",
				"localField":   "_id",
				"foreignField": "taskId",
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$lookup", Value: bson.M{
						"from": tblUsers,
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
						},
					}}},
					bson.D{{Key: "$set", Value: bson.M{"worker": bson.M{"$first": "$usera"}}}},
				},
			}}},
		},
	}}})

	if input.Sort != nil && len(input.Sort) > 0 {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
		}
		pipe = append(pipe, bson.D{{"$sort", sortParam}})
	}

	skip := 0
	limit := 10
	dataOptions := bson.A{}
	if input.Skip != nil {
		skip = *input.Skip
		dataOptions = append(dataOptions, bson.D{{"$skip", skip}})
	}
	if input.Limit != nil {
		limit = *input.Limit
		dataOptions = append(dataOptions, bson.D{{"$limit", limit}})
	}
	if input.Sort != nil {
		sortParam := bson.D{}
		for i := range input.Sort {
			sortParam = append(sortParam, bson.E{input.Sort[i].Key, input.Sort[i].Value})
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

	cursor, err := r.db.Collection(TblArchiveOrder).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	resultMap := []bson.M{}
	if er := cursor.All(ctx, &resultMap); er != nil {
		return response, er
	}
	resultFacetOne := domain.ResultFacetArchiveOrder{}
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

	response = domain.Response[domain.ArchiveOrder]{
		Total: total,
		Skip:  skip,
		Limit: limit,
		Data:  resultFacetOne.Data,
	}
	return response, nil
}

func (r *ArchiveOrderMongo) DeleteArchiveOrder(id string) (*domain.ArchiveOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &domain.ArchiveOrder{}
	collection := r.db.Collection(TblArchiveOrder)

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
