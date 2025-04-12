package repository

import (
	"reflect"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Analytic interface {
	GetAnalytic() (domain.Analytic, error)
}

type Authorization interface {
	CreateAuth(auth *domain.AuthInputMongo) (string, error)
	GetAuth(id string) (domain.Auth, error)
	CheckExistAuth(auth *domain.AuthInput) (domain.Auth, error)
	GetByCredentials(auth *domain.AuthInput) (domain.Auth, error)
	SetSession(authID primitive.ObjectID, session domain.Session) error
	VerificationCode(userID string, code string) error
	RefreshToken(refreshToken string) (domain.Auth, error)
	RemoveRefreshToken(refreshToken string) (string, error)
	UpdateAuth(id string, auth *domain.AuthInput) (domain.Auth, error)
	DeleteAuth(id string) (domain.Auth, error)
}

type Message interface {
	CreateMessage(userID string, message *domain.MessageInput) (*domain.Message, error)
	FindMessage(params *domain.MessageFilter) (domain.Response[domain.Message], error)
	UpdateMessage(id string, userID string, data *domain.MessageInput) (*domain.Message, error)
	DeleteMessage(id string) (domain.Message, error)
	GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error)
}

type ArchiveMessage interface {
	CreateArchiveMessage(userID string, message *domain.Message) (*domain.ArchiveMessage, error)
	FindArchiveMessage(params *domain.ArchiveMessageFilter) (domain.Response[domain.ArchiveMessage], error)
	DeleteArchiveMessage(id string) (domain.ArchiveMessage, error)
}

type MessageStatus interface {
	CreateMessageStatus(userID string, message *domain.MessageStatus) (*domain.MessageStatus, error)
	FindMessageStatus(params *domain.MessageStatusFilter) (domain.Response[domain.MessageStatus], error)
	UpdateMessageStatus(id string, userID string, data *domain.MessageStatus) (*domain.MessageStatus, error)
	DeleteMessageStatus(id string) (domain.MessageStatus, error)
	// GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error)
}

type Notify interface {
	CreateNotify(userID string, data *domain.NotifyInput) (*domain.Notify, error)
	FindNotifyPopulate(input *domain.NotifyFilter) (domain.Response[domain.Notify], error)
	UpdateNotify(id string, userID string, data *domain.NotifyInput) (*domain.Notify, error)
	DeleteNotify(id string) (*domain.Notify, error)
	ClearNotify(userID string) error
}

type ArchiveNotify interface {
	CreateArchiveNotify(userID string, data *domain.Notify) (*domain.ArchiveNotify, error)
	FindArchiveNotifyPopulate(input *domain.ArchiveNotifyFilter) (domain.Response[domain.ArchiveNotify], error)
	DeleteArchiveNotify(id string, userID string) (*domain.ArchiveNotify, error)
}

type Order interface {
	FindOrder(input *domain.OrderFilter) (domain.Response[domain.Order], error)
	// GetAllOrder(params domain.RequestParams) (domain.Response[domain.Order], error)
	CreateOrder(userID string, Order *domain.Order) (*domain.Order, error)
	UpdateOrder(id string, userID string, data *domain.OrderInput) (*domain.Order, error)
	DeleteOrder(id string) (*domain.Order, error)
}

type ArchiveOrder interface {
	CreateArchiveOrder(userID string, Order *domain.Order) (*domain.ArchiveOrder, error)
	FindArchiveOrder(input *domain.ArchiveOrderFilter) (domain.Response[domain.ArchiveOrder], error)
	DeleteArchiveOrder(id string) (*domain.ArchiveOrder, error)
}

type Operation interface {
	FindOperation(params domain.RequestParams) (domain.Response[domain.Operation], error)
	CreateOperation(userID string, data *domain.Operation) (*domain.Operation, error)
	UpdateOperation(id string, userID string, data *domain.OperationInput) (*domain.Operation, error)
	DeleteOperation(id string) (*domain.Operation, error)
}

