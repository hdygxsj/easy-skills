package cli

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var staticAssets embed.FS

// GetFS returns the embedded filesystem for the web UI
func GetFS() (fs.FS, error) {
	return fs.Sub(staticAssets, "dist")
}

// ServeFileServer returns an HTTP handler for the embedded files
func ServeFileServer() http.Handler {
	fsys, err := GetFS()
	if err != nil {
		// Fallback to empty handler if no dist
		return http.NotFoundHandler()
	}
	return http.FileServer(http.FS(fsys))
}
