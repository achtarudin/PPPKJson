package handlers

import (
	"cutbray/pppk-json/web"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type frontendHandler struct {
}

func NewFrontendHandler() *frontendHandler {
	return &frontendHandler{}
}

func (h *frontendHandler) RegisterRoutes(router *gin.Engine) {
	distFS, err := web.GetDistFS()

	if err != nil {
		log.Fatal("Gagal load frontend:", err)
	}

	httpFS := http.FS(distFS)

	router.StaticFS("/assets", http.FS(mustSub(distFS, "assets")))

	indexHTML, err := fs.ReadFile(distFS, "index.html")

	if err != nil {
		log.Fatalf("Gagal membaca index.html: %v", err)
	}

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}

		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if fileExists(distFS, path) {
			c.FileFromFS(path, httpFS)
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, string(indexHTML))
	})

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
