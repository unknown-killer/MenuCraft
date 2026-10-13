package theme

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

// shadowLevels maps the 0..5 slider to concrete elevation shadows.
var shadowLevels = []string{
	"none",
	"0 1px 2px rgba(0,0,0,.18)",
	"0 4px 14px rgba(0,0,0,.22)",
	"0 10px 30px rgba(0,0,0,.28)",
	"0 18px 48px rgba(0,0,0,.34)",
	"0 28px 70px rgba(0,0,0,.42)",
}

// densityPadding maps the density switch to card/section padding in px.
var densityPadding = map[string][2]int{
	"compact": {12, 28},
	"normal":  {18, 44},
	"relaxed": {26, 64},
}

// aspectRatios maps the ratio choice to a CSS aspect-ratio value.
var aspectRatios = map[string]string{
	"square": "1 / 1",
	"4:3":    "4 / 3",
	"16:9":   "16 / 9",
	"3:4":    "3 / 4",
	"circle": "1 / 1",
}

var animDurations = map[string]string{
	"slow":   "620ms",
	"normal": "320ms",
	"fast":   "160ms",
}

// Variables renders the config as a CSS custom-property block.
//
// The return type is template.CSS, which tells html/template the string is
// already trusted. That is only sound because Sanitize has run: callers must
// pass a config that came from Parse or has had Sanitize called on it.
func (c Config) Variables(locale string) template.CSS {
	c.Sanitize()

	var b strings.Builder
	b.Grow(2048)
	w := func(name, value string) {
		b.WriteString("  ")
		b.WriteString(name)
		b.WriteString(": ")
		b.WriteString(value)
		b.WriteString(";\n")
	}

	b.WriteString(":root{\n")

	// -- palette ------------------------------------------------------------
	// The base palette is whichever mode the venue publishes in; the other one
	// is emitted afterwards inside media queries and [data-mode] overrides.
	base := c.Colors
	if c.Mode == "dark" {
		base = c.DarkColors
	}

	w("--c-primary", base.Primary)
	w("--c-secondary", base.Secondary)
	w("--c-accent", base.Accent)
	w("--c-bg", base.Background)
	w("--c-surface", base.Surface)
	w("--c-surface-alt", base.SurfaceAlt)
	w("--c-text", base.Text)
	w("--c-muted", base.TextMuted)
	w("--c-border", base.Border)
	w("--c-price", base.Price)
	w("--c-badge", base.Badge)

	// Derived translucent tints. color-mix keeps these in sync with the base
	// colours automatically instead of hard-coding a second palette.
	w("--c-primary-08", "color-mix(in srgb, var(--c-primary) 8%, transparent)")
	w("--c-primary-16", "color-mix(in srgb, var(--c-primary) 16%, transparent)")
	w("--c-primary-30", "color-mix(in srgb, var(--c-primary) 30%, transparent)")
	w("--c-surface-70", "color-mix(in srgb, var(--c-surface) 70%, transparent)")
	w("--c-surface-88", "color-mix(in srgb, var(--c-surface) 88%, transparent)")
	w("--c-text-10", "color-mix(in srgb, var(--c-text) 10%, transparent)")

	// -- typography ---------------------------------------------------------
	font := stackFor(faFontList, c.Typography.FontFa)
	if locale == "en" {
		font = stackFor(enFontList, c.Typography.FontEn)
	}
	w("--font-body", font)
	w("--font-fa", stackFor(faFontList, c.Typography.FontFa))
	w("--font-en", stackFor(enFontList, c.Typography.FontEn))
	w("--fs-base", itoaPx(c.Typography.BaseSize))
	w("--fs-sm", fmt.Sprintf("%.3frem", 0.875))
	w("--h-scale", ftoa(c.Typography.HeadingScale))
	w("--h1", fmt.Sprintf("%.3frem", pow(c.Typography.HeadingScale, 3)))
	w("--h2", fmt.Sprintf("%.3frem", pow(c.Typography.HeadingScale, 2)))
	w("--h3", fmt.Sprintf("%.3frem", pow(c.Typography.HeadingScale, 1)))
	w("--h-weight", strconv.Itoa(c.Typography.HeadingWeight))
	w("--body-weight", strconv.Itoa(c.Typography.BodyWeight))
	w("--lh", ftoa(c.Typography.LineHeight))
	w("--ls", ftoa(c.Typography.LetterSpacing)+"px")

	// -- layout -------------------------------------------------------------
	w("--gap", itoaPx(c.Layout.Gap))
	w("--max-w", itoaPx(c.Layout.MaxWidth))
	w("--cols", strconv.Itoa(c.Layout.Columns))
	pad := densityPadding[c.Layout.Density]
	w("--pad-card", itoaPx(pad[0]))
	w("--pad-section", itoaPx(pad[1]))

	// -- cards --------------------------------------------------------------
	w("--radius", itoaPx(c.Cards.Radius))
	w("--radius-sm", itoaPx(maxInt(0, c.Cards.Radius/2)))
	w("--radius-lg", itoaPx(c.Cards.Radius+8))
	w("--shadow", shadowLevels[clampIdx(c.Cards.Shadow, len(shadowLevels))])
	w("--shadow-hover", shadowLevels[clampIdx(c.Cards.Shadow+1, len(shadowLevels))])
	w("--border-w", itoaPx(c.Cards.BorderWidth))
	w("--img-ratio", aspectRatios[c.Cards.ImageRatio])
	if c.Cards.ImageRatio == "circle" {
		w("--img-radius", "50%")
	} else {
		w("--img-radius", "calc(var(--radius) - var(--border-w))")
	}

	// -- header -------------------------------------------------------------
	w("--cover-h", itoaPx(c.Header.CoverHeight))
	w("--logo-size", itoaPx(c.Header.LogoSize))
	switch c.Header.LogoShape {
	case "circle":
		w("--logo-radius", "50%")
	case "square":
		w("--logo-radius", "0px")
	default:
		w("--logo-radius", "18px")
	}
	w("--cover-overlay", ftoa(c.Header.Overlay))

	// -- effects ------------------------------------------------------------
	dur := animDurations[c.Effects.AnimSpeed]
	if dur == "" {
		dur = animDurations["normal"]
	}
	if !c.Effects.Animations {
		dur = "1ms"
	}
	w("--anim", dur)
	if c.Effects.RoundedBtns {
		w("--btn-radius", "999px")
	} else {
		w("--btn-radius", "calc(var(--radius) / 2)")
	}
	if c.Effects.Glass {
		w("--glass-blur", "blur(18px)")
	} else {
		w("--glass-blur", "none")
	}

	// -- background layers --------------------------------------------------
	baseIsDark := c.Mode == "dark"
	w("--bg-base", c.backgroundBaseFor(base, c.gradientIsDark() == baseIsDark))
	w("--bg-image", c.backgroundImage())
	w("--bg-blur", itoaPx(c.Background.ImageBlur))
	w("--bg-overlay", fmt.Sprintf("rgba(0,0,0,%s)", ftoa(c.Background.ImageOverlay)))
	w("--bg-pattern", c.patternImage())
	w("--bg-pattern-op", ftoa(c.Background.PatternOpac))
	w("--bg-pattern-size", itoaPx(c.Background.PatternScale))
	w("--bg-attach", c.Background.Attachment)

	b.WriteString("}\n")

	// -- the other mode -----------------------------------------------------
	//
	// In auto mode the visitor's system preference decides, and an explicit
	// [data-mode] on <html> (set by the toggle) wins over both. Emitting the
	// overrides as plain custom-property blocks means every component picks up
	// the switch without any component-level dark styling.
	writePalette := func(selector string, p Colors, dark bool) {
		b.WriteString(selector)
		b.WriteString("{\n")
		// The page background travels with the palette, or a mode switch would
		// leave the text and the backdrop from different themes.
		w("--bg-base", c.backgroundBaseFor(p, c.gradientIsDark() == dark))
		w("--c-primary", p.Primary)
		w("--c-secondary", p.Secondary)
		w("--c-accent", p.Accent)
		w("--c-bg", p.Background)
		w("--c-surface", p.Surface)
		w("--c-surface-alt", p.SurfaceAlt)
		w("--c-text", p.Text)
		w("--c-muted", p.TextMuted)
		w("--c-border", p.Border)
		w("--c-price", p.Price)
		w("--c-badge", p.Badge)
		b.WriteString("}\n")
	}

	switch c.Mode {
	case "auto":
		b.WriteString("@media (prefers-color-scheme: dark){\n")
		writePalette(":root:not([data-mode=\"light\"])", c.DarkColors, true)
		b.WriteString("}\n")
		writePalette(":root[data-mode=\"dark\"]", c.DarkColors, true)
		writePalette(":root[data-mode=\"light\"]", c.Colors, false)
	case "dark":
		// Published dark, but a visitor can still ask for light.
		writePalette(":root[data-mode=\"light\"]", c.Colors, false)
	default:
		writePalette(":root[data-mode=\"dark\"]", c.DarkColors, true)
	}

	return template.CSS(b.String())
}

