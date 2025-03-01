package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/MobasirSarkar/pass-manage/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
)

var (
	instance *sql.DB
	once     sync.Once
)

const (
	MAX_IDLE_TIME  = 10 * time.Minute
	MAX_OPEN_CONNS = 1
	MAX_IDLE_CONNS = 1
	DB_TIMEOUT     = 5 * time.Second
)

// DbConnection initializes the singleton SQLITE database connection...
func DbConnection() *sql.DB {

	once.Do(func() {
		dbPath := os.Getenv("SQLITE_DB_PATH")
		if dbPath == "" {
			dbPath = "data/main.db" // default path
		}

		// ensure directory exists
		if err := ensureDir(filepath.Dir(dbPath)); err != nil {
			logger.Fatal("Failed to create database directory: %v", err)
			return
		}

		var err error
		instance, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			logger.Error("Failed to Open SQLITE database: %v", err)
			return
		}

		// configure connection pooling
		instance.SetMaxOpenConns(MAX_OPEN_CONNS)
		instance.SetMaxIdleConns(MAX_IDLE_CONNS)
		instance.SetConnMaxIdleTime(MAX_IDLE_TIME)

		ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
		defer cancel()
		if err := instance.PingContext(ctx); err != nil {
			logger.Fatal("Database connection test failed\n details:%s", err)
			instance.Close()
			return
		}
		logger.Good("Connected to SQLITE database at %s", dbPath)
	})

	if instance == nil {
		log.Fatal("Failed to initializes database. check logs for details.")
	}

	return instance
}

func CloseDb() {
	if instance != nil {
		if err := instance.Close(); err != nil {
			logger.Error("error closing database: %v", err)
		} else {
			logger.Info("Database connection closed.")
		}
	}
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Info("Creating database directory: %s", dir)
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
