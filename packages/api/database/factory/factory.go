package factory

import (
	"fmt"
	"hue-n-gnk/csv-to-gitisue/database/dao"
	"hue-n-gnk/csv-to-gitisue/pkg/logger"

	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// dbManager is a default database manager.
	dbManager *gorm.DB

	// dbMutex prevents dup of default database manager.
	dbMutex sync.Mutex
)

type Config struct {
	User               string
	Password           string
	Host               string
	Database           string
	Port               string
	ConnMaxLifeTime    time.Duration
	ConnMaxIdle        int
	ConnMaxOpen        int
	TransactionTimeout time.Duration
}

// NewConnection is the way to connect to databse
func NewConnection(cfg Config) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if dbManager != nil {
		return errors.New("database: already initialized")
	}
	db, err := connectDatabase(cfg)
	if err != nil {
		return errors.Wrapf(err, "connection failed")
	}
	dbManager = db

	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "error getting generic database interface")
	}
	sqlDB.SetMaxIdleConns(cfg.ConnMaxIdle)
	sqlDB.SetMaxOpenConns(cfg.ConnMaxOpen)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	return nil
}

func connectDatabase(cfg Config) (*gorm.DB, error) {
	var dsn string
	lg, err := logger.NewLogger()
	if err != nil {
		return nil, errors.Wrap(err, "setup Logger failed")
	}
	gormConfig := &gorm.Config{
		Logger: logger.NewGormLogger(lg),
	}
	dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Database, cfg.Password)
	return gorm.Open(postgres.Open(dsn), gormConfig)
}

func Query() (*dao.Query, error) {
	if dbManager == nil {
		return nil, errors.New("Don't have any connection to Database")
	}
	return dao.Use(dbManager), nil
}
