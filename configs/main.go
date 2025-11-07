package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Common CommonConfig

	Version VersionConfig

	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type CommonConfig struct {
	APISecretKey string `env:"API_SECRET_KEY"`
	BaseURL      string `env:"BASE_URL"`
	Name         string `env:"NAME"`
	Mode         string `env:"MODE"`
	Port         string `env:"PORT"`
	Log          bool   `env:"LOG"`
	LogTypes     string `env:"LOG_TYPES"`
}

type VersionConfig struct {
	Code int    `env:"VERSION_CODE"`
	Name string `env:"VERSION_NAME"`
	Path string `env:"VERSION_PATH"`
}

type DatabaseConfig struct {
	Host     string `env:"POSTGES_HOST"`
	Port     int    `env:"POSTGES_PORT"`
	User     string `env:"POSTGES_USER"`
	Password string `env:"POSTGES_PASSWORD"`
	DBName   string `env:"POSTGES_DB"`
}

type RedisConfig struct {
	URL     string `env:"REDIS_URL"`
	SkipTLS bool   `env:"REDIS_SKIP_TLS"`
}

type JWTConfig struct {
	Secret    string `env:"JWT_SECRET"`
	ExpiresIn int    `env:"JWT_EXPIRES_IN"`
}

var instance *Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Can not load form env file (there is no problem if we set from deployment)")
	}
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	instance = cfg
}

func GetConfig() Config {
	if instance == nil {
		cfg := &Config{}
		if err := env.Parse(cfg); err != nil {
			panic(err)
		}
		instance = cfg
	}
	return *instance
}
