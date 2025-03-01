package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/MobasirSarkar/pass-manage/pkg/logger"
)

var CONTEXT_TIMEOUT = 5 * time.Second

func SetUpDb(conn *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancel()

	queries := []string{`
   CREATE TABLE IF NOT EXISTS vaults (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      site_name TEXT NOT NULL,
      user_name TEXT NOT NULL,
      password BLOB NOT NULL,
      iv BLOB NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );`,

		`CREATE TABLE IF NOT EXISTS settings (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      master_key BLOB NOT NULL,
      salt BLOB NOT NULL,
      iterations INTEGER DEFAULT 1000000,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );`,

		`CREATE TABLE IF NOT EXISTS audit_logs (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      action TEXT NOT NULL,
      site_name TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );`,
	}

	for _, query := range queries {
		_, err := conn.ExecContext(ctx, query)
		if err != nil {
			logger.Error("error setting up the schema:\n details: %s", err)
			return
		}
	}
	logger.Good("Schema is setup successfully...")
}
