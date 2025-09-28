package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/handler"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/server"
	"github.com/mikalai2006/kingwood-api/internal/service"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/auths"
	"github.com/mikalai2006/kingwood-api/pkg/hasher"
	"github.com/mikalai2006/kingwood-api/pkg/logger"
	"github.com/sirupsen/logrus"
)

// @title Template API
// @version 1.0
// @description API Server for Template App

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func Run(configPath string) {
	// setting logrus
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// read config file
	cfg, err := config.Init(fmt.Sprintf("%sconfigs", configPath), fmt.Sprintf("%s.env", configPath))
	if err != nil {
		logger.Error(err)
		return
	}

	// initialize mongoDB
	mongoClient, err := repository.NewMongoDB(&repository.ConfigMongoDB{
		Host:     cfg.Mongo.Host,
		Port:     cfg.Mongo.Port,
		DBName:   cfg.Mongo.Dbname,
		Username: cfg.Mongo.User,
		SSL:      cfg.Mongo.SslMode,
		Password: cfg.Mongo.Password,
	})

	if err != nil {
		logger.Error(err)
	}

	mongoDB := mongoClient.Database(cfg.Mongo.Dbname)

	// if cfg.Environment != "prod" {
	// 	logger.Info(mongoDB)
	// }

	// initialize hasher
	hasherP := hasher.NewSHA1Hasher(cfg.Auth.Salt)

	// initialize token manager
	tokenManager, err := auths.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	// intiale opt
	otpGenerator := utils.NewGOTPGenerator()

	hub := service.NewHub()
	go hub.Run()

	repositories := repository.NewRepositories(mongoDB, cfg.I18n)
	services := service.NewServices(&service.ConfigServices{
		Repositories:           repositories,
		Hasher:                 hasherP,
		TokenManager:           tokenManager,
		OtpGenerator:           otpGenerator,
		AccessTokenTTL:         cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.RefreshTokenTTL,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		I18n:                   cfg.I18n,
		ImageService:           cfg.IImage,
		Hub:                    hub,
	})

	handlers := handler.NewHandler(services, repositories, mongoDB, &cfg.Oauth, &cfg.Auth, &cfg.I18n, &cfg.IImage, *hub)

	// initialize server
	srv := server.NewServer(cfg, handlers.InitRoutes(cfg, mongoDB))

	go func() {

		// // create a scheduler
		// s, err := gocron.NewScheduler()
		// if err != nil {
		// 	// handle error
		// 	logger.Errorf("Error gocron: %s", err.Error())
		// }

		// // add a job to the scheduler
		// j, err := s.NewJob(
		// 	gocron.DurationJob(
		// 		1*time.Second,
		// 	),
		// 	gocron.NewTask(
		// 		func(a string, b int) {
		// 			// do things
		// 			langs, err := services.Lang.FindLanguage(domain.RequestParams{Filter: bson.D{}})
		// 			if err != nil {
		// 				logger.Errorf("Error gocron services: %s", err.Error())
		// 			}
		// 			fmt.Println(langs.Total)
		// 		},
		// 		"hello",
		// 		1,
		// 	),
		// )
		// if err != nil {
		// 	// handle error
		// 	logger.Errorf("Error gocron task: %s", err.Error())
		// }
		// // each job has a unique id
		// fmt.Println(j.ID())
		// // start the scheduler
		// s.Start()

		if er := srv.Run(); !errors.Is(er, http.ErrServerClosed) {
			logger.Errorf("Error starting server: %s", er.Error())
		}
	}()

	logger.Infof("API service start on port: %s", cfg.HTTP.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("API service shutdown")
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if er := srv.Stop(ctx); er != nil {
		logger.Errorf("failed to stop server: %v", er)
	}

	if er := mongoClient.Disconnect(context.Background()); er != nil {
		logger.Error(er.Error())
	}
}
