package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/oswaldom-code/api-template-gin/pkg/log"
	"github.com/spf13/pflag"
)

var environment string

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

func SetEnvironment(env string) {
	environment = env
}

func (s *ServerConfig) Validate() error {
	if s.Host == "" || s.Port == "0" || s.Scheme == "" || s.Mode == "" {
		return fmt.Errorf(`ServerConfig is invalid: 
        env: %s
        host: %s
        port: %s
        scheme: %s
        mode: %s`,
			environment, s.Host, s.Port, s.Scheme, s.Mode)
	}
	return nil
}

func (s ServerConfig) AsUri() string {
	return s.Host + ":" + s.Port
}

type LoggingConfig struct {
	Level        string
	ErrorLogFile string
}

type AuthenticateKeyConfig struct {
	Secret string
}

func LoadConfiguration() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found, proceeding with default values")
	}

	// Definir banderas de configuración
	pflag.String("db.user", "", "Database user")
	pflag.String("db.password", "", "Database password")
	pflag.String("db.host", "", "Database host")
	pflag.Int("db.port", 0, "Database port")
	pflag.String("server.host", "localhost", "Server host")
	pflag.String("server.port", "9000", "Server port")
	pflag.String("server.scheme", "http", "Server scheme")
	pflag.String("server.mode", "debug", "Server mode")
	pflag.String("auth.secret", "", "Authentication secret")

	pflag.Parse()

	loadEnvVariables()

	log.SetLogLevel(GetLogConfig().Level)
}

func loadEnvVariables() {
	envVars := []struct {
		key   string
		value string
	}{
		{"DB_USER", "db.user"},
		{"DB_PASSWORD", "db.password"},
		{"DB_HOST", "db.host"},
		{"DB_PORT", "db.port"},
		{"SERVER_HOST", "server.host"},
		{"SERVER_PORT", "server.port"},
		{"SERVER_SCHEME", "server.scheme"},
		{"SERVER_MODE", "server.mode"},
		{"AUTH_SECRET", "auth.secret"},
	}

	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar.key); exists {
			if err := os.Setenv(envVar.value, value); err != nil {
				log.Warn(fmt.Sprintf("Failed to set environment variable %s: %v", envVar.value, err))
			}
		}
	}
}

func GetDBConfig() DBConfig {
	port, _ := strconv.Atoi(os.Getenv("db.port"))
	config := DBConfig{
		User:               getEnv("db.user", "default_user"),     // Valor predeterminado
		Password:           getEnv("db.password", "default_pass"), // Valor predeterminado
		Host:               getEnv("db.host", "localhost"),        // Valor predeterminado
		Port:               port,
		Database:           getEnv("db.database", "default_db"), // Valor predeterminado
		MaxOpenConnections: getIntEnv("db.max_connections", 10), // Valor predeterminado
		SSLMode:            getEnv("db.ssl_mode", "disable"),    // Valor predeterminado
		LogMode:            getEnv("db.log_mode", "info"),       // Valor predeterminado
		Engine:             getEnv("db.engine", "postgres"),     // Valor predeterminado
	}
	log.Debug("DBConfig", log.Fields{"config": config})
	return config
}

func GetServerConfig() ServerConfig {
	config := ServerConfig{
		Host:              getEnv("server.host", "localhost"), // Valor predeterminado
		Port:              getEnv("server.port", "9000"),      // Valor predeterminado
		Scheme:            getEnv("server.scheme", "http"),    // Valor predeterminado
		Mode:              getEnv("server.mode", "debug"),     // Valor predeterminado
		PathToSSLKeyFile:  os.Getenv("server.ssl.key"),        // Podría ser opcional
		PathToSSLCertFile: os.Getenv("server.ssl.cert"),       // Podría ser opcional
		Static:            os.Getenv("server.static"),         // Podría ser opcional
	}
	log.Debug("ServerConfig", log.Fields{"config": config})
	return config
}

func GetLogConfig() LoggingConfig {
	return LoggingConfig{
		Level:        getEnv("log.level", "info"),
		ErrorLogFile: os.Getenv("log.errorLogFile"),
	}
}

func GetEnvironmentConfig() EnvironmentConfig {
	return EnvironmentConfig{
		Environment: getEnv("environment", "development"),
	}
}

func GetAuthenticationKey() AuthenticateKeyConfig {
	return AuthenticateKeyConfig{
		Secret: getEnv("auth.secret", "default_secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
