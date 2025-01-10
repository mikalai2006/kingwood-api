package service

import (
	"time"

	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/auths"
	"github.com/mikalai2006/kingwood-api/pkg/hasher"
)

type Action interface {
	FindAction(params domain.RequestParams) (domain.Response[model.Action], error)
	GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error)
	CreateAction(userID string, data *model.ActionInput) (*model.Action, error)
	UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error)
	DeleteAction(id string) (model.Action, error)
}

type Authorization interface {
	CreateAuth(auth *domain.AuthInput) (string, error)
	GetAuth(id string) (domain.Auth, error)
	SignIn(input *domain.AuthInput) (domain.ResponseTokens, error)
	ExistAuth(auth *domain.AuthInput) (domain.Auth, error)
	CreateSession(auth *domain.Auth) (domain.ResponseTokens, error)
	VerificationCode(userID string, code string) error
	RefreshTokens(refreshToken string) (domain.ResponseTokens, error)
	RemoveRefreshTokens(refreshToken string) (string, error)
	UpdateAuth(id string, auth *domain.AuthInput) (domain.Auth, error)
	ResetPassword(authID string) (string, error)
}

type Product interface {
	FindProduct(params *model.ProductFilter) (domain.Response[model.Product], error)
	CreateProduct(userID string, node *model.ProductInputData) (*model.Product, error)
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
	CreateOffer(userID string, data *model.OfferInput) (*model.Offer, error)
	FindOffer(params *model.OfferFilter) (domain.Response[model.Offer], error)
	UpdateOffer(id string, userID string, data *model.Offer) (*model.Offer, error)
	DeleteOffer(id string) (model.Offer, error)
}

type Order interface {
	CreateOrder(userID string, data *domain.Order) (*domain.Order, error)
	FindOrder(input *domain.OrderFilter) (domain.Response[domain.Order], error)
	UpdateOrder(id string, userID string, data *domain.OrderInput) (*domain.Order, error)
	DeleteOrder(id string) (*domain.Order, error)

	GetAllOrder(params domain.RequestParams) (domain.Response[domain.Order], error)
}

type Operation interface {
	CreateOperation(userID string, data *domain.Operation) (*domain.Operation, error)
	FindOperation(params domain.RequestParams) (domain.Response[domain.Operation], error)
	UpdateOperation(id string, userID string, data *domain.OperationInput) (*domain.Operation, error)
	DeleteOperation(id string) (*domain.Operation, error)
}

type Object interface {
	CreateObject(userID string, data *domain.Object) (*domain.Object, error)
	FindObject(input *domain.ObjectFilter) (domain.Response[domain.Object], error)
	UpdateObject(id string, userID string, data *domain.ObjectInput) (*domain.Object, error)
	DeleteObject(id string) (*domain.Object, error)
}

type Task interface {
	CreateTask(userID string, data *domain.Task) (*domain.Task, error)
	FindTask(params domain.RequestParams) (domain.Response[domain.Task], error)
	FindTaskPopulate(filter domain.TaskFilter) (domain.Response[domain.Task], error)
	UpdateTask(id string, userID string, data *domain.TaskInput) (*domain.Task, error)
	DeleteTask(id string) (*domain.Task, error)
}

type TaskMontaj interface {
	CreateTaskMontaj(userID string, data *domain.TaskMontaj) (*domain.TaskMontaj, error)
	FindTaskMontaj(input domain.TaskMontajFilter) (domain.Response[domain.TaskMontaj], error)
	UpdateTaskMontaj(id string, userID string, data *domain.TaskMontajInput) (*domain.TaskMontaj, error)
	DeleteTaskMontaj(id string) (*domain.TaskMontaj, error)
}

type WorkTime interface {
	CreateWorkTime(userID string, data *domain.WorkTime) (*domain.WorkTime, error)
	FindWorkTime(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error)
	FindWorkTimePopulate(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error)
	UpdateWorkTime(id string, userID string, data *domain.WorkTimeInput) (*domain.WorkTime, error)
	DeleteWorkTime(id string) (*domain.WorkTime, error)
}

