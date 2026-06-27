package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
		Host string
		ReadTimeout time.Duration
		WriteTimeout time.Duration
	}

	Database struct {
		Host string
		Port string
		User string
		Password string
		DBName string
		SSLMode string
	}

	JWT struct {
		Secret string
		TokenExpiry time.Duration
		RefreshExpiry time.Duration
	}

	Environment string
}

var Conf *Config

func Load() error {
	if err := godotenv.Load("../.env"); err != nil {
		return err
	}

	cfg := &Config{}

	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.Server.Host = os.Getenv("SERVER_HOST")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15

	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.DBName = os.Getenv("DB_NAME")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.SSLMode = os.Getenv("DB_SSLMODE")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.JWT.TokenExpiry = time.Hour * 24
	cfg.JWT.RefreshExpiry = time.Hour * 168

	cfg.Environment = os.Getenv("ENV")

	Conf = cfg

	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}