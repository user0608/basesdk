package connection

import (
	"basesdk/configs"
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConfigParams struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUsername string
	DBPassword string
	DBLogLevel string
}

type StorageManager interface {
	Conn(ctx context.Context) *gorm.DB
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type storageManager struct {
	db *gorm.DB
}

var _ StorageManager = (*storageManager)(nil)

func NewConnection(config configs.DatabaseConfig) (StorageManager, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.GetHost(),
		config.GetUsername(),
		config.GetPassword(),
		config.GetDBName(),
		config.GetPort(),
	)

	db, err := openConnection(dsn, config.GetLogLevel())
	if err != nil {
		return nil, err
	}

	return &storageManager{db: db}, nil
}

func openConnection(dsn string, logLevel string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(parseLogLevel(logLevel)),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully.")

	return db, nil
}

func parseLogLevel(value string) logger.LogLevel {
	switch value {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	default:
		return logger.Silent
	}
}

type txContextKey struct{}

var transactionKey txContextKey

func (s *storageManager) Conn(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}

	return s.db.WithContext(ctx)
}

func (s *storageManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if fn == nil {
		return nil
	}

	if _, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return fn(ctx)
	}

	return s.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, transactionKey, tx)
		return fn(txCtx)
	})
}

// skipStorage is a StorageManager that disables persistence using
// the null object pattern.
type skipStorage struct{}

var _ StorageManager = (*skipStorage)(nil)

func SkipStorage() (StorageManager, error) {
	return &skipStorage{}, nil
}

func (*skipStorage) Conn(ctx context.Context) *gorm.DB {
	return nil
}

func (*skipStorage) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if fn == nil {
		return nil
	}

	return fn(ctx)
}