type TaskHistory interface {
	CreateTaskHistory(userID string, data *domain.TaskHistory) (*domain.TaskHistory, error)
	FindTaskHistory(input domain.TaskHistoryFilter) (domain.Response[domain.TaskHistory], error)
	FindTaskHistoryPopulate(input domain.TaskHistoryFilter) (domain.Response[domain.TaskHistory], error)
	UpdateTaskHistory(id string, userID string, data *domain.TaskHistoryInput) (*domain.TaskHistory, error)
	DeleteTaskHistory(id string) (*domain.TaskHistory, error)
}

type TaskMontajWorker interface {
	FindTaskMontajWorkerPopulate(input *domain.TaskMontajWorkerFilter) (domain.Response[domain.TaskMontajWorker], error)
	CreateTaskMontajWorker(userID string, data *domain.TaskMontajWorker) (*domain.TaskMontajWorker, error)
	FindTaskMontajWorker(params domain.RequestParams) (domain.Response[domain.TaskMontajWorker], error)
	UpdateTaskMontajWorker(id string, userID string, data *domain.TaskMontajWorkerInput) (*domain.TaskMontajWorker, error)
	DeleteTaskMontajWorker(id string) (*domain.TaskMontajWorker, error)
}

type TaskWorker interface {
	CreateTaskWorker(userID string, data *domain.TaskWorker, autoCreate int) (*domain.TaskWorker, error)
	FindTaskWorkerPopulate(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorker], error)
	// FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error)
	UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput, autoUpdate int) (*domain.TaskWorker, error)
	DeleteTaskWorker(id string) (*domain.TaskWorker, error)
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

type Pay interface {
	CreatePay(userID string, data *domain.Pay) (*domain.Pay, error)
	FindPay(params domain.RequestParams) (domain.Response[domain.Pay], error)
	UpdatePay(id string, userID string, data *domain.PayInput) (*domain.Pay, error)
	DeletePay(id string) (*domain.Pay, error)
}

type Image interface {
	CreateImage(userID string, data *domain.ImageInput) (domain.Image, error)
	GetImage(id string) (domain.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[domain.Image], error)
	DeleteImage(id string) (domain.Image, error)
}
type Country interface {
	CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error)
	GetCountry(id string) (domain.Country, error)
	FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error)
	UpdateCountry(id string, data interface{}) (domain.Country, error)
	DeleteCountry(id string) (domain.Country, error)
}

type Role interface {
	CreateRole(userID string, data *domain.RoleInput) (domain.Role, error)
	GetRole(id string) (domain.Role, error)
	FindRole(params domain.RequestParams) (domain.Response[domain.Role], error)
	UpdateRole(id string, data interface{}) (domain.Role, error)
	DeleteRole(id string) (domain.Role, error)
}

type Lang interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Question interface {
	FindQuestion(params *model.QuestionFilter) (domain.Response[model.Question], error)
	CreateQuestion(userID string, question *model.QuestionInput) (*model.Question, error)
	UpdateQuestion(id string, userID string, data *model.QuestionInput) (*model.Question, error)
	DeleteQuestion(id string, userID string) (model.Question, error)
}
type Ticket interface {
	FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error)
	CreateTicketMessage(userID string, message *model.TicketMessage) (*model.TicketMessage, error)
	DeleteTicket(id string) (model.Ticket, error)
}

type Post interface {
	FindPost(params domain.RequestParams) (domain.Response[domain.Post], error)
	GetAllPost(params domain.RequestParams) (domain.Response[domain.Post], error)
	CreatePost(userID string, Post *domain.Post) (*domain.Post, error)
	UpdatePost(id string, userID string, data *domain.PostInput) (*domain.Post, error)
	DeletePost(id string) (domain.Post, error)
}

type Services struct {
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
	TaskWorker
	TaskStatus
	TaskMontaj
	TaskMontajWorker
	Operation
	Pay
	Object
	WorkTime
	TaskHistory
}

