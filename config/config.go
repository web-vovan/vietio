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
		log.Println(".env file not found, using environment variables")
	}

	// db config
	dbHost := getEnvVar("DB_HOST")
	dbPort := getEnvVar("DB_PORT")
	dbName := getEnvVar("DB_NAME")
	dbUser := getEnvVar("DB_USER")
	dbPassword := getEnvVar("DB_PASSWORD")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	httpPort := getEnvVar("HTTP_PORT")
	publicUrl := getEnvVar("PUBLIC_URL")
	env := getEnvVar("APP_ENV")
	s3Bucket := getEnvVar("S3_BUCKET")
	s3Key := getEnvVar("S3_KEY")
	s3Secret := getEnvVar("S3_SECRET")
	s3PublicUrl := getEnvVar("S3_PUBLIC_URL")
	storageType := getEnvVar("STORAGE_TYPE")
	botToken := getEnvVar("BOT_TOKEN")
	jwtSecret := getEnvVar("JWT_SECRET")

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

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set", key)
	}
	return value
}
