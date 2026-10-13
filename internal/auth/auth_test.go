package auth

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestPasswordsAndMalformedHashes(t *testing.T) {
	hash, err := HashPassword("CorrectPassword123")
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyPassword("CorrectPassword123", hash); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(VerifyPassword("wrong", hash), ErrMismatch) {
		t.Fatal("wrong password accepted")
	}
	for _, value := range []string{"pbkdf2-sha256$999999999$AA$AA", "pbkdf2-sha256$1000$$", "not-a-hash"} {
		if !errors.Is(VerifyPassword("anything", value), ErrInvalidHash) {
			t.Fatal(value)
		}
	}
}

func TestSessionSignatureExpiryAndFlags(t *testing.T) {
	m := NewManager(bytes.Repeat([]byte("k"), 32), time.Hour, true)
	w := httptest.NewRecorder()
	m.Issue(w, 12, 3)
	cookie := w.Result().Cookies()[0]
	if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatal("unsafe cookie")
	}
	r := httptest.NewRequest("GET", "https://example.com/", nil)
	r.AddCookie(cookie)
	sess, err := m.Read(r)
	if err != nil || sess.UserID != 12 || sess.Version != 3 {
		t.Fatal(sess, err)
	}
	if _, err := m.verify(cookie.Value + "x"); !errors.Is(err, ErrInvalidSession) {
		t.Fatal("signature bypass")
	}
	for _, s := range []Session{{UserID: 12, Version: 3, Expires: time.Now().Add(-time.Minute)}, {UserID: 0, Version: 1, Expires: time.Now().Add(time.Hour)}} {
		if _, err := m.verify(m.sign(s)); err == nil {
			t.Fatal("invalid session accepted")
		}
	}
}

func TestCSRFRejectsQueryAndMismatch(t *testing.T) {
	m := NewManager(bytes.Repeat([]byte("k"), 32), time.Hour, false)
	w := httptest.NewRecorder()
	token := m.IssueCSRF(w, httptest.NewRequest("GET", "/", nil))
	if raw, err := base64.RawURLEncoding.DecodeString(token); err != nil || len(raw) != 24 {
		t.Fatal("invalid token")
	}
	r := httptest.NewRequest("POST", "/?csrf_token="+url.QueryEscape(token), strings.NewReader(""))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(w.Result().Cookies()[0])
	if m.CheckCSRF(r) {
		t.Fatal("query token accepted")
	}
	r.Header.Set("X-CSRF-Token", token)
	if !m.CheckCSRF(r) {
		t.Fatal("valid token rejected")
	}
	r.Header.Set("X-CSRF-Token", strings.Repeat("x", 32))
	if m.CheckCSRF(r) {
		t.Fatal("mismatch accepted")
	}
}
