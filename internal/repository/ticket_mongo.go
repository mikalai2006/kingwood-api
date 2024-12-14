package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTicketMongo(db *mongo.Database, i18n config.I18nConfig) *TicketMongo {
	return &TicketMongo{db: db, i18n: i18n}
}

func (r *TicketMongo) FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Ticket
	var response domain.Response[model.Ticket]
	// filter, opts, err := CreateFilterAndOptions(params)
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Ticket]{}, err
	}
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usera",
		"let":  bson.D{{Key: "userId", Value: "$user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": tblImage,
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})

	// get messages for ticket.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": TblTicketMessage,
			"as":   "messages",
			"let":  bson.D{{Key: "ticketId", Value: "$_id"}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$ticket_id", "$$ticketId"}}}}},
				bson.D{{
					Key: "$lookup",
					Value: bson.M{
						"from": tblImage,
						"as":   "images",
						"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
						"pipeline": mongo.Pipeline{
							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
						},
					},
				}},
				bson.D{{
					Key: "$lookup",
					Value: bson.M{
						"from": "users",
						"as":   "userx",
						"let":  bson.D{{Key: "userId", Value: "$user_id"}},
						"pipeline": mongo.Pipeline{
							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
							bson.D{{"$limit", 1}},
						},
					},
				}},
				bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userx"}}}},
			},
		},
	}})

	cursor, err := r.db.Collection(TblTicket).Aggregate(ctx, pipe) //.Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Ticket, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTicket).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Ticket]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TicketMongo) GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Ticket
	var response domain.Response[model.Ticket]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Ticket]{}, err
	}

	cursor, err := r.db.Collection(TblTicket).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Ticket, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTicket).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Ticket]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TicketMongo) CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error) {
	var result *model.Ticket

	collection := r.db.Collection(TblTicket)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newTicket := model.TicketMongo{
		UserID:   userIDPrimitive,
		Title:    ticket.Title,
		Status:   ticket.Status,
		Progress: ticket.Progress,
		// Description: Ticket.Description,
		// Props:     Ticket.Props,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newTicket)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblTicket).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TicketMongo) CreateTicketMessage(userID string, ticket *model.TicketMessage) (*model.TicketMessage, error) {
	var result *model.TicketMessage

	collection := r.db.Collection(TblTicketMessage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newTicket := model.TicketMessageMongo{
		UserID:    userIDPrimitive,
		Text:      ticket.Text,
		Status:    ticket.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newTicket)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblTicketMessage).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TicketMongo) GqlGetTickets(params domain.RequestParams) ([]*model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Ticket
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	cursor, err := r.db.Collection(TblTicket).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Ticket, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *TicketMongo) DeleteTicket(id string) (model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Ticket{}
	collection := r.db.Collection(TblTicket)

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
