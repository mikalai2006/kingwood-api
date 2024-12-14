package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	db *mongo.Database
}

func NewAuthMongo(db *mongo.Database) *AuthMongo {
	return &AuthMongo{db: db}
}

func (r *AuthMongo) CreateAuth(user *domain.AuthInputMongo) (string, error) {
	var id string

	collection := r.db.Collection(TblAuth)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// data := domain.Auth{
	// 	Login:    user.Login,
	// 	Email:    user.Email,
	// 	Strategy: user.Strategy,
	// 	Password: user.Password,
	// }

	// if !user.Role.IsZero() {
	// 	roleID, err := primitive.ObjectIDFromHex(user.Role)
	// 	if err != nil {
	// 		return id, err
	// 	}
	// 	data.Role = roleID
	// }
	// newAuth := domain.AuthInputMongo{
	// 	RoleId:       user.RoleId,
	// 	PostId:       user.PostId,
	// 	Login:        user.Login,
	// 	Email:        user.Email,
	// 	Password:     user.Password,
	// 	Strategy:     user.Strategy,
	// 	Verification: user.Verification,
	// 	Session:      user.Session,
	// 	CreatedAt:    time.Now(),
	// 	UpdatedAt:    time.Now(),
	// }

	auth, err := collection.InsertOne(ctx, user)
	if err != nil {
		return id, err
	}
	id = auth.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func chooseProvider(auth *domain.AuthInput) bson.D {
	if auth.Strategy == "local" {
		var filter bson.A
		if auth.Email != "" {
			// add email filter
			filter = append(filter, bson.M{"email": auth.Email}, bson.M{"email": auth.Email})
		}
		if auth.Login != "" {
			// add email filter
			filter = append(filter, bson.M{"login": auth.Login}, bson.M{"email": auth.Login})
		}

		return bson.D{
			{Key: "$or", Value: filter},
			// {Key: "password", Value: auth.Password},
		}
	}

	// if auth.VkID != "" {
	// 	return bson.D{{Key: "vk_id", Value: auth.VkID}}
	// } else if auth.GoogleID != "" {
	// 	return bson.D{{Key: "google_id", Value: auth.GoogleID}}
	// }
	return bson.D{{Key: "vk_id", Value: "none"}}
}

func (r *AuthMongo) CheckExistAuth(auth *domain.AuthInput) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	filter := chooseProvider(auth)
	err := r.db.Collection(TblAuth).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = nil
		}
	}
	return user, err
}

func (r *AuthMongo) GetAuth(id string) (domain.Auth, error) {
	var auth domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	idUser, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return auth, err
	}

	// query := bson.M{"_id": idUser}
	// // fmt.Println("")
	// // fmt.Printf("GetAuth: query=%s", query)
	// err = r.db.Collection(TblAuth).FindOne(ctx, query).Decode(&auth)
	// if err != nil {
	// 	return auth, err
	// }

	q := bson.M{"_id": idUser}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})
	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         TblRole,
			"as":           "roles",
			"localField":   "roleId",
			"foreignField": "_id",
			// "let": bson.D{{Key: "roleId", Value: bson.D{{"$toString", "$roleId"}}}},
			// "pipeline": mongo.Pipeline{
			// 	bson.D{{Key: "$match", Value: bson.M{"$_id": bson.M{"$eq": [2]string{"$roleId", "$$_id"}}}}},
			// },
		},
	}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"roleObject": bson.M{"$first": "$roles"}}}})

	cursor, err := r.db.Collection(TblAuth).Aggregate(ctx, pipe) // .FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return auth, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if er := cursor.Decode(&auth); er != nil {
			return auth, er
		}
	}

	return auth, err
}

func (r *AuthMongo) GetByCredentials(auth *domain.AuthInput) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	filter := chooseProvider(auth)

	// pipe := mongo.Pipeline{}
	// pipe = append(pipe, bson.D{{"$match", filter}})
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from":         "users",
	// 	"as":           "user_data",
	// 	"localField":   "_id",
	// 	"foreignField": "user_id",
	// }}})
	// pipe = append(pipe, bson.D{{Key: "$unwind", Value: "$user_data"}})

	err := r.db.Collection(TblAuth).FindOne(ctx, filter).Decode(&user) // .Aggregate(ctx, pipe) //
	if err != nil {
		return user, err
	}

	if err := r.db.Collection(tblUsers).FindOneAndUpdate(ctx, bson.M{
		"userId": user.ID,
	}, bson.M{"$set": bson.M{"online": true}}).Decode(&user.UserData); err != nil {
		return user, err
	}

	// defer cursor.Close(ctx)

	// for cursor.Next(context.TODO()) {
	// 	if err := cursor.Decode(&user); err != nil {
	// 		return user, err
	// 	}
	// 	// fmt.Printf("userData %+v\n", user.UserData)
	// }
	// if err := cursor.Err(); err != nil {
	// 	return user, err
	// }
	return user, err
}

func (r *AuthMongo) SetSession(authID primitive.ObjectID, session domain.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	_, err := r.db.Collection(TblAuth).UpdateOne(
		ctx,
		bson.M{"_id": authID},
		bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}},
	)

	return err
}

func (r *AuthMongo) VerificationCode(userID, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	res, err := r.db.Collection(TblAuth).UpdateOne(ctx,
		bson.M{"verification.code": code, "_id": id},
		bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("failed verify")
	}

	return nil
}

func (r *AuthMongo) RefreshToken(refreshToken string) (domain.Auth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Auth

	pipe := mongo.Pipeline{}

	pipe = append(pipe)

	if err := r.db.Collection(TblAuth).FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&result); err != nil {
		return result, err
	}

	if err := r.db.Collection(tblUsers).FindOne(ctx, bson.M{
		"userId": result.ID,
	}).Decode(&result.UserData); err != nil {
		return result, err
	}

	return result, nil
}

func (r *AuthMongo) RemoveRefreshToken(refreshToken string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var auth domain.Auth
	var result = ""

	if err := r.db.Collection(TblAuth).FindOneAndUpdate(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}, bson.M{"$set": bson.M{"session.refreshToken": "", "session.expiresAt": time.Now()}}).Decode(&auth); err != nil {
		return result, err
	}

	if err := r.db.Collection(tblUsers).FindOneAndUpdate(ctx, bson.M{
		"userId": auth.ID,
	}, bson.M{"$set": bson.M{"online": false}}).Decode(&auth.UserData); err != nil {
		return result, err
	}

	return result, nil
}

func (r *AuthMongo) UpdateAuth(id string, auth *domain.AuthInput) (domain.User, error) {
	var result domain.User

	return result, nil
}
