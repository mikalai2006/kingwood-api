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

type TaskStatusMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTaskStatusMongo(db *mongo.Database, i18n config.I18nConfig) *TaskStatusMongo {
	return &TaskStatusMongo{db: db, i18n: i18n}
}

func (r *TaskStatusMongo) FindTaskStatus(params domain.RequestParams) (domain.Response[domain.TaskStatus], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.TaskStatus
	var response domain.Response[domain.TaskStatus]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.TaskStatus]{}, err
	// }
	// cursor, err := r.db.Collection(TblTaskStatus).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(tblTaskStatus).Aggregate(ctx, pipe)
	// fmt.Println("filter TaskStatus:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.TaskStatus, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblTaskStatus).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.TaskStatus]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TaskStatusMongo) CreateTaskStatus(userID string, data *domain.TaskStatus) (*domain.TaskStatus, error) {
	var result *domain.TaskStatus
	collection := r.db.Collection(tblTaskStatus)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIdPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	newTaskStatus := domain.TaskStatus{
		UserID:      userIdPrimitive,
		Name:        data.Name,
		Description: data.Description,
		Props:       data.Props,
		Color:       data.Color,
		Enabled:     data.Enabled,
		Icon:        data.Icon,
		Animate:     data.Animate,
		Status:      data.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newTaskStatus)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblTaskStatus).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TaskStatusMongo) UpdateTaskStatus(id string, userID string, data *domain.TaskStatusInput) (*domain.TaskStatus, error) {
	var result *domain.TaskStatus
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblTaskStatus)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}
	// _, err = collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
	// 	"seo":         data.Seo,
	// 	"title":       data.Title,
	// 	"description": data.Description,
	// 	"props":       data.Props,
	// 	"locale":      data.Locale,
	// 	"parent":      data.Parent,
	// 	"status":      data.Status,
	// 	"sort_order":  data.SortOrder,
	// 	"updated_at":  time.Now(),
	// }})
	// if err != nil {
	// 	return result, err
	// }
	newData := bson.M{}

	if data.Name != "" {
		newData["name"] = data.Name
	}
	if data.Color != "" {
		newData["color"] = data.Color
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if data.Icon != "" {
		newData["icon"] = data.Icon
	}
	if data.Animate != nil {
		newData["animate"] = data.Animate
	}
	if data.Start != nil {
		newData["start"] = data.Start
	}
	if data.Finish != nil {
		newData["finish"] = data.Finish
	}
	if data.Process != nil {
		newData["process"] = data.Process
	}
	if data.Status != "" {
		newData["status"] = data.Status
	}

	if data.Props != nil {
		//newProps := make(map[string]interface{})
		newProps := data.Props
		if val, ok := data.Props["status"]; ok {
			if val == -1.0 {
				newDel := make(map[string]interface{})
				newDel["user_id"] = userID
				newDel["del_at"] = time.Now()
				newProps["del"] = newDel
			}
		}
		newData["props"] = newProps
	}
	newData["updatedAt"] = time.Now()
	// test := model.ProductLike{}
	// if data.ProductLike != test {
	// 	newData["Product_like"] = data.ProductLike
	// }
	// if data.Status != 0 {
	// 	newData["status"] = data.Status
	// }
	// bson.M{
	// 	"lon":        data.Lon,
	// 	"lat":        data.Lat,
	// 	"type":       data.Type,
	// 	"osm_id":     data.OsmID,
	// 	"amenity_id": data.AmenityID,
	// 	"props":      data.Props,
	// 	"name":       data.Name,
	// 	"updated_at": time.Now(),
	// }
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *TaskStatusMongo) DeleteTaskStatus(id string) (domain.TaskStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.TaskStatus{}
	collection := r.db.Collection(tblTaskStatus)

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
