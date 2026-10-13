// Package web embeds the HTML templates and static assets so the compiled
// binary is self-contained and can be copied to a server on its own.
package web

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed templates static
var embedded embed.FS

// Templates returns the template tree.
func Templates() fs.FS {
	if live, ok := liveDir("templates"); ok {
		return live
	}
	sub, err := fs.Sub(embedded, "templates")
	if err != nil {
		panic("web: templates missing from embed: " + err.Error())
	}
	return sub
}

// Static returns the static asset tree.
func Static() fs.FS {
	if live, ok := liveDir("static"); ok {
		return live
	}
	sub, err := fs.Sub(embedded, "static")
	if err != nil {
		panic("web: static missing from embed: " + err.Error())
	}
	return sub
}

// liveDir serves files straight from disk when the source tree is present and
// MENU_ENV is not production, so editing a template or stylesheet only needs a
// browser refresh instead of a rebuild.
func liveDir(name string) (fs.FS, bool) {
	if os.Getenv("MENU_ENV") == "production" {
		return nil, false
	}
	for _, base := range []string{"web", filepath.Join("..", "web")} {
		dir := filepath.Join(base, name)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return os.DirFS(dir), true
		}
	}
	return nil, false
}
