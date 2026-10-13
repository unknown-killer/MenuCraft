package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ResturantMenue/internal/auth"
	"ResturantMenue/internal/i18n"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
)

type dashboardPage struct {
	Categories     []*models.Category
	TotalItems     int
	AvailableItems int
	PhotoCount     int
	MenuURL        string
	IsEmpty        bool
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		// A signed-in account should always have a venue; recreate one rather
		// than dead-ending the user.
		if err := s.provisionRestaurant(v.User, v.User.DisplayName, string(models.BusinessRestaurant), v.Locale.Code); err != nil {
			logf("dashboard: reprovision for user %d: %v", v.User.ID, err)
			s.renderError(w, r, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	cats, err := s.store.FullMenu(v.Restaurant.ID, false)
	if err != nil {
		logf("dashboard: load menu %d: %v", v.Restaurant.ID, err)
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}

	page := &dashboardPage{
		Categories: cats,
		MenuURL:    "/m/" + v.Restaurant.Slug,
	}
	for _, c := range cats {
		page.TotalItems += len(c.Items)
		for _, it := range c.Items {
			if it.IsAvailable {
				page.AvailableItems++
			}
		}
	}
	page.IsEmpty = page.TotalItems == 0

	if media, err := s.store.Media.List(v.Restaurant.ID, 500); err == nil {
		page.PhotoCount = len(media)
	}

	v.Title = v.T("nav.dashboard")
	v.Page = page
	s.render.Render(w, http.StatusOK, "panel", "dashboard", v)
}

// --- settings --------------------------------------------------------------

type settingsPage struct {
	Hours         []*models.WorkingHour
	BusinessTypes []string
	Currencies    []string
	MenuURL       string
}

func (s *Server) handleSettingsPage(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	hours, err := s.store.Restaurants.WorkingHours(v.Restaurant.ID)
	if err != nil {
		logf("settings: hours for venue %d: %v", v.Restaurant.ID, err)
	}
	// Always show seven rows, even before anything has been saved.
	byDay := map[int]*models.WorkingHour{}
	for _, h := range hours {
		byDay[h.DayOfWeek] = h
	}
	full := make([]*models.WorkingHour, 7)
	for d := 0; d < 7; d++ {
		if h, ok := byDay[d]; ok {
			full[d] = h
		} else {
			full[d] = &models.WorkingHour{RestaurantID: v.Restaurant.ID, DayOfWeek: d, OpensAt: "09:00", ClosesAt: "23:00"}
		}
	}

	v.Title = v.T("set.title")
	v.Page = &settingsPage{
		Hours:         full,
		BusinessTypes: businessTypes,
		Currencies:    i18n.Currencies,
		MenuURL:       "/m/" + v.Restaurant.Slug,
	}
	s.render.Render(w, http.StatusOK, "panel", "settings", v)
}

func (s *Server) handleSettingsSave(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := parseForm(r); err != nil {
		s.setFlash(w, "error", v.T("err.generic"))
		http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
		return
	}

	venue := v.Restaurant

	// The slug is the public URL, so validate it before anything else.
	if slug := strings.ToLower(formStr(r, "slug")); slug != "" && slug != venue.Slug {
		if !validSlug(slug) {
			s.setFlash(w, "error", v.T("set.slug_hint"))
			http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
			return
		}
		taken, err := s.store.Restaurants.SlugTaken(slug, venue.ID)
		if err != nil {
			logf("settings: slug check: %v", err)
			s.setFlash(w, "error", v.T("err.generic"))
			http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
			return
		}
		if taken {
			s.setFlash(w, "error", v.T("err.slug_taken"))
			http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
			return
		}
		venue.Slug = slug
	}

	venue.NameFa = formStr(r, "name_fa")
	venue.NameEn = formStr(r, "name_en")
	venue.TaglineFa = formStr(r, "tagline_fa")
	venue.TaglineEn = formStr(r, "tagline_en")
	venue.DescriptionFa = formStr(r, "description_fa")
	venue.DescriptionEn = formStr(r, "description_en")
	venue.Phone = formStr(r, "phone")
	venue.AddressFa = formStr(r, "address_fa")
	venue.AddressEn = formStr(r, "address_en")
	venue.Instagram = formStr(r, "instagram")
	venue.Telegram = formStr(r, "telegram")
	venue.WhatsApp = formStr(r, "whatsapp")
	venue.Website = formStr(r, "website")
	venue.LogoPath = s.sanitizeUploadPath(venue.ID, formStr(r, "logo_path"))
	venue.CoverPath = s.sanitizeUploadPath(venue.ID, formStr(r, "cover_path"))
	venue.IsPublished = formBool(r, "is_published")

	if bt := formStr(r, "business_type"); validBusinessType(bt) {
		venue.BusinessType = bt
	}
	if cur := strings.ToUpper(formStr(r, "currency")); validCurrency(cur) {
		venue.Currency = cur
	}

	// Working hours arrive as seven parallel rows.
	hours := make([]*models.WorkingHour, 0, 7)
	for d := 0; d < 7; d++ {
		suffix := strconv.Itoa(d)
		hours = append(hours, &models.WorkingHour{
			RestaurantID: venue.ID,
			DayOfWeek:    d,
			OpensAt:      sanitizeTime(formStr(r, "opens_"+suffix)),
			ClosesAt:     sanitizeTime(formStr(r, "closes_"+suffix)),
			IsClosed:     formBool(r, "closed_"+suffix),
		})
	}
	if err := s.store.Transaction(func(tx *repository.Store) error {
		if err := tx.Restaurants.Update(venue); err != nil {
			return err
		}
		return tx.Restaurants.SaveWorkingHours(venue.ID, hours)
	}); err != nil {
		key := "err.generic"
		if errors.Is(err, repository.ErrDuplicate) {
			key = "err.slug_taken"
		}
		logf("settings: save venue %d: %v", venue.ID, err)
		s.setFlash(w, "error", v.T(key))
		http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
		return
	}

	s.setFlash(w, "success", v.T("msg.saved"))
	http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
}

func (s *Server) handlePasswordChange(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if err := parseForm(r); err != nil {
		s.setFlash(w, "error", v.T("err.generic"))
		http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
		return
	}

	current := r.FormValue("old_password")
	next := r.FormValue("new_password")

	if err := auth.VerifyPassword(current, v.User.PasswordHash); err != nil {
		s.setFlash(w, "error", v.T("err.bad_login"))
		http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
		return
	}
	if problem := auth.PasswordProblem(next); problem != "" {
		s.setFlash(w, "error", v.T(problem))
		http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
		return
	}

	hash, err := auth.HashPassword(next)
	if err != nil {
		logf("password: hash for user %d: %v", v.User.ID, err)
		s.setFlash(w, "error", v.T("err.generic"))
		http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
		return
	}
	if err := s.store.Users.UpdatePassword(v.User.ID, hash); err != nil {
		logf("password: update for user %d: %v", v.User.ID, err)
		s.setFlash(w, "error", v.T("err.generic"))
		http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
		return
	}

	// UpdatePassword bumped session_version, so the current cookie is now
	// stale: issue a fresh one instead of bouncing the user to the login page.
	s.sessions.Issue(w, v.User.ID, v.User.SessionVersion+1)
	s.setFlash(w, "success", v.T("set.pass_changed"))
	http.Redirect(w, r, "/panel/settings", http.StatusSeeOther)
}

// --- shared ----------------------------------------------------------------

// renderError shows a styled error page.
func (s *Server) renderError(w http.ResponseWriter, r *http.Request, status int) {
	v := s.newView(w, r)
	switch status {
	case http.StatusNotFound:
		v.Title = v.T("err.404_title")
	default:
		v.Title = v.T("err.500_title")
	}
	v.Page = map[string]any{"Status": status}
	s.render.Render(w, status, "auth", "error", v)
}

func validCurrency(code string) bool {
	for _, c := range i18n.Currencies {
		if c == code {
			return true
		}
	}
	return false
}

// sanitizeTime accepts only HH:MM, so nothing odd reaches the template.
func sanitizeTime(v string) string {
	v = normaliseDigits(strings.TrimSpace(v))
	if len(v) != 5 || v[2] != ':' {
		return ""
	}
	for i, ch := range v {
		if i == 2 {
			continue
		}
		if ch < '0' || ch > '9' {
			return ""
		}
	}
	return v
}
