package server

import (
	"html/template"
	"strconv"
	"strings"

	"ResturantMenue/internal/models"
)

// canvasShadows maps the 0..5 slider to concrete shadows.
var canvasShadows = []string{
	"none",
	"0 1px 3px rgba(0,0,0,.2)",
	"0 4px 12px rgba(0,0,0,.25)",
	"0 10px 26px rgba(0,0,0,.3)",
	"0 18px 44px rgba(0,0,0,.36)",
	"0 28px 66px rgba(0,0,0,.44)",
}

// themeVar maps a "@token" colour reference onto the theme's CSS variable, so a
// canvas designed with theme colours follows the palette (and dark mode) rather
// than freezing one set of hex values.
func themeVar(v string) string {
	switch v {
	case "@text":
		return "var(--c-text)"
	case "@muted":
		return "var(--c-muted)"
	case "@primary":
		return "var(--c-primary)"
	case "@secondary":
		return "var(--c-secondary)"
	case "@accent":
		return "var(--c-accent)"
	case "@bg":
		return "var(--c-bg)"
	case "@surface":
		return "var(--c-surface)"
	case "@surfaceAlt":
		return "var(--c-surface-alt)"
	case "@border":
		return "var(--c-border)"
	case "@price":
		return "var(--c-price)"
	case "@badge":
		return "var(--c-badge)"
	case "transparent":
		return "transparent"
	}
	// Anything else has already passed canvasColor, so it is a hex literal.
	return v
}

// CanvasView renders one canvas block.
type CanvasView struct {
	Data   models.CanvasData
	Locale string
	Edit   bool
}

// StageStyle sets the canvas box: an aspect ratio plus its background.
//
// container-type is what makes the whole design responsive — element font
// sizes are expressed in cqw, so they scale with the canvas rather than the
// viewport.
func (c CanvasView) StageStyle() template.CSS {
	var b strings.Builder

	b.WriteString("aspect-ratio:")
	b.WriteString(ftoa(c.Data.Ratio))
	b.WriteString(";")

	b.WriteString("background-color:")
	b.WriteString(themeVar(c.Data.Bg))
	b.WriteString(";")

	if c.Data.BgImage != "" {
		b.WriteString("background-image:url(\"")
		b.WriteString(c.Data.BgImage)
		b.WriteString("\");background-size:")
		b.WriteString(c.Data.BgFit)
		b.WriteString(";background-position:center;")
	}

	b.WriteString("--cv-overlay:")
	b.WriteString(ftoa(c.Data.Overlay))
	b.WriteString(";")

	return template.CSS(b.String())
}

// HasOverlay reports whether the darkening veil should be drawn.
func (c CanvasView) HasOverlay() bool {
	return c.Data.BgImage != "" && c.Data.Overlay > 0
}

// Elements exposes the placed items to the template.
func (c CanvasView) Elements() []models.CanvasElement { return c.Data.Elements }

