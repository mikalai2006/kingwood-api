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

type PostMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewPostMongo(db *mongo.Database, i18n config.I18nConfig) *PostMongo {
	return &PostMongo{db: db, i18n: i18n}
}

func (r *PostMongo) FindPost(params domain.RequestParams) (domain.Response[domain.Post], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Post
	var response domain.Response[domain.Post]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.Post]{}, err
	// }
	// cursor, err := r.db.Collection(TblPost).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblPost).Aggregate(ctx, pipe)
	// fmt.Println("filter Post:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Post, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblPost).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Post]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *PostMongo) GetAllPost(params domain.RequestParams) (domain.Response[domain.Post], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Post
	var response domain.Response[domain.Post]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Post]{}, err
	}

	cursor, err := r.db.Collection(TblPost).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Post, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblPost).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Post]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *PostMongo) CreatePost(userID string, post *domain.Post) (*domain.Post, error) {
	var result *domain.Post
	collection := r.db.Collection(TblPost)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	idPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	newPost := domain.Post{
		UserID:      idPrimitive,
		Name:        post.Name,
		Description: post.Description,
		Props:       post.Props,
		Hidden:      post.Hidden,
		Color:       post.Color,
		SortOrder:   post.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newPost)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblPost).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PostMongo) GqlGetPosts(params domain.RequestParams) ([]*domain.Post, error) {
	fmt.Println("GqlGetPosts")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*domain.Post
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(TblPost).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*domain.Post, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *PostMongo) UpdatePost(id string, userID string, data *domain.PostInput) (*domain.Post, error) {
	var result *domain.Post
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblPost)

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

	if data.SortOrder != 0 {
		newData["sortOrder"] = data.SortOrder
	}
	if data.Name != "" {
		newData["name"] = data.Name
	}
	if data.Color != "" {
		newData["color"] = data.Color
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if data.Hidden != nil {
		newData["hidden"] = data.Hidden
	}

	if data.Props != nil {
		//newProps := make(map[string]interface{})
		newProps := data.Props
		if val, ok := data.Props["status"]; ok {
			if val == -1.0 {
				newDel := make(map[string]interface{})
				newDel["userId"] = userID
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

func (r *PostMongo) DeletePost(id string) (domain.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Post{}
	collection := r.db.Collection(TblPost)

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
