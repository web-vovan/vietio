package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	BotToken    string
	Server      Server
	S3Storage   S3Storage
	StorageType string
	Db          DbConfig
	JwtSecret   string
}

type Server struct {
	HttpPort  string
	PublicUrl string
}

type S3Storage struct {
	Key       string
	Secret    string
	Bucket    string
	PublicUrl string
}

type DbConfig struct {
	Dsn string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("not found .env file")
	}

	// db config
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatal("HTTP_PORT is not set")
	}

	publicUrl := os.Getenv("PUBLIC_URL")
	if publicUrl == "" {
		log.Fatal("PUBLIC_URL is not set")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("APP_ENV is not set")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_BUCKET is not set")
	}

	s3Key := os.Getenv("S3_KEY")
	if s3Key == "" {
		log.Fatal("S3_KEY is not set")
	}

	s3Secret := os.Getenv("S3_SECRET")
	if s3Secret == "" {
		log.Fatal("S3_SECRET is not set")
	}

	s3PublicUrl := os.Getenv("S3_PUBLIC_URL")
	if s3PublicUrl == "" {
		log.Fatal("S3_PUBLIC_URL is not set")
	}

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		log.Fatal("STORAGE_TYPE is not set")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	return &Config{
		Env: env,
		Server: Server{
			HttpPort:  httpPort,
			PublicUrl: publicUrl,
		},
		S3Storage: S3Storage{
			Key:       s3Key,
			Secret:    s3Secret,
			Bucket:    s3Bucket,
			PublicUrl: s3PublicUrl,
		},
		StorageType: storageType,
		BotToken:    botToken,
		Db: DbConfig{
			Dsn: dsn,
		},
		JwtSecret: jwtSecret,
	}
}
