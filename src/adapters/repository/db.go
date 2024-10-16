package repository

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/oswaldom-code/api-template-gin/pkg/log"
	"github.com/oswaldom-code/api-template-gin/src/aplication/system_services/ports"
)

// repository handles the database context
type repository struct {
	db *gorm.DB
}

var repositoryInstance *repository

// New returns a new instance of a Store
func NewConnection(dsn config.DBConfig) ports.Store {
	var dsnStrConnection string
	log.Debug("Creating new database connection", log.Fields{"dsn": dsn})

	switch dsn.Engine {
	case "postgre":
		dsnStrConnection = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=America/Lima",
			dsn.Host, dsn.User, dsn.Password, dsn.Database, dsn.Port, dsn.SSLMode)
	case "mysql":
		dsnStrConnection = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
	default:
		log.Error("Invalid database engine", log.Fields{"engine": dsn.Engine})
	}

	// configure connection
	config := &gorm.Config{
		// SkipDefaultTransaction: (default false) - skip default transaction for each request
		// (useful for performance) 30 % faster but you need to handle transactions manually (begin, commit, rollback)
		SkipDefaultTransaction: true,
		FullSaveAssociations:   false, // default is true
	}
	db, err := gorm.Open(postgres.Open(dsnStrConnection), config)
	if err != nil {
		log.Error("error connecting to db ", log.Fields{
			"engine":   dsn.Engine,
			"host":     dsn.Host,
			"port":     dsn.Port,
			"database": dsn.Database,
			"username": dsn.User,
			"err":      err,
		})
		os.Exit(1)
	}

	return &repository{db: db.Set("gorm:auto_preload", true)}
}

func NewRepository() ports.Store {
	log.Debug("Creating new database connection", log.Fields{"dsn": config.GetDBConfig()})
	if repositoryInstance == nil {
		NewConnection(config.GetDBConfig())
		return repositoryInstance
	}
	return repositoryInstance
}