// ColorScheme is the value for the color-scheme meta tag, so form controls and
// scrollbars match the palette.
func (c Config) ColorScheme() string {
	switch c.Mode {
	case "auto":
		return "light dark"
	case "dark":
		return "dark"
	default:
		return "light"
	}
}

// ShowModeToggle reports whether the public page should offer a switch.
func (c Config) ShowModeToggle() bool { return c.Mode == "auto" }

// backgroundBaseFor returns the page background to paint behind a palette.
//
// Solid, pattern and image backgrounds simply follow the palette. A gradient
// keeps its configured stops only in the mode it was authored for; the other
// mode gets a gradient built from that palette's own neutrals. Without this,
// switching to light mode would leave dark text sitting on the dark gradient
// the owner designed, which is unreadable.
func (c Config) backgroundBaseFor(p Colors, authored bool) string {
	if c.Background.Type != "gradient" {
		return p.Background
	}

	if authored {
		return fmt.Sprintf("linear-gradient(%ddeg, %s 0%%, %s 52%%, %s 100%%)",
			c.Background.GradientDeg, c.Background.GradientFrom,
			c.Background.GradientVia, c.Background.GradientTo)
	}
	return fmt.Sprintf("linear-gradient(%ddeg, %s 0%%, %s 52%%, %s 100%%)",
		c.Background.GradientDeg, p.Background, p.Surface, p.SurfaceAlt)
}

