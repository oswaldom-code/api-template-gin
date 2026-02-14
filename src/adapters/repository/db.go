package repository

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/oswaldom-code/api-template-gin/pkg/log"
	"github.com/oswaldom-code/api-template-gin/src/application/system_services/ports"
)

// repository handles the database context
type repository struct {
	db *gorm.DB
}

var (
	repositoryInstance *repository
	once               sync.Once
	initErr            error
)

// NewConnection creates a new database connection based on the provided config.
func NewConnection(dsn config.DBConfig) (ports.Store, error) {
	var dsnStrConnection string
	log.Debug("Creating new database connection", log.Fields{"dsn": dsn})

	switch dsn.Engine {
	case "postgres":
		dsnStrConnection = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=America/Lima",
			dsn.Host, dsn.User, dsn.Password, dsn.Database, dsn.Port, dsn.SSLMode)
	case "mysql":
		dsnStrConnection = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
	default:
		return nil, fmt.Errorf("invalid database engine: %s", dsn.Engine)
	}

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		FullSaveAssociations:   false,
	}
	db, err := gorm.Open(postgres.Open(dsnStrConnection), gormConfig)
	if err != nil {
		log.Error("error connecting to db", log.Fields{
			"engine":   dsn.Engine,
			"host":     dsn.Host,
			"port":     dsn.Port,
			"database": dsn.Database,
			"username": dsn.User,
			"err":      err,
		})
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &repository{db: db.Set("gorm:auto_preload", true)}, nil
}

func NewRepository() (ports.Store, error) {
	once.Do(func() {
		store, err := NewConnection(config.GetDBConfig())
		if err != nil {
			initErr = err
			return
		}
		repositoryInstance = store.(*repository)
	})
	if initErr != nil {
		return nil, initErr
	}
	return repositoryInstance, nil
}
