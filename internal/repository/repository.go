package repository

import (
	"reflect"

	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Action interface {
	FindAction(params domain.RequestParams) (domain.Response[model.Action], error)
	GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error)
	CreateAction(userID string, tag *model.ActionInput) (*model.Action, error)
	UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error)
	DeleteAction(id string) (model.Action, error)
	GqlGetActions(params domain.RequestParams) ([]*model.Action, error)
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
}

type Product interface {
	FindProduct(params *model.ProductFilter) (domain.Response[model.Product], error)
	CreateProduct(userID string, product *model.ProductInputData) (*model.Product, error)
	UpdateProduct(id string, userID string, data *model.Product) (*model.Product, error)
	DeleteProduct(id string) (model.Product, error)
}

type Message interface {
	CreateMessage(userID string, message *model.MessageInput) (*model.Message, error)
	FindMessage(params *model.MessageFilter) (domain.Response[model.Message], error)
	UpdateMessage(id string, userID string, data *model.MessageInput) (*model.Message, error)
	DeleteMessage(id string) (model.Message, error)
	GetGroupForUser(userID string) ([]model.MessageGroupForUser, error)
}
type MessageRoom interface {
	CreateMessageRoom(userID string, message *model.MessageRoom) (*model.MessageRoom, error)
	FindMessageRoom(params *model.MessageRoomFilter) (domain.Response[model.MessageRoom], error)
	UpdateMessageRoom(id string, userID string, data *model.MessageRoom) (*model.MessageRoom, error)
	DeleteMessageRoom(id string) (model.MessageRoom, error)
	// GetGroupForUser(userID string) ([]model.MessageGroupForUser, error)
}

type Offer interface {
	FindOffer(params *model.OfferFilter) (domain.Response[model.Offer], error)
	GetOffer(id string) (*model.Offer, error)
	CreateOffer(userID string, data *model.OfferInput) (*model.Offer, error)
	UpdateOffer(id string, userID string, data *model.Offer) (*model.Offer, error)
	DeleteOffer(id string) (model.Offer, error)
}

type Order interface {
	FindOrder(input *domain.OrderFilter) (domain.Response[domain.Order], error)
	// GetAllOrder(params domain.RequestParams) (domain.Response[domain.Order], error)
	CreateOrder(userID string, Order *domain.Order) (*domain.Order, error)
	UpdateOrder(id string, userID string, data *domain.OrderInput) (*domain.Order, error)
	DeleteOrder(id string) (*domain.Order, error)
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
type TaskMontajWorker interface {
	FindTaskMontajWorkerPopulate(input *domain.TaskMontajWorkerFilter) (domain.Response[domain.TaskMontajWorker], error)
	FindTaskMontajWorker(params domain.RequestParams) (domain.Response[domain.TaskMontajWorker], error)
	CreateTaskMontajWorker(userID string, Order *domain.TaskMontajWorker) (*domain.TaskMontajWorker, error)
	UpdateTaskMontajWorker(id string, userID string, data *domain.TaskMontajWorkerInput) (*domain.TaskMontajWorker, error)
	DeleteTaskMontajWorker(id string) (*domain.TaskMontajWorker, error)
}

type TaskMontaj interface {
	FindTaskMontaj(input domain.TaskMontajFilter) (domain.Response[domain.TaskMontaj], error)
	FindTaskPopulate(params domain.TaskMontajFilter) (domain.Response[domain.TaskMontaj], error)
	CreateTaskMontaj(userID string, Order *domain.TaskMontaj) (*domain.TaskMontaj, error)
	UpdateTaskMontaj(id string, userID string, data *domain.TaskMontajInput) (*domain.TaskMontaj, error)
	DeleteTaskMontaj(id string) (*domain.TaskMontaj, error)
}

type WorkTime interface {
	FindWorkTime(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error)
	FindWorkTimePopulate(params domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error)
	CreateWorkTime(userID string, Order *domain.WorkTime) (*domain.WorkTime, error)
	UpdateWorkTime(id string, userID string, data *domain.WorkTimeInput) (*domain.WorkTime, error)
	DeleteWorkTime(id string) (*domain.WorkTime, error)
}

type TaskHistory interface {
	FindTaskHistory(input domain.TaskHistoryFilter) (domain.Response[domain.TaskHistory], error)
	FindTaskHistoryPopulate(params domain.TaskHistoryFilter) (domain.Response[domain.TaskHistory], error)
	CreateTaskHistory(userID string, Order *domain.TaskHistory) (*domain.TaskHistory, error)
	UpdateTaskHistory(id string, userID string, data *domain.TaskHistoryInput) (*domain.TaskHistory, error)
	DeleteTaskHistory(id string) (*domain.TaskHistory, error)
}

type Object interface {
	FindObject(input *domain.ObjectFilter) (domain.Response[domain.Object], error)
	CreateObject(userID string, Order *domain.Object) (*domain.Object, error)
	UpdateObject(id string, userID string, data *domain.ObjectInput) (*domain.Object, error)
	DeleteObject(id string) (*domain.Object, error)
}

type TaskWorker interface {
	FindTaskWorkerPopulate(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorker], error)
	// FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error)
	CreateTaskWorker(userID string, Order *domain.TaskWorker) (*domain.TaskWorker, error)
	UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput) (*domain.TaskWorker, error)
	DeleteTaskWorker(id string) (*domain.TaskWorker, error)
}

