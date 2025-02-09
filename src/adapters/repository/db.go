package repository

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/oswaldom-code/api-template-gin/pkg/log"
	"github.com/oswaldom-code/api-template-gin/src/aplication/system_services/ports"
)

type repository struct {
	db *gorm.DB
}

var repositoryInstance *repository

func retryOnError(maxRetries int, delay time.Duration, action func() error) error {
	for i := 1; i <= maxRetries; i++ {
		err := action()
		if err == nil {
			return nil
		}
		log.Error(fmt.Sprintf("Attempt %d/%d failed: %v", i, maxRetries, err))
		if i < maxRetries {
			time.Sleep(delay * time.Duration(1<<i))
		}
	}
	return fmt.Errorf("failed after %d attempts", maxRetries)
}

// NewConnection establece una nueva conexión a la base de datos.
func NewConnection(dbConfig config.DBConfig) (ports.Repository, error) {
	var dsnStrConnection string
	log.Debug("Creating new database connection", log.Fields{
		"engine":   dbConfig.Engine,
		"host":     dbConfig.Host,
		"port":     dbConfig.Port,
		"database": dbConfig.Database,
		"username": dbConfig.User,
	})

	switch dbConfig.Engine {
	case "postgres":
		timeZone := os.Getenv("DB_TIMEZONE")
		if timeZone == "" {
			timeZone = "UTC"
		}
		dsnStrConnection = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
			dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.SSLMode, timeZone)
	case "mysql":
		dsnStrConnection = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	default:
		log.Error(fmt.Sprintf("Invalid database engine: %s", dbConfig.Engine), log.Fields{"engine": dbConfig.Engine})
		return &repository{}, fmt.Errorf("unsupported database engine: %s", dbConfig.Engine)
	}

	var dialector gorm.Dialector
	switch dbConfig.Engine {
	case "postgres":
		dialector = postgres.Open(dsnStrConnection)
	case "mysql":
		dialector = mysql.Open(dsnStrConnection)
	}

	err := retryOnError(5, time.Second, func() error {
		db, err := gorm.Open(dialector, &gorm.Config{
			SkipDefaultTransaction: true,
			FullSaveAssociations:   false,
		})
		if err != nil {
			return err
		}
		repositoryInstance = &repository{db: db}
		return nil
	})

	if err != nil {
		log.Error("Error connecting to db", log.Fields{"err": err})
		return nil, err
	}

	return repositoryInstance, nil
}

func NewRepository(config config.DBConfig) (*repository, error) {
	log.Debug("Creating new database connection")
	if repositoryInstance == nil {
		instance, err := NewConnection(config)
		if err != nil {
			log.Error("Error creating new database connection", log.Fields{"err": err})
			return nil, err
		}
		repositoryInstance = instance.(*repository)
	}
	return repositoryInstance, nil
}
