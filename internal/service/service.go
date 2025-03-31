package service

import (
	"time"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/auths"
	"github.com/mikalai2006/kingwood-api/pkg/hasher"
)

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
	ResetPassword(authID string, userID string, input *domain.ResetPassword) (string, error)
}

type Message interface {
	CreateMessage(userID string, message *domain.MessageInput) (*domain.Message, error)
	FindMessage(params *domain.MessageFilter) (domain.Response[domain.Message], error)
	UpdateMessage(id string, userID string, data *domain.MessageInput) (*domain.Message, error)
	DeleteMessage(id string, userID string) (domain.Message, error)
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

type Order interface {
	CreateOrder(userID string, data *domain.Order) (*domain.Order, error)
	FindOrder(input *domain.OrderFilter) (domain.Response[domain.Order], error)
	UpdateOrder(id string, userID string, data *domain.OrderInput) (*domain.Order, error)
	DeleteOrder(id string, userID string) (*domain.Order, error)
}

type ArchiveOrder interface {
	CreateArchiveOrder(userID string, Order *domain.Order) (*domain.ArchiveOrder, error)
	FindArchiveOrder(input *domain.ArchiveOrderFilter) (domain.Response[domain.ArchiveOrder], error)
	DeleteArchiveOrder(id string, userID string) (*domain.ArchiveOrder, error)
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
	DeleteObject(id string, userID string) (*domain.Object, error)
}

type ArchiveObject interface {
	CreateArchiveObject(userID string, data *domain.Object) (*domain.ArchiveObject, error)
	FindArchiveObject(input *domain.ArchiveObjectFilter) (domain.Response[domain.ArchiveObject], error)
	DeleteArchiveObject(id string, userID string) (*domain.ArchiveObject, error)
}

type Task interface {
	CreateTask(userID string, data *domain.Task) (*domain.Task, error)
	FindTask(params domain.RequestParams) (domain.Response[domain.Task], error)
	FindTaskPopulate(filter domain.TaskFilter) (domain.Response[domain.Task], error)
	UpdateTask(id string, userID string, data *domain.TaskInput) (*domain.Task, error)
	DeleteTask(id string, userID string, checkStatus bool) (*domain.Task, error)
}

type ArchiveTask interface {
	CreateArchiveTask(userID string, data *domain.Task) (*domain.ArchiveTask, error)
	FindArchiveTask(params domain.ArchiveTaskFilter) (domain.Response[domain.ArchiveTask], error)
	DeleteArchiveTask(id string, userID string) (*domain.ArchiveTask, error)
}

type WorkHistory interface {
	CreateWorkHistory(userID string, data *domain.WorkHistory) (*domain.WorkHistory, error)
	FindWorkHistory(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error)
	FindWorkHistoryPopulate(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error)
	UpdateWorkHistory(id string, userID string, data *domain.WorkHistoryInput) (*domain.WorkHistory, error)
	DeleteWorkHistory(id string, userID string, createNotify bool) (*domain.WorkHistory, error)
	GetStatByOrder(input domain.WorkHistoryFilter) ([]domain.WorkHistoryStatByOrder, error)
}

type ArchiveWorkHistory interface {
	CreateArchiveWorkHistory(userID string, data *domain.WorkHistory) (*domain.ArchiveWorkHistory, error)
	FindArchiveWorkHistory(input domain.ArchiveWorkHistoryFilter) (domain.Response[domain.ArchiveWorkHistory], error)
	DeleteArchiveWorkHistory(id string, userID string) (*domain.ArchiveWorkHistory, error)
}

type TaskWorker interface {
	CreateTaskWorker(userID string, data *domain.TaskWorker, autoCreate int) (*domain.TaskWorker, error)
	FindTaskWorkerPopulate(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorker], error)
	// FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error)
	UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput, autoUpdate int) (*domain.TaskWorker, error)
	DeleteTaskWorker(id string, userID string, checkStatus bool) (*domain.TaskWorker, error)
}

type ArchiveTaskWorker interface {
	CreateArchiveTaskWorker(userID string, data *domain.TaskWorker) (*domain.ArchiveTaskWorker, error)
	FindArchiveTaskWorker(input *domain.ArchiveTaskWorkerFilter) (domain.Response[domain.ArchiveTaskWorker], error)
	DeleteArchiveTaskWorker(id string, userID string) (*domain.ArchiveTaskWorker, error)
}

type Notify interface {
	CreateNotify(userID string, data *domain.NotifyInput) (*domain.Notify, error)
	FindNotifyPopulate(input *domain.NotifyFilter) (domain.Response[domain.Notify], error)
	UpdateNotify(id string, userID string, data *domain.NotifyInput) (*domain.Notify, error)
	DeleteNotify(id string) (*domain.Notify, error)
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
	DeleteUser(id string, userID string) (domain.User, error)
	UpdateUser(id string, user *domain.UserInput) (domain.User, error)
	Iam(userID string) (domain.User, error)
}

type Pay interface {
	CreatePay(userID string, data *domain.Pay) (*domain.Pay, error)
	FindPay(input *domain.PayFilter) (domain.Response[domain.Pay], error)
	UpdatePay(id string, userID string, data *domain.PayInput) (*domain.Pay, error)
	DeletePay(id string, userID string) (*domain.Pay, error)
}

type AppError interface {
	CreateAppError(userID string, data *domain.AppError) (*domain.AppError, error)
	FindAppError(input *domain.AppErrorFilter) (domain.Response[domain.AppError], error)
	UpdateAppError(id string, userID string, data *domain.AppErrorInput) (*domain.AppError, error)
	DeleteAppError(id string, userID string) (*domain.AppError, error)
}

