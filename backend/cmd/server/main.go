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
	"cutbray/pppk-json/web"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	logger.New()

	port := utils.GetEnvOrDefault("PORT", "8080")
	ginMode := utils.GetEnvOrDefault("GIN_MODE", gin.ReleaseMode)

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
	examHandler := handlers.NewGinExamHandler(db)
	examHandler.RegisterRoutes(ginEngine)

	distFS, err := web.GetDistFS()

	if err != nil {
		log.Fatal("Gagal load frontend:", err)
	}

	httpFS := http.FS(distFS)

	ginEngine.StaticFS("/assets", http.FS(mustSub(distFS, "assets")))

	indexHTML, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		log.Fatalf("Gagal membaca index.html: %v", err)
	}

	ginEngine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}

		// 2. Cek apakah request adalah FILE FISIK (js, css, png)
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if fileExists(distFS, path) {
			c.FileFromFS(path, httpFS)
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, string(indexHTML))
	})

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
