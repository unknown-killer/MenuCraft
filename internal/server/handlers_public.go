package server

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
	"ResturantMenue/internal/theme"
)

type publicPage struct {
	Venue       *models.Restaurant
	Categories  []*models.Category
	Hours       []*models.WorkingHour
	Config      theme.Config
	CSSVars     template.CSS
	BodyClass   string
	IsPreview   bool
	MenuURL     string
	HasItems    bool
	TableNumber string

	// Sections is the page in the owner's drag-and-drop order. The template
	// walks this instead of hard-coding the layout.
	Sections []PageSection
}

// PageSection is one entry of the ordered public page.
type PageSection struct {
	Key   string     // built-in section key, or "block"
	Block *BlockView // set when Key == "block"
}

// IsBlock reports whether this entry renders custom content.
func (p PageSection) IsBlock() bool { return p.Key == "block" }

type landingPage struct {
	Presets []theme.Preset
}

func (s *Server) handleLanding(w http.ResponseWriter, r *http.Request) {
	if s.cfg.VenueSlug != "" {
		s.handleVenueLanding(w, r)
		return
	}
	v := s.newView(w, r)
	// A signed-in owner goes straight to their panel.
	if s.currentUser(r) != nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	v.Title = v.T("app.name")
	v.Page = &landingPage{Presets: theme.Presets()}
	s.render.Render(w, http.StatusOK, "public", "landing", v)
}

func (s *Server) handlePublicMenu(w http.ResponseWriter, r *http.Request) {
	slug := strings.ToLower(r.PathValue("slug"))

	venue, err := s.store.Restaurants.BySlug(slug)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			logf("public: lookup slug %q: %v", truncate(slug, 60), err)
		}
		s.renderError(w, r, http.StatusNotFound)
		return
	}
	s.renderPublicMenu(w, r, venue, false)
}

// renderPublicMenu draws the themed menu. preview mode is used by the editor
// iframe: it bypasses the published check and injects the live-update hook.
func (s *Server) renderPublicMenu(w http.ResponseWriter, r *http.Request, venue *models.Restaurant, preview bool) {
	v := s.newView(w, r)
	if !preview {
		s.publicLocale(v, r, venue)
	}
	if preview || !venue.IsPublished {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	}

	if !venue.IsPublished && !preview {
		// The owner can still see their own draft.
		owner := s.currentUser(r)
		if owner == nil || owner.ID != venue.UserID {
			v.Title = v.T("pub.not_published")
			v.Page = map[string]any{"Status": http.StatusNotFound, "Message": v.T("pub.not_published")}
			s.render.Render(w, http.StatusNotFound, "auth", "error", v)
			return
		}
	}

	_, cfg, err := s.store.Themes.Get(venue.ID)
	if err != nil {
		logf("public: theme for venue %d: %v", venue.ID, err)
		cfg = theme.Default()
	}

	cats, err := s.store.FullMenu(venue.ID, true)
	if err != nil {
		logf("public: menu for venue %d: %v", venue.ID, err)
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}

	// Drop categories that ended up with nothing visible in them.
	shown := make([]*models.Category, 0, len(cats))
	total := 0
	for _, c := range cats {
		if len(c.Items) == 0 {
			continue
		}
		shown = append(shown, c)
		total += len(c.Items)
	}

	var hours []*models.WorkingHour
	if cfg.Display.ShowHours {
		if h, err := s.store.Restaurants.WorkingHours(venue.ID); err == nil {
			hours = h
		}
	}

	// Custom blocks: hidden ones stay out of the public page, but the editor
	// preview shows them dimmed so the owner can still find them.
	blocks, err := s.store.Blocks.List(venue.ID, !preview)
	if err != nil {
		logf("public: blocks for venue %d: %v", venue.ID, err)
	}
	blocksByID := make(map[int64]*models.Block, len(blocks))
	ids := make([]int64, 0, len(blocks))
	for _, b := range blocks {
		blocksByID[b.ID] = b
		ids = append(ids, b.ID)
	}

	order := theme.ReconcileSections(cfg.Layout.Sections, ids)
	sections := make([]PageSection, 0, len(order))
	for _, key := range order {
		if id, ok := theme.IsBlockSection(key); ok {
			if b, exists := blocksByID[id]; exists {
				sections = append(sections, PageSection{
					Key:   "block",
					Block: &BlockView{Block: b, Locale: v.Locale.Code},
				})
			}
			continue
		}
		sections = append(sections, PageSection{Key: key})
	}

	// The venue's own currency and language drive formatting here, not the
	// panel user's.
	v.Restaurant = venue
	v.Title = venue.Name(v.Locale.Code)
	v.Theme = cfg
	v.Page = &publicPage{
		Venue:       venue,
		Categories:  shown,
		Hours:       hours,
		Config:      cfg,
		CSSVars:     cfg.Variables(v.Locale.Code),
		BodyClass:   cfg.BodyClasses(),
		IsPreview:   preview,
		MenuURL:     s.menuPath(venue.Slug),
		HasItems:    total > 0,
		TableNumber: validTableNumber(r.URL.Query().Get("table")),
		Sections:    sections,
	}
	if !preview {
		s.setVenueSEO(v, true)
	}
	s.render.Render(w, http.StatusOK, "menu", "public_menu", v)
}

// handleQRCode renders a PNG QR code pointing at the venue's public menu.
func (s *Server) handleQRCode(w http.ResponseWriter, r *http.Request) {
	slug := strings.ToLower(r.PathValue("slug"))

	venue, err := s.store.Restaurants.BySlug(slug)
	if err != nil {
		s.renderError(w, r, http.StatusNotFound)
		return
	}
	if !venue.IsPublished {
		owner := s.currentUser(r)
		if owner == nil || owner.ID != venue.UserID {
			http.NotFound(w, r)
			return
		}
	}

	size := 512
	if n, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil && n >= 128 && n <= 2048 {
		size = n
	}

	table := validTableNumber(r.URL.Query().Get("table"))
	if r.URL.Query().Get("table") != "" && table == "" {
		http.Error(w, "Invalid table number (1-9999)", http.StatusBadRequest)
		return
	}
	png, err := qrcode.Encode(tableMenuURL(s.absoluteMenuURL(r, venue.Slug), table), qrcode.Medium, size)
	if err != nil {
		logf("qr: encode for %q: %v", slug, err)
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "private, no-store")
	if r.URL.Query().Get("download") != "" {
		suffix := ""
		if table != "" {
			suffix = "-table-" + table
		}
		w.Header().Set("Content-Disposition", `attachment; filename="menu-qr-`+venue.Slug+suffix+`.png"`)
	}
	_, _ = w.Write(png)
}

// absoluteMenuURL builds the public URL a QR code should point at.
func (s *Server) absoluteMenuURL(r *http.Request, slug string) string {
	if s.cfg.PublicURL != "" {
		return s.cfg.PublicURL + s.menuPath(slug)
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	return scheme + "://" + host + s.menuPath(slug)
}
