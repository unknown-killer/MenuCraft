package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"ResturantMenue/internal/i18n"
	"ResturantMenue/internal/models"
)

type SEOData struct {
	Title, Description, Canonical, Persian, English, Image, Locale string
	Structured                                                     template.JS
}

// Single-venue public URLs have a stable language independent of cookies,
// account preferences and crawler Accept-Language headers.
func (s *Server) publicLocale(v *View, r *http.Request, venue *models.Restaurant) {
	if venue.Slug != s.cfg.VenueSlug {
		return
	}
	lang := "fa"
	if r.URL.Query().Get("lang") == "en" {
		lang = "en"
	}
	v.Locale = i18n.Get(lang)
}

func (s *Server) setVenueSEO(v *View, menu bool) {
	venue := v.Restaurant
	if venue == nil || !venue.IsPublished || s.cfg.PublicURL == "" || venue.Slug != s.cfg.VenueSlug {
		return
	}
	origin := strings.TrimRight(s.cfg.PublicURL, "/")
	path := "/"
	if menu {
		path = s.menuPath(venue.Slug)
	}
	fa, en := origin+path, origin+path+"?lang=en"
	seo := &SEOData{Persian: fa, English: en, Canonical: fa, Locale: "fa_IR"}
	seo.Title = venue.Name("fa") + " | " + venue.Address("fa")
	seo.Description = venue.Description("fa") + " نشانی: " + venue.Address("fa") + ". مشاهدهٔ منو، قیمت‌ها، ساعت کاری و مسیر دسترسی."
	if menu {
		seo.Title = "منوی " + venue.Name("fa") + " | غذاها و قیمت‌ها"
		seo.Description = "منوی " + venue.Name("fa") + "؛ غذاها، محصولات و قیمت‌ها. انتخاب غذا و یادداشت برای گارسون. " + venue.Address("fa")
	}
	if v.Lang() == "en" {
		seo.Canonical = en
		seo.Locale = "en_US"
		seo.Title = venue.Name("en") + " | " + venue.Address("en")
		seo.Description = venue.Description("en") + " Find our menu, opening hours and directions. " + venue.Address("en")
		if menu {
			seo.Title = venue.Name("en") + " Menu | Food & Prices"
			seo.Description = "Explore the menu at " + venue.Name("en") + ". Save your choices for the waiter. " + venue.Address("en")
		}
	}
	if strings.HasPrefix(venue.LogoPath, "/") && !strings.HasPrefix(venue.LogoPath, "//") {
		seo.Image = origin + venue.LogoPath
	}
	business := map[string]any{"@context": "https://schema.org", "@type": "Restaurant", "@id": origin + "/#restaurant", "name": venue.Name("fa"), "alternateName": venue.Name("en"), "url": origin + "/", "menu": origin + s.menuPath(venue.Slug), "image": seo.Image, "logo": seo.Image, "description": venue.Description(v.Lang()), "address": map[string]any{"@type": "PostalAddress", "streetAddress": venue.Address("fa")}}
	if phones := venue.Phones(); len(phones) > 0 {
		business["telephone"] = phones
	}
	// Marshal escapes HTML-sensitive characters before marking the result as JS;
	// user-provided text cannot close the JSON-LD script element.
	if b, err := json.Marshal(business); err == nil {
		seo.Structured = template.JS(b)
	}
	v.SEO = seo
	v.Title = seo.Title
}

type sitemapURL struct {
	Loc string `xml:"loc"`
}
type sitemapDocument struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

func (s *Server) handleSitemap(w http.ResponseWriter, r *http.Request) {
	if s.cfg.VenueSlug == "" || s.cfg.PublicURL == "" {
		http.NotFound(w, r)
		return
	}
	venue, err := s.store.Restaurants.BySlug(s.cfg.VenueSlug)
	if err != nil || !venue.IsPublished {
		http.NotFound(w, r)
		return
	}
	origin := strings.TrimRight(s.cfg.PublicURL, "/")
	doc := sitemapDocument{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, path := range []string{"/", "/?lang=en", s.menuPath(venue.Slug), s.menuPath(venue.Slug) + "?lang=en"} {
		doc.URLs = append(doc.URLs, sitemapURL{Loc: origin + path})
	}
	body, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		http.Error(w, "Unable to build sitemap", 500)
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write(append([]byte(xml.Header), body...))
}

func (s *Server) handleRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprint(w, "User-agent: *\nAllow: /\nDisallow: /panel/\nDisallow: /logout\nDisallow: /lang/\n")
	if s.cfg.VenueSlug != "" && s.cfg.PublicURL != "" {
		_, _ = fmt.Fprintf(w, "\nSitemap: %s/sitemap.xml\n", strings.TrimRight(s.cfg.PublicURL, "/"))
	}
}
