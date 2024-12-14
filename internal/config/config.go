package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		Oauth       OauthConfig
		I18n        I18nConfig
		IImage      IImageConfig
		// FileStorage FileStorageConfig
		// Email       EmailConfig
		// Payment     PaymentConfig
		// Limiter     LimiterConfig
		// CacheTTL    time.Duration `mapstructure:"ttl"`
		// SMTP        SMTPConfig
		// Cloudflare  CloudflareConfig
	}

	MongoConfig struct {
		Host     string
		User     string
		Password string
		Port     string
		Dbname   string `mapstructure:"dbname"`
		SslMode  bool   `mapstructure:"sslmode"`
		DBTest   string
	}

	OauthConfig struct {
		// TimeExpireCookie int

		VkAuthURI      string
		VkTokenURI     string
		VkUserinfoURI  string
		VkClientID     string
		VkClientSecret string
		VkRedirectURI  string
		VkScopes       []string

		GoogleAuthURI      string
		GoogleTokenURI     string
		GoogleUserinfoURI  string
		GoogleRedirectURI  string
		GoogleClientID     string
		GoogleClientSecret string
		GoogleScopes       []string
	}

	AuthConfig struct {
		Salt              string
		SigningKey        string
		NameCookieRefresh string
		AccessTokenTTL    time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL   time.Duration `mapstructure:"refreshTokenTTL"`

		VerificationCodeLength int
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	I18nConfig struct {
		Prefix  string
		Default string
		Locale  string
	}

	IIModeImage struct {
		Size    int
		Prefix  string
		Quality int
	}
	IImageConfig struct {
		Sizes []IIModeImage
	}
)

func Init(configsDir, envPath string) (*Config, error) {
	var cfg Config
	setDefaultConfigs(&cfg)

	// read env configs
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("mongodb", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("i18n", &cfg.I18n); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("images", &cfg.IImage); err != nil {
		return err
	}

	return viper.UnmarshalKey("oauth", &cfg.Oauth)
}

func setFromEnv(cfg *Config) {
	cfg.HTTP.Host = os.Getenv("HOST")
	cfg.HTTP.Port = os.Getenv("PORT")
	// cfg.HTTP.ReadTimeout = 10 * time.Second
	// cfg.HTTP.WriteTimeout = 10 * time.Second

	cfg.Mongo.Host = os.Getenv("MONGODB_HOST")
	cfg.Mongo.Port = os.Getenv("MONGODB_PORT")
	cfg.Mongo.User = os.Getenv("MONGODB_USER")
	cfg.Mongo.Password = os.Getenv("MONGODB_PASSWORD")

	cfg.Auth.Salt = os.Getenv("SALT")
	cfg.Auth.SigningKey = os.Getenv("SIGNING_KEY")
	cfg.Auth.NameCookieRefresh = os.Getenv("NAME_COOKIE_REFRESH")

	cfg.Environment = os.Getenv("APP_ENV")

	cfg.Oauth.VkClientID = os.Getenv("VK_CLIENT_ID")
	cfg.Oauth.VkClientSecret = os.Getenv("VK_CLIENT_SECRET")

	cfg.Oauth.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.Oauth.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
}

func setDefaultConfigs(cfg *Config) {
	cfg.Oauth.VkAuthURI = "https://oauth.vk.com/authorize"
	cfg.Oauth.VkRedirectURI = "https://poihub.ru/api/v1/oauth/vk/me"
	cfg.Oauth.VkScopes = []string{"account"}
	cfg.Oauth.VkTokenURI = "https://oauth.vk.com/access_token"
	cfg.Oauth.VkUserinfoURI = "https://api.vk.com/method/users.get"
	cfg.Oauth.GoogleAuthURI = "https://accounts.google.com/o/oauth2/auth"
	cfg.Oauth.GoogleRedirectURI = "http://localhost:8000/api/v1/oauth/google/me"
	cfg.Oauth.GoogleTokenURI = "https://accounts.google.com/o/oauth2/token"
	cfg.Oauth.GoogleUserinfoURI = "https://www.googleapis.com/oauth2/v3/userinfo"
	cfg.Oauth.GoogleScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	cfg.Auth.VerificationCodeLength = 10
}