type Pay interface {
	FindPay(params domain.RequestParams) (domain.Response[domain.Pay], error)
	CreatePay(userID string, Order *domain.Pay) (*domain.Pay, error)
	UpdatePay(id string, userID string, data *domain.PayInput) (*domain.Pay, error)
	DeletePay(id string) (*domain.Pay, error)
}

type TaskStatus interface {
	FindTaskStatus(params domain.RequestParams) (domain.Response[domain.TaskStatus], error)
	CreateTaskStatus(userID string, data *domain.TaskStatus) (*domain.TaskStatus, error)
	UpdateTaskStatus(id string, userID string, data *domain.TaskStatusInput) (*domain.TaskStatus, error)
	DeleteTaskStatus(id string) (domain.TaskStatus, error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(params domain.RequestParams) (domain.Response[domain.User], error)
	CreateUser(userID string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user *domain.UserInput) (domain.User, error)
	Iam(userID string) (domain.User, error)
}

type Image interface {
	CreateImage(userID string, data *domain.ImageInput) (domain.Image, error)
	GetImage(id string) (domain.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[domain.Image], error)
	DeleteImage(id string) (domain.Image, error)
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
	FindRole(params domain.RequestParams) (domain.Response[domain.Role], error)
	UpdateRole(id string, data interface{}) (domain.Role, error)
	DeleteRole(id string) (domain.Role, error)
}

type Country interface {
	CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error)
	GetCountry(id string) (domain.Country, error)
	FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error)
	UpdateCountry(id string, data interface{}) (domain.Country, error)
	DeleteCountry(id string) (domain.Country, error)
}

type Question interface {
	FindQuestion(params *model.QuestionFilter) (domain.Response[model.Question], error)
	CreateQuestion(userID string, data *model.QuestionInput) (*model.Question, error)
	UpdateQuestion(id string, userID string, data *model.QuestionInput) (*model.Question, error)
	DeleteQuestion(id string) (model.Question, error)
}

type Ticket interface {
	FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error)
	CreateTicketMessage(userID string, message *model.TicketMessage) (*model.TicketMessage, error)
	DeleteTicket(id string) (model.Ticket, error)
	GqlGetTickets(params domain.RequestParams) ([]*model.Ticket, error)
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
	Action
	Post
	Authorization
	Lang
	Country
	Role

	Image
	Order
	User
	Product
	Message
	MessageRoom
	Offer
	Question
	Ticket
	Task
	TaskMontaj
	TaskStatus
	TaskWorker
	TaskMontajWorker
	Operation
	Pay
	Object
	WorkTime
	TaskHistory
}

func NewRepositories(mongodb *mongo.Database, i18n config.I18nConfig) *Repositories {
	return &Repositories{
		Action:        NewActionMongo(mongodb, i18n),
		Post:          NewPostMongo(mongodb, i18n),
		Authorization: NewAuthMongo(mongodb),
		Lang:          NewLangMongo(mongodb, i18n),
		Country:       NewCountryMongo(mongodb, i18n),
		Role:          NewRoleMongo(mongodb, i18n),
		Image:         NewImageMongo(mongodb, i18n),
		Order:         NewOrderMongo(mongodb, i18n),
		User:          NewUserMongo(mongodb, i18n),
		Product:       NewProductMongo(mongodb, i18n),
		Message:       NewMessageMongo(mongodb, i18n),
		MessageRoom:   NewMessageRoomMongo(mongodb, i18n),
		// MessageImage:  NewMessageImageMongo(mongodb, i18n),
		Offer:            NewOfferMongo(mongodb, i18n),
		Question:         NewQuestionMongo(mongodb, i18n),
		Ticket:           NewTicketMongo(mongodb, i18n),
		Task:             NewTaskMongo(mongodb, i18n),
		TaskWorker:       NewTaskWorkerMongo(mongodb, i18n),
		TaskStatus:       NewTaskStatusMongo(mongodb, i18n),
		TaskMontaj:       NewTaskMontajMongo(mongodb, i18n),
		TaskMontajWorker: NewTaskMontajWorkerMongo(mongodb, i18n),
		Operation:        NewOperationMongo(mongodb, i18n),
		Pay:              NewPayMongo(mongodb, i18n),
		Object:           NewObjectMongo(mongodb, i18n),
		WorkTime:         NewWorkTimeMongo(mongodb, i18n),
		TaskHistory:      NewTaskHistoryMongo(mongodb, i18n),
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