type PayTemplate interface {
	FindPayTemplate(params domain.RequestParams) (domain.Response[domain.PayTemplate], error)
	CreatePayTemplate(userID string, data *domain.PayTemplate) (*domain.PayTemplate, error)
	UpdatePayTemplate(id string, userID string, data *domain.PayTemplateInput) (*domain.PayTemplate, error)
	DeletePayTemplate(id string) (domain.PayTemplate, error)
}

type Image interface {
	CreateImage(userID string, data *domain.ImageInput) (domain.Image, error)
	GetImage(id string) (domain.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[domain.Image], error)
	DeleteImage(userID string, id string) (domain.Image, error)
}

type ArchiveImage interface {
	CreateArchiveImage(userID string, data *domain.Image) (domain.ArchiveImage, error)
	FindArchiveImage(params domain.RequestParams) (domain.Response[domain.ArchiveImage], error)
	DeleteArchiveImage(id string) (domain.ArchiveImage, error)
}

type Role interface {
	CreateRole(userID string, data *domain.RoleInput) (domain.Role, error)
	GetRole(id string) (domain.Role, error)
	FindRole(filter *domain.RoleFilter) (domain.Response[domain.Role], error)
	UpdateRole(id string, data *domain.RoleInput) (domain.Role, error)
	DeleteRole(id string) (domain.Role, error)
}

type Lang interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Post interface {
	FindPost(params domain.RequestParams) (domain.Response[domain.Post], error)
	GetAllPost(params domain.RequestParams) (domain.Response[domain.Post], error)
	CreatePost(userID string, Post *domain.Post) (*domain.Post, error)
	UpdatePost(id string, userID string, data *domain.PostInput) (*domain.Post, error)
	DeletePost(id string) (domain.Post, error)
}

type Services struct {
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
	TaskWorker
	TaskStatus
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
	Post := NewPostService(cfgService.Repositories.Post, cfgService.I18n)
	TaskStatus := NewTaskStatusService(cfgService.Repositories.TaskStatus, cfgService.I18n)
	Lang := NewLangService(cfgService.Repositories, cfgService.I18n)
	Role := NewRoleService(cfgService.Repositories, cfgService.I18n)
	Image := NewImageService(cfgService.Repositories.Image, cfgService.ImageService)
	MessageStatus := NewMessageStatusService(cfgService.Repositories.MessageStatus, cfgService.Hub)
	Message := NewMessageService(cfgService.Repositories.Message, cfgService.Hub, cfgService.ImageService)
	Operation := NewOperationService(cfgService.Repositories.Operation, User)
	Order := NewOrderService(cfgService.Repositories.Order, User, cfgService.Hub, Operation)
	Task := NewTaskService(cfgService.Repositories.Task, cfgService.Hub, User, TaskStatus, Order)
	TaskWorker := NewTaskWorkerService(cfgService.Repositories.TaskWorker, User, TaskStatus, Task, cfgService.Hub)
	WorkHistory := NewWorkHistoryService(cfgService.Repositories.WorkHistory, cfgService.Hub)
	Notify := NewNotifyService(cfgService.Repositories.Notify, cfgService.Hub)
	Pay := NewPayService(cfgService.Repositories.Pay, cfgService.Hub)
	Object := NewObjectService(cfgService.Repositories.Object, cfgService.Hub)

	ArchiveOrder := NewArchiveOrderService(cfgService.Repositories.ArchiveOrder)
	ArchiveTask := NewArchiveTaskService(cfgService.Repositories.ArchiveTask)
	ArchiveTaskWorker := NewArchiveTaskWorkerService(cfgService.Repositories.ArchiveTaskWorker, cfgService.Hub)
	ArchiveWorkHistory := NewArchiveWorkHistoryService(cfgService.Repositories.ArchiveWorkHistory, cfgService.Hub)
	ArchiveImage := NewArchiveImageService(cfgService.Repositories.ArchiveImage, Image.imageConfig)
	ArchiveMessage := NewArchiveMessageService(cfgService.Repositories.ArchiveMessage, cfgService.Hub, Image.imageConfig)
	ArchiveObject := NewArchiveObjectService(cfgService.Repositories.ArchiveObject, cfgService.Hub)

	services := &Services{
		AppError:      NewAppErrorService(cfgService.Repositories.AppError, cfgService.Hub),
		Authorization: Authorization,
		Post:          Post,
		User:          User,
		Lang:          Lang,
		Image:         Image,
		Message:       Message,
		MessageStatus: MessageStatus,
		Order:         Order,
		Task:          Task,
		TaskWorker:    TaskWorker,
		Operation:     Operation,
		Role:          Role,
		TaskStatus:    TaskStatus,
		Pay:           Pay,
		PayTemplate:   NewPayTemplateService(cfgService.Repositories.PayTemplate, cfgService.I18n),
		Object:        Object,
		WorkHistory:   WorkHistory,
		Notify:        Notify,

		ArchiveOrder:       ArchiveOrder,
		ArchiveTask:        ArchiveTask,
		ArchiveTaskWorker:  ArchiveTaskWorker,
		ArchiveWorkHistory: ArchiveWorkHistory,
		ArchiveImage:       ArchiveImage,
		ArchiveMessage:     ArchiveMessage,
		ArchiveObject:      ArchiveObject,
	}
	Task.Services = services
	TaskWorker.Services = services
	Order.Services = services
	Notify.Services = services
	Pay.Services = services
	User.Services = services
	MessageStatus.Services = services
	Message.Services = services
	WorkHistory.Services = services
	Image.Services = services
	Object.Services = services

	ArchiveOrder.Services = services
	ArchiveTask.Services = services
	ArchiveTaskWorker.Services = services
	ArchiveWorkHistory.Services = services
	ArchiveMessage.Services = services
	ArchiveObject.Services = services

	return services
}
