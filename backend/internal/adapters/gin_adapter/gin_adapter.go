package gin_adapter

import (
	"context"
	"cutbray/pppk-json/internal/ports"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var _ ports.AdapterPort = &ginAdapter{}

type ginAdapter struct {
	server *http.Server
	engine *gin.Engine
	port   string
}

// GinConfig holds configuration for Gin server
type GinConfig struct {
	Port string
	Mode string // debug, release, test
}

func New(config GinConfig) *ginAdapter {
	if config.Port == "" {
		config.Port = "8080"
	}

	if config.Mode == "" {
		config.Mode = gin.ReleaseMode
	}

	gin.SetMode(config.Mode)

	engine := gin.New()

	// Add middlewares
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware())

	adapter := &ginAdapter{
		engine: engine,
		port:   config.Port,
	}

	adapter.setupRoutes()

	return adapter
}

func (g *ginAdapter) Connect(ctx context.Context) error {

	lc := net.ListenConfig{}

	ln, err := lc.Listen(ctx, "tcp", ":"+g.port)

	if err != nil {
		return err
	}
	g.server = &http.Server{
		Addr:         ":" + g.port,
		Handler:      g.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := g.server.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("[Error] Gin server failed to start: %v", err)
		}
	}()

	return nil
}

func (g *ginAdapter) Disconnect(ctx context.Context) error {
	if g.server == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := g.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[Error] Gin server shutdown error: %v", err)
		return err
	}

	return nil
}

func (g *ginAdapter) IsReady() bool {
	return g.server != nil
}

func (g *ginAdapter) Value() any {
	return g.engine
}

func (g *ginAdapter) setupRoutes() {
	// Health check endpoint
	g.engine.GET("/health", g.healthCheck)

	// Swagger documentation
	g.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := g.engine.Group("/api/v1")
	{
		v1.GET("/health", g.healthCheck)
	}
}

// @Summary Health Check
// @Description Get the health status of the application
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (g *ginAdapter) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "PPPKJson Exam API",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With, X-HTTP-Method-Override, Accept-Language, Accept-Encoding")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
