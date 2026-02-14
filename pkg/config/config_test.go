package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDBConfig_Defaults(t *testing.T) {
	envKeys := []string{"db.user", "db.password", "db.host", "db.port", "db.database",
		"db.max_connections", "db.ssl_mode", "db.log_mode", "db.engine"}
	for _, key := range envKeys {
		os.Unsetenv(key)
	}

	cfg := GetDBConfig()

	assert.Equal(t, "default_user", cfg.User)
	assert.Equal(t, "default_pass", cfg.Password)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 0, cfg.Port)
	assert.Equal(t, "default_db", cfg.Database)
	assert.Equal(t, 10, cfg.MaxOpenConnections)
	assert.Equal(t, "disable", cfg.SSLMode)
	assert.Equal(t, "info", cfg.LogMode)
	assert.Equal(t, "postgres", cfg.Engine)
}

func TestGetDBConfig_FromEnv(t *testing.T) {
	os.Setenv("db.user", "testuser")
	os.Setenv("db.password", "testpass")
	os.Setenv("db.host", "testhost")
	os.Setenv("db.port", "5432")
	os.Setenv("db.database", "testdb")
	os.Setenv("db.engine", "postgres")
	defer func() {
		os.Unsetenv("db.user")
		os.Unsetenv("db.password")
		os.Unsetenv("db.host")
		os.Unsetenv("db.port")
		os.Unsetenv("db.database")
		os.Unsetenv("db.engine")
	}()

	cfg := GetDBConfig()

	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpass", cfg.Password)
	assert.Equal(t, "testhost", cfg.Host)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "testdb", cfg.Database)
	assert.Equal(t, "postgres", cfg.Engine)
}

func TestGetServerConfig_Defaults(t *testing.T) {
	envKeys := []string{"server.host", "server.port", "server.scheme", "server.mode"}
	for _, key := range envKeys {
		os.Unsetenv(key)
	}

	cfg := GetServerConfig()

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "9000", cfg.Port)
	assert.Equal(t, "http", cfg.Scheme)
	assert.Equal(t, "debug", cfg.Mode)
}

func TestServerConfig_Validate_Valid(t *testing.T) {
	cfg := ServerConfig{
		Host:   "localhost",
		Port:   "9000",
		Scheme: "http",
		Mode:   "debug",
	}
	assert.NoError(t, cfg.Validate())
}

func TestServerConfig_Validate_Invalid(t *testing.T) {
	cfg := ServerConfig{
		Host:   "",
		Port:   "9000",
		Scheme: "http",
		Mode:   "debug",
	}
	assert.Error(t, cfg.Validate())
}

func TestServerConfig_AsUri(t *testing.T) {
	cfg := ServerConfig{Host: "localhost", Port: "9000"}
	assert.Equal(t, "localhost:9000", cfg.AsUri())
}

func TestGetLogConfig_Defaults(t *testing.T) {
	os.Unsetenv("log.level")
	os.Unsetenv("log.errorLogFile")

	cfg := GetLogConfig()
	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "", cfg.ErrorLogFile)
}

func TestLoadEnvVariables_MapsCorrectly(t *testing.T) {
	os.Setenv("DB_USER", "envuser")
	os.Setenv("DB_DATABASE", "envdb")
	os.Setenv("DB_ENGINE", "postgres")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENVIRONMENT", "production")
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_DATABASE")
		os.Unsetenv("DB_ENGINE")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("db.user")
		os.Unsetenv("db.database")
		os.Unsetenv("db.engine")
		os.Unsetenv("log.level")
		os.Unsetenv("environment")
	}()

	loadEnvVariables()

	assert.Equal(t, "envuser", os.Getenv("db.user"))
	assert.Equal(t, "envdb", os.Getenv("db.database"))
	assert.Equal(t, "postgres", os.Getenv("db.engine"))
	assert.Equal(t, "debug", os.Getenv("log.level"))
	assert.Equal(t, "production", os.Getenv("environment"))
}

func TestGetEnvironmentConfig_Default(t *testing.T) {
	os.Unsetenv("environment")
	cfg := GetEnvironmentConfig()
	assert.Equal(t, "development", cfg.Environment)
}

func TestGetAuthenticationKey_Default(t *testing.T) {
	os.Unsetenv("auth.secret")
	cfg := GetAuthenticationKey()
	assert.Equal(t, "default_secret", cfg.Secret)
}