// ElementStyle turns one element's geometry and styling into inline CSS.
//
// Every number here has been clamped by CanvasElement.Normalize and every
// colour validated by canvasColor, which is what makes marking the result as
// template.CSS sound.
func (c CanvasView) ElementStyle(e models.CanvasElement) template.CSS {
	var b strings.Builder

	// --- geometry ---
	b.WriteString("left:")
	b.WriteString(ftoa(e.X))
	b.WriteString("%;top:")
	b.WriteString(ftoa(e.Y))
	b.WriteString("%;width:")
	b.WriteString(ftoa(e.W))
	b.WriteString("%;height:")
	b.WriteString(ftoa(e.H))
	b.WriteString("%;z-index:")
	b.WriteString(strconv.Itoa(e.Z))
	b.WriteString(";")

	if e.Rotate != 0 {
		b.WriteString("transform:rotate(")
		b.WriteString(ftoa(e.Rotate))
		b.WriteString("deg);")
	}

	// --- typography, sized relative to the canvas width ---
	b.WriteString("font-size:")
	b.WriteString(ftoa(e.FontSize))
	b.WriteString("cqw;font-weight:")
	b.WriteString(strconv.Itoa(e.Weight))
	b.WriteString(";line-height:")
	b.WriteString(ftoa(e.LineHeight))
	b.WriteString(";")

	if e.Letter != 0 {
		b.WriteString("letter-spacing:")
		b.WriteString(ftoa(e.Letter / 10))
		b.WriteString("cqw;")
	}

	// --- colour ---
	b.WriteString("color:")
	b.WriteString(themeVar(e.Color))
	b.WriteString(";")

	if e.Bg != "transparent" {
		b.WriteString("background-color:")
		b.WriteString(themeVar(e.Bg))
		b.WriteString(";")
	}

	if e.BorderW > 0 {
		b.WriteString("border:")
		b.WriteString(ftoa(e.BorderW))
		b.WriteString("px solid ")
		b.WriteString(themeVar(e.BorderCol))
		b.WriteString(";")
	}

	// --- shape ---
	switch e.Shape {
	case "circle":
		b.WriteString("border-radius:50%;")
	case "pill":
		b.WriteString("border-radius:999px;")
	case "triangle":
		b.WriteString("clip-path:polygon(50% 0%,100% 100%,0% 100%);")
	case "ring":
		b.WriteString("border-radius:50%;background-color:transparent;")
	default:
		if e.Radius > 0 {
			b.WriteString("border-radius:")
			b.WriteString(ftoa(e.Radius / 4))
			b.WriteString("cqw;")
		}
	}

	if e.Opacity < 1 {
		b.WriteString("opacity:")
		b.WriteString(ftoa(e.Opacity))
		b.WriteString(";")
	}
	if e.Shadow > 0 && e.Shadow < len(canvasShadows) {
		b.WriteString("box-shadow:")
		b.WriteString(canvasShadows[e.Shadow])
		b.WriteString(";")
	}
	if e.Blur > 0 {
		b.WriteString("backdrop-filter:blur(")
		b.WriteString(ftoa(e.Blur))
		b.WriteString("px);")
	}

	// --- text layout ---
	b.WriteString("text-align:")
	b.WriteString(e.Align)
	b.WriteString(";justify-content:")
	b.WriteString(flexPos(e.Align))
	b.WriteString(";align-items:")
	b.WriteString(flexPos(e.VAlign))
	b.WriteString(";")

	if e.Italic {
		b.WriteString("font-style:italic;")
	}
	if e.Underline {
		b.WriteString("text-decoration:underline;")
	}
	if e.Uppercase {
		b.WriteString("text-transform:uppercase;")
	}
	if e.FontRef == "heading" {
		b.WriteString("font-family:var(--font-body);")
	}

	return template.CSS(b.String())
}

func flexPos(v string) string {
	switch v {
	case "start":
		return "flex-start"
	case "end":
		return "flex-end"
	default:
		return "center"
	}
}

// ElementClasses returns the CSS hooks for an element.
func (c CanvasView) ElementClasses(e models.CanvasElement) string {
	parts := []string{"cv-el", "cv-" + string(e.Kind)}
	if e.Effect != "none" {
		parts = append(parts, "cv-fx-"+e.Effect)
	}
	if e.Kind == models.CanvasShape {
		parts = append(parts, "cv-shape-"+e.Shape)
	}
	if e.Kind == models.CanvasImage {
		parts = append(parts, "cv-fit-"+e.ObjFit)
	}
	return strings.Join(parts, " ")
}

// Text returns the localised copy for an element.
func (c CanvasView) Text(e models.CanvasElement) string {
	if c.Locale == "en" {
		if strings.TrimSpace(e.TextEn) != "" {
			return e.TextEn
		}
		return e.TextFa
	}
	if strings.TrimSpace(e.TextFa) != "" {
		return e.TextFa
	}
	return e.TextEn
}

// Link returns an element's link only when its scheme is one we will emit.
func (c CanvasView) Link(e models.CanvasElement) string {
	u := strings.TrimSpace(e.LinkURL)
	if u == "" {
		return ""
	}
	lower := strings.ToLower(u)
	switch {
	case strings.HasPrefix(lower, "https://"),
		strings.HasPrefix(lower, "http://"),
		strings.HasPrefix(lower, "tel:"),
		strings.HasPrefix(lower, "mailto:"),
		strings.HasPrefix(u, "/"),
		strings.HasPrefix(u, "#"):
		return u
	}
	return ""
}

// ftoa formats a float without a trailing ".000".
func ftoa(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
