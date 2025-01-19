package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerV1 struct {
	db           *mongo.Database
	repositories *repository.Repositories
	Services     *service.Services
	oauth        config.OauthConfig
	auth         config.AuthConfig
	i18n         config.I18nConfig
	imageConfig  config.IImageConfig
	hub          service.Hub
}

func NewHandler(services *service.Services, repositories *repository.Repositories, db *mongo.Database, oauth *config.OauthConfig, auth *config.AuthConfig, i18n *config.I18nConfig, imageConfig *config.IImageConfig, hub service.Hub) *HandlerV1 {
	return &HandlerV1{
		repositories: repositories,
		db:           db,
		Services:     services,
		oauth:        *oauth,
		auth:         *auth,
		i18n:         *i18n,
		imageConfig:  *imageConfig,
		hub:          hub,
	}
}

func (h *HandlerV1) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{

		h.registerAuth(v1)
		h.registerPost(v1)
		h.RegisterLang(v1)
		h.RegisterRole(v1)
		h.registerWs(v1)
		h.registerOperation(v1)
		h.registerTaskStatus(v1)

		authenticated := v1.Group("", h.SetUserFromRequest)
		{
			// h.registerAction(authenticated)
			h.registerAppError(authenticated)
			h.RegisterImage(authenticated)
			h.registerMessage(authenticated)
			h.registerMessageRoom(authenticated)
			h.registerOrder(authenticated)
			h.registerObject(authenticated)
			h.registerPay(authenticated)
			h.registerTask(authenticated)
			h.registerTaskWorker(authenticated)
			h.registerNotify(authenticated)
			h.RegisterUser(authenticated)
			h.registerWorkTime(authenticated)
			h.registerWorkHistory(authenticated)
			h.registerPayTemplate(authenticated)
		}

		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": "v1",
			})
		})
	}
}
