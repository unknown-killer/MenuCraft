// Package theme models everything a restaurant owner can restyle, and renders
// it into CSS custom properties consumed by the public menu.
//
// Every value that reaches the stylesheet is validated here rather than in the
// template: the config is user-controlled, and it is interpolated into a
// <style> block, so unvalidated values would be a CSS-injection vector.
package theme

import (
	"encoding/json"
	"math"
	"strings"
)

// Config is the complete visual configuration for one venue.
type Config struct {
	// Mode selects which palette the public page uses:
	//   light — always Colors
	//   dark  — always DarkColors
	//   auto  — follows the visitor's system setting, with a manual toggle
	Mode string `json:"mode"`

	Colors Colors `json:"colors"`
	// DarkColors is the palette used in dark mode. It is kept separate rather
	// than derived, because auto-darkening a hand-picked brand palette rarely
	// produces something an owner is happy to publish.
	DarkColors Colors     `json:"darkColors"`
	Typography Typography `json:"typography"`
	Background Background `json:"background"`
	Layout     Layout     `json:"layout"`
	Cards      Cards      `json:"cards"`
	Header     Header     `json:"header"`
	Display    Display    `json:"display"`
	Effects    Effects    `json:"effects"`

	// Custom is the owner's own stylesheet, appended after the generated one so
	// it can override anything. Sanitised on save and again on render.
	Custom string `json:"custom"`
}

type Colors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Accent     string `json:"accent"`
	Background string `json:"background"`
	Surface    string `json:"surface"`
	SurfaceAlt string `json:"surfaceAlt"`
	Text       string `json:"text"`
	TextMuted  string `json:"textMuted"`
	Border     string `json:"border"`
	Price      string `json:"price"`
	Badge      string `json:"badge"`
}

type Typography struct {
	FontFa        string  `json:"fontFa"`
	FontEn        string  `json:"fontEn"`
	BaseSize      int     `json:"baseSize"`      // px, 12..22
	HeadingScale  float64 `json:"headingScale"`  // 1.05..1.6
	HeadingWeight int     `json:"headingWeight"` // 400..900
	BodyWeight    int     `json:"bodyWeight"`    // 300..600
	LineHeight    float64 `json:"lineHeight"`    // 1.2..2.2
	LetterSpacing float64 `json:"letterSpacing"` // -1..3 (px)
}

type Background struct {
	Type         string  `json:"type"` // solid | gradient | image | pattern
	Color        string  `json:"color"`
	GradientFrom string  `json:"gradientFrom"`
	GradientTo   string  `json:"gradientTo"`
	GradientVia  string  `json:"gradientVia"`
	GradientDeg  int     `json:"gradientDeg"` // 0..360
	ImagePath    string  `json:"imagePath"`
	ImageOverlay float64 `json:"imageOverlay"` // 0..1 darkening veil
	ImageBlur    int     `json:"imageBlur"`    // 0..30 px
	Pattern      string  `json:"pattern"`      // dots | grid | diagonal | crosshatch | none
	PatternOpac  float64 `json:"patternOpacity"`
	PatternScale int     `json:"patternScale"` // 8..80 px tile
	Attachment   string  `json:"attachment"`   // scroll | fixed
}

type Layout struct {
	Mode        string `json:"mode"`        // grid | list | magazine | compact
	Columns     int    `json:"columns"`     // 1..4 (desktop)
	Gap         int    `json:"gap"`         // 4..48 px
	MaxWidth    int    `json:"maxWidth"`    // 600..1600 px
	Density     string `json:"density"`     // compact | normal | relaxed
	CategoryNav string `json:"categoryNav"` // tabs | pills | sidebar | dropdown | none
	StickyNav   bool   `json:"stickyNav"`

	// Sections is the drag-ordered page layout: built-in section keys mixed
	// with "block:<id>" references to custom content.
	Sections []string `json:"sections"`
}

