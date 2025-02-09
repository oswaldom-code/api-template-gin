package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

var (
	environment string
)

type EnvironmentConfig struct {
	Environment string
}

type DBConfig struct {
	User               string
	Password           string
	Host               string
	Port               int
	Database           string
	MaxOpenConnections int
	SSLMode            string
	LogMode            string
	Engine             string
}

type ServerConfig struct {
	Host              string
	Port              string
	Scheme            string
	Mode              string
	PathToSSLKeyFile  string
	PathToSSLCertFile string
	Static            string
}

type LoggingConfig struct {
	Level        string
	ErrorLogFile string
}

type AuthenticateKeyConfig struct {
	Secret string
}

func SetEnvironment(env string) {
	environment = env
}

func (s *ServerConfig) Validate() error {
	requiredFields := map[string]string{
		"HOST":   s.Host,
		"PORT":   s.Port,
		"SCHEME": s.Scheme,
		"MODE":   s.Mode,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("ServerConfig is invalid: %s is empty", field)
		}
	}
	return nil
}

func (s ServerConfig) AsUri() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func LoadConfigFromFlagsAndEnv() error {
	pflag.String("DB_HOST", "localhost", "Database host")
	pflag.Int("DB_PORT", 5432, "Database port")
	pflag.String("DB_USER", "default_user", "Database user")
	pflag.String("DB_PASSWORD", "default_pass", "Database password")
	pflag.String("SERVER_HOST", "0.0.0.0", "Server host")
	pflag.String("SERVER_PORT", "9000", "Server port")
	pflag.String("SERVER_SCHEME", "http", "Server scheme")
	pflag.String("SERVER_MODE", "release", "Server mode")
	pflag.String("AUTH_SECRET", "default_secret", "Authentication secret")

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, proceeding with default values")
	} else {
		fmt.Println(".env file loaded successfully")
	}

	if err := applyEnvVarsToFlags(); err != nil {
		return fmt.Errorf("failed to apply environment variables to flags: %v", err)
	}
	pflag.Parse()

	return nil
}

func applyEnvVarsToFlags() error {
	// Mapa de las variables de entorno que queremos cargar y los flags correspondientes
	envVars := map[string]string{
		"DB_HOST":       "DB_HOST",
		"DB_PORT":       "DB_PORT",
		"DB_USER":       "DB_USER",
		"DB_PASSWORD":   "DB_PASSWORD",
		"SERVER_HOST":   "SERVER_HOST",
		"SERVER_PORT":   "SERVER_PORT",
		"SERVER_SCHEME": "SERVER_SCHEME",
		"SERVER_MODE":   "SERVER_MODE",
		"AUTH_SECRET":   "AUTH_SECRET",
	}

	// Recorremos las variables de entorno y asignamos los valores a los flags correspondientes
	for envVar, flag := range envVars {
		if value, exists := os.LookupEnv(envVar); exists {
			if err := pflag.Set(flag, value); err != nil {
				return fmt.Errorf("failed to set flag %s with value %s: %v", flag, value, err)
			}
		}
	}
	return nil
}

func GetDBConfig() DBConfig {
	return DBConfig{
		User:               getEnv("DB_USER", "default_user"),
		Password:           getEnv("DB_PASSWORD", "default_pass"),
		Host:               getEnv("DB_HOST", "localhost"),
		Port:               getIntEnv("DB_PORT", 5432),
		Database:           getEnv("DB_DATABASE", "default_db"),
		MaxOpenConnections: getIntEnv("DB_MAX_CONNECTIONS", 10),
		SSLMode:            getEnv("DB_SSL_MODE", "disable"),
		LogMode:            getEnv("DB_LOG_MODE", "info"),
		Engine:             getEnv("DB_ENGINE", "postgres"),
	}
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		Host:              getEnv("SERVER_HOST", "0.0.0.0"),
		Port:              getEnv("SERVER_PORT", "9000"),
		Scheme:            getEnv("SERVER_SCHEME", "http"),
		Mode:              getEnv("SERVER_MODE", "release"),
		PathToSSLKeyFile:  os.Getenv("SERVER_SSL_KEY"),
		PathToSSLCertFile: os.Getenv("SERVER_SSL_CERT"),
		Static:            os.Getenv("SERVER_STATIC"),
	}
}

func GetLogConfig() LoggingConfig {
	return LoggingConfig{
		Level:        getEnv("LOG_LEVEL", "info"),
		ErrorLogFile: os.Getenv("LOG_ERROR_LOG_FILE"),
	}
}

func GetEnvironmentConfig() EnvironmentConfig {
	return EnvironmentConfig{
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func GetAuthenticationKey() AuthenticateKeyConfig {
	return AuthenticateKeyConfig{
		Secret: getEnv("AUTH_SECRET", "default_secret"),
	}
}

// Función genérica para obtener variables de entorno
func getEnv(key, defaultValue string) string {
	// Primero verifica si el flag está definido
	if value := pflag.Lookup(key); value != nil && value.Value.String() != "" {
		return value.Value.String()
	}
	// Si el flag no está definido, busca en las variables de entorno
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	// Si no hay flag ni variable de entorno, devuelve el valor por defecto
	return defaultValue
}

// Función genérica para obtener variables enteras
func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if value := pflag.Lookup(key); value != nil {
		if intValue, err := strconv.Atoi(value.Value.String()); err == nil {
			return intValue
		}
	}
	return defaultValue
}
