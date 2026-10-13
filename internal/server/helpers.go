package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"ResturantMenue/internal/i18n"
)

// tr is a shorthand for i18n.T used where no View is available.
func tr(locale, key string, args ...any) string { return i18n.T(locale, key, args...) }

// --- flash messages --------------------------------------------------------

// Flash is a one-shot message shown after a redirect.
type Flash struct {
	Kind string // success | error | info
	Text string
}

const flashCookie = "menu_flash"

// setFlash stores a message for the next request.
func (s *Server) setFlash(w http.ResponseWriter, kind, text string) {
	payload := kind + "|" + text
	http.SetCookie(w, &http.Cookie{
		Name:     flashCookie,
		Value:    base64.RawURLEncoding.EncodeToString([]byte(payload)),
		Path:     "/",
		MaxAge:   30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// takeFlash reads and immediately clears the flash cookie.
func takeFlash(w http.ResponseWriter, r *http.Request) *Flash {
	c, err := r.Cookie(flashCookie)
	if err != nil || c.Value == "" {
		return nil
	}
	http.SetCookie(w, &http.Cookie{Name: flashCookie, Value: "", Path: "/", MaxAge: -1})

	raw, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return nil
	}
	kind, text, ok := strings.Cut(string(raw), "|")
	if !ok {
		return nil
	}
	return &Flash{Kind: kind, Text: text}
}

// --- request helpers -------------------------------------------------------

// parseForm populates r.Form for both urlencoded and multipart submissions.
//
// http.Request.ParseForm only decodes application/x-www-form-urlencoded
// bodies; a multipart body (what fetch() sends when handed a FormData) leaves
// r.Form holding just the query string. Handlers that inspect r.Form directly
// — rather than going through r.FormValue, which parses lazily — must call
// this instead of ParseForm.
func parseForm(r *http.Request) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return r.ParseMultipartForm(8 << 20)
	}
	return r.ParseForm()
}

func wantsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json") ||
		strings.Contains(r.Header.Get("Content-Type"), "application/json") ||
		r.Header.Get("X-Requested-With") == "fetch"
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// formStr trims a form value.
func formStr(r *http.Request, key string) string {
	return strings.TrimSpace(r.FormValue(key))
}

// formInt parses an integer form value, tolerating thousands separators and
// Persian digits, which are both common in pasted prices.
func formInt(r *http.Request, key string) int64 {
	n, _ := parseWholeNumber(formStr(r, key))
	return n
}

func parseWholeNumber(raw string) (int64, error) {
	raw = normaliseDigits(strings.TrimSpace(raw))
	raw = strings.NewReplacer(",", "", "٬", "", " ", "").Replace(raw)
	if raw == "" {
		return 0, nil
	}
	for _, c := range raw {
		if c < '0' || c > '9' {
			return 0, errors.New("expected non-negative whole number")
		}
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	return n, err
}

func formBool(r *http.Request, key string) bool {
	v := strings.ToLower(formStr(r, key))
	return v == "on" || v == "1" || v == "true" || v == "yes"
}

func pathInt(r *http.Request, key string) int64 {
	n, err := strconv.ParseInt(r.PathValue(key), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// normaliseDigits converts Persian and Arabic-Indic digits to ASCII so numeric
// input works regardless of the user's keyboard.
func normaliseDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= '۰' && r <= '۹': // Persian
			b.WriteRune('0' + (r - '۰'))
		case r >= '٠' && r <= '٩': // Arabic-Indic
			b.WriteRune('0' + (r - '٠'))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// --- validation ------------------------------------------------------------

var (
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)
	slugRe     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	emailRe    = regexp.MustCompile(`^[^@\s]+@[^@\s.]+\.[^@\s]+$`)
)

func validUsername(s string) bool { return usernameRe.MatchString(s) }
func validEmail(s string) bool    { return len(s) <= 254 && emailRe.MatchString(s) }
func validSlug(s string) bool     { return len(s) >= 2 && len(s) <= 60 && slugRe.MatchString(s) }

// Slugify turns a venue name into a URL-safe identifier. Persian text has no
// ASCII transliteration here, so a non-Latin name falls back to a generated
// slug the owner can edit in settings.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_' || r == '.':
			if !lastDash && b.Len() > 0 {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-")
	if len(out) > 60 {
		out = strings.Trim(out[:60], "-")
	}
	return out
}

// safeNext decodes the ?next parameter, accepting only same-site paths so the
// login form cannot be used as an open redirect.
func safeNext(encoded string) string {
	if encoded == "" {
		return "/panel/"
	}
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "/panel/"
	}
	next := string(raw)
	if !safeLocalPath(next) {
		return "/panel/"
	}
	if u, err := url.Parse(next); err != nil || u.Host != "" || u.Scheme != "" {
		return "/panel/"
	}
	return next
}

func safeLocalPath(path string) bool {
	decoded, err := url.PathUnescape(path)
	return err == nil && strings.HasPrefix(decoded, "/") && !strings.HasPrefix(decoded, "//") &&
		!strings.ContainsAny(decoded, "\\\r\n\x00")
}

// truncate shortens a string to n runes, for log lines and labels.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
