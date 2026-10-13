package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSigningKeyPersistsAndEnvironmentOverrides(t *testing.T) {
	t.Setenv("MENU_SECRET", "")
	path := filepath.Join(t.TempDir(), ".secret")
	first := loadOrCreateSecret(path)
	if !bytes.Equal(first, loadOrCreateSecret(path)) {
		t.Fatal("key changed")
	}
	t.Setenv("MENU_SECRET", "explicit-key-that-is-at-least-32-bytes-long")
	if string(loadOrCreateSecret(path)) != os.Getenv("MENU_SECRET") {
		t.Fatal("environment did not override persisted key")
	}
}
func TestOriginsAreExactAndCannotInjectCSP(t *testing.T) {
	for _, origin := range []string{"https://example.com", "http://localhost:8080"} {
		if !validOrigin(origin) {
			t.Fatal(origin)
		}
	}
	for _, origin := range []string{"https://*.example.com", "https://example.com; script-src *", "javascript:alert(1)", "https://user:pw@example.com", "https://example.com/path", "https://example.com?q=1"} {
		if validOrigin(origin) {
			t.Fatal(origin)
		}
	}
}
