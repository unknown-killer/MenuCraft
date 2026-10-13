// Package config loads runtime settings from environment variables with
// sensible defaults, so the app runs with zero configuration in development.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr            string
	DatabasePath    string
	UploadDir       string
	SecretKey       []byte
	SessionDuration time.Duration
	MaxUploadBytes  int64
	DefaultLocale   string
	Env             string
}

// Load reads configuration from the environment. Any missing value falls back
// to a development-friendly default.
func Load() *Config {
	root := workingRoot()

	cfg := &Config{
		Addr:            env("MENU_ADDR", ":8080"),
		DatabasePath:    env("MENU_DB", filepath.Join(root, "data", "menu.db")),
		UploadDir:       env("MENU_UPLOADS", filepath.Join(root, "uploads")),
		SessionDuration: envDuration("MENU_SESSION_TTL", 30*24*time.Hour),
		MaxUploadBytes:  envInt64("MENU_MAX_UPLOAD", 8<<20), // 8 MiB
		DefaultLocale:   env("MENU_LOCALE", "fa"),
		Env:             env("MENU_ENV", "development"),
	}
	for _, dir := range []string{filepath.Dir(cfg.DatabasePath), cfg.UploadDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatalf("config: cannot create directory %s: %v", dir, err)
		}
	}

	// The signing key lives beside the database rather than under the working
	// directory. In a container the process starts at "/" while the data sits
	// on a mounted volume; anchoring to the database keeps the key on that same
	// volume, so restarts do not silently sign every session out.
	cfg.SecretKey = loadOrCreateSecret(filepath.Join(filepath.Dir(cfg.DatabasePath), ".secret"))

	return cfg
}

func (c *Config) IsProduction() bool { return strings.EqualFold(c.Env, "production") }

// loadOrCreateSecret keeps the session signing key stable across restarts so
// users are not logged out every time the server reboots.
func loadOrCreateSecret(path string) []byte {
	if raw, err := os.ReadFile(path); err == nil {
		if key, err := hex.DecodeString(strings.TrimSpace(string(raw))); err == nil && len(key) >= 32 {
			return key
		}
	}
	if v := os.Getenv("MENU_SECRET"); len(v) >= 32 {
		return []byte(v)
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("config: cannot generate secret key: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err == nil {
		_ = os.WriteFile(path, []byte(hex.EncodeToString(key)), 0o600)
	}
	return key
}

func workingRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envInt64(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return fallback
}
