package server

import (
	"net/http"
	"strconv"
	"strings"

	"ResturantMenue/internal/models"
	"ResturantMenue/internal/theme"
)

// sectionRow is one draggable entry in the page designer.
type sectionRow struct {
	Key     string
	Label   string
	Icon    string
	IsBlock bool
	Block   *models.Block
	Hidden  bool
}

type designPage struct {
	Sections   []sectionRow
	Blocks     []*models.Block
	Kinds      []models.BlockKind
	BlockCount int
}

// handleDesignPage renders the drag-and-drop page designer.
func (s *Server) handleDesignPage(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	blocks, err := s.store.Blocks.List(v.Restaurant.ID, false)
	if err != nil {
		logf("design: list blocks %d: %v", v.Restaurant.ID, err)
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}

	ids := make([]int64, 0, len(blocks))
	byID := make(map[int64]*models.Block, len(blocks))
	for _, b := range blocks {
		ids = append(ids, b.ID)
		byID[b.ID] = b
	}

	order := theme.ReconcileSections(v.Theme.Layout.Sections, ids)

	rows := make([]sectionRow, 0, len(order))
	for _, key := range order {
		if id, ok := theme.IsBlockSection(key); ok {
			b, exists := byID[id]
			if !exists {
				continue
			}
			label := b.Title(v.Locale.Code)
			if strings.TrimSpace(label) == "" {
				label = v.T("blk.kind." + string(b.Kind))
			}
			rows = append(rows, sectionRow{
				Key: key, Label: label, Icon: blockIcon(b.Kind),
				IsBlock: true, Block: b, Hidden: !b.IsVisible,
			})
			continue
		}
		rows = append(rows, sectionRow{
			Key: key, Label: v.T("sect." + key), Icon: sectionIcon(key),
		})
	}

	v.Title = v.T("design.title")
	v.Page = &designPage{
		Sections:   rows,
		Blocks:     blocks,
		Kinds:      models.BlockKinds,
		BlockCount: len(blocks),
	}
	s.render.Render(w, http.StatusOK, "panel", "design", v)
}

func sectionIcon(key string) string {
	switch key {
	case theme.SectionHeader:
		return "🖼"
	case theme.SectionSearch:
		return "⌕"
	case theme.SectionNav:
		return "☰"
	case theme.SectionMenu:
		return "🍽"
	case theme.SectionHours:
		return "🕐"
	case theme.SectionFooter:
		return "⌄"
	}
	return "▫"
}

func blockIcon(kind models.BlockKind) string {
	switch kind {
	case models.BlockHeading:
		return "🅷"
	case models.BlockText:
		return "¶"
	case models.BlockImage:
		return "🖼"
	case models.BlockBanner:
		return "🎏"
	case models.BlockGallery:
		return "🖼"
	case models.BlockQuote:
		return "❝"
	case models.BlockButton:
		return "🔘"
	case models.BlockDivider:
		return "―"
	case models.BlockSpacer:
		return "␣"
	}
	return "▫"
}

// handleBlockSave creates or updates a custom content block.
func (s *Server) handleBlockSave(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := parseForm(r); err != nil {
		s.blockRedirect(w, r, "error", v.T("err.generic"))
		return
	}

	kind := formStr(r, "kind")
	if !models.ValidBlockKind(kind) {
		s.blockRedirect(w, r, "error", v.T("err.not_found"))
		return
	}

	b := &models.Block{
		RestaurantID: v.Restaurant.ID,
		Kind:         models.BlockKind(kind),
		TitleFa:      formStr(r, "title_fa"),
		TitleEn:      formStr(r, "title_en"),
		BodyFa:       formStr(r, "body_fa"),
		BodyEn:       formStr(r, "body_en"),
		ImagePath:    s.sanitizeUploadPath(v.Restaurant.ID, formStr(r, "image_path")),
		LinkURL:      formStr(r, "link_url"),
		LinkLabelFa:  formStr(r, "link_label_fa"),
		LinkLabelEn:  formStr(r, "link_label_en"),
		IsVisible:    formBool(r, "is_visible"),
		Style:        parseBlockStyle(r, models.BlockKind(kind)),
	}

	// Gallery images arrive as a comma-separated list of upload paths.
	for _, p := range strings.Split(formStr(r, "images"), ",") {
		if clean := s.sanitizeUploadPath(v.Restaurant.ID, strings.TrimSpace(p)); clean != "" {
			b.Images = append(b.Images, clean)
		}
	}

	id, _ := strconv.ParseInt(formStr(r, "id"), 10, 64)
	if id > 0 {
		b.ID = id
		if err := s.store.Blocks.Update(b); err != nil {
			logf("block: update %d: %v", id, err)
			s.blockRedirect(w, r, "error", v.T("err.generic"))
			return
		}
		s.blockRedirect(w, r, "success", v.T("msg.saved"))
		return
	}

	if err := s.store.Blocks.Create(b); err != nil {
		logf("block: create for venue %d: %v", v.Restaurant.ID, err)
		s.blockRedirect(w, r, "error", v.T("err.generic"))
		return
	}

	// Place the new block just before the footer so it is visible immediately.
	s.appendBlockToSections(v.Restaurant.ID, b.ID)
	s.blockRedirect(w, r, "success", v.T("msg.created"))
}