type Cards struct {
	Style       string `json:"style"`       // elevated | flat | outline | glass | gradient
	Radius      int    `json:"radius"`      // 0..40 px
	Shadow      int    `json:"shadow"`      // 0..5
	BorderWidth int    `json:"borderWidth"` // 0..4 px
	ImageRatio  string `json:"imageRatio"`  // square | 4:3 | 16:9 | 3:4 | circle
	ImagePos    string `json:"imagePos"`    // top | side | background | none
	HoverEffect string `json:"hoverEffect"` // lift | zoom | glow | tilt | none
}

type Header struct {
	Style       string  `json:"style"`       // hero | centered | split | minimal
	CoverHeight int     `json:"coverHeight"` // 120..600 px
	LogoSize    int     `json:"logoSize"`    // 40..200 px
	LogoShape   string  `json:"logoShape"`   // circle | rounded | square
	Overlay     float64 `json:"overlay"`     // 0..1
	ShowTagline bool    `json:"showTagline"`
}

// Display holds the on/off switches for individual page elements.
type Display struct {
	ShowImages       bool `json:"showImages"`
	ShowDescriptions bool `json:"showDescriptions"`
	ShowPrices       bool `json:"showPrices"`
	ShowCategoryIcon bool `json:"showCategoryIcon"`
	ShowSearch       bool `json:"showSearch"`
	ShowTags         bool `json:"showTags"`
	ShowCalories     bool `json:"showCalories"`
	ShowPrepTime     bool `json:"showPrepTime"`
	ShowSocial       bool `json:"showSocial"`
	ShowHours        bool `json:"showHours"`
	ShowFooter       bool `json:"showFooter"`
	ShowLangSwitch   bool `json:"showLangSwitch"`
}

type Effects struct {
	Animations    bool   `json:"animations"`
	AnimSpeed     string `json:"animSpeed"` // slow | normal | fast
	Glass         bool   `json:"glass"`
	Grain         bool   `json:"grain"`
	RoundedBtns   bool   `json:"roundedBtns"`
	SmoothScroll  bool   `json:"smoothScroll"`
	ImageZoomOnHV bool   `json:"imageZoom"`
}

// ---------------------------------------------------------------------------
// Defaults
// ---------------------------------------------------------------------------

// Default returns a complete, sane configuration. Every preset is built by
// overlaying changes on top of this, so no field is ever left at a zero value.
func Default() Config {
	return Config{
		Mode: "dark",
		DarkColors: Colors{
			Primary:    "#e8b04b",
			Secondary:  "#2a2f45",
			Accent:     "#ff7a59",
			Background: "#0f1117",
			Surface:    "#181b25",
			SurfaceAlt: "#1f2331",
			Text:       "#f5f6fa",
			TextMuted:  "#a2a8bd",
			Border:     "#2c3142",
			Price:      "#e8b04b",
			Badge:      "#ff7a59",
		},
		// The light counterpart of the midnight palette, so switching modes
		// lands somewhere designed rather than on an inverted guess.
		Colors: Colors{
			Primary:    "#a8761d",
			Secondary:  "#42485f",
			Accent:     "#d1542f",
			Background: "#fbf9f5",
			Surface:    "#ffffff",
			SurfaceAlt: "#f2eee6",
			Text:       "#1b1d24",
			TextMuted:  "#6b7183",
			Border:     "#e2ddd2",
			Price:      "#a8761d",
			Badge:      "#d1542f",
		},
		Typography: Typography{
			FontFa: "vazir", FontEn: "inter",
			BaseSize: 16, HeadingScale: 1.28, HeadingWeight: 800,
			BodyWeight: 400, LineHeight: 1.75, LetterSpacing: 0,
		},
		Background: Background{
			Type: "gradient", Color: "#0f1117",
			GradientFrom: "#0f1117", GradientVia: "#161a26", GradientTo: "#1b1430",
			GradientDeg: 160, ImageOverlay: 0.55, ImageBlur: 0,
			Pattern: "none", PatternOpac: 0.5, PatternScale: 22, Attachment: "fixed",
		},
		Layout: Layout{
			Mode: "grid", Columns: 2, Gap: 20, MaxWidth: 1180,
			Density: "normal", CategoryNav: "pills", StickyNav: true,
			Sections: DefaultSections(),
		},
		Cards: Cards{
			Style: "elevated", Radius: 20, Shadow: 2, BorderWidth: 1,
			ImageRatio: "4:3", ImagePos: "top", HoverEffect: "lift",
		},
		Header: Header{
			Style: "hero", CoverHeight: 320, LogoSize: 104,
			LogoShape: "circle", Overlay: 0.5, ShowTagline: true,
		},
		Display: Display{
			ShowImages: true, ShowDescriptions: true, ShowPrices: true,
			ShowCategoryIcon: true, ShowSearch: true, ShowTags: true,
			ShowCalories: false, ShowPrepTime: false, ShowSocial: true,
			ShowHours: false, ShowFooter: true, ShowLangSwitch: true,
		},
		Effects: Effects{
			Animations: true, AnimSpeed: "normal", Glass: false, Grain: false,
			RoundedBtns: true, SmoothScroll: true, ImageZoomOnHV: true,
		},
	}
}

