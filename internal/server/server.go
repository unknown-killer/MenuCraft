package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"ResturantMenue/internal/auth"
	"ResturantMenue/internal/config"
	"ResturantMenue/internal/i18n"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
	"ResturantMenue/internal/theme"
	"ResturantMenue/web"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func logf(format string, args ...any) { logger.Printf(format, args...) }

// Server holds every dependency the handlers need.
type Server struct {
	cfg        *config.Config
	store      *repository.Store
	sessions   *auth.Manager
	render     *Renderer
	mux        *http.ServeMux
	assetVer   string
	authLimits *loginLimiter
}

func New(cfg *config.Config, store *repository.Store) *Server {
	s := &Server{
		cfg:        cfg,
		store:      store,
		sessions:   auth.NewManager(cfg.SecretKey, cfg.SessionDuration, cfg.IsProduction()),
		render:     NewRenderer(web.Templates(), !cfg.IsProduction()),
		mux:        http.NewServeMux(),
		assetVer:   assetVersion(),
		authLimits: newLoginLimiter(),
	}
	s.routes()
	return s
}

// assetVersion produces the ?v= stamp appended to static asset URLs. It is
// derived from the executable's build time so every release busts the cache;
// in development it changes on each restart.
func assetVersion() string {
	if exe, err := os.Executable(); err == nil {
		if info, err := os.Stat(exe); err == nil {
			return strconv.FormatInt(info.ModTime().Unix(), 36)
		}
	}
	return strconv.FormatInt(time.Now().Unix(), 36)
}

// Handler returns the fully wrapped root handler.
func (s *Server) Handler() http.Handler {
	return s.recoverPanic(s.securityHeaders(s.limitBody(s.mux)))
}

func (s *Server) routes() {
	// --- static assets ---
	static := http.FileServer(http.FS(web.Static()))
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", s.cacheControl(static)))

	// User uploads are served from disk, never executed.
	uploads := http.FileServer(http.Dir(s.cfg.UploadDir))
	s.mux.Handle("GET /uploads/", s.onlyUploadedFiles(http.StripPrefix("/uploads/", noSniff(uploads))))

	// --- public ---
	s.mux.HandleFunc("GET /{$}", s.handleLanding)
	s.mux.HandleFunc("GET /sitemap.xml", s.handleSitemap)
	s.mux.HandleFunc("GET /robots.txt", s.handleRobots)
	s.mux.HandleFunc("GET /menu", s.handleVenueMenu)
	s.mux.HandleFunc("GET /menu/qr", s.handleVenueQR)
	s.mux.HandleFunc("GET /m/{slug}", s.handlePublicMenu)
	s.mux.HandleFunc("GET /m/{slug}/qr", s.handleQRCode)
	s.mux.HandleFunc("GET /lang/{code}", s.handleSwitchLocale)
	s.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	s.mux.HandleFunc("GET /readyz", s.handleReadiness)

	// --- auth ---
	s.mux.HandleFunc("GET /login", s.handleLoginForm)
	s.mux.Handle("POST /login", s.limitAuthentication(http.HandlerFunc(s.handleLogin)))
	s.mux.HandleFunc("GET /register", s.handleRegisterForm)
	s.mux.Handle("POST /register", s.limitAuthentication(http.HandlerFunc(s.handleRegister)))
	s.mux.HandleFunc("POST /logout", s.handleLogout)

	// --- panel (all require a session) ---
	p := func(pattern string, h http.HandlerFunc) {
		s.mux.Handle(pattern, s.requireAuth(h))
	}

	p("GET /panel/{$}", s.handleDashboard)
	p("GET /panel/menu", s.handleMenuPage)
	p("GET /panel/transfer", s.handleTransferPage)
	p("GET /panel/menu/export", s.handleMenuExport)
	p("POST /panel/menu/import", s.handleMenuImport)
	p("POST /panel/menu/category", s.handleCategorySave)
	p("POST /panel/menu/category/{id}/delete", s.handleCategoryDelete)
	p("POST /panel/menu/category/reorder", s.handleCategoryReorder)
	p("POST /panel/menu/item", s.handleItemSave)
	p("POST /panel/menu/item/{id}/delete", s.handleItemDelete)
	p("POST /panel/menu/item/{id}/toggle", s.handleItemToggle)
	p("POST /panel/menu/item/reorder", s.handleItemReorder)
	p("POST /panel/menu/template", s.handleApplyMenuTemplate)

	p("GET /panel/gallery", s.handleGallery)
	p("POST /panel/gallery/apply", s.handleApplyPack)

	p("GET /panel/design", s.handleDesignPage)
	p("POST /panel/design/block", s.handleBlockSave)
	p("POST /panel/design/block/{id}/delete", s.handleBlockDelete)
	p("POST /panel/design/sections", s.handleSectionOrder)

	// Free-form canvas designer.
	p("GET /panel/design/canvas/new", s.handleCanvasEdit)
	p("GET /panel/design/canvas/{id}", s.handleCanvasEdit)
	p("POST /panel/design/canvas", s.handleCanvasSave)
	p("POST /panel/design/canvas/render", s.handleCanvasRender)
	p("POST /panel/design/canvas/element", s.handleCanvasElement)

	p("GET /panel/theme", s.handleThemePage)
	p("POST /panel/theme", s.handleThemeSave)
	p("POST /panel/theme/css", s.handleThemeCSS)
	p("POST /panel/theme/reset", s.handleThemeReset)
	p("GET /panel/theme/preview", s.handleThemePreview)

	p("GET /panel/settings", s.handleSettingsPage)
	p("GET /panel/qr", s.handleQRStands)
	p("POST /panel/settings", s.handleSettingsSave)
	p("POST /panel/settings/password", s.handlePasswordChange)

	p("GET /panel/media", s.handleMediaPage)
	p("POST /panel/media/upload", s.handleUpload)
	p("POST /panel/media/{id}/delete", s.handleMediaDelete)
}

