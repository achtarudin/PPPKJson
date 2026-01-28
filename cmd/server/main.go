package main

import (
	"context"
	"cutbray/pppk-json/cmd/config"
	"cutbray/pppk-json/internal/adapters/db_adapter"
	"cutbray/pppk-json/internal/adapters/logger"
	"cutbray/pppk-json/internal/utils"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"gorm.io/gorm"
)

func main() {
	logger.New()

	dbHost := utils.GetEnvOrDefault("DB_HOST", "localhost")
	dbUser := utils.GetEnvOrDefault("DB_USER", "encang_cutbray")
	dbPassword := utils.GetEnvOrDefault("DB_PASSWORD", "encang_cutbray")
	dbName := utils.GetEnvOrDefault("DB_NAME", "togotestgo")
	dbPort := utils.GetEnvOrDefault("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize Adapters
	dbAdapter := db_adapter.New(dsn, &gorm.Config{})

	connectManagers := []config.ConnectManager{
		{Name: "Postgres DB", Adapter: dbAdapter},
	}

	if err := config.ConnectAdapters(shutdown, connectManagers...); err != nil {
		log.Fatalf("%v", err)
	}

	// Safety Check Initialize Connected Adapters
	_, ok := dbAdapter.Value().(*gorm.DB)
	if !ok {
		log.Fatalf("Database adapter is not properly initialized")
	}

	<-shutdown.Done()

	config.DisconnectAdapters(connectManagers...)

}