// Parse decodes stored JSON on top of the defaults, then sanitises the result.
// Unknown or invalid values silently fall back rather than failing the render.
func Parse(raw string) Config {
	cfg := Default()
	if strings.TrimSpace(raw) != "" && raw != "{}" {
		_ = json.Unmarshal([]byte(raw), &cfg)
	}
	cfg.Sanitize()
	return cfg
}

// JSON serialises the config for storage.
func (c Config) JSON() string {
	b, err := json.Marshal(c)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

// Sanitize clamps every numeric field and rejects any string not on its
// whitelist. Call this before rendering; it is what makes the config safe to
// interpolate into CSS.
func (c *Config) Sanitize() {
	d := Default()

	c.Mode = oneOf(c.Mode, []string{"light", "dark", "auto"}, d.Mode)

	c.repairPalettes(d)
	c.DarkColors.sanitize(d.DarkColors)

	c.Colors.Primary = color(c.Colors.Primary, d.Colors.Primary)
	c.Colors.Secondary = color(c.Colors.Secondary, d.Colors.Secondary)
	c.Colors.Accent = color(c.Colors.Accent, d.Colors.Accent)
	c.Colors.Background = color(c.Colors.Background, d.Colors.Background)
	c.Colors.Surface = color(c.Colors.Surface, d.Colors.Surface)
	c.Colors.SurfaceAlt = color(c.Colors.SurfaceAlt, d.Colors.SurfaceAlt)
	c.Colors.Text = color(c.Colors.Text, d.Colors.Text)
	c.Colors.TextMuted = color(c.Colors.TextMuted, d.Colors.TextMuted)
	c.Colors.Border = color(c.Colors.Border, d.Colors.Border)
	c.Colors.Price = color(c.Colors.Price, d.Colors.Price)
	c.Colors.Badge = color(c.Colors.Badge, d.Colors.Badge)

	c.Typography.FontFa = oneOf(c.Typography.FontFa, FontsFa, d.Typography.FontFa)
	c.Typography.FontEn = oneOf(c.Typography.FontEn, FontsEn, d.Typography.FontEn)
	c.Typography.BaseSize = clampInt(c.Typography.BaseSize, 12, 22, d.Typography.BaseSize)
	c.Typography.HeadingScale = clampFloat(c.Typography.HeadingScale, 1.05, 1.6, d.Typography.HeadingScale)
	c.Typography.HeadingWeight = clampInt(c.Typography.HeadingWeight, 400, 900, d.Typography.HeadingWeight)
	c.Typography.BodyWeight = clampInt(c.Typography.BodyWeight, 300, 600, d.Typography.BodyWeight)
	c.Typography.LineHeight = clampFloat(c.Typography.LineHeight, 1.2, 2.2, d.Typography.LineHeight)
	c.Typography.LetterSpacing = clampFloat(c.Typography.LetterSpacing, -1, 3, 0)

	c.Background.Type = oneOf(c.Background.Type, []string{"solid", "gradient", "image", "pattern"}, d.Background.Type)
	c.Background.Color = color(c.Background.Color, d.Background.Color)
	c.Background.GradientFrom = color(c.Background.GradientFrom, d.Background.GradientFrom)
	c.Background.GradientVia = color(c.Background.GradientVia, d.Background.GradientVia)
	c.Background.GradientTo = color(c.Background.GradientTo, d.Background.GradientTo)
	c.Background.GradientDeg = clampInt(c.Background.GradientDeg, 0, 360, d.Background.GradientDeg)
	c.Background.ImagePath = safePath(c.Background.ImagePath)
	c.Background.ImageOverlay = clampFloat(c.Background.ImageOverlay, 0, 1, d.Background.ImageOverlay)
	c.Background.ImageBlur = clampInt(c.Background.ImageBlur, 0, 30, 0)
	c.Background.Pattern = oneOf(c.Background.Pattern, []string{"none", "dots", "grid", "diagonal", "crosshatch"}, d.Background.Pattern)
	// Border-coloured patterns are subtle by construction, so the slider can
	// run the full range without the overlay ever becoming harsh.
	c.Background.PatternOpac = clampFloat(c.Background.PatternOpac, 0, 1, d.Background.PatternOpac)
	c.Background.PatternScale = clampInt(c.Background.PatternScale, 8, 80, d.Background.PatternScale)
	c.Background.Attachment = oneOf(c.Background.Attachment, []string{"scroll", "fixed"}, d.Background.Attachment)

	c.Layout.Mode = oneOf(c.Layout.Mode, []string{"grid", "list", "magazine", "compact"}, d.Layout.Mode)
	c.Layout.Columns = clampInt(c.Layout.Columns, 1, 4, d.Layout.Columns)
	c.Layout.Gap = clampInt(c.Layout.Gap, 4, 48, d.Layout.Gap)
	c.Layout.MaxWidth = clampInt(c.Layout.MaxWidth, 600, 1600, d.Layout.MaxWidth)
	c.Layout.Density = oneOf(c.Layout.Density, []string{"compact", "normal", "relaxed"}, d.Layout.Density)
	c.Layout.CategoryNav = oneOf(c.Layout.CategoryNav, []string{"tabs", "pills", "sidebar", "dropdown", "none"}, d.Layout.CategoryNav)
	c.Layout.Sections = SanitizeSections(c.Layout.Sections)

	c.Cards.Style = oneOf(c.Cards.Style, []string{"elevated", "flat", "outline", "glass", "gradient"}, d.Cards.Style)
	c.Cards.Radius = clampInt(c.Cards.Radius, 0, 40, d.Cards.Radius)
	c.Cards.Shadow = clampInt(c.Cards.Shadow, 0, 5, d.Cards.Shadow)
	c.Cards.BorderWidth = clampInt(c.Cards.BorderWidth, 0, 4, d.Cards.BorderWidth)
	c.Cards.ImageRatio = oneOf(c.Cards.ImageRatio, []string{"square", "4:3", "16:9", "3:4", "circle"}, d.Cards.ImageRatio)
	c.Cards.ImagePos = oneOf(c.Cards.ImagePos, []string{"top", "side", "background", "none"}, d.Cards.ImagePos)
	c.Cards.HoverEffect = oneOf(c.Cards.HoverEffect, []string{"lift", "zoom", "glow", "tilt", "none"}, d.Cards.HoverEffect)

	c.Header.Style = oneOf(c.Header.Style, []string{"hero", "centered", "split", "minimal"}, d.Header.Style)
	c.Header.CoverHeight = clampInt(c.Header.CoverHeight, 120, 600, d.Header.CoverHeight)
	c.Header.LogoSize = clampInt(c.Header.LogoSize, 40, 200, d.Header.LogoSize)
	c.Header.LogoShape = oneOf(c.Header.LogoShape, []string{"circle", "rounded", "square"}, d.Header.LogoShape)
	c.Header.Overlay = clampFloat(c.Header.Overlay, 0, 1, d.Header.Overlay)

	c.Effects.AnimSpeed = oneOf(c.Effects.AnimSpeed, []string{"slow", "normal", "fast"}, d.Effects.AnimSpeed)

	c.Custom, _ = SanitizeCustomCSS(c.Custom)
}

// repairPalettes enforces the invariant that Colors holds the light palette and
// DarkColors the dark one.
//
// Themes saved before dark mode existed have their designed palette — usually a
// dark one — sitting in Colors, with DarkColors still at the factory default.
// Left alone, such a theme would show the dark design in light mode and an
// unrelated default in dark mode. The untouched-default check is what keeps this
// from clobbering a palette the owner has actually customised.
func (c *Config) repairPalettes(d Config) {
	if c.DarkColors.Background == "" {
		c.DarkColors = DeriveDark(c.Colors)
		return
	}

	untouched := c.DarkColors == d.DarkColors
	designedDark := luminance(c.Colors.Background) < 0.45

	if untouched && designedDark {
		designed := c.Colors
		c.DarkColors = designed
		c.Colors = DeriveLight(designed)
	}
}

// sanitize validates one palette against a fallback palette.
func (p *Colors) sanitize(d Colors) {
	p.Primary = color(p.Primary, d.Primary)
	p.Secondary = color(p.Secondary, d.Secondary)
	p.Accent = color(p.Accent, d.Accent)
	p.Background = color(p.Background, d.Background)
	p.Surface = color(p.Surface, d.Surface)
	p.SurfaceAlt = color(p.SurfaceAlt, d.SurfaceAlt)
	p.Text = color(p.Text, d.Text)
	p.TextMuted = color(p.TextMuted, d.TextMuted)
	p.Border = color(p.Border, d.Border)
	p.Price = color(p.Price, d.Price)
	p.Badge = color(p.Badge, d.Badge)
}

// DeriveDark builds a dark counterpart for a palette.
//
// Brand colours (primary, accent, price, badge) are kept and only lifted if
// they are too dark to read on a dark surface; the neutrals are rebuilt around
// the brand hue so the result looks designed rather than merely inverted.
func DeriveDark(light Colors) Colors {
	hue := hueOf(light.Primary)

	dark := Colors{
		Primary:    ensureBright(light.Primary, 0.62),
		Secondary:  mixHex(light.Secondary, "#0d0f16", 0.55),
		Accent:     ensureBright(light.Accent, 0.60),
		Background: shadeFromHue(hue, 0.07),
		Surface:    shadeFromHue(hue, 0.11),
		SurfaceAlt: shadeFromHue(hue, 0.15),
		Text:       "#f3f4f8",
		TextMuted:  shadeFromHue(hue, 0.62),
		Border:     shadeFromHue(hue, 0.22),
		Price:      ensureBright(light.Price, 0.62),
		Badge:      ensureBright(light.Badge, 0.60),
	}
	return dark
}

// DeriveLight builds a light counterpart, for owners who design in dark first.
func DeriveLight(dark Colors) Colors {
	hue := hueOf(dark.Primary)

	return Colors{
		Primary:    ensureDark(dark.Primary, 0.46),
		Secondary:  mixHex(dark.Secondary, "#ffffff", 0.4),
		Accent:     ensureDark(dark.Accent, 0.48),
		Background: tintFromHue(hue, 0.985),
		Surface:    "#ffffff",
		SurfaceAlt: tintFromHue(hue, 0.94),
		Text:       "#1a1c22",
		TextMuted:  tintFromHue(hue, 0.42),
		Border:     tintFromHue(hue, 0.88),
		Price:      ensureDark(dark.Price, 0.46),
		Badge:      ensureDark(dark.Badge, 0.48),
	}
}

var hexRunes = "0123456789abcdefABCDEF"

// color accepts only #RGB / #RRGGBB / #RRGGBBAA and returns the fallback
// otherwise, so nothing arbitrary can reach the stylesheet.
func color(v, fallback string) string {
	v = strings.TrimSpace(v)
	if len(v) != 4 && len(v) != 7 && len(v) != 9 {
		return fallback
	}
	if v[0] != '#' {
		return fallback
	}
	for _, r := range v[1:] {
		if !strings.ContainsRune(hexRunes, r) {
			return fallback
		}
	}
	return v
}

func oneOf(v string, allowed []string, fallback string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	for _, a := range allowed {
		if v == a {
			return a
		}
	}
	return fallback
}

func clampInt(v, min, max, fallback int) int {
	if v == 0 && fallback != 0 {
		return fallback
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampFloat(v, min, max, fallback float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fallback
	}
	if v == 0 && fallback != 0 {
		return fallback
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// safePath allows only uploads served by this app: a relative path under
// /uploads/ with no traversal segments.
func safePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	if strings.Contains(p, "..") || strings.ContainsAny(p, `\"'()<>;{}`) {
		return ""
	}
	if !strings.HasPrefix(p, "/uploads/") {
		return ""
	}
	return p
}

// ---------------------------------------------------------------------------
// Font registry
// ---------------------------------------------------------------------------

// Font describes a selectable typeface and the CSS stack it maps to.
type Font struct {
	Key     string
	LabelFa string
	LabelEn string
	Stack   string
}

var faFontList = []Font{
	{"vazir", "وزیرمتن", "Vazirmatn", `"Vazirmatn","Vazir",Tahoma,"Segoe UI",sans-serif`},
	{"iransans", "ایران‌سنس", "IRANSans", `"IRANSans","IRANSansX",Tahoma,"Segoe UI",sans-serif`},
	{"shabnam", "شبنم", "Shabnam", `"Shabnam","Vazirmatn",Tahoma,sans-serif`},
	{"sahel", "ساحل", "Sahel", `"Sahel","Vazirmatn",Tahoma,sans-serif`},
	{"tahoma", "تاهوما", "Tahoma", `Tahoma,"Segoe UI",Arial,sans-serif`},
	{"nazanin", "بی‌نازنین", "B Nazanin", `"B Nazanin","Nazanin","Vazirmatn",serif`},
	{"lalezar", "لاله‌زار", "Lalezar", `"Lalezar","Vazirmatn",cursive,sans-serif`},
	{"system", "پیش‌فرض سیستم", "System", `system-ui,-apple-system,"Segoe UI",Tahoma,sans-serif`},
}

var enFontList = []Font{
	{"inter", "اینتر", "Inter", `"Inter","Segoe UI",system-ui,sans-serif`},
	{"poppins", "پاپینز", "Poppins", `"Poppins","Segoe UI",system-ui,sans-serif`},
	{"playfair", "پلی‌فر", "Playfair Display", `"Playfair Display",Georgia,"Times New Roman",serif`},
	{"georgia", "جورجیا", "Georgia", `Georgia,"Times New Roman",serif`},
	{"montserrat", "مونتسرات", "Montserrat", `"Montserrat","Segoe UI",sans-serif`},
	{"courier", "کوریر", "Courier", `"Courier New",Courier,monospace`},
	{"system", "پیش‌فرض سیستم", "System", `system-ui,-apple-system,"Segoe UI",sans-serif`},
}

// FontsFa / FontsEn are the accepted keys, derived from the registries.
var FontsFa = fontKeys(faFontList)
var FontsEn = fontKeys(enFontList)

func fontKeys(list []Font) []string {
	out := make([]string, len(list))
	for i, f := range list {
		out[i] = f.Key
	}
	return out
}

// FaFonts and EnFonts expose the registries to templates.
func FaFonts() []Font { return faFontList }
func EnFonts() []Font { return enFontList }

func stackFor(list []Font, key string) string {
	for _, f := range list {
		if f.Key == key {
			return f.Stack
		}
	}
	return list[len(list)-1].Stack
}
