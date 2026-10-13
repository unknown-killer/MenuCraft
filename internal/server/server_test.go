package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"ResturantMenue/internal/auth"
	"ResturantMenue/internal/config"
	"ResturantMenue/internal/database"
	"ResturantMenue/internal/menudata"
	"ResturantMenue/internal/menupack"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
)

func fixture(t *testing.T) (*Server, *models.User, *models.Restaurant) {
	t.Helper()
	root := t.TempDir()
	db, err := database.Open(filepath.Join(root, "menu # sample.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	store := repository.New(db)
	cfg := &config.Config{DatabasePath: filepath.Join(root, "menu.db"), UploadDir: filepath.Join(root, "uploads"),
		SecretKey: bytes.Repeat([]byte("k"), 32), SessionDuration: time.Hour, MaxUploadBytes: 1024, DefaultLocale: "fa",
		PublicURL: "https://menu.example.com", EmbedOrigins: []string{"https://wp.example.com"}}
	s := New(cfg, store)
	u := &models.User{Username: "owner", Email: "owner@example.com", PasswordHash: dummyHash, DisplayName: "Owner", Locale: "fa"}
	if err := store.Users.Create(u); err != nil {
		t.Fatal(err)
	}
	if err := s.provisionRestaurant(u, "Test Kitchen", "restaurant", "fa"); err != nil {
		t.Fatal(err)
	}
	v, err := store.Restaurants.ByUser(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	v.Phone = "+98 21 12345678"
	if err := store.Restaurants.Update(v); err != nil {
		t.Fatal(err)
	}
	cats, err := store.Categories.List(v.ID)
	if err != nil || len(cats) == 0 {
		t.Fatal(err)
	}
	it := &models.Item{RestaurantID: v.ID, CategoryID: cats[0].ID, NameFa: "سالاد", NameEn: "Salad", Price: 100, IsAvailable: true, Tags: []string{"vegan"}}
	if err := store.Items.Create(it); err != nil {
		t.Fatal(err)
	}
	return s, u, v
}
func session(s *Server, u *models.User) *http.Cookie {
	w := httptest.NewRecorder()
	s.sessions.Issue(w, u.ID, u.SessionVersion)
	return w.Result().Cookies()[0]
}
func request(s *Server, method, path string, form url.Values, cookie *http.Cookie) *httptest.ResponseRecorder {
	var body *strings.Reader
	if form == nil {
		body = strings.NewReader("")
	} else {
		body = strings.NewReader(form.Encode())
	}
	r := httptest.NewRequest(method, path, body)
	if form != nil {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if cookie != nil {
		r.AddCookie(cookie)
	}
	token := strings.Repeat("a", 32)
	r.AddCookie(&http.Cookie{Name: auth.CSRFCookie, Value: token})
	r.Header.Set("X-CSRF-Token", token)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	return w
}
func TestAllPagesRender(t *testing.T) {
	s, u, v := fixture(t)
	for _, lang := range []string{"fa", "en"} {
		for _, path := range []string{"/panel/", "/panel/menu", "/panel/gallery", "/panel/theme", "/panel/design", "/panel/design/canvas/new", "/panel/settings", "/panel/media", "/panel/transfer", "/panel/theme/preview", "/m/" + v.Slug} {
			t.Run(lang+path, func(t *testing.T) {
				w := request(s, "GET", path+"?lang="+lang, nil, session(s, u))
				if w.Code != 200 {
					t.Fatalf("status %d: %s", w.Code, w.Body.String())
				}
				if strings.Contains(w.Body.String(), "ZgotmplZ") {
					t.Fatal("unsafe template URL")
				}
			})
		}
	}
	w := request(s, "GET", "/m/"+v.Slug, nil, nil)
	if !strings.Contains(html.UnescapeString(w.Body.String()), v.Phone) {
		t.Fatal("phone number was corrupted")
	}
}
func TestRevokedSessionDoesNotRedirectFromLogin(t *testing.T) {
	s, u, _ := fixture(t)
	old := session(s, u)
	if err := s.store.Users.UpdatePassword(u.ID, dummyHash); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/login", "/register"} {
		if w := request(s, "GET", path, nil, old); w.Code != 200 {
			t.Fatalf("%s: %d", path, w.Code)
		}
	}
	if w := request(s, "GET", "/panel/", nil, old); w.Code != 303 {
		t.Fatal(w.Code)
	}
}
func TestDraftPrivacyAndOwnerPreview(t *testing.T) {
	s, u, v := fixture(t)
	v.IsPublished = false
	s.store.Restaurants.Update(v)
	for _, path := range []string{"/m/" + v.Slug, "/m/" + v.Slug + "/qr"} {
		if w := request(s, "GET", path, nil, nil); w.Code != 404 {
			t.Fatalf("draft leaked: %d", w.Code)
		}
		if w := request(s, "GET", path, nil, session(s, u)); w.Code != 200 {
			t.Fatalf("owner cannot preview: %d", w.Code)
		}
	}
}
func TestRegistrationRollsBackOnProvisionFailure(t *testing.T) {
	s, _, _ := fixture(t)
	_, err := s.store.DB().Exec("CREATE TRIGGER fail_venue BEFORE INSERT ON restaurants BEGIN SELECT RAISE(ABORT, 'test'); END")
	if err != nil {
		t.Fatal(err)
	}
	w := request(s, "POST", "/register", url.Values{"username": {"newuser"}, "email": {"new@example.com"}, "restaurant_name": {"New"}, "password": {"Password123"}, "password_confirm": {"Password123"}}, nil)
	if w.Code != 422 {
		t.Fatalf("status %d", w.Code)
	}
	if _, err := s.store.Users.ByLogin("newuser"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatal("orphan account remains", err)
	}
}
func TestRegistrationCreatesUsableSession(t *testing.T) {
	s, _, _ := fixture(t)
	w := request(s, "POST", "/register", url.Values{"username": {"newuser"}, "email": {"new@example.com"}, "restaurant_name": {"New"}, "password": {"Password123"}, "password_confirm": {"Password123"}}, nil)
	if w.Code != 303 {
		t.Fatalf("status %d", w.Code)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == auth.SessionCookie {
			if got := request(s, "GET", "/panel/", nil, c); got.Code != 200 {
				t.Fatal("new session unusable")
			}
			return
		}
	}
	t.Fatal("no session")
}
func TestBodyLimitsApplyBeforeMultipartCSRF(t *testing.T) {
	s, u, _ := fixture(t)
	for _, length := range []int64{-1, 3 << 20} {
		r := httptest.NewRequest("POST", "/panel/media/upload", strings.NewReader(strings.Repeat("x", 3<<20)))
		r.ContentLength = length
		r.Header.Set("Content-Type", "multipart/form-data; boundary=test")
		r.AddCookie(session(s, u))
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		if w.Code != 413 && w.Code != 400 {
			t.Fatalf("unbounded body: %d", w.Code)
		}
	}
}
func TestCSRFAndTenantIsolation(t *testing.T) {
	s, u, v := fixture(t)
	r := httptest.NewRequest("POST", "/panel/menu/category/1/delete", nil)
	r.AddCookie(session(s, u))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal(w.Code)
	}
	other := &models.User{Username: "other", Email: "other@example.com", PasswordHash: dummyHash}
	s.store.Users.Create(other)
	s.provisionRestaurant(other, "Other", "cafe", "en")
	cats, _ := s.store.Categories.List(v.ID)
	request(s, "POST", "/panel/menu/category/"+strconvID(cats[0].ID)+"/delete", url.Values{}, session(s, other))
	if _, err := s.store.Categories.ByID(cats[0].ID, v.ID); err != nil {
		t.Fatal("cross-tenant deletion", err)
	}
}
func strconvID(id int64) string { return strconv.FormatInt(id, 10) }
func TestPackFailureKeepsExistingMenu(t *testing.T) {
	s, u, v := fixture(t)
	before, _ := s.store.FullMenu(v.ID, false)
	_, err := s.store.DB().Exec("CREATE TRIGGER fail_item BEFORE INSERT ON items BEGIN SELECT RAISE(ABORT, 'test'); END")
	if err != nil {
		t.Fatal(err)
	}
	packs := menupack.ForBusiness("restaurant")
	if len(packs) == 0 {
		t.Fatal("no packs")
	}
	request(s, "POST", "/panel/gallery/apply", url.Values{"pack": {packs[0].Key}, "replace": {"1"}, "with_theme": {"1"}}, session(s, u))
	after, _ := s.store.FullMenu(v.ID, false)
	if len(after) != len(before) || len(after[0].Items) != 1 || after[0].Items[0].NameEn != "Salad" {
		t.Fatal("replace was not rolled back")
	}
}
func TestTransferRoundTripAndValidation(t *testing.T) {
	s, u, v := fixture(t)
	w := request(s, "GET", "/panel/menu/export", nil, session(s, u))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	if strings.Contains(w.Body.String(), "password") || strings.Contains(w.Body.String(), u.Email) {
		t.Fatal("private data exported")
	}
	var doc menudata.Document
	if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := s.importMenu(v.ID, doc, true); err != nil {
		t.Fatal(err)
	}
	cats, _ := s.store.FullMenu(v.ID, false)
	if len(cats[0].Items) != 1 || cats[0].Items[0].NameEn != "Salad" {
		t.Fatal("round trip failed")
	}
	var body bytes.Buffer
	mp := multipart.NewWriter(&body)
	file, _ := mp.CreateFormFile("menu_file", "bad.json")
	file.Write([]byte("{bad"))
	mp.WriteField("replace", "1")
	mp.Close()
	r := httptest.NewRequest("POST", "/panel/menu/import", &body)
	r.Header.Set("Content-Type", mp.FormDataContentType())
	r.AddCookie(session(s, u))
	r.AddCookie(&http.Cookie{Name: auth.CSRFCookie, Value: strings.Repeat("a", 32)})
	r.Header.Set("X-CSRF-Token", strings.Repeat("a", 32))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, r)
	cats, _ = s.store.FullMenu(v.ID, false)
	if len(cats[0].Items) != 1 {
		t.Fatal("invalid import changed menu")
	}
	if got := request(s, "GET", "/panel/menu/export", nil, nil); got.Code != 303 {
		t.Fatal("anonymous export")
	}
}
func TestSecurityHeadersCanonicalURLAndReadiness(t *testing.T) {
	s, u, v := fixture(t)
	w := request(s, "GET", "/m/"+v.Slug, nil, nil)
	if !strings.Contains(w.Header().Get("Content-Security-Policy"), "https://wp.example.com") {
		t.Fatal("embed origin missing")
	}
	if w.Header().Get("X-Frame-Options") != "" {
		t.Fatal("embed blocked")
	}
	w = request(s, "GET", "/panel/", nil, session(s, u))
	if w.Header().Get("X-Frame-Options") != "SAMEORIGIN" || strings.Contains(w.Header().Get("Content-Security-Policy"), "wp.example") {
		t.Fatal("admin frame boundary")
	}
	r := httptest.NewRequest("GET", "http://internal/m/test", nil)
	r.Header.Set("X-Forwarded-Host", "attacker.invalid")
	if got := s.absoluteMenuURL(r, v.Slug); got != "https://menu.example.com/m/"+v.Slug {
		t.Fatal(got)
	}
	if w := request(s, "GET", "/readyz", nil, nil); w.Code != 200 {
		t.Fatal(w.Code)
	}
	s.store.DB().Close()
	if w := request(s, "GET", "/readyz", nil, nil); w.Code != 503 {
		t.Fatal(w.Code)
	}
}
func TestRateLimiterAndSafeInputs(t *testing.T) {
	l := newLoginLimiter()
	now := time.Now()
	for i := 0; i < 20; i++ {
		if !l.allow("client", now) {
			t.Fatal(i)
		}
	}
	if l.allow("client", now) {
		t.Fatal("limit bypass")
	}
	if !l.allow("client", now.Add(61*time.Second)) {
		t.Fatal("window did not reset")
	}
	for _, input := range []string{"//evil.test", "/\\evil.test", "/%2f%2fevil.test", "/\r\nevil"} {
		if safeNext(base64.RawURLEncoding.EncodeToString([]byte(input))) != "/panel/" {
			t.Fatal(input)
		}
	}
	for _, input := range []string{"-100", "1.5", "NaN", "999999999999999999999999999", "100USD"} {
		if _, err := parseWholeNumber(input); err == nil {
			t.Fatal(input)
		}
	}
	if n, err := parseWholeNumber("۱۲۳٬۰۰۰"); err != nil || n != 123000 {
		t.Fatal(n, err)
	}
	if safeCSVCell(" =1+1") != "' =1+1" {
		t.Fatal("formula injection")
	}
}

func TestUploadDirectoriesAndArbitraryFilesAreNotPublic(t *testing.T) {
	s, _, v := fixture(t)
	dir := filepath.Join(s.cfg.UploadDir, strconv.FormatInt(v.ID, 10))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "private.txt"), []byte("not an image"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/uploads/", "/uploads/" + strconv.FormatInt(v.ID, 10) + "/", "/uploads/" + strconv.FormatInt(v.ID, 10) + "/private.txt"} {
		if w := request(s, "GET", path, nil, nil); w.Code != http.StatusNotFound {
			t.Fatalf("%s: %d", path, w.Code)
		}
	}
}
