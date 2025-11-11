package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/r0manch1k/umbrella/signature-server/pkg/urlbuilder"
)

type (
	Config struct {
		App       App
		HTTP      HTTP
		Log       Log
		Signature Signature
		DB        DB
	}

	App struct {
		Name string `env:"APP_NAME,required"`
		TZ   string `env:"TZ,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	Signature struct {
		PrivateKeyPath string `env:"SIGNATURE_PRIVATE_KEY_PATH,required"`
		PublicKeyPath  string `env:"SIGNATURE_PUBLIC_KEY_PATH,required"`
		Product        string `env:"SIGNATURE_PRODUCT,required"`
		RSAKeyBits     int    `env:"SIGNATURE_RSA_KEY_BITS" envDefault:"2048"`
	}

	DB struct {
		Host        string `env:"DB_HOST,required"`
		Port        string `env:"DB_PORT" envDefault:"5432"`
		User        string `env:"DB_USERNAME" envDefault:"postgres"`
		Password    string `env:"DB_PASSWORD" envDefault:"postgres"`
		Database    string `env:"DB_DATABASE" envDefault:"postgres"`
		MaxPoolSize int32  `env:"DB_MAX_POOL_SIZE" envDefault:"10"`
	}

	HTTP struct {
		Host string `env:"HTTP_HOST,required"`
		Port string `env:"HTTP_PORT,required"`
	}
)

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("file .env not found: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (db *DB) URL() string {
	return urlbuilder.New().
		Scheme("postgres").
		User(db.User, db.Password).
		Host(db.Host).
		Port(db.Port).
		Path("/" + db.Database).
		Build()
}
