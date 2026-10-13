package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SessionCookie = "menu_session"
	CSRFCookie    = "menu_csrf"
	CSRFField     = "csrf_token"
	LocaleCookie  = "menu_locale"
)

var ErrInvalidSession = errors.New("auth: invalid session")

// Session is the payload carried by the signed cookie. Sessions are stateless:
// there is no server-side table, and revocation works by bumping the user's
// session_version column.
type Session struct {
	UserID  int64
	Version int
	Expires time.Time
}

// Manager issues and verifies signed cookies.
type Manager struct {
	key    []byte
	ttl    time.Duration
	secure bool
}

func NewManager(key []byte, ttl time.Duration, secure bool) *Manager {
	return &Manager{key: key, ttl: ttl, secure: secure}
}

// Issue writes a fresh session cookie for the user.
func (m *Manager) Issue(w http.ResponseWriter, userID int64, version int) {
	s := Session{UserID: userID, Version: version, Expires: time.Now().Add(m.ttl)}
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    m.sign(s),
		Path:     "/",
		Expires:  s.Expires,
		MaxAge:   int(m.ttl.Seconds()),
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// Clear expires the session cookie.
func (m *Manager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// Read verifies the cookie signature and expiry, returning the session.
func (m *Manager) Read(r *http.Request) (Session, error) {
	c, err := r.Cookie(SessionCookie)
	if err != nil {
		return Session{}, ErrInvalidSession
	}
	return m.verify(c.Value)
}

// sign encodes "userID.version.expiry" and appends an HMAC over that payload.
func (m *Manager) sign(s Session) string {
	payload := strconv.FormatInt(s.UserID, 10) + "." +
		strconv.Itoa(s.Version) + "." +
		strconv.FormatInt(s.Expires.Unix(), 10)

	mac := hmac.New(sha256.New, m.key)
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}

func (m *Manager) verify(token string) (Session, error) {
	dot := strings.LastIndex(token, ".")
	if dot <= 0 {
		return Session{}, ErrInvalidSession
	}
	rawPayload, sig := token[:dot], token[dot+1:]

	payloadBytes, err := base64.RawURLEncoding.DecodeString(rawPayload)
	if err != nil {
		return Session{}, ErrInvalidSession
	}

	mac := hmac.New(sha256.New, m.key)
	mac.Write(payloadBytes)
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sig), []byte(expected)) != 1 {
		return Session{}, ErrInvalidSession
	}

	parts := strings.Split(string(payloadBytes), ".")
	if len(parts) != 3 {
		return Session{}, ErrInvalidSession
	}
	userID, err1 := strconv.ParseInt(parts[0], 10, 64)
	version, err2 := strconv.Atoi(parts[1])
	expUnix, err3 := strconv.ParseInt(parts[2], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil || userID < 1 || version < 1 {
		return Session{}, ErrInvalidSession
	}

	expires := time.Unix(expUnix, 0)
	if !time.Now().Before(expires) {
		return Session{}, ErrInvalidSession
	}
	return Session{UserID: userID, Version: version, Expires: expires}, nil
}

// --- CSRF ------------------------------------------------------------------

// IssueCSRF returns the request's CSRF token, creating and setting one when
// the cookie is missing. This is the double-submit pattern: the attacker's
// site can trigger a request but cannot read the cookie to echo it back.
func (m *Manager) IssueCSRF(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(CSRFCookie); err == nil && len(c.Value) == 32 {
		return c.Value
	}

	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	token := base64.RawURLEncoding.EncodeToString(buf)

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(m.ttl.Seconds()),
		HttpOnly: false, // readable by JS so fetch() calls can send the header
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
	return token
}

// CheckCSRF compares the submitted token against the cookie in constant time.
// The token may arrive as a form field or as the X-CSRF-Token header.
func (m *Manager) CheckCSRF(r *http.Request) bool {
	c, err := r.Cookie(CSRFCookie)
	if err != nil || len(c.Value) != 32 {
		return false
	}

	submitted := r.Header.Get("X-CSRF-Token")
	if submitted == "" {
		// Query parameters must never stand in for a token in the request body.
		submitted = r.PostFormValue(CSRFField)
	}
	if submitted == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(submitted), []byte(c.Value)) == 1
}
