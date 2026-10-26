package server

import (
	"net/http"
	"strconv"
	"strings"

	"ResturantMenue/internal/menudata"
	"ResturantMenue/internal/menutemplate"
	"ResturantMenue/internal/models"
)

// knownTags is the whitelist of item badges the editor offers.
var knownTags = []string{"spicy", "vegetarian", "vegan", "new", "popular", "chef", "glutenfree", "cold", "hot"}

type menuPage struct {
	Categories []*models.Category
	Tags       []string
	Templates  []menutemplate.Template
	MenuURL    string
}

func (s *Server) handleMenuPage(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	cats, err := s.store.FullMenu(v.Restaurant.ID, false)
	if err != nil {
		logf("menu: load for venue %d: %v", v.Restaurant.ID, err)
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}

	v.Title = v.T("menu.title")
	v.Page = &menuPage{
		Categories: cats,
		Tags:       knownTags,
		Templates:  menutemplate.All(),
		MenuURL:    "/m/" + v.Restaurant.Slug,
	}
	s.render.Render(w, http.StatusOK, "panel", "menu", v)
}

// --- categories ------------------------------------------------------------

// handleCategorySave creates or updates a category depending on whether an id
// was submitted.
func (s *Server) handleCategorySave(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := parseForm(r); err != nil {
		s.menuRedirect(w, r, "error", v.T("err.generic"))
		return
	}

	nameFa := formStr(r, "name_fa")
	nameEn := formStr(r, "name_en")
	if nameFa == "" && nameEn == "" {
		s.menuRedirect(w, r, "error", v.T("err.required"))
		return
	}

	cat := &models.Category{
		RestaurantID:  v.Restaurant.ID,
		NameFa:        nameFa,
		NameEn:        nameEn,
		DescriptionFa: formStr(r, "description_fa"),
		DescriptionEn: formStr(r, "description_en"),
		Icon:          firstRunes(formStr(r, "icon"), 4),
		ImagePath:     s.sanitizeUploadPath(v.Restaurant.ID, formStr(r, "image_path")),
		IsVisible:     formBool(r, "is_visible"),
	}

	id, _ := strconv.ParseInt(formStr(r, "id"), 10, 64)
	if id > 0 {
		cat.ID = id
		if err := s.store.Categories.Update(cat); err != nil {
			logf("category: update %d: %v", id, err)
			s.menuRedirect(w, r, "error", v.T("err.generic"))
			return
		}
		s.menuRedirect(w, r, "success", v.T("msg.saved"))
		return
	}

	if err := s.store.Categories.Create(cat); err != nil {
		logf("category: create for venue %d: %v", v.Restaurant.ID, err)
		s.menuRedirect(w, r, "error", v.T("err.generic"))
		return
	}
	s.menuRedirect(w, r, "success", v.T("msg.created"))
}

func (s *Server) handleCategoryDelete(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := s.store.Categories.Delete(pathInt(r, "id"), v.Restaurant.ID); err != nil {
		logf("category: delete: %v", err)
		s.menuRedirect(w, r, "error", v.T("err.generic"))
		return
	}
	s.menuRedirect(w, r, "success", v.T("msg.deleted"))
}

func (s *Server) handleCategoryReorder(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no venue"})
		return
	}
	ids := parseIDList(r.FormValue("order"))
	if len(ids) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty order"})
		return
	}
	if err := s.store.Categories.Reorder(v.Restaurant.ID, ids); err != nil {
		logf("category: reorder: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": v.T("err.generic")})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// --- items -----------------------------------------------------------------

func (s *Server) handleItemSave(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := parseForm(r); err != nil {
		s.menuRedirect(w, r, "error", v.T("err.generic"))
		return
	}

	categoryID, _ := strconv.ParseInt(formStr(r, "category_id"), 10, 64)
	if categoryID <= 0 {
		s.menuRedirect(w, r, "error", v.T("err.required"))
		return
	}
	// Confirm the category belongs to this venue before attaching an item.
	if _, err := s.store.Categories.ByID(categoryID, v.Restaurant.ID); err != nil {
		s.menuRedirect(w, r, "error", v.T("err.forbidden"))
		return
	}

	nameFa := formStr(r, "name_fa")
	nameEn := formStr(r, "name_en")
	if nameFa == "" && nameEn == "" {
		s.menuRedirect(w, r, "error", v.T("err.required"))
		return
	}

	item := &models.Item{
		RestaurantID:  v.Restaurant.ID,
		CategoryID:    categoryID,
		NameFa:        nameFa,
		NameEn:        nameEn,
		DescriptionFa: formStr(r, "description_fa"),
		DescriptionEn: formStr(r, "description_en"),
		Price:         formInt(r, "price"),
		PricePending:  formBool(r, "price_pending"),
		DiscountPrice: formInt(r, "discount_price"),
		ImagePath:     s.sanitizeUploadPath(v.Restaurant.ID, formStr(r, "image_path")),
		Tags:          selectedTags(r),
		IsAvailable:   formBool(r, "is_available"),
		IsFeatured:    formBool(r, "is_featured"),
		Calories:      int(formInt(r, "calories")),
		PrepMinutes:   int(formInt(r, "prep_minutes")),
	}
	for _, key := range []string{"price", "discount_price", "calories", "prep_minutes"} {
		if _, err := parseWholeNumber(formStr(r, key)); err != nil {
			s.menuRedirect(w, r, "error", v.T("err.price"))
			return
		}
	}
	if item.PricePending {
		item.Price = 0
		item.DiscountPrice = 0
	}
	if item.Price > menudata.MaxPrice || item.DiscountPrice > menudata.MaxPrice ||
		(item.DiscountPrice > 0 && item.DiscountPrice >= item.Price) || item.Calories > 100000 || item.PrepMinutes > 10080 {
		s.menuRedirect(w, r, "error", v.T("err.price"))
		return
	}

	id, _ := strconv.ParseInt(formStr(r, "id"), 10, 64)
	if id > 0 {
		item.ID = id
		if err := s.store.Items.Update(item); err != nil {
			logf("item: update %d: %v", id, err)
			s.menuRedirect(w, r, "error", v.T("err.generic"))
			return
		}
		s.menuRedirect(w, r, "success", v.T("msg.saved"))
		return
	}

	if err := s.store.Items.Create(item); err != nil {
		logf("item: create for venue %d: %v", v.Restaurant.ID, err)
		s.menuRedirect(w, r, "error", v.T("err.generic"))
		return
	}
	s.menuRedirect(w, r, "success", v.T("msg.created"))
}

