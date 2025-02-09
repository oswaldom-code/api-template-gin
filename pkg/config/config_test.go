package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestGetDBConfig(t *testing.T) {
	// set environment variables
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_DATABASE", "testdb")
	os.Setenv("DB_MAX_CONNECTIONS", "20")
	os.Setenv("DB_SSL_MODE", "enable")
	os.Setenv("DB_LOG_MODE", "debug")
	os.Setenv("DB_ENGINE", "postgres")

	// clean up the environment variables after the test
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_DATABASE")
		os.Unsetenv("DB_MAX_CONNECTIONS")
		os.Unsetenv("DB_SSL_MODE")
		os.Unsetenv("DB_LOG_MODE")
		os.Unsetenv("DB_ENGINE")
	}()
	// reset pflag to avoid contamination between tests
	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	// load the configuration
	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "LoadConfigFromFlagsAndEnv should not return an error")

	// Get the configuration values
	dbConfig := GetDBConfig()

	assert.Equal(t, "testuser", dbConfig.User)
	assert.Equal(t, "testpass", dbConfig.Password)
	assert.Equal(t, "localhost", dbConfig.Host)
	assert.Equal(t, 5432, dbConfig.Port)
	assert.Equal(t, "testdb", dbConfig.Database)
	assert.Equal(t, 20, dbConfig.MaxOpenConnections)
	assert.Equal(t, "enable", dbConfig.SSLMode)
	assert.Equal(t, "debug", dbConfig.LogMode)
	assert.Equal(t, "postgres", dbConfig.Engine)
}

func TestGetServerConfig(t *testing.T) {
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_SCHEME", "https")
	os.Setenv("SERVER_MODE", "development")

	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_SCHEME")
		os.Unsetenv("SERVER_MODE")
	}()

	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "LoadConfigFromFlagsAndEnv should not return an error")

	serverConfig := GetServerConfig()

	assert.Equal(t, "127.0.0.1", serverConfig.Host)
	assert.Equal(t, "8080", serverConfig.Port)
	assert.Equal(t, "https", serverConfig.Scheme)
	assert.Equal(t, "development", serverConfig.Mode)
}

func TestLoadConfigFromFlagsAndEnv(t *testing.T) {
	os.Setenv("DB_USER", "envuser")
	os.Setenv("SERVER_HOST", "envhost")

	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("SERVER_HOST")
	}()

	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "LoadConfigFromFlagsAndEnv should not return an error")

	dbUser := GetDBConfig().User
	serverHost := GetServerConfig().Host

	assert.Equal(t, "envuser", dbUser)
	assert.Equal(t, "envhost", serverHost)
}

func TestGetLogConfig(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")

	defer func() {
		os.Unsetenv("LOG_LEVEL")
	}()

	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "LoadConfigFromFlagsAndEnv should not return an error")

	logConfig := GetLogConfig()

	assert.Equal(t, "debug", logConfig.Level)
}

func TestGetAuthenticationKey(t *testing.T) {
	os.Setenv("AUTH_SECRET", "mysecretkey")

	defer func() {
		os.Unsetenv("AUTH_SECRET")
	}()

	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "LoadConfigFromFlagsAndEnv should not return an error")

	authKey := GetAuthenticationKey()

	assert.Equal(t, "mysecretkey", authKey.Secret)
}

func TestGetEnvironmentConfig(t *testing.T) {
	os.Setenv("ENVIRONMENT", "production")

	defer func() {
		os.Unsetenv("ENVIRONMENT")
	}()

	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "LoadConfigFromFlagsAndEnv should not return an error")

	envConfig := GetEnvironmentConfig()

	assert.Equal(t, "production", envConfig.Environment)
}

func TestGetEnvFallback(t *testing.T) {
	os.Setenv("DB_USER", "envuser")
	defer os.Unsetenv("DB_USER")

	// Simulates a command line argument
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--DB_USER=flaguser"} // Simulates "--DB_USER=flaguser"

	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "Error al cargar la configuración")

	dbConfig := GetDBConfig()

	assert.Equal(t, "flaguser", dbConfig.User, "El flag DB_USER debería sobrescribir la variable de entorno")
}

func TestGetIntEnv(t *testing.T) {
	// Configuramos las variables de entorno
	os.Setenv("DB_PORT", "3306")

	// Limpiamos las variables de entorno después de la prueba
	defer func() {
		os.Unsetenv("DB_PORT")
	}()

	// Reiniciamos pflag para evitar contaminación entre pruebas
	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Cargamos la configuración
	err := LoadConfigFromFlagsAndEnv()
	assert.NoError(t, err, "LoadConfigFromFlagsAndEnv should not return an error")

	// Verificamos los valores
	port := getIntEnv("DB_PORT", 5432)

	assert.Equal(t, 3306, port)
}
