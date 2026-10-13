package server

import (
	"html/template"
	"strconv"
	"strings"

	"ResturantMenue/internal/models"
)

// BlockView wraps a Block with the presentation helpers its template needs, so
// the class and style logic lives in Go rather than in template conditionals.
type BlockView struct {
	*models.Block
	Locale string
}

// Classes returns the CSS hooks for the block's style options.
func (b BlockView) Classes() string {
	s := b.Style
	parts := []string{
		"blk",
		"blk-" + string(b.Kind),
		"blk-size-" + s.Size,
		"blk-align-" + s.Align,
		"blk-color-" + s.Color,
		"blk-bg-" + s.Bg,
		"blk-w-" + s.Width,
		"blk-pad-" + s.Padding,
		"blk-fx-" + s.Effect,
	}
	if s.Italic {
		parts = append(parts, "blk-italic")
	}
	if s.Uppercase {
		parts = append(parts, "blk-upper")
	}
	if s.Border {
		parts = append(parts, "blk-border")
	}
	if b.Kind == models.BlockDivider {
		parts = append(parts, "blk-div-"+s.Divider)
	}
	return strings.Join(parts, " ")
}

// InlineStyle emits the per-block custom properties.
//
// Every value here has been through BlockStyle.Normalize, which is what makes
// it safe to mark as template.CSS.
func (b BlockView) InlineStyle() template.CSS {
	s := b.Style
	var sb strings.Builder

	sb.WriteString("--blk-radius:")
	sb.WriteString(strconv.Itoa(s.Radius))
	sb.WriteString("px;")

	sb.WriteString("--blk-weight:")
	sb.WriteString(strconv.Itoa(s.Weight))
	sb.WriteString(";")

	sb.WriteString("--blk-shadow-i:")
	sb.WriteString(strconv.Itoa(s.Shadow))
	sb.WriteString(";")

	sb.WriteString("--blk-h:")
	sb.WriteString(strconv.Itoa(s.Height))
	sb.WriteString("px;")

	sb.WriteString("--blk-cols:")
	sb.WriteString(strconv.Itoa(s.Columns))
	sb.WriteString(";")

	sb.WriteString("--blk-overlay:")
	sb.WriteString(strconv.FormatFloat(s.Overlay, 'f', 2, 64))
	sb.WriteString(";")

	sb.WriteString("--blk-ratio:")
	sb.WriteString(ratioValue(s.Ratio))
	sb.WriteString(";")

	if s.Color == "custom" {
		sb.WriteString("--blk-custom:")
		sb.WriteString(s.CustomCol)
		sb.WriteString(";")
	}

	if s.Bg == "image" && b.ImagePath != "" {
		// ImagePath is validated by sanitizeUploadPath before it is stored, so
		// it cannot contain quotes or parentheses that would break out of url().
		sb.WriteString(`--blk-img:url("`)
		sb.WriteString(b.ImagePath)
		sb.WriteString(`");`)
	}

	return template.CSS(sb.String())
}

func ratioValue(r string) string {
	switch r {
	case "square":
		return "1 / 1"
	case "4:3":
		return "4 / 3"
	case "16:9":
		return "16 / 9"
	case "21:9":
		return "21 / 9"
	case "3:4":
		return "3 / 4"
	default:
		return "auto"
	}
}

// HasText reports whether the block carries any copy worth rendering.
func (b BlockView) HasText() bool {
	return strings.TrimSpace(b.Title(b.Locale)) != "" || strings.TrimSpace(b.Body(b.Locale)) != ""
}

// Paragraphs splits the body on blank lines so multi-paragraph copy renders
// properly without allowing any HTML through.
func (b BlockView) Paragraphs() []string {
	body := strings.ReplaceAll(b.Body(b.Locale), "\r\n", "\n")
	var out []string
	for _, p := range strings.Split(body, "\n") {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// SafeLink returns the block's link only when it is a scheme we are willing to
// emit, so a stored "javascript:" URL can never become a live link.
func (b BlockView) SafeLink() string {
	u := strings.TrimSpace(b.LinkURL)
	if u == "" {
		return ""
	}
	lower := strings.ToLower(u)
	switch {
	case strings.HasPrefix(lower, "https://"),
		strings.HasPrefix(lower, "http://"),
		strings.HasPrefix(lower, "tel:"),
		strings.HasPrefix(lower, "mailto:"),
		strings.HasPrefix(u, "/"):
		return u
	}
	return ""
}

// CanvasView builds the renderer for a canvas block.
func (b BlockView) CanvasView(locale string) CanvasView {
	return CanvasView{Data: b.Canvas, Locale: locale}
}

// GalleryImages returns the images to show in a gallery block, falling back to
// the single image field when the list is empty.
func (b BlockView) GalleryImages() []string {
	if len(b.Images) > 0 {
		return b.Images
	}
	if b.ImagePath != "" {
		return []string{b.ImagePath}
	}
	return nil
}