type Task interface {
	FindTask(params domain.RequestParams) (domain.Response[domain.Task], error)
	FindTaskPopulate(input domain.TaskFilter) (domain.Response[domain.Task], error)
	CreateTask(userID string, Order *domain.Task) (*domain.Task, error)
	UpdateTask(id string, userID string, data *domain.TaskInput) (*domain.Task, error)
	DeleteTask(id string) (*domain.Task, error)
}

type ArchiveTask interface {
	CreateArchiveTask(userID string, Order *domain.Task) (*domain.ArchiveTask, error)
	FindArchiveTask(params domain.ArchiveTaskFilter) (domain.Response[domain.ArchiveTask], error)
	DeleteArchiveTask(id string) (*domain.ArchiveTask, error)
}

type WorkHistory interface {
	FindWorkHistory(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error)
	FindWorkHistoryPopulate(params domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error)
	CreateWorkHistory(userID string, Order *domain.WorkHistory) (*domain.WorkHistory, error)
	UpdateWorkHistory(id string, userID string, data *domain.WorkHistoryInput) (*domain.WorkHistory, error)
	DeleteWorkHistory(id string) (*domain.WorkHistory, error)
	GetStatByOrder(input domain.WorkHistoryFilter) ([]domain.WorkHistoryStatByOrder, error)
}

type ArchiveWorkHistory interface {
	CreateArchiveWorkHistory(userID string, Order *domain.WorkHistory) (*domain.ArchiveWorkHistory, error)
	FindArchiveWorkHistory(input domain.ArchiveWorkHistoryFilter) (domain.Response[domain.ArchiveWorkHistory], error)
	DeleteArchiveWorkHistory(id string) (*domain.ArchiveWorkHistory, error)
}

type Object interface {
	FindObject(input *domain.ObjectFilter) (domain.Response[domain.Object], error)
	CreateObject(userID string, Order *domain.Object) (*domain.Object, error)
	UpdateObject(id string, userID string, data *domain.ObjectInput) (*domain.Object, error)
	DeleteObject(id string, userID string) (*domain.Object, error)
}

type ArchiveObject interface {
	CreateArchiveObject(userID string, Order *domain.Object) (*domain.ArchiveObject, error)
	FindArchiveObject(input *domain.ArchiveObjectFilter) (domain.Response[domain.ArchiveObject], error)
	DeleteArchiveObject(id string, userID string) (*domain.ArchiveObject, error)
}

type TaskWorker interface {
	FindTaskWorkerPopulate(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorker], error)
	// FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error)
	CreateTaskWorker(userID string, Order *domain.TaskWorker) (*domain.TaskWorker, error)
	UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput) (*domain.TaskWorker, error)
	DeleteTaskWorker(id string) (*domain.TaskWorker, error)
}

type ArchiveTaskWorker interface {
	CreateArchiveTaskWorker(userID string, Order *domain.TaskWorker) (*domain.ArchiveTaskWorker, error)
	FindArchiveTaskWorker(input *domain.ArchiveTaskWorkerFilter) (domain.Response[domain.ArchiveTaskWorker], error)
	DeleteArchiveTaskWorker(id string) (*domain.ArchiveTaskWorker, error)
}

type Pay interface {
	FindPay(input *domain.PayFilter) (domain.Response[domain.Pay], error)
	CreatePay(userID string, Order *domain.Pay) (*domain.Pay, error)
	UpdatePay(id string, userID string, data *domain.PayInput) (*domain.Pay, error)
	DeletePay(id string, userID string) (*domain.Pay, error)
}

type ArchivePay interface {
	FindArchivePay(input *domain.ArchivePayFilter) (domain.Response[domain.ArchivePay], error)
	CreateArchivePay(userID string, Order *domain.Pay) (*domain.ArchivePay, error)
	DeleteArchivePay(id string, userID string) (*domain.ArchivePay, error)
}

type AppError interface {
	FindAppError(input *domain.AppErrorFilter) (domain.Response[domain.AppError], error)
	CreateAppError(userID string, Order *domain.AppError) (*domain.AppError, error)
	UpdateAppError(id string, userID string, data *domain.AppErrorInput) (*domain.AppError, error)
	DeleteAppError(id string, userID string) (*domain.AppError, error)
	ClearAppError(userID string) error
}

