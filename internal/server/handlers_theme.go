package server

import (
	"net/http"
	"strconv"

	"ResturantMenue/internal/theme"
)

type themePage struct {
	Preset string
	Config theme.Config
	// Presets is ordered so the looks suited to this venue type come first.
	Presets []theme.Preset
}

func (s *Server) handleThemePage(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	v.Title = v.T("theme.title")
	v.Page = &themePage{
		Preset:  v.ThemePreset,
		Config:  v.Theme,
		Presets: theme.PresetsFor(v.Restaurant.BusinessType),
	}
	s.render.Render(w, http.StatusOK, "panel", "theme", v)
}

// handleThemeSave persists the submitted configuration.
func (s *Server) handleThemeSave(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Error(w, "no venue", http.StatusBadRequest)
		return
	}
	if err := parseForm(r); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	preset, cfg := parseThemeForm(r, v.Theme)
	if err := s.store.Themes.Save(v.Restaurant.ID, preset, cfg); err != nil {
		logf("theme: save for venue %d: %v", v.Restaurant.ID, err)
		if wantsJSON(r) {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": v.T("err.generic")})
			return
		}
		s.setFlash(w, "error", v.T("err.generic"))
		http.Redirect(w, r, "/panel/theme", http.StatusSeeOther)
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": v.T("theme.saved")})
		return
	}
	s.setFlash(w, "success", v.T("theme.saved"))
	http.Redirect(w, r, "/panel/theme", http.StatusSeeOther)
}

// handleThemeCSS renders the submitted configuration to CSS without saving it.
// The editor calls this on every change to drive the live preview, which keeps
// the CSS generation logic in exactly one place (Go) rather than duplicating it
// in JavaScript.
func (s *Server) handleThemeCSS(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if err := parseForm(r); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	_, cfg := parseThemeForm(r, v.Theme)
	locale := r.FormValue("preview_locale")
	if locale == "" {
		locale = v.Locale.Code
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"css":       string(cfg.Variables(locale)),
		"bodyClass": cfg.BodyClasses(),
		"display":   cfg.Display,
	})
}

func (s *Server) handleThemeReset(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := s.store.Themes.Save(v.Restaurant.ID, "midnight", theme.Default()); err != nil {
		logf("theme: reset for venue %d: %v", v.Restaurant.ID, err)
		s.setFlash(w, "error", v.T("err.generic"))
	} else {
		s.setFlash(w, "success", v.T("msg.saved"))
	}
	http.Redirect(w, r, "/panel/theme", http.StatusSeeOther)
}

// handleThemePreview renders the owner's public menu inside the editor iframe.
func (s *Server) handleThemePreview(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Error(w, "no venue", http.StatusBadRequest)
		return
	}
	s.renderPublicMenu(w, r, v.Restaurant, true)
}

// --- form parsing ----------------------------------------------------------

