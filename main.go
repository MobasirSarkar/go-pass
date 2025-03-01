package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MobasirSarkar/pass-manage/database"
	"github.com/MobasirSarkar/pass-manage/pkg/logger"
)

func main() {
	logger.Good("starting application....")
	conn := database.DbConnection()
	if conn == nil {
		logger.Fatal("Database initialization failed. Existing application.")
		os.Exit(1)
	}
	defer database.CloseDb()

	database.SetUpDb(conn)
}

func gracefulShutdown() {
	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigChan // Block until signal received
		logger.Info("Received signal: %s. Shutting down gracefully...", sig)
		database.CloseDb()
	}()

	wg.Wait() // Wait for shutdown routines to complete
}