// gradientIsDark reports whether the configured gradient was authored for the
// dark palette, judged by the brightness of its first stop.
func (c Config) gradientIsDark() bool {
	return luminance(c.Background.GradientFrom) < 0.45
}

// backgroundImage returns a url() for the image layer, or none.
func (c Config) backgroundImage() string {
	if c.Background.Type != "image" || c.Background.ImagePath == "" {
		return "none"
	}
	// ImagePath already passed safePath: no quotes, parens or traversal.
	return fmt.Sprintf("url(%q)", c.Background.ImagePath)
}

// patternImage builds the decorative overlay out of pure CSS gradients, which
// avoids shipping or encoding any image asset.
func (c Config) patternImage() string {
	if c.Background.Type != "pattern" || c.Background.Pattern == "none" {
		return "none"
	}
	// Painted with the border colour, not the text colour. Border is the one
	// palette entry defined as "a line you can barely see against the
	// background", which is exactly what a decorative pattern needs; using the
	// text colour made a 22px grid read as graph paper in both modes.
	col := "var(--c-border)"
	switch c.Background.Pattern {
	case "dots":
		return fmt.Sprintf("radial-gradient(%s 1.2px, transparent 1.2px)", col)
	case "grid":
		return fmt.Sprintf("linear-gradient(%s 1px, transparent 1px), linear-gradient(90deg, %s 1px, transparent 1px)", col, col)
	case "diagonal":
		return fmt.Sprintf("repeating-linear-gradient(45deg, %s 0 1px, transparent 1px 14px)", col)
	case "crosshatch":
		return fmt.Sprintf("repeating-linear-gradient(45deg, %s 0 1px, transparent 1px 12px), repeating-linear-gradient(-45deg, %s 0 1px, transparent 1px 12px)", col, col)
	}
	return "none"
}

// BodyClasses returns the layout/style hooks the stylesheet keys off, so the
// template does not need to reimplement this mapping.
func (c Config) BodyClasses() string {
	c.Sanitize()
	parts := []string{
		"layout-" + c.Layout.Mode,
		"cards-" + c.Cards.Style,
		"imgpos-" + strings.ReplaceAll(c.Cards.ImagePos, ":", ""),
		"hover-" + c.Cards.HoverEffect,
		"nav-" + c.Layout.CategoryNav,
		"header-" + c.Header.Style,
		"bg-" + c.Background.Type,
		"density-" + c.Layout.Density,
	}
	if c.Effects.Glass {
		parts = append(parts, "fx-glass")
	}
	if c.Effects.Grain {
		parts = append(parts, "fx-grain")
	}
	if c.Effects.Animations {
		parts = append(parts, "fx-anim")
	}
	if c.Effects.ImageZoomOnHV {
		parts = append(parts, "fx-imgzoom")
	}
	if c.Layout.StickyNav {
		parts = append(parts, "nav-sticky")
	}
	if !c.Display.ShowImages {
		parts = append(parts, "no-images")
	}
	return strings.Join(parts, " ")
}

func itoaPx(v int) string { return strconv.Itoa(v) + "px" }

func ftoa(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func pow(base float64, n int) float64 {
	out := 1.0
	for i := 0; i < n; i++ {
		out *= base
	}
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampIdx(v, n int) int {
	if v < 0 {
		return 0
	}
	if v >= n {
		return n - 1
	}
	return v
}
