package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type Conf struct {
	Server ConfServer
	DB     ConfDB
}

type ConfServer struct {
	Port         string        `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
}

type ConfDB struct {
	Host        string        `env:"DB_HOST,required"`
	Port        int           `env:"DB_PORT,required"`
	Username    string        `env:"DB_USER,required"`
	Password    string        `env:"DB_PASS,required"`
	DBName      string        `env:"DB_NAME,required"`
	Debug       bool          `env:"DB_DEBUG,required"`
	MaxIdleConn int           `env:"DB_MAX_IDLE_CONNS,required"`
	MaxOpenCOnn int           `env:"DB_MAX_OPEN_CONNS,required"`
	MaxIdleTime time.Duration `env:"DB_CONN_MAX_LIFETIME,required"`
}

func New() *Conf {
	var c Conf
	if err := env.Parse(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func NewDB() *ConfDB {
	var c ConfDB
	if err := env.Parse(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
