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
	DeleteMessage(id string) (domain.Message, error)
	GetGroupForUser(userID string) ([]domain.MessageGroupForUser, error)
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
	DeleteTask(id string, userID string, checkStatus bool) (*domain.Task, error)
}

type WorkTime interface {
	CreateWorkTime(userID string, data *domain.WorkTime) (*domain.WorkTime, error)
	FindWorkTime(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error)
	FindWorkTimePopulate(input domain.WorkTimeFilter) (domain.Response[domain.WorkTime], error)
	UpdateWorkTime(id string, userID string, data *domain.WorkTimeInput) (*domain.WorkTime, error)
	DeleteWorkTime(id string) (*domain.WorkTime, error)
}

type WorkHistory interface {
	CreateWorkHistory(userID string, data *domain.WorkHistory) (*domain.WorkHistory, error)
	FindWorkHistory(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error)
	FindWorkHistoryPopulate(input domain.WorkHistoryFilter) (domain.Response[domain.WorkHistory], error)
	UpdateWorkHistory(id string, userID string, data *domain.WorkHistoryInput) (*domain.WorkHistory, error)
	DeleteWorkHistory(id string) (*domain.WorkHistory, error)
	GetStatByOrder(input domain.WorkHistoryFilter) ([]domain.WorkHistoryStatByOrder, error)
}

type TaskWorker interface {
	CreateTaskWorker(userID string, data *domain.TaskWorker, autoCreate int) (*domain.TaskWorker, error)
	FindTaskWorkerPopulate(input *domain.TaskWorkerFilter) (domain.Response[domain.TaskWorker], error)
	// FindTaskWorker(params domain.RequestParams) (domain.Response[domain.TaskWorker], error)
	UpdateTaskWorker(id string, userID string, data *domain.TaskWorkerInput, autoUpdate int) (*domain.TaskWorker, error)
	DeleteTaskWorker(id string, userID string, checkStatus bool) (*domain.TaskWorker, error)
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
	DeleteImage(id string) (domain.Image, error)
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
	WorkTime
	WorkHistory
	Notify
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
	WorkTime := NewWorkTimeService(cfgService.Repositories.WorkTime, cfgService.Hub, User, TaskStatus)
	Task := NewTaskService(cfgService.Repositories.Task, cfgService.Hub, User, TaskStatus, Order)
	TaskWorker := NewTaskWorkerService(cfgService.Repositories.TaskWorker, User, TaskStatus, Task, cfgService.Hub)
	TaskHistory := NewWorkHistoryService(cfgService.Repositories.WorkHistory, cfgService.Hub, User, TaskStatus)
	Notify := NewNotifyService(cfgService.Repositories.Notify, cfgService.Hub)
	Pay := NewPayService(cfgService.Repositories.Pay, cfgService.Hub)

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
		Object:        NewObjectService(cfgService.Repositories.Object, cfgService.Hub, User),
		WorkTime:      WorkTime,
		WorkHistory:   TaskHistory,
		Notify:        Notify,
	}
	Task.Services = services
	TaskWorker.Services = services
	Order.Services = services
	Notify.Services = services
	Pay.Services = services
	User.Services = services
	MessageStatus.Services = services
	Message.Services = services
	WorkTime.Services = services

	return services
}
