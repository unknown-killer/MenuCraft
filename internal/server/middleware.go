package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
)

type ctxKey string

const (
	userCtxKey  ctxKey = "user"
	nonceCtxKey ctxKey = "nonce"
)

// nonceFrom returns the per-request CSP nonce.
func nonceFrom(r *http.Request) string {
	n, _ := r.Context().Value(nonceCtxKey).(string)
	return n
}

// userFrom pulls the authenticated user out of the request context.
func userFrom(r *http.Request) *models.User {
	u, _ := r.Context().Value(userCtxKey).(*models.User)
	return u
}

// requireAuth verifies the session cookie, loads the user, and rejects any
// unsafe method whose CSRF token does not match.
func (s *Server) requireAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.sessions.Read(r)
		if err != nil {
			s.redirectToLogin(w, r)
			return
		}

		user, err := s.store.Users.ByID(sess.UserID)
		if err != nil {
			if !errors.Is(err, repository.ErrNotFound) {
				logf("auth: load user %d: %v", sess.UserID, err)
			}
			s.sessions.Clear(w)
			s.redirectToLogin(w, r)
			return
		}

		// A password change bumps session_version, retiring older cookies.
		if user.SessionVersion != sess.Version {
			s.sessions.Clear(w)
			s.redirectToLogin(w, r)
			return
		}

		if isUnsafeMethod(r.Method) && !s.sessions.CheckCSRF(r) {
			s.failCSRF(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isUnsafeMethod(m string) bool {
	switch m {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

func (s *Server) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	if wantsJSON(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	next := r.URL.RequestURI()
	http.Redirect(w, r, "/login?next="+base64.RawURLEncoding.EncodeToString([]byte(next)), http.StatusSeeOther)
}

func (s *Server) failCSRF(w http.ResponseWriter, r *http.Request) {
	loc := s.newView(w, r).Locale.Code
	if wantsJSON(r) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": tr(loc, "err.csrf")})
		return
	}
	http.Error(w, tr(loc, "err.csrf"), http.StatusForbidden)
}

// securityHeaders applies a conservative baseline to every response.
//
// The CSP allows inline styles because the theme engine emits a per-venue
// <style> block, but keeps script-src to 'self' so no injected markup can run
// JavaScript.
func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A per-request nonce lets the public menu run the one inline script it
		// needs (applying the saved light/dark choice before first paint)
		// without opening script-src up to 'unsafe-inline'.
		buf := make([]byte, 16)
		nonce := ""
		if _, err := rand.Read(buf); err == nil {
			nonce = base64.RawStdEncoding.EncodeToString(buf)
		}

		script := "script-src 'self'"
		if nonce != "" {
			script += " 'nonce-" + nonce + "'"
		}

		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-Frame-Options", "SAMEORIGIN")
		h.Set("Content-Security-Policy", strings.Join([]string{
			"default-src 'self'",
			"img-src 'self' data: blob:",
			"style-src 'self' 'unsafe-inline'",
			script,
			"font-src 'self' data:",
			"frame-ancestors 'self'",
			"base-uri 'self'",
			"form-action 'self'",
		}, "; "))

		ctx := context.WithValue(r.Context(), nonceCtxKey, nonce)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// recoverPanic keeps one bad request from taking the process down.
func (s *Server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logf("panic on %s %s: %v", r.Method, r.URL.Path, rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// cacheControl caches static assets aggressively in production and not at all
// in development.
//
// Production URLs carry a ?v= build stamp (see View.AssetVer), so a new release
// changes the URL and the immutable cache is bypassed automatically. Without
// that stamp a long max-age would serve stale CSS and JS after every deploy.
func (s *Server) cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.IsProduction() {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}

// noSniff guards the uploads directory: everything is served as a download-safe
// type and never inherits the page's origin privileges.
func noSniff(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'none'; script-src 'none'; sandbox")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		next.ServeHTTP(w, r)
	})
}
