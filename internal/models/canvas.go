package models

import (
	"encoding/json"
	"strings"
)

// CanvasElementKind enumerates what can be dropped onto a canvas.
type CanvasElementKind string

const (
	CanvasText    CanvasElementKind = "text"
	CanvasHeading CanvasElementKind = "heading"
	CanvasImage   CanvasElementKind = "image"
	CanvasButton  CanvasElementKind = "button"
	CanvasShape   CanvasElementKind = "shape"
	CanvasIcon    CanvasElementKind = "icon"
	CanvasPrice   CanvasElementKind = "price"
)

// CanvasElementKinds is the whitelist, in toolbar order.
var CanvasElementKinds = []CanvasElementKind{
	CanvasHeading, CanvasText, CanvasImage, CanvasButton,
	CanvasShape, CanvasIcon, CanvasPrice,
}

func ValidCanvasKind(k string) bool {
	for _, v := range CanvasElementKinds {
		if string(v) == k {
			return true
		}
	}
	return false
}

// CanvasElement is one freely positioned item.
//
// Geometry is stored as a percentage of the canvas box rather than in pixels.
// Combined with container query units for the font size, that makes a layout
// designed on a desktop scale down to a phone with its proportions intact —
// no separate mobile coordinates to maintain.
type CanvasElement struct {
	ID   string            `json:"id"`
	Kind CanvasElementKind `json:"kind"`

	X float64 `json:"x"` // 0..100, left edge
	Y float64 `json:"y"` // 0..100, top edge
	W float64 `json:"w"` // 1..100
	H float64 `json:"h"` // 1..100
	Z int     `json:"z"`

	Rotate float64 `json:"rotate"` // -180..180 degrees

	TextFa  string `json:"textFa"`
	TextEn  string `json:"textEn"`
	Image   string `json:"image"`
	LinkURL string `json:"link"`

	// FontSize is expressed in cqw: 1 unit = 1% of the canvas width.
	FontSize   float64 `json:"fontSize"`
	Weight     int     `json:"weight"`     // 100..900
	Align      string  `json:"align"`      // start | center | end
	VAlign     string  `json:"vAlign"`     // start | center | end
	Color      string  `json:"color"`      // hex, or a theme token
	Bg         string  `json:"bg"`         // hex, "transparent", or a theme token
	BorderCol  string  `json:"borderCol"`  // hex
	BorderW    float64 `json:"borderW"`    // 0..12 px
	Radius     float64 `json:"radius"`     // 0..100 (% of the shorter side)
	Opacity    float64 `json:"opacity"`    // 0..1
	Shadow     int     `json:"shadow"`     // 0..5
	Blur       float64 `json:"blur"`       // 0..30 px, backdrop blur
	LineHeight float64 `json:"lineHeight"` // 0.8..3
	Letter     float64 `json:"letter"`     // -2..8, in cqw/100
	Italic     bool    `json:"italic"`
	Underline  bool    `json:"underline"`
	Uppercase  bool    `json:"uppercase"`
	Shape      string  `json:"shape"`   // rect | circle | pill | triangle | ring
	ObjFit     string  `json:"objFit"`  // cover | contain | fill
	Effect     string  `json:"effect"`  // none | fade | slide | zoom | float | pulse
	FontRef    string  `json:"fontRef"` // body | heading
}

// DefaultCanvasElement returns a ready-to-place element of the given kind,
// sized and styled so it looks deliberate the moment it lands on the canvas.
func DefaultCanvasElement(kind CanvasElementKind, id string) CanvasElement {
	e := CanvasElement{
		ID: id, Kind: kind,
		X: 25, Y: 40, W: 50, H: 14, Z: 1,
		FontSize: 5, Weight: 500, Align: "center", VAlign: "center",
		Color: "@text", Bg: "transparent", BorderCol: "@border",
		BorderW: 0, Radius: 8, Opacity: 1, Shadow: 0, Blur: 0,
		LineHeight: 1.3, Letter: 0, Shape: "rect", ObjFit: "cover",
		Effect: "none", FontRef: "body",
	}

	switch kind {
	case CanvasHeading:
		e.FontSize, e.Weight, e.H = 9, 800, 16
		e.TextFa, e.TextEn = "عنوان شما", "Your heading"
		e.Color, e.FontRef = "@primary", "heading"
	case CanvasText:
		e.FontSize, e.H = 4, 18
		e.TextFa, e.TextEn = "متن دلخواه شما اینجا می‌آید.", "Your text goes here."
		e.Color = "@muted"
	case CanvasButton:
		e.FontSize, e.Weight, e.W, e.H = 4, 700, 30, 10
		e.TextFa, e.TextEn = "کلیک کنید", "Click here"
		e.Color, e.Bg, e.Radius = "@bg", "@primary", 50
	case CanvasImage:
		e.W, e.H, e.Radius = 40, 40, 12
	case CanvasShape:
		e.W, e.H, e.Bg, e.Radius = 30, 30, "@primary", 12
		e.Opacity = 0.9
	case CanvasIcon:
		e.W, e.H, e.FontSize = 14, 14, 9
		e.TextFa, e.TextEn = "🍽", "🍽"
		e.Color = "@primary"
	case CanvasPrice:
		e.FontSize, e.Weight, e.W, e.H = 6, 800, 26, 10
		e.TextFa, e.TextEn = "۱۲۰٬۰۰۰", "120,000"
		e.Color = "@price"
	}
	return e
}

