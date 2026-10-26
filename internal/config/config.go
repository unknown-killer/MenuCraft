// Package config loads runtime settings from environment variables with
// sensible defaults, so the app runs with zero configuration in development.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr                string
	DatabasePath        string
	UploadDir           string
	SecretKey           []byte
	SessionDuration     time.Duration
	MaxUploadBytes      int64
	DefaultLocale       string
	Env                 string
	PublicURL           string
	EmbedOrigins        []string
	DisableRegistration bool
	VenueSlug           string
}

// Load reads configuration from the environment. Any missing value falls back
// to a development-friendly default.
func Load() *Config {
	root := workingRoot()

	cfg := &Config{
		Addr:                env("MENU_ADDR", ":8080"),
		DatabasePath:        env("MENU_DB", filepath.Join(root, "data", "menu.db")),
		UploadDir:           env("MENU_UPLOADS", filepath.Join(root, "uploads")),
		SessionDuration:     envDuration("MENU_SESSION_TTL", 30*24*time.Hour),
		MaxUploadBytes:      envInt64("MENU_MAX_UPLOAD", 8<<20), // 8 MiB
		DefaultLocale:       env("MENU_LOCALE", "fa"),
		Env:                 env("MENU_ENV", "development"),
		PublicURL:           strings.TrimRight(env("MENU_PUBLIC_URL", ""), "/"),
		DisableRegistration: strings.EqualFold(env("MENU_DISABLE_REGISTRATION", "false"), "true"),
		VenueSlug:           env("MENU_VENUE_SLUG", ""),
	}
	if cfg.PublicURL != "" && !validOrigin(cfg.PublicURL) {
		log.Fatal("config: MENU_PUBLIC_URL must be an http(s) origin without a path, credentials or query")
	}
	for _, origin := range strings.Split(env("MENU_EMBED_ORIGINS", ""), ",") {
		origin = strings.TrimRight(strings.TrimSpace(origin), "/")
		if origin == "" {
			continue
		}
		if !validOrigin(origin) {
			log.Fatal("config: MENU_EMBED_ORIGINS must contain exact http(s) origins")
		}
		cfg.EmbedOrigins = append(cfg.EmbedOrigins, origin)
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
	if v := os.Getenv("MENU_SECRET"); v != "" {
		if len(v) < 32 {
			log.Fatal("config: MENU_SECRET must contain at least 32 bytes")
		}
		return []byte(v)
	}
	if raw, err := os.ReadFile(path); err == nil {
		if key, err := hex.DecodeString(strings.TrimSpace(string(raw))); err == nil && len(key) >= 32 {
			return key
		}
		log.Fatal("config: stored signing key is invalid; restore it or explicitly set MENU_SECRET")
	} else if !os.IsNotExist(err) {
		log.Fatalf("config: cannot read signing key: %v", err)
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("config: cannot generate secret key: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if os.IsExist(err) {
		return loadOrCreateSecret(path)
	}
	if err != nil {
		log.Fatalf("config: cannot create signing key: %v", err)
	}
	_, writeErr := f.WriteString(hex.EncodeToString(key))
	closeErr := f.Close()
	if writeErr != nil || closeErr != nil {
		log.Fatal("config: cannot persist signing key")
	}
	return key
}

func validOrigin(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Hostname() != "" &&
		u.User == nil && u.Path == "" && u.RawQuery == "" && u.Fragment == "" &&
		!strings.ContainsAny(raw, " \t\r\n;*\\")
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
