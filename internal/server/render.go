package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"sync"

	"ResturantMenue/internal/i18n"
	"ResturantMenue/internal/menutemplate"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/theme"
)

// Renderer parses and caches the page templates.
type Renderer struct {
	fsys  fs.FS
	dev   bool
	mu    sync.RWMutex
	cache map[string]*template.Template
}

func NewRenderer(fsys fs.FS, dev bool) *Renderer {
	return &Renderer{fsys: fsys, dev: dev, cache: map[string]*template.Template{}}
}

// page composes a layout with a page template plus every partial. In dev mode
// templates are re-parsed on each request so edits show up on refresh.
func (r *Renderer) page(layout, name string) (*template.Template, error) {
	key := layout + "|" + name

	if !r.dev {
		r.mu.RLock()
		if t, ok := r.cache[key]; ok {
			r.mu.RUnlock()
			return t, nil
		}
		r.mu.RUnlock()
	}

	patterns := []string{
		"layouts/" + layout + ".html",
		"pages/" + name + ".html",
	}
	// Partials are optional: a fresh checkout may not have any yet.
	if matches, err := fs.Glob(r.fsys, "partials/*.html"); err == nil {
		patterns = append(patterns, matches...)
	}

	t, err := template.New(layout+".html").Funcs(funcMap()).ParseFS(r.fsys, patterns...)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", key, err)
	}

	if !r.dev {
		r.mu.Lock()
		r.cache[key] = t
		r.mu.Unlock()
	}
	return t, nil
}

// Render writes a page, buffering first so a template error cannot leave a
// half-written response with a 200 status.
func (r *Renderer) Render(w http.ResponseWriter, status int, layout, name string, data any) {
	t, err := r.page(layout, name)
	if err != nil {
		httpError(w, err)
		return
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		httpError(w, fmt.Errorf("execute %s: %w", name, err))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}

func httpError(w http.ResponseWriter, err error) {
	// The detail goes to the log, not the browser.
	logf("render error: %v", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// funcMap is the helper set available inside every template.
func funcMap() template.FuncMap {
	return template.FuncMap{
		// --- i18n ---
		"t": i18n.T,
		"num": func(locale string, n any) string {
			return i18n.Number(locale, toInt64(n))
		},
		"price": func(locale, currency string, amount any) string {
			return i18n.Price(locale, currency, toInt64(amount))
		},
		"currencyLabel": i18n.CurrencyLabel,
		"weekdays":      i18n.WeekdayNames,
		"locales":       func() []i18n.Locale { return i18n.Supported },

		// customCSS re-sanitises on the way out, so a row written by an older
		// build can never reach the page unfiltered.
		"customCSS": theme.CustomCSS,

		// --- theme metadata for the editor ---
		"presets": theme.Presets,
		"faFonts": theme.FaFonts,
		"enFonts": theme.EnFonts,

		// --- menu templates ---
		"menuTemplates": menutemplate.All,

		// canvasView wraps raw canvas data for the editor stage, which renders
		// through the same partial the public page uses.
		"canvasView": func(data models.CanvasData, locale string) CanvasView {
			return CanvasView{Data: data, Locale: locale, Edit: true}
		},

		// --- small utilities ---
		"dict":      dict,
		"list":      func(items ...any) []any { return items },
		"jsonAttr":  jsonAttr,
		"seq":       seq,
		"add":       func(a, b int) int { return a + b },
		"sub":       func(a, b int) int { return a - b },
		"hasTag":    hasTag,
		"hasPrefix": strings.HasPrefix,
		"deref":     deref,
		"nonEmpty":  func(vals ...string) string { return firstNonEmpty(vals...) },
		"pct":       func(f float64) int { return int(f * 100) },

		// attr marks a pre-validated value as safe for an inline style
		// attribute. Only ever called with theme values that have been
		// through theme.Sanitize.
		"styleAttr": func(s string) template.CSS { return template.CSS(s) },
	}
}

// jsonAttr serialises a value for a data-* attribute. It returns a plain
// string, so html/template still HTML-escapes it for the attribute context;
// the browser decodes the entities before JSON.parse sees them.
func jsonAttr(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("jsonAttr: %w", err)
	}
	return string(b), nil
}

// dict builds a map so partials can be called with named arguments.
func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict: odd number of arguments")
	}
	out := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key %d is not a string", i)
		}
		out[key] = values[i+1]
	}
	return out, nil
}

// seq produces [from, to] for range loops in templates.
func seq(from, to int) []int {
	if to < from {
		return nil
	}
	out := make([]int, 0, to-from+1)
	for i := from; i <= to; i++ {
		out = append(out, i)
	}
	return out
}

func hasTag(tags []string, want string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, want) {
			return true
		}
	}
	return false
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	default:
		return 0
	}
}