func (s *Server) handleBlockDelete(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := s.store.Blocks.Delete(pathInt(r, "id"), v.Restaurant.ID); err != nil {
		logf("block: delete: %v", err)
		s.blockRedirect(w, r, "error", v.T("err.generic"))
		return
	}
	s.blockRedirect(w, r, "success", v.T("msg.deleted"))
}

// handleSectionOrder persists the drag-and-drop page order.
func (s *Server) handleSectionOrder(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no venue"})
		return
	}
	if err := parseForm(r); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad form"})
		return
	}

	order := strings.Split(formStr(r, "order"), ",")
	if len(order) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty order"})
		return
	}

	preset, cfg, err := s.store.Themes.Get(v.Restaurant.ID)
	if err != nil {
		logf("section order: load theme %d: %v", v.Restaurant.ID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": v.T("err.generic")})
		return
	}

	blocks, err := s.store.Blocks.List(v.Restaurant.ID, false)
	if err != nil {
		logf("section order: list blocks %d: %v", v.Restaurant.ID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": v.T("err.generic")})
		return
	}
	ids := make([]int64, 0, len(blocks))
	for _, b := range blocks {
		ids = append(ids, b.ID)
	}

	cfg.Layout.Sections = theme.ReconcileSections(theme.SanitizeSections(order), ids)

	if err := s.store.Themes.Save(v.Restaurant.ID, preset, cfg); err != nil {
		logf("section order: save theme %d: %v", v.Restaurant.ID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": v.T("err.generic")})
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "sections": cfg.Layout.Sections})
		return
	}
	s.blockRedirect(w, r, "success", v.T("msg.saved"))
}

// appendBlockToSections adds a freshly created block to the saved page order.
func (s *Server) appendBlockToSections(restaurantID, blockID int64) {
	preset, cfg, err := s.store.Themes.Get(restaurantID)
	if err != nil {
		logf("block: load theme %d: %v", restaurantID, err)
		return
	}

	blocks, err := s.store.Blocks.List(restaurantID, false)
	if err != nil {
		logf("block: list for section merge %d: %v", restaurantID, err)
		return
	}
	ids := make([]int64, 0, len(blocks))
	for _, b := range blocks {
		ids = append(ids, b.ID)
	}

	// ReconcileSections inserts anything unknown before the footer, which is
	// exactly where a new block belongs.
	cfg.Layout.Sections = theme.ReconcileSections(cfg.Layout.Sections, ids)
	if err := s.store.Themes.Save(restaurantID, preset, cfg); err != nil {
		logf("block: save section order %d: %v", restaurantID, err)
	}
}

// parseBlockStyle reads the style fields, layering them over the kind's defaults.
func parseBlockStyle(r *http.Request, kind models.BlockKind) models.BlockStyle {
	s := models.DefaultBlockStyle(kind)

	if v := formStr(r, "s_size"); v != "" {
		s.Size = v
	}
	if v := formStr(r, "s_align"); v != "" {
		s.Align = v
	}
	if v := formStr(r, "s_color"); v != "" {
		s.Color = v
	}
	if v := formStr(r, "s_customCol"); v != "" {
		s.CustomCol = v
	}
	if v := formStr(r, "s_bg"); v != "" {
		s.Bg = v
	}
	if v := formStr(r, "s_width"); v != "" {
		s.Width = v
	}
	if v := formStr(r, "s_padding"); v != "" {
		s.Padding = v
	}
	if v := formStr(r, "s_effect"); v != "" {
		s.Effect = v
	}
	if v := formStr(r, "s_ratio"); v != "" {
		s.Ratio = v
	}
	if v := formStr(r, "s_divider"); v != "" {
		s.Divider = v
	}

	if n, err := strconv.Atoi(normaliseDigits(formStr(r, "s_weight"))); err == nil {
		s.Weight = n
	}
	if n, err := strconv.Atoi(normaliseDigits(formStr(r, "s_radius"))); err == nil {
		s.Radius = n
	}
	if n, err := strconv.Atoi(normaliseDigits(formStr(r, "s_shadow"))); err == nil {
		s.Shadow = n
	}
	if n, err := strconv.Atoi(normaliseDigits(formStr(r, "s_height"))); err == nil {
		s.Height = n
	}
	if n, err := strconv.Atoi(normaliseDigits(formStr(r, "s_columns"))); err == nil {
		s.Columns = n
	}
	if f, err := strconv.ParseFloat(normaliseDigits(formStr(r, "s_overlay")), 64); err == nil {
		s.Overlay = f
	}

	// The style form always submits these, so absence means unchecked.
	s.Italic = formBool(r, "s_italic")
	s.Uppercase = formBool(r, "s_uppercase")
	s.Border = formBool(r, "s_border")

	s.Normalize(kind)
	return s
}

func (s *Server) blockRedirect(w http.ResponseWriter, r *http.Request, kind, msg string) {
	s.setFlash(w, kind, msg)
	http.Redirect(w, r, "/panel/design", http.StatusSeeOther)
}
