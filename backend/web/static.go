package web

import (
	"embed"
	"io/fs"
)

// Embed folder dist yang ada di sebelah file ini
//
//go:embed dist/*
var distFS embed.FS

// Helper function agar kita langsung dapat subfolder "dist"
func GetDistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