// Normalize clamps geometry and style so an element can never escape its
// canvas or inject arbitrary CSS.
func (e *CanvasElement) Normalize() {
	if !ValidCanvasKind(string(e.Kind)) {
		e.Kind = CanvasText
	}

	e.W = clampF(e.W, 1, 100, 20)
	e.H = clampF(e.H, 1, 100, 10)
	// Position is clamped so at least part of the element stays on the canvas.
	e.X = clampF(e.X, -50, 100, 0)
	e.Y = clampF(e.Y, -50, 100, 0)

	if e.Z < 0 {
		e.Z = 0
	}
	if e.Z > 999 {
		e.Z = 999
	}

	e.Rotate = clampF(e.Rotate, -180, 180, 0)
	e.FontSize = clampF(e.FontSize, 0.5, 30, 4)
	e.LineHeight = clampF(e.LineHeight, 0.8, 3, 1.3)
	e.Letter = clampF(e.Letter, -2, 8, 0)
	e.BorderW = clampF(e.BorderW, 0, 12, 0)
	e.Radius = clampF(e.Radius, 0, 100, 0)
	e.Blur = clampF(e.Blur, 0, 30, 0)

	if e.Opacity <= 0 || e.Opacity > 1 {
		e.Opacity = 1
	}
	e.Weight = clampI(e.Weight, 100, 900, 400)
	e.Shadow = clampI(e.Shadow, 0, 5, 0)

	e.Align = oneOfStr(e.Align, []string{"start", "center", "end"}, "center")
	e.VAlign = oneOfStr(e.VAlign, []string{"start", "center", "end"}, "center")
	e.Shape = oneOfStr(e.Shape, []string{"rect", "circle", "pill", "triangle", "ring"}, "rect")
	e.ObjFit = oneOfStr(e.ObjFit, []string{"cover", "contain", "fill"}, "cover")
	e.Effect = oneOfStr(e.Effect, []string{"none", "fade", "slide", "zoom", "float", "pulse"}, "none")
	e.FontRef = oneOfStr(e.FontRef, []string{"body", "heading"}, "body")

	e.Color = canvasColor(e.Color, "@text")
	e.Bg = canvasColor(e.Bg, "transparent")
	e.BorderCol = canvasColor(e.BorderCol, "@border")

	e.TextFa = truncRunes(e.TextFa, 600)
	e.TextEn = truncRunes(e.TextEn, 600)
	e.ID = safeID(e.ID)
}

// canvasColor accepts a hex literal, the keyword "transparent", or a theme
// token like "@primary" that resolves to a CSS variable at render time.
func canvasColor(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "transparent" {
		return v
	}
	if strings.HasPrefix(v, "@") {
		switch v {
		case "@text", "@muted", "@primary", "@secondary", "@accent",
			"@bg", "@surface", "@surfaceAlt", "@border", "@price", "@badge":
			return v
		}
		return fallback
	}
	if hexOr(v, "") != "" {
		return v
	}
	return fallback
}

// safeID keeps element ids to a harmless character set; they end up in DOM
// attributes and CSS selectors.
func safeID(v string) string {
	var b strings.Builder
	for _, r := range v {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		}
		if b.Len() >= 24 {
			break
		}
	}
	if b.Len() == 0 {
		return "el"
	}
	return b.String()
}

func truncRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func clampF(v, min, max, fallback float64) float64 {
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

// CanvasData is the whole free-form section.
type CanvasData struct {
	// Ratio drives the canvas height: width / ratio. Storing a ratio rather
	// than a pixel height is what keeps the design intact on a phone.
	Ratio    float64         `json:"ratio"`   // 0.4..3.0
	Bg       string          `json:"bg"`      // hex or theme token
	BgImage  string          `json:"bgImage"` // upload path
	BgFit    string          `json:"bgFit"`   // cover | contain
	Overlay  float64         `json:"overlay"` // 0..1
	Grid     int             `json:"grid"`    // snap grid in %, 0 = off
	Elements []CanvasElement `json:"elements"`
}

// DefaultCanvas returns an empty canvas with sensible proportions.
func DefaultCanvas() CanvasData {
	return CanvasData{
		Ratio: 1.9, Bg: "@surface", BgFit: "cover",
		Overlay: 0.3, Grid: 1,
	}
}

// Normalize validates the canvas and every element on it.
func (c *CanvasData) Normalize() {
	c.Ratio = clampF(c.Ratio, 0.4, 3.0, 1.9)
	c.Overlay = clampF(c.Overlay, 0, 1, 0.3)
	c.Bg = canvasColor(c.Bg, "@surface")
	c.BgFit = oneOfStr(c.BgFit, []string{"cover", "contain"}, "cover")
	c.Grid = clampI(c.Grid, 0, 10, 1)

	if len(c.Elements) > 60 {
		c.Elements = c.Elements[:60]
	}

	seen := make(map[string]bool, len(c.Elements))
	for i := range c.Elements {
		c.Elements[i].Normalize()

		// Duplicate ids would break selection in the editor.
		id := c.Elements[i].ID
		for seen[id] {
			id += "x"
		}
		c.Elements[i].ID = id
		seen[id] = true
	}
}

// CanvasJSON serialises the canvas for storage.
func (c CanvasData) CanvasJSON() string {
	c.Normalize()
	out, err := json.Marshal(c)
	if err != nil {
		return "{}"
	}
	return string(out)
}

// ParseCanvas decodes stored canvas JSON on top of the defaults.
func ParseCanvas(raw string) CanvasData {
	c := DefaultCanvas()
	if strings.TrimSpace(raw) != "" && raw != "{}" {
		_ = json.Unmarshal([]byte(raw), &c)
	}
	c.Normalize()
	return c
}