// parseThemeForm builds a config from the editor's flat form fields.
//
// When "preset" names a known preset and "preset_only" is set, the preset's
// configuration is returned wholesale — that is the "apply preset" action.
// Otherwise the submitted fields are layered over the current config, so a
// partial form never wipes settings it did not include.
func parseThemeForm(r *http.Request, current theme.Config) (string, theme.Config) {
	preset := r.FormValue("preset")

	if r.FormValue("preset_only") != "" && theme.IsPreset(preset) {
		return preset, theme.ApplyPreset(preset)
	}

	cfg := current

	setStr(r, "mode", &cfg.Mode)

	// Regenerating a palette is an explicit button, not something that happens
	// silently whenever a colour changes.
	switch r.FormValue("derive") {
	case "dark":
		cfg.DarkColors = theme.DeriveDark(cfg.Colors)
	case "light":
		cfg.Colors = theme.DeriveLight(cfg.DarkColors)
	}

	// dark palette
	setStr(r, "dc_primary", &cfg.DarkColors.Primary)
	setStr(r, "dc_secondary", &cfg.DarkColors.Secondary)
	setStr(r, "dc_accent", &cfg.DarkColors.Accent)
	setStr(r, "dc_background", &cfg.DarkColors.Background)
	setStr(r, "dc_surface", &cfg.DarkColors.Surface)
	setStr(r, "dc_surfaceAlt", &cfg.DarkColors.SurfaceAlt)
	setStr(r, "dc_text", &cfg.DarkColors.Text)
	setStr(r, "dc_textMuted", &cfg.DarkColors.TextMuted)
	setStr(r, "dc_border", &cfg.DarkColors.Border)
	setStr(r, "dc_price", &cfg.DarkColors.Price)
	setStr(r, "dc_badge", &cfg.DarkColors.Badge)

	// colours
	setStr(r, "c_primary", &cfg.Colors.Primary)
	setStr(r, "c_secondary", &cfg.Colors.Secondary)
	setStr(r, "c_accent", &cfg.Colors.Accent)
	setStr(r, "c_background", &cfg.Colors.Background)
	setStr(r, "c_surface", &cfg.Colors.Surface)
	setStr(r, "c_surfaceAlt", &cfg.Colors.SurfaceAlt)
	setStr(r, "c_text", &cfg.Colors.Text)
	setStr(r, "c_textMuted", &cfg.Colors.TextMuted)
	setStr(r, "c_border", &cfg.Colors.Border)
	setStr(r, "c_price", &cfg.Colors.Price)
	setStr(r, "c_badge", &cfg.Colors.Badge)

	// typography
	setStr(r, "t_fontFa", &cfg.Typography.FontFa)
	setStr(r, "t_fontEn", &cfg.Typography.FontEn)
	setInt(r, "t_baseSize", &cfg.Typography.BaseSize)
	setFloat(r, "t_headingScale", &cfg.Typography.HeadingScale)
	setInt(r, "t_headingWeight", &cfg.Typography.HeadingWeight)
	setInt(r, "t_bodyWeight", &cfg.Typography.BodyWeight)
	setFloat(r, "t_lineHeight", &cfg.Typography.LineHeight)
	setFloat(r, "t_letterSpacing", &cfg.Typography.LetterSpacing)

	// background
	setStr(r, "b_type", &cfg.Background.Type)
	setStr(r, "b_color", &cfg.Background.Color)
	setStr(r, "b_gradientFrom", &cfg.Background.GradientFrom)
	setStr(r, "b_gradientVia", &cfg.Background.GradientVia)
	setStr(r, "b_gradientTo", &cfg.Background.GradientTo)
	setInt(r, "b_gradientDeg", &cfg.Background.GradientDeg)
	setStr(r, "b_imagePath", &cfg.Background.ImagePath)
	setFloat(r, "b_imageOverlay", &cfg.Background.ImageOverlay)
	setInt(r, "b_imageBlur", &cfg.Background.ImageBlur)
	setStr(r, "b_pattern", &cfg.Background.Pattern)
	setFloat(r, "b_patternOpacity", &cfg.Background.PatternOpac)
	setStr(r, "b_attachment", &cfg.Background.Attachment)
	// An explicit clear must win over "field absent means keep".
	if r.FormValue("b_clearImage") != "" {
		cfg.Background.ImagePath = ""
	}

	// layout
	setStr(r, "l_mode", &cfg.Layout.Mode)
	setInt(r, "l_columns", &cfg.Layout.Columns)
	setInt(r, "l_gap", &cfg.Layout.Gap)
	setInt(r, "l_maxWidth", &cfg.Layout.MaxWidth)
	setStr(r, "l_density", &cfg.Layout.Density)
	setStr(r, "l_categoryNav", &cfg.Layout.CategoryNav)

	// cards
	setStr(r, "card_style", &cfg.Cards.Style)
	setInt(r, "card_radius", &cfg.Cards.Radius)
	setInt(r, "card_shadow", &cfg.Cards.Shadow)
	setInt(r, "card_borderWidth", &cfg.Cards.BorderWidth)
	setStr(r, "card_imageRatio", &cfg.Cards.ImageRatio)
	setStr(r, "card_imagePos", &cfg.Cards.ImagePos)
	setStr(r, "card_hoverEffect", &cfg.Cards.HoverEffect)

	// header
	setStr(r, "h_style", &cfg.Header.Style)
	setInt(r, "h_coverHeight", &cfg.Header.CoverHeight)
	setInt(r, "h_logoSize", &cfg.Header.LogoSize)
	setStr(r, "h_logoShape", &cfg.Header.LogoShape)
	setFloat(r, "h_overlay", &cfg.Header.Overlay)

	// effects
	setStr(r, "e_animSpeed", &cfg.Effects.AnimSpeed)

	// custom stylesheet
	setStr(r, "custom_css", &cfg.Custom)

	// Checkboxes only appear in the payload when ticked, so they are read
	// unconditionally from the submitted form rather than via setBool.
	if r.FormValue("checkboxes_present") != "" {
		cfg.Layout.StickyNav = formBool(r, "l_stickyNav")
		cfg.Header.ShowTagline = formBool(r, "h_showTagline")

		cfg.Display.ShowImages = formBool(r, "d_showImages")
		cfg.Display.ShowDescriptions = formBool(r, "d_showDescriptions")
		cfg.Display.ShowPrices = formBool(r, "d_showPrices")
		cfg.Display.ShowCategoryIcon = formBool(r, "d_showCategoryIcon")
		cfg.Display.ShowSearch = formBool(r, "d_showSearch")
		cfg.Display.ShowTags = formBool(r, "d_showTags")
		cfg.Display.ShowCalories = formBool(r, "d_showCalories")
		cfg.Display.ShowPrepTime = formBool(r, "d_showPrepTime")
		cfg.Display.ShowSocial = formBool(r, "d_showSocial")
		cfg.Display.ShowHours = formBool(r, "d_showHours")
		cfg.Display.ShowFooter = formBool(r, "d_showFooter")
		cfg.Display.ShowLangSwitch = formBool(r, "d_showLangSwitch")

		cfg.Effects.Animations = formBool(r, "e_animations")
		cfg.Effects.Glass = formBool(r, "e_glass")
		cfg.Effects.Grain = formBool(r, "e_grain")
		cfg.Effects.RoundedBtns = formBool(r, "e_roundedBtns")
		cfg.Effects.SmoothScroll = formBool(r, "e_smoothScroll")
		cfg.Effects.ImageZoomOnHV = formBool(r, "e_imageZoom")
	}

	cfg.Sanitize()

	if !theme.IsPreset(preset) {
		preset = "custom"
	}
	return preset, cfg
}

// setStr assigns only when the field was actually submitted, so a partial form
// leaves the rest of the configuration untouched.
func setStr(r *http.Request, key string, dst *string) {
	if !r.Form.Has(key) {
		return
	}
	*dst = formStr(r, key)
}

func setInt(r *http.Request, key string, dst *int) {
	if !r.Form.Has(key) {
		return
	}
	if n, err := strconv.Atoi(normaliseDigits(formStr(r, key))); err == nil {
		*dst = n
	}
}

func setFloat(r *http.Request, key string, dst *float64) {
	if !r.Form.Has(key) {
		return
	}
	if f, err := strconv.ParseFloat(normaliseDigits(formStr(r, key)), 64); err == nil {
		*dst = f
	}
}
