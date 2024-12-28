package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	v1 "github.com/mikalai2006/kingwood-api/internal/handler/v1"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/service"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/auths"
	"github.com/mikalai2006/kingwood-api/pkg/hasher"
	"github.com/mikalai2006/kingwood-api/pkg/logger"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestSuite struct {
	suite.Suite

	db       *mongo.Database
	handler  *v1.HandlerV1
	services *service.Services
	repos    *repository.Repositories

	hasher       hasher.PasswordHasher
	tokenManager auths.TokenManager

	i18n config.I18nConfig
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	// read config file
	cfg, err := config.Init("../configs", "../.env")
	if err != nil {
		logger.Error(err)
		return
	}

	// initialize mongoDB
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	client, err := repository.NewMongoDB(&repository.ConfigMongoDB{
		Host:     cfg.Mongo.Host,
		Port:     cfg.Mongo.Port,
		DBName:   cfg.Mongo.DBTest,
		Username: cfg.Mongo.User,
		SSL:      cfg.Mongo.SslMode,
		Password: cfg.Mongo.Password,
	})

	if err != nil {
		s.FailNow("Failed to connect to mongo", err)
	}
	// defer client.Disconnect(ctx)

	// client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	// if err != nil {
	// 	s.FailNow("Failed to connect to mongo", err)
	// }

	s.db = client.Database(cfg.Mongo.DBTest)

	// Init domain deps
	repos := repository.NewRepositories(s.db, cfg.I18n)
	hasherP := hasher.NewSHA1Hasher(cfg.Auth.Salt)

	tokenManager, err := auths.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}

	// intiale opt
	otpGenerator := utils.NewGOTPGenerator()

	services := service.NewServices(&service.ConfigServices{

		Repositories:           repos,
		Hasher:                 hasherP,
		TokenManager:           tokenManager,
		AccessTokenTTL:         time.Minute * 15,
		RefreshTokenTTL:        time.Minute * 15,
		OtpGenerator:           otpGenerator,
		VerificationCodeLength: 8,
		I18n:                   cfg.I18n,
	})

	hub := service.NewHub()
	go hub.Run()

	s.repos = repos
	s.services = services
	s.handler = v1.NewHandler(services, repos, s.db, &cfg.Oauth, &cfg.Auth, &cfg.I18n, &cfg.IImage, *hub)
	s.hasher = hasherP
	s.tokenManager = tokenManager
}

func (s *TestSuite) TearDownSuite() {
	if er := s.db.Client().Disconnect(context.Background()); er != nil {
		s.FailNow("Failed disconnect DB", er)
	}
}

func (s *TestSuite) Make() {

}

func (s *TestSuite) Auth(router *gin.Engine) (domain.ResponseTokens, error) {
	testUser := domain.AuthInput{
		Email:    "mail@mail.com",
		Password: "pass12345",
	}
	// user := strings.NewReader(fmt.Sprintf("%#v", testUser))
	// data := `{"email": "mail@mail.com", "password": "pass12345"}`

	err := s.db.Collection(repository.TblAuth).Drop(context.Background())
	s.NoError(err)
	r := s.Require()

	dataJSON, err := json.Marshal(testUser)
	s.NoError(err)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewReader(dataJSON))
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var bodies map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&bodies)
	s.NoError(err)

	req, err = http.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewReader(dataJSON))
	req.Header.Set("Content-type", "application/json")
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	var auth domain.ResponseTokens
	err = json.NewDecoder(response.Body).Decode(&auth)
	s.NoError(err)

	r.Equal(http.StatusOK, response.StatusCode)

	return auth, nil
}

func (s *TestSuite) TestHomePageVersion1() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/", http.NoBody)
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	r.Equal(http.StatusOK, response.StatusCode)
}
