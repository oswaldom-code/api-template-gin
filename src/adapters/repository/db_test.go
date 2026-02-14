package repository

import (
	"testing"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConnection_InvalidEngine(t *testing.T) {
	dsn := config.DBConfig{
		Engine: "invalid_engine",
	}

	store, err := NewConnection(dsn)
	assert.Nil(t, store)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid database engine")
}

func TestNewConnection_PostgresConnectionFailure(t *testing.T) {
	dsn := config.DBConfig{
		Engine:   "postgres",
		Host:     "invalid-host",
		User:     "test",
		Password: "test",
		Database: "test",
		Port:     5432,
		SSLMode:  "disable",
	}

	store, err := NewConnection(dsn)
	assert.Nil(t, store)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}
