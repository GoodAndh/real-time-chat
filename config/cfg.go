package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type config struct {
	PublicHost   string
	Port         string
	DBUser       string
	DBPassword   string
	DBAddress    string
	DBName       string
	DBTest       string
	JWTSecretKey string
}

var Env config

func InitConfig(path string) config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envPath := ""
	if path == "main" {
		envPath = filepath.Join(dir, ".", ".env")
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(getEnv("DB_HOST", "gada1"), getEnv("DB_PORT", "gada2"), getEnv("PORT","gada3"))

	return config{
		PublicHost:   getEnv("PUBLIC_HOST", "http://localhost"),
		Port:         getEnv("PORT", "3000"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", "r23p"),
		DBAddress:    fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:       getEnv("DB_NAME", "real_time"),
		DBTest:       getEnv("DB_NAME_TEST", "real_time_test"),
		JWTSecretKey: getEnv("SECRET_KEY", "SECRET_KEY"),
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