type PayTemplate interface {
	FindPayTemplate(params domain.RequestParams) (domain.Response[domain.PayTemplate], error)
	CreatePayTemplate(userID string, data *domain.PayTemplate) (*domain.PayTemplate, error)
	UpdatePayTemplate(id string, userID string, data *domain.PayTemplateInput) (*domain.PayTemplate, error)
	DeletePayTemplate(id string) (domain.PayTemplate, error)
}

type TaskStatus interface {
	FindTaskStatus(params domain.RequestParams) (domain.Response[domain.TaskStatus], error)
	CreateTaskStatus(userID string, data *domain.TaskStatus) (*domain.TaskStatus, error)
	UpdateTaskStatus(id string, userID string, data *domain.TaskStatusInput) (*domain.TaskStatus, error)
	DeleteTaskStatus(id string) (domain.TaskStatus, error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(filter *domain.UserFilter) (domain.Response[domain.User], error)
	CreateUser(userID string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user *domain.UserInput) (domain.User, error)
	Iam(userID string) (domain.User, error)
}

type ArchiveUser interface {
	CreateArchiveUser(userID string, user *domain.User) (*domain.ArchiveUser, error)
	FindArchiveUser(filter *domain.ArchiveUserFilter) (domain.Response[domain.ArchiveUser], error)
	DeleteArchiveUser(id string) (domain.ArchiveUser, error)
}

type Image interface {
	CreateImage(userID string, data *domain.ImageInput) (domain.Image, error)
	GetImage(id string) (domain.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params *domain.ImageFilter) (domain.Response[domain.Image], error)
	DeleteImage(id string) (domain.Image, error)
}

type ArchiveImage interface {
	CreateArchiveImage(userID string, data *domain.Image) (domain.ArchiveImage, error)
	FindArchiveImage(params *domain.ArchiveImageFilter) (domain.Response[domain.ArchiveImage], error)
	DeleteArchiveImage(id string) (domain.ArchiveImage, error)
}

type Lang interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Role interface {
	CreateRole(userID string, data *domain.RoleInput) (domain.Role, error)
	GetRole(id string) (domain.Role, error)
	FindRole(filter *domain.RoleFilter) (domain.Response[domain.Role], error)
	UpdateRole(id string, data *domain.RoleInput) (domain.Role, error)
	DeleteRole(id string) (domain.Role, error)
}

type Post interface {
	FindPost(params domain.RequestParams) (domain.Response[domain.Post], error)
	GetAllPost(params domain.RequestParams) (domain.Response[domain.Post], error)
	CreatePost(userID string, Post *domain.Post) (*domain.Post, error)
	UpdatePost(id string, userID string, data *domain.PostInput) (*domain.Post, error)
	DeletePost(id string) (domain.Post, error)
	GqlGetPosts(params domain.RequestParams) ([]*domain.Post, error)
}

type Repositories struct {
	AppError
	Post
	Authorization
	Lang
	Role
	Image
	Order
	User
	Message
	MessageStatus
	Task
	TaskStatus
	TaskWorker
	Operation
	Pay
	PayTemplate
	Object

	WorkHistory
	Notify

	ArchiveOrder
	ArchiveTask
	ArchiveTaskWorker
	ArchiveWorkHistory
	ArchiveImage
	ArchiveMessage
	ArchiveObject
	ArchiveNotify
	ArchiveUser
	ArchivePay

	Analytic
}

func NewRepositories(mongodb *mongo.Database, i18n config.I18nConfig) *Repositories {
	return &Repositories{
		AppError:      NewAppErrorMongo(mongodb, i18n),
		Post:          NewPostMongo(mongodb, i18n),
		Authorization: NewAuthMongo(mongodb),
		Lang:          NewLangMongo(mongodb, i18n),
		Role:          NewRoleMongo(mongodb, i18n),
		Image:         NewImageMongo(mongodb, i18n),
		Order:         NewOrderMongo(mongodb, i18n),
		User:          NewUserMongo(mongodb, i18n),
		Message:       NewMessageMongo(mongodb, i18n),
		MessageStatus: NewMessageStatusMongo(mongodb, i18n),
		Task:          NewTaskMongo(mongodb, i18n),
		TaskWorker:    NewTaskWorkerMongo(mongodb, i18n),
		TaskStatus:    NewTaskStatusMongo(mongodb, i18n),
		Operation:     NewOperationMongo(mongodb, i18n),
		Pay:           NewPayMongo(mongodb, i18n),
		PayTemplate:   NewPayTemplateMongo(mongodb, i18n),
		Object:        NewObjectMongo(mongodb, i18n),
		WorkHistory:   NewWorkHistoryMongo(mongodb, i18n),
		Notify:        NewNotifyMongo(mongodb, i18n),

		ArchiveOrder:       NewArchiveOrderMongo(mongodb, i18n),
		ArchiveTask:        NewArchiveTaskMongo(mongodb, i18n),
		ArchiveTaskWorker:  NewArchiveTaskWorkerMongo(mongodb, i18n),
		ArchiveWorkHistory: NewArchiveWorkHistoryMongo(mongodb, i18n),
		ArchiveImage:       NewArchiveImageMongo(mongodb, i18n),
		ArchiveMessage:     NewArchiveMessageMongo(mongodb, i18n),
		ArchiveObject:      NewArchiveObjectMongo(mongodb, i18n),
		ArchiveNotify:      NewArchiveNotifyMongo(mongodb, i18n),
		ArchiveUser:        NewArchiveUserMongo(mongodb, i18n),
		ArchivePay:         NewArchivePayMongo(mongodb, i18n),

		Analytic: NewAnalyticMongo(mongodb, i18n),
	}
}

// func getPaginationOpts(pagination *domain.PaginationQuery) *options.FindOptions {
// 	var opts *options.FindOptions
// 	if pagination != nil {
// 		opts = &options.FindOptions{
// 			Skip:  pagination.GetSkip(),
// 			Limit: pagination.GetLimit(),
// 		}
// 	}

// 	return opts
// }

func createFilter[V any](filterData *V) any {
	var filter V

	filterReflect := reflect.ValueOf(filterData)
	// fmt.Println("========== filterReflect ===========")
	// fmt.Println("struct > ", filterReflect)
	// fmt.Println("struct type > ", filterReflect.Type())
	filterIndirectData := reflect.Indirect(filterReflect)
	// fmt.Println("filter data > ", filterIndirectData)
	// fmt.Println("filter numField > ", filterIndirectData.NumField())
	dataFilter := bson.M{}

	var tagJSON, tagPrimitive string
	for i := 0; i < filterIndirectData.NumField(); i++ {
		field := filterIndirectData.Field(i)
		if field.Kind() == reflect.Ptr {
			field = reflect.Indirect(field)
		}
		typeField := filterIndirectData.Type().Field(i)
		tag := typeField.Tag
		// tagBson = tag.Get("bson")
		tagJSON = tag.Get("json")
		tagPrimitive = tag.Get("primitive")
		switch field.Kind() {
		case reflect.String:
			value := field.String()
			if tagPrimitive == "true" {
				id, _ := primitive.ObjectIDFromHex(value)
				// fmt.Println("===== string add ", tag, value)
				dataFilter[tagJSON] = id
			} else {
				dataFilter[tagJSON] = value
			}

		case reflect.Bool:
			value := field.Bool()
			dataFilter[tagJSON] = value

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := field.Int()
			dataFilter[tagJSON] = value

		default:

		}

		// fmt.Println(tagBson, tagJSON, tagPrimitive, fmt.Sprintf("[%s]", field), field.Kind(), field)
	}

	// structure := reflect.ValueOf(&filter)
	// fmt.Println("========== filter ===========")
	// fmt.Println("struct > ", structure)
	// fmt.Println("struct type > ", structure.Type())
	// fmt.Println("filter data > ", reflect.Indirect(structure))
	// fmt.Println("filter numField > ", reflect.Indirect(structure).NumField())

	// fmt.Println("========== result ===========")
	// fmt.Println("dataFilter > ", dataFilter)
	return filter
}
