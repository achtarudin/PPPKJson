// Package main provides the entry point for the PPPKJson Exam API application.
//
// @title           PPPKJson Exam API
// @version         1.0.0
// @description     PPPK Examination Management System API with randomized questions
// @termsOfService  http://swagger.io/terms/
//
// @contact.name    API Support
// @contact.url     http://www.example.com/support
// @contact.email   support@example.com
//
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
//
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http https
//
// @securityDefinitions.basic BasicAuth
package main

import (
	"context"
	"cutbray/pppk-json/cmd/config"
	_ "cutbray/pppk-json/docs" // for swagger docs
	"cutbray/pppk-json/internal/adapters/db_adapter"
	"cutbray/pppk-json/internal/adapters/gin_adapter"
	"cutbray/pppk-json/internal/adapters/logger"
	"cutbray/pppk-json/internal/handlers"
	"cutbray/pppk-json/internal/utils"
	"fmt"
	"io/fs"
	"log"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	logger.New()

	err := config.LoadEnvFile()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := utils.GetEnvOrDefault("APP_PORT", "8080")
	ginMode := utils.GetEnvOrDefault("APP_MODE", gin.ReleaseMode)

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
	ginAdapter := gin_adapter.New(gin_adapter.GinConfig{
		Port: port,
		Mode: ginMode,
	})

	connectManagers := []config.ConnectManager{
		{Name: "Postgres DB", Adapter: dbAdapter},
		{Name: "Gin HTTP Server", Adapter: ginAdapter},
	}

	if err := config.ConnectAdapters(shutdown, connectManagers...); err != nil {
		log.Fatalf("%v", err)
	}

	// Safety Check Initialize Connected Adapters
	db, ok := dbAdapter.Value().(*gorm.DB)
	if !ok {
		log.Fatalf("Database adapter is not properly initialized")
	}

	// Setup handlers and routes
	ginEngine, ok := ginAdapter.Value().(*gin.Engine)
	if !ok {
		log.Fatalf("Gin adapter is not properly initialized")
	}

	// Register exam handlers
	handlers.NewGinExamHandler(db).RegisterRoutes(ginEngine)
	handlers.NewFrontendHandler().RegisterRoutes(ginEngine)
	<-shutdown.Done()

	config.DisconnectAdapters(connectManagers...)

}

func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

func fileExists(f fs.FS, path string) bool {
	file, err := f.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return !stat.IsDir()
}