// --- view model ------------------------------------------------------------

// View is the data every template receives. Page-specific values go in Page.
type View struct {
	Locale      i18n.Locale
	CSRF        string
	User        *models.User
	Restaurant  *models.Restaurant
	Theme       theme.Config
	ThemePreset string
	Path        string
	Title       string
	SEO         *SEOData
	Flash       *Flash
	AssetVer    string
	// Nonce authorises the single inline script the public menu needs.
	Nonce string
	Page  any
}

// Asset stamps a static URL with the build version so caches update on deploy.
func (v *View) Asset(path string) string { return "/static/" + path + "?v=" + v.AssetVer }

// T translates a key in the request's locale.
func (v *View) T(key string, args ...any) string { return i18n.T(v.Locale.Code, key, args...) }

// N formats a number in the request's locale.
func (v *View) N(n any) string { return i18n.Number(v.Locale.Code, toInt64(n)) }

// Money formats an amount with the venue's currency.
func (v *View) Money(amount any) string {
	cur := "IRT"
	if v.Restaurant != nil && v.Restaurant.Currency != "" {
		cur = v.Restaurant.Currency
	}
	return i18n.Price(v.Locale.Code, cur, toInt64(amount))
}

// IsRTL lets templates branch on direction.
func (v *View) IsRTL() bool { return v.Locale.IsRTL }

// Lang is the short locale code.
func (v *View) Lang() string { return v.Locale.Code }

// Active marks the current nav item.
func (v *View) Active(prefix string) bool { return strings.HasPrefix(v.Path, prefix) }

// newView assembles the common view data for a request.
func (s *Server) newView(w http.ResponseWriter, r *http.Request) *View {
	loc := i18n.Detect(r, s.cfg.DefaultLocale)
	v := &View{
		Locale:   loc,
		CSRF:     s.sessions.IssueCSRF(w, r),
		Path:     r.URL.Path,
		Flash:    takeFlash(w, r),
		Theme:    theme.Default(),
		AssetVer: s.assetVer,
		Nonce:    nonceFrom(r),
	}

	if u := userFrom(r); u != nil {
		v.User = u
		// The panel follows the account's preferred language unless the
		// request explicitly asked for another one.
		if r.URL.Query().Get("lang") == "" {
			if _, err := r.Cookie("menu_locale"); err != nil && u.Locale != "" {
				v.Locale = i18n.Get(u.Locale)
			}
		}
		if rest, err := s.store.Restaurants.ByUser(u.ID); err == nil {
			v.Restaurant = rest
			if preset, cfg, err := s.store.Themes.Get(rest.ID); err == nil {
				v.Theme = cfg
				v.ThemePreset = preset
			}
		}
	}
	return v
}

// --- graceful shutdown -----------------------------------------------------

// Shutdown is a hook for main to close resources owned by the server.
func (s *Server) Shutdown(ctx context.Context) error { return nil }
