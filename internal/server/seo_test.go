package server

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVenueSEOAndStablePublicLanguage(t *testing.T) {
	s, _, v := fixture(t)
	v.Slug = "sample-venue"
	if err := s.store.Restaurants.Update(v); err != nil {
		t.Fatal(err)
	}
	s.cfg.VenueSlug = v.Slug
	for _, path := range []string{"/", "/menu", "/menu?table=7", "/m/sample-venue?lang=fa&table=2"} {
		r := httptest.NewRequest("GET", path, nil)
		r.Header.Set("Accept-Language", "en-US,en;q=0.9")
		r.AddCookie(&http.Cookie{Name: "menu_locale", Value: "en"})
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		if w.Code != 200 || !strings.Contains(w.Body.String(), `lang="fa"`) || !strings.Contains(w.Body.String(), `dir="rtl"`) {
			t.Fatalf("default Persian for %s: %d", path, w.Code)
		}
		canonical := "https://menu.example.com/"
		if path != "/" {
			canonical += "menu"
		}
		for _, want := range []string{`rel="canonical" href="` + canonical + `"`, `hreflang="en"`, `hreflang="x-default"`, `"@type":"Restaurant"`, `application/ld+json`} {
			if !strings.Contains(w.Body.String(), want) {
				t.Fatalf("%s missing %s", path, want)
			}
		}
	}
	w := request(s, "GET", "/menu?lang=en&table=7", nil, nil)
	if !strings.Contains(w.Body.String(), `lang="en"`) || !strings.Contains(w.Body.String(), `rel="canonical" href="https://menu.example.com/menu?lang=en"`) {
		t.Fatal("explicit English canonical lost")
	}
	for _, path := range []string{"/login", "/panel/", "/menu/qr", "/readyz"} {
		if w := request(s, "GET", path, nil, nil); !strings.Contains(w.Header().Get("X-Robots-Tag"), "noindex") {
			t.Fatal("private/utility URL indexable: " + path)
		}
	}
	robots := request(s, "GET", "/robots.txt", nil, nil)
	if !strings.Contains(robots.Body.String(), "Sitemap: https://menu.example.com/sitemap.xml") {
		t.Fatal("missing sitemap discovery")
	}
	w = request(s, "GET", "/sitemap.xml", nil, nil)
	var sitemap sitemapDocument
	if err := xml.Unmarshal(w.Body.Bytes(), &sitemap); err != nil {
		t.Fatal(err)
	}
	if len(sitemap.URLs) != 4 {
		t.Fatalf("sitemap: %d", len(sitemap.URLs))
	}
	for _, u := range sitemap.URLs {
		if strings.Contains(u.Loc, "table=") || strings.Contains(u.Loc, "panel") {
			t.Fatal("noncanonical sitemap URL")
		}
	}
	v.IsPublished = false
	_ = s.store.Restaurants.Update(v)
	if w := request(s, "GET", "/sitemap.xml", nil, nil); w.Code != 404 {
		t.Fatal("draft venue exposed in sitemap")
	}
}

func TestStructuredDataEscapesUserContent(t *testing.T) {
	s, _, v := fixture(t)
	s.cfg.VenueSlug = v.Slug
	v.DescriptionFa = `</script><script>alert(1)</script>`
	view := &View{Restaurant: v}
	s.setVenueSEO(view, false)
	if view.SEO == nil || strings.Contains(string(view.SEO.Structured), "</script>") {
		t.Fatal("unsafe structured data script")
	}
}
