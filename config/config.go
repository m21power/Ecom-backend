package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost    string
	Port          string
	DBUser        string
	DBPassword    string
	DBAddress     string
	DBName        string
	JWTExpiration int64
	JWTSecret     string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:    getEnv("PUBLIC_HOST", "http://localhost"),
		Port:          getEnv("PORT", "8080"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "abate"),
		DBName:        getEnv("DB_NAME", "ecom"),
		DBAddress:     getEnv("DB_ADDRESS", "localhost:3306"),
		JWTExpiration: getEnvInt64("JWT_EXP", 3600*24*7),
		JWTSecret:     getEnv("JWT_SECRET", "not_secret_secret_anymoer_secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback // if u didn't find the key, return the fallback, the defaultd
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return val
	}
	return fallback
}