func (s *Server) handleItemDelete(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := s.store.Items.Delete(pathInt(r, "id"), v.Restaurant.ID); err != nil {
		logf("item: delete: %v", err)
		s.menuRedirect(w, r, "error", v.T("err.generic"))
		return
	}
	s.menuRedirect(w, r, "success", v.T("msg.deleted"))
}

func (s *Server) handleItemToggle(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no venue"})
		return
	}
	if err := s.store.Items.ToggleAvailable(pathInt(r, "id"), v.Restaurant.ID); err != nil {
		logf("item: toggle: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": v.T("err.generic")})
		return
	}
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	s.menuRedirect(w, r, "success", v.T("msg.saved"))
}

func (s *Server) handleItemReorder(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no venue"})
		return
	}
	ids := parseIDList(r.FormValue("order"))
	if len(ids) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty order"})
		return
	}
	if err := s.store.Items.Reorder(v.Restaurant.ID, ids); err != nil {
		logf("item: reorder: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": v.T("err.generic")})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleApplyMenuTemplate appends a starter category set to the current menu.
func (s *Server) handleApplyMenuTemplate(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	tpl, ok := menutemplate.ByKey(formStr(r, "template"))
	if !ok {
		s.menuRedirect(w, r, "error", v.T("err.not_found"))
		return
	}

	// Skip categories the venue already has, so applying twice is harmless.
	existing, err := s.store.Categories.List(v.Restaurant.ID)
	if err != nil {
		logf("template: list categories: %v", err)
		s.menuRedirect(w, r, "error", v.T("err.generic"))
		return
	}
	have := make(map[string]bool, len(existing))
	for _, c := range existing {
		have[strings.ToLower(c.NameFa)] = true
		have[strings.ToLower(c.NameEn)] = true
	}

	added := 0
	for _, c := range tpl.Categories {
		if have[strings.ToLower(c.NameFa)] || have[strings.ToLower(c.NameEn)] {
			continue
		}
		cat := &models.Category{
			RestaurantID: v.Restaurant.ID,
			NameFa:       c.NameFa,
			NameEn:       c.NameEn,
			Icon:         c.Icon,
			IsVisible:    true,
		}
		if err := s.store.Categories.Create(cat); err != nil {
			logf("template: create category %q: %v", c.NameEn, err)
			continue
		}
		added++
	}

	s.menuRedirect(w, r, "success", v.T("menu.template_applied"))
	_ = added
}

// --- helpers ---------------------------------------------------------------

func (s *Server) menuRedirect(w http.ResponseWriter, r *http.Request, kind, msg string) {
	s.setFlash(w, kind, msg)
	http.Redirect(w, r, "/panel/menu", http.StatusSeeOther)
}

// selectedTags keeps only whitelisted tag values.
func selectedTags(r *http.Request) []string {
	submitted := r.Form["tags"]
	if len(submitted) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(knownTags))
	for _, t := range knownTags {
		allowed[t] = true
	}

	var out []string
	seen := map[string]bool{}
	for _, t := range submitted {
		t = strings.ToLower(strings.TrimSpace(t))
		if allowed[t] && !seen[t] {
			out = append(out, t)
			seen[t] = true
		}
	}
	return out
}

// parseIDList turns "3,1,2" into []int64{3,1,2}.
func parseIDList(raw string) []int64 {
	var out []int64
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if n, err := strconv.ParseInt(part, 10, 64); err == nil && n > 0 {
			out = append(out, n)
		}
	}
	return out
}

// firstRunes truncates to n runes, keeping multi-byte emoji intact.
func firstRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
