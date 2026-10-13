package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"ResturantMenue/internal/auth"
	"ResturantMenue/internal/i18n"
	"ResturantMenue/internal/menutemplate"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
	"ResturantMenue/internal/theme"
)

// authPage is the payload for the login and register templates.
type authPage struct {
	Mode          string // "login" | "register"
	Error         string
	Next          string
	Values        map[string]string
	BusinessTypes []string
}

var businessTypes = []string{
	string(models.BusinessRestaurant),
	string(models.BusinessFastFood),
	string(models.BusinessCafe),
	string(models.BusinessJigaraki),
	string(models.BusinessTeahouse),
	string(models.BusinessBakery),
	string(models.BusinessIceCream),
}

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	if _, err := s.sessions.Read(r); err == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	v := s.newView(w, r)
	v.Title = v.T("auth.login.title")
	v.Page = &authPage{Mode: "login", Next: r.URL.Query().Get("next"), Values: map[string]string{}}
	s.render.Render(w, http.StatusOK, "auth", "login", v)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	v.Title = v.T("auth.login.title")

	if err := parseForm(r); err != nil {
		s.renderAuthError(w, v, "login", v.T("err.generic"), nil)
		return
	}
	if !s.sessions.CheckCSRF(r) {
		s.renderAuthError(w, v, "login", v.T("err.csrf"), nil)
		return
	}

	login := formStr(r, "login")
	password := r.FormValue("password")
	values := map[string]string{"login": login}

	user, err := s.store.Users.ByLogin(login)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			logf("login: lookup %q: %v", truncate(login, 40), err)
		}
		// Same message and roughly the same cost whether the account exists
		// or not, so the form does not reveal which usernames are registered.
		_ = auth.VerifyPassword(password, dummyHash)
		s.renderAuthError(w, v, "login", v.T("err.bad_login"), values)
		return
	}

	if err := auth.VerifyPassword(password, user.PasswordHash); err != nil {
		s.renderAuthError(w, v, "login", v.T("err.bad_login"), values)
		return
	}

	s.sessions.Issue(w, user.ID, user.SessionVersion)
	if user.Locale != "" {
		i18n.SetCookie(w, user.Locale, s.cfg.IsProduction())
	}
	http.Redirect(w, r, safeNext(r.FormValue("next")), http.StatusSeeOther)
}

// dummyHash is a real PBKDF2 hash of an unguessable value. Verifying against
// it on the "user not found" path keeps login timing roughly uniform.
var dummyHash = mustHash("timing-equaliser-not-a-real-password-9f3c")

func mustHash(pw string) string {
	h, err := auth.HashPassword(pw)
	if err != nil {
		panic("auth: cannot build dummy hash: " + err.Error())
	}
	return h
}

func (s *Server) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	if _, err := s.sessions.Read(r); err == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	v := s.newView(w, r)
	v.Title = v.T("auth.register.title")
	v.Page = &authPage{Mode: "register", Values: map[string]string{}, BusinessTypes: businessTypes}
	s.render.Render(w, http.StatusOK, "auth", "register", v)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	v.Title = v.T("auth.register.title")

	if err := parseForm(r); err != nil {
		s.renderAuthError(w, v, "register", v.T("err.generic"), nil)
		return
	}
	if !s.sessions.CheckCSRF(r) {
		s.renderAuthError(w, v, "register", v.T("err.csrf"), nil)
		return
	}

	username := strings.ToLower(formStr(r, "username"))
	email := strings.ToLower(formStr(r, "email"))
	displayName := formStr(r, "display_name")
	venueName := formStr(r, "restaurant_name")
	bizType := formStr(r, "business_type")
	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")

	values := map[string]string{
		"username": username, "email": email,
		"display_name": displayName, "restaurant_name": venueName,
		"business_type": bizType,
	}
	fail := func(msg string) { s.renderAuthError(w, v, "register", msg, values) }

	switch {
	case !validUsername(username):
		fail(v.T("err.username"))
		return
	case !validEmail(email):
		fail(v.T("err.email"))
		return
	case venueName == "":
		fail(v.T("err.required"))
		return
	case password != confirm:
		fail(v.T("err.password_match"))
		return
	}
	if problem := auth.PasswordProblem(password); problem != "" {
		fail(v.T(problem))
		return
	}
	if !validBusinessType(bizType) {
		bizType = string(models.BusinessRestaurant)
	}

	taken, err := s.store.Users.Exists(username, email)
	if err != nil {
		logf("register: exists check: %v", err)
		fail(v.T("err.generic"))
		return
	}
	if taken {
		fail(v.T("err.taken"))
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		logf("register: hash: %v", err)
		fail(v.T("err.generic"))
		return
	}

	if displayName == "" {
		displayName = username
	}
	user := &models.User{
		Username: username, Email: email, PasswordHash: hash,
		DisplayName: displayName, Locale: v.Locale.Code,
	}
	if err := s.store.Users.Create(user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			fail(v.T("err.taken"))
			return
		}
		logf("register: create user: %v", err)
		fail(v.T("err.generic"))
		return
	}

	if err := s.provisionRestaurant(user, venueName, bizType, v.Locale.Code); err != nil {
		logf("register: provision venue for user %d: %v", user.ID, err)
		fail(v.T("err.generic"))
		return
	}

	s.sessions.Issue(w, user.ID, user.SessionVersion)
	i18n.SetCookie(w, user.Locale, s.cfg.IsProduction())
	s.setFlash(w, "success", i18n.T(v.Locale.Code, "msg.created"))
	http.Redirect(w, r, "/panel/", http.StatusSeeOther)
}

