package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/docs"
	"github.com/mikalai2006/kingwood-api/internal/config"
	v1 "github.com/mikalai2006/kingwood-api/internal/handler/v1"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/service"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	db           *mongo.Database
	repositories *repository.Repositories
	services     *service.Services
	oauth        config.OauthConfig
	auth         config.AuthConfig
	i18n         config.I18nConfig
	imageConfig  config.IImageConfig
	hub          service.Hub
}

func NewHandler(services *service.Services, repositories *repository.Repositories, mongoDB *mongo.Database, oauth *config.OauthConfig, auth *config.AuthConfig, i18n *config.I18nConfig, imageConfig *config.IImageConfig, hub service.Hub) *Handler {
	return &Handler{
		repositories: repositories,
		db:           mongoDB,
		services:     services,
		oauth:        *oauth,
		auth:         *auth,
		i18n:         *i18n,
		imageConfig:  *imageConfig,
		hub:          hub,
	}
}

func (h *Handler) InitRoutes(cfg *config.Config, mongoDB *mongo.Database) *gin.Engine {
	// appG := app.Gin{C: *gin.Context}
	router := gin.New() // New() // Default
	router.Use(
		// gzip.Gzip(gzip.DefaultCompression),
		gin.Recovery(),
		gin.Logger(),
		middleware.Cors,
		// middleware.JSONAppErrorReporter(),
	)
	// add swagger route
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}
	if cfg.Environment != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.GET("/", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, "API")
	})

	// router.GET("/room/:roomId", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", nil)
	// })
	// router.GET("/ws", func(c *gin.Context) {
	// 	socket.HandleConnections(c.Writer, c.Request)
	// })

	// router.GET("/ws/:roomId", func(c *gin.Context) {
	// 	roomId := c.Param("roomId")
	// 	websocket.ServeWS(c, roomId, &h.hub)
	// })

	// create session
	// store := cookie.NewStore([]byte(os.Getenv("secret")))
	// router.Use(sessions.Sessions("mysession", store))

	h.initAPI(router)

	router.NoRoute(func(c *gin.Context) {
		appG := app.Gin{C: c}
		// c.AbortWithError(http.StatusNotFound, errors.New("page not found"))
		appG.ResponseError(http.StatusNotFound, errors.New("page not found"), nil)
		// .SetMeta(gin.H{
		// 	"code": http.StatusNotFound,
		// 	"status": "error",
		// 	"message": "hello",
		// })
	})
	router.Static("/images", "./public")
	router.Static("/css", "./public/css")
	router.Static("/js", "./public/js")
	router.Static("/files", "./public/files")
	// router.GET("/debug/vars", expvar.Handler())

	// var upgrader = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }
	// router.GET("/ws", func(c *gin.Context) {
	// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	// 	if err != nil {
	// 		return
	// 	}
	// 	defer conn.Close()
	// 	var input interface{}
	// 	if er := c.BindJSON(&input); er != nil {
	// 		fmt.Println("Error: ", er)
	// 		return
	// 	}
	// 	conn.WriteMessage(websocket.TextMessage, []byte(input.(string)))
	// 	// for {
	// 	// 	conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	// 	// 	time.Sleep(time.Second)
	// 	// }
	// })
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	fmt.Println("IImageConfig", &h.imageConfig)
	api := router.Group("/api")
	api.Use(GetLang(&h.i18n))

	handlerV1 := v1.NewHandler(h.services, h.repositories, h.db, &h.oauth, &h.auth, &h.i18n, &h.imageConfig, h.hub)
	handlerV1.Init(api)
}

func GetLang(i18n *config.I18nConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			lang = i18n.Default
		}
		// i18n.Locale = lang
		c.Set("i18nLocale", lang)
		// fmt.Println("middleware h.i18n.Locale=", i18n.Locale)
		c.Next()
	}
}
