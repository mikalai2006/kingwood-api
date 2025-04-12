package repository

import (
	"context"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewAnalyticMongo(db *mongo.Database, i18n config.I18nConfig) *AnalyticMongo {
	return &AnalyticMongo{db: db, i18n: i18n}
}

func (r *AnalyticMongo) GetAnalytic() (domain.Analytic, error) {
	var result domain.Analytic

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// order.
	countOrder, err := r.db.Collection(TblOrder).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountOrder = countOrder

	// user.
	countUser, err := r.db.Collection(tblUsers).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountUser = countUser

	// workHistory.
	countWorkHistory, err := r.db.Collection(tblWorkHistory).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountWorkHistory = countWorkHistory

	// task.
	countTask, err := r.db.Collection(tblTask).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountTask = countTask

	// taskWorker.
	countTaskWorker, err := r.db.Collection(tblTaskWorker).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountTaskWorker = countTaskWorker

	// message.
	countMessage, err := r.db.Collection(TblMessage).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountMessage = countMessage

	// notify.
	countNotify, err := r.db.Collection(tblNotify).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountNotify = countNotify

	// pay.
	countPay, err := r.db.Collection(tblPay).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountPay = countPay

	// image.
	countImage, err := r.db.Collection(tblImage).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.CountImage = countImage

	// archive user.
	countArchiveUser, err := r.db.Collection(TblArchiveUser).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveUser = countArchiveUser

	// archive image.
	countArchiveImage, err := r.db.Collection(TblArchiveImage).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveImage = countArchiveImage

	// archive pay.
	countArchivePay, err := r.db.Collection(TblArchivePay).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchivePay = countArchivePay

	// archive notify.
	countArchiveNotify, err := r.db.Collection(TblArchiveNotify).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveNotify = countArchiveNotify

	// archive message.
	countArchiveMessage, err := r.db.Collection(TblArchiveMessage).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveMessage = countArchiveMessage

	// archive taskWorker.
	countArchiveTaskWorker, err := r.db.Collection(TblArchiveTaskWorker).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveTaskWorker = countArchiveTaskWorker

	// archive task.
	countArchiveTask, err := r.db.Collection(TblArchiveTask).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveTask = countArchiveTask

	// archive workHistory.
	countArchiveWorkHistory, err := r.db.Collection(TblArchiveWorkHistory).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveWorkHistory = countArchiveWorkHistory

	// archive order.
	countArchiveOrder, err := r.db.Collection(TblArchiveOrder).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Archive.CountArchiveOrder = countArchiveOrder

	// app error.
	appErrors, err := r.db.Collection(tblAppError).CountDocuments(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	result.Active.AppError = appErrors
	return result, nil
}