type ConfigServices struct {
	Repositories           *repository.Repositories
	Hasher                 hasher.PasswordHasher
	TokenManager           auths.TokenManager
	OtpGenerator           utils.Generator
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeLength int
	I18n                   config.I18nConfig
	ImageService           config.IImageConfig
	Hub                    *Hub
}

func NewServices(cfgService *ConfigServices) *Services {
	User := NewUserService(cfgService.Repositories.User, cfgService.Hub)
	Authorization := NewAuthService(
		cfgService.Repositories.Authorization,
		cfgService.Hasher,
		cfgService.TokenManager,
		cfgService.RefreshTokenTTL,
		cfgService.AccessTokenTTL,
		cfgService.OtpGenerator,
		cfgService.VerificationCodeLength,
		User,
		cfgService.Hub,
	)
	Action := NewActionService(cfgService.Repositories.Action, cfgService.I18n)
	Post := NewPostService(cfgService.Repositories.Post, cfgService.I18n)
	TaskStatus := NewTaskStatusService(cfgService.Repositories.TaskStatus, cfgService.I18n)
	// Review := NewReviewService(cfgService.Repositories.Review)
	Lang := NewLangService(cfgService.Repositories, cfgService.I18n)
	Country := NewCountryService(cfgService.Repositories, cfgService.I18n)
	Role := NewRoleService(cfgService.Repositories, cfgService.I18n)
	Image := NewImageService(cfgService.Repositories.Image, cfgService.ImageService)
	Product := NewProductService(cfgService.Repositories.Product, User, cfgService.Hub)
	MessageRoom := NewMessageRoomService(cfgService.Repositories.MessageRoom, cfgService.Hub)
	Message := NewMessageService(cfgService.Repositories.Message, cfgService.Hub, MessageRoom)
	Question := NewQuestionService(cfgService.Repositories.Question, cfgService.Hub)
	Ticket := NewTicketService(cfgService.Repositories.Ticket)
	Operation := NewOperationService(cfgService.Repositories.Operation, User)
	Order := NewOrderService(cfgService.Repositories.Order, User, cfgService.Hub, Operation)
	WorkTime := NewWorkTimeService(cfgService.Repositories.WorkTime, cfgService.Hub, User, TaskStatus)
	Task := NewTaskService(cfgService.Repositories.Task, cfgService.Hub, User, TaskStatus, Order)
	TaskWorker := NewTaskWorkerService(cfgService.Repositories.TaskWorker, User, TaskStatus, Task, cfgService.Hub)
	TaskMontaj := NewTaskMontajService(cfgService.Repositories.TaskMontaj, cfgService.Hub, User, TaskStatus)
	TaskHistory := NewTaskHistoryService(cfgService.Repositories.TaskHistory, cfgService.Hub, User, TaskStatus)

	services := &Services{
		Authorization:    Authorization,
		Action:           Action,
		Post:             Post,
		User:             User,
		Lang:             Lang,
		Country:          Country,
		Image:            Image,
		Product:          Product,
		Message:          Message,
		MessageRoom:      MessageRoom,
		Offer:            NewOfferService(cfgService.Repositories.Offer, User, cfgService.Hub, Message, MessageRoom),
		Question:         Question,
		Ticket:           Ticket,
		Order:            Order,
		Task:             Task,
		TaskWorker:       TaskWorker,
		Operation:        Operation,
		Role:             Role,
		TaskStatus:       TaskStatus,
		TaskMontaj:       TaskMontaj,
		TaskMontajWorker: NewTaskMontajWorkerService(cfgService.Repositories.TaskMontajWorker, User, TaskStatus, Task, cfgService.Hub),
		Pay:              NewPayService(cfgService.Repositories.Pay, User, cfgService.Hub),
		Object:           NewObjectService(cfgService.Repositories.Object, cfgService.Hub, User),
		WorkTime:         WorkTime,
		TaskHistory:      TaskHistory,
	}
	Task.Services = services
	TaskWorker.Services = services
	return services
}
