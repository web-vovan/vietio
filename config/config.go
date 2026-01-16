package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SeedFlag bool
	Env      string
	Server   Server
	Db       DbConfig
}

type Server struct {
	HttpPort string
}

type DbConfig struct {
	Dsn string
}

func Load() *Config {
	seedFlag := flag.Bool("seed", false, "наполнение БД тестовыми данными")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("not found .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatal("HTTP_PORT is not set")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("APP_ENV is not set")
	}

	return &Config{
		SeedFlag: *seedFlag,
		Env:      env,
		Server: Server{
			HttpPort: httpPort,
		},
		Db: DbConfig{
			Dsn: dsn,
		},
	}
}
