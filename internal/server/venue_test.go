package server

import (
	"ResturantMenue/internal/models"
	"bytes"
	qrcode "github.com/skip2/go-qrcode"
	"strings"
	"testing"
)

func TestSingleVenueRoutesAndQR(t *testing.T) {
	s, u, v := fixture(t)
	if w := request(s, "GET", "/menu", nil, nil); w.Code != 404 {
		t.Fatal("alias must be disabled by default")
	}
	s.cfg.VenueSlug = v.Slug
	for _, p := range []string{"/", "/menu?table=۷&lang=en"} {
		w := request(s, "GET", p, nil, session(s, u))
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", p, w.Code, w.Body.String())
		}
	}
	want, _ := qrcode.Encode("https://menu.example.com/menu?table=7", qrcode.Medium, 512)
	if w := request(s, "GET", "/menu/qr?table=7", nil, nil); w.Code != 200 || !bytes.Equal(w.Body.Bytes(), want) {
		t.Fatal("alias QR lost canonical path/table")
	}
	w := request(s, "GET", "/panel/qr?start=1&end=7", nil, session(s, u))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "https://menu.example.com/menu?table=7") {
		t.Fatal("stands must use canonical menu alias")
	}
	v.IsPublished = false
	if err := s.store.Restaurants.Update(v); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"/", "/menu", "/menu/qr"} {
		if w := request(s, "GET", p, nil, nil); w.Code != 404 {
			t.Fatalf("draft exposed via %s", p)
		}
	}
}

func TestVenueLandingUsesSavedDetailsAndMenuEmbed(t *testing.T) {
	s, _, v := fixture(t)
	v.NameEn = "Example Dining"
	v.Phone = "+1 202 555 0100 / +1 202 555 0101"
	v.AddressEn = "Example street"
	if err := s.store.Restaurants.Update(v); err != nil {
		t.Fatal(err)
	}
	if err := s.store.Restaurants.SaveWorkingHours(v.ID, []*models.WorkingHour{{DayOfWeek: 1, OpensAt: "09:15", ClosesAt: "18:45"}}); err != nil {
		t.Fatal(err)
	}
	s.cfg.VenueSlug = v.Slug
	s.cfg.EmbedOrigins = []string{"https://wordpress.example.com"}
	w := request(s, "GET", "/?lang=en", nil, nil)
	for _, want := range []string{"Example Dining", "Example street", "09:15", "18:45", "202 555 0100", "202 555 0101"} {
		if !strings.Contains(w.Body.String(), want) {
			t.Fatalf("landing missing %s", want)
		}
	}
	if strings.Count(w.Body.String(), `href="tel:`) != 2 {
		t.Fatal("expected two phone links")
	}
	w = request(s, "GET", "/menu", nil, nil)
	if !strings.Contains(w.Header().Get("Content-Security-Policy"), "frame-ancestors 'self' https://wordpress.example.com") || w.Header().Get("X-Frame-Options") != "" {
		t.Fatal("short menu URL must support configured WordPress embeds")
	}
	for _, path := range []string{"/", "/login", "/panel/"} {
		w := request(s, "GET", path, nil, nil)
		if strings.Contains(w.Header().Get("Content-Security-Policy"), "wordpress.example.com") {
			t.Fatal("non-menu page became embeddable")
		}
	}
}
