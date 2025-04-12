package domain

type AnalyticArchive struct {
	CountArchiveUser        int64 `json:"countArchiveUser" bson:"countArchiveUser"`
	CountArchiveWorkHistory int64 `json:"countArchiveWorkHistory" bson:"countArchiveWorkHistory"`
	CountArchiveTask        int64 `json:"countArchiveTask" bson:"countArchiveTask"`
	CountArchiveTaskWorker  int64 `json:"countArchiveTaskWorker" bson:"countArchiveTaskWorker"`
	CountArchiveOrder       int64 `json:"countArchiveOrder" bson:"countArchiveOrder"`
	CountArchiveNotify      int64 `json:"countArchiveNotify" bson:"countArchiveNotify"`
	CountArchiveMessage     int64 `json:"countArchiveMessage" bson:"countArchiveMessage"`
	CountArchivePay         int64 `json:"countArchivePay" bson:"countArchivePay"`
	CountArchiveImage       int64 `json:"countArchiveImage" bson:"countArchiveImage"`
}
type AnalyticActive struct {
	CountUser        int64 `json:"countUser" bson:"countUser"`
	CountOrder       int64 `json:"countOrder" bson:"countOrder"`
	CountTask        int64 `json:"countTask" bson:"countTask"`
	CountTaskWorker  int64 `json:"countTaskWorker" bson:"countTaskWorker"`
	CountNotify      int64 `json:"countNotify" bson:"countNotify"`
	CountMessage     int64 `json:"countMessage" bson:"countMessage"`
	CountPay         int64 `json:"countPay" bson:"countPay"`
	CountImage       int64 `json:"countImage" bson:"countImage"`
	CountWorkHistory int64 `json:"countWorkHistory" bson:"countWorkHistory"`
	AppError         int64 `json:"appError" bson:"appError"`
}
type Analytic struct {
	Active  AnalyticActive  `json:"active" bson:"active"`
	Archive AnalyticArchive `json:"archive" bson:"archive"`
}
