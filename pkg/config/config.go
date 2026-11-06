package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Name            string
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // in minutes
}

type Config struct {
	AppName            string
	AppEnv             string
	AppPort            string
	GinMode            string
	APIPrefix          string
	JWTSecret          string
	JWTExpirationHours int
	CORSAllowedOrigins []string
	RateLimitRPS       float64
	RateLimitBurst     int
	StoragePath        string
	StorageURL         string
	Databases          map[string]DatabaseConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("[Config] .env file not found, using environment variables")
	}

	cfg := &Config{
		AppName:            getEnv("APP_NAME", "BarongFramework"),
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		GinMode:            getEnv("GIN_MODE", "debug"),
		APIPrefix:          getEnv("API_PREFIX", "/api/v1"),
		JWTSecret:          getEnv("JWT_SECRET", "default-insecure-secret-key-32-chars"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		CORSAllowedOrigins: parseCommaSeparated(getEnv("CORS_ALLOWED_ORIGINS", "*")),
		RateLimitRPS:       getEnvAsFloat("RATE_LIMIT_RPS", 20.0),
		RateLimitBurst:     getEnvAsInt("RATE_LIMIT_BURST", 40),
		StoragePath:        getEnv("STORAGE_PATH", "./storage/uploads"),
		StorageURL:         getEnv("STORAGE_URL", "/storage"),
		Databases:          make(map[string]DatabaseConfig),
	}

	// 1. Primary/Default Database
	cfg.Databases["default"] = DatabaseConfig{
		Name:            "default",
		Host:            getEnv("DB_HOST", "127.0.0.1"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "barong_main"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 5),
	}

	// 2. Secondary DB (Optional User DB)
	if userDBName := os.Getenv("DB_USER_NAME"); userDBName != "" {
		cfg.Databases["user"] = DatabaseConfig{
			Name:            "user",
			Host:            getEnv("DB_USER_HOST", cfg.Databases["default"].Host),
			Port:            getEnv("DB_USER_PORT", cfg.Databases["default"].Port),
			User:            getEnv("DB_USER_USER", cfg.Databases["default"].User),
			Password:        getEnv("DB_USER_PASSWORD", cfg.Databases["default"].Password),
			DBName:          userDBName,
			SSLMode:         getEnv("DB_USER_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_USER_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_USER_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvAsInt("DB_USER_CONN_MAX_LIFETIME_MIN", 5),
		}
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseCommaSeparated(val string) []string {
	if strings.TrimSpace(val) == "" {
		return []string{}
	}
	parts := strings.Split(val, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