// provisionRestaurant creates the venue, a unique slug, a starter theme and
// the category skeleton for the chosen business type.
func (s *Server) provisionRestaurant(user *models.User, name, bizType, locale string) error {
	slug := Slugify(name)
	if slug == "" {
		slug = Slugify(user.Username)
	}
	if slug == "" {
		slug = "menu"
	}

	// Append a counter until the slug is free.
	base := slug
	for i := 2; ; i++ {
		taken, err := s.store.Restaurants.SlugTaken(slug, 0)
		if err != nil {
			return fmt.Errorf("slug check: %w", err)
		}
		if !taken {
			break
		}
		slug = fmt.Sprintf("%s-%d", base, i)
		if i > 200 {
			return errors.New("could not find a free slug")
		}
	}

	venue := &models.Restaurant{
		UserID:       user.ID,
		Slug:         slug,
		BusinessType: bizType,
		Currency:     "IRT",
		IsPublished:  true,
	}
	if locale == "en" {
		venue.NameEn = name
	} else {
		venue.NameFa = name
	}
	if err := s.store.Restaurants.Create(venue); err != nil {
		return fmt.Errorf("create venue: %w", err)
	}

	// Give every new venue the default theme so the editor starts populated.
	if err := s.store.Themes.Save(venue.ID, "midnight", theme.Default()); err != nil {
		return fmt.Errorf("seed theme: %w", err)
	}

	tpl := menutemplate.ForBusiness(bizType)
	for _, c := range tpl.Categories {
		cat := &models.Category{
			RestaurantID: venue.ID,
			NameFa:       c.NameFa,
			NameEn:       c.NameEn,
			Icon:         c.Icon,
			IsVisible:    true,
		}
		if err := s.store.Categories.Create(cat); err != nil {
			return fmt.Errorf("seed category %q: %w", c.NameEn, err)
		}
	}
	return nil
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Logout is a POST with a CSRF token, so a stray link cannot sign a user out.
	if !s.sessions.CheckCSRF(r) {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	s.sessions.Clear(w)
	s.setFlash(w, "info", i18n.T(i18n.Detect(r, s.cfg.DefaultLocale).Code, "auth.logout_done"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleSwitchLocale(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	i18n.SetCookie(w, code, s.cfg.IsProduction())

	// Persist the choice on the account when someone is signed in.
	if sess, err := s.sessions.Read(r); err == nil {
		if u, err := s.store.Users.ByID(sess.UserID); err == nil {
			u.Locale = i18n.Get(code).Code
			if err := s.store.Users.UpdateProfile(u); err != nil {
				logf("locale: persist for user %d: %v", u.ID, err)
			}
		}
	}

	back := r.Referer()
	if back == "" || !strings.HasPrefix(back, "/") {
		if u, err := parseSameSiteReferer(back); err == nil {
			back = u
		} else {
			back = "/"
		}
	}
	http.Redirect(w, r, back, http.StatusSeeOther)
}

// parseSameSiteReferer reduces a Referer header to its path, discarding any
// cross-site origin so the redirect stays local.
func parseSameSiteReferer(ref string) (string, error) {
	if ref == "" {
		return "", errors.New("empty referer")
	}
	i := strings.Index(ref, "://")
	if i < 0 {
		return "", errors.New("not absolute")
	}
	rest := ref[i+3:]
	slash := strings.Index(rest, "/")
	if slash < 0 {
		return "/", nil
	}
	return rest[slash:], nil
}

func (s *Server) renderAuthError(w http.ResponseWriter, v *View, mode, msg string, values map[string]string) {
	if values == nil {
		values = map[string]string{}
	}
	v.Page = &authPage{Mode: mode, Error: msg, Values: values, BusinessTypes: businessTypes}
	s.render.Render(w, http.StatusUnprocessableEntity, "auth", mode, v)
}

func validBusinessType(t string) bool {
	for _, b := range businessTypes {
		if b == t {
			return true
		}
	}
	return false
}
