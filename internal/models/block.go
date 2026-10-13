package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// BlockKind enumerates the free-form content blocks an owner can place on the
// public page.
type BlockKind string

const (
	BlockText    BlockKind = "text"
	BlockHeading BlockKind = "heading"
	BlockImage   BlockKind = "image"
	BlockBanner  BlockKind = "banner"
	BlockGallery BlockKind = "gallery"
	BlockQuote   BlockKind = "quote"
	BlockButton  BlockKind = "button"
	BlockDivider BlockKind = "divider"
	BlockSpacer  BlockKind = "spacer"
	// BlockCanvas is the free-form designer surface: elements placed anywhere,
	// at any size, rather than flowing in the normal document order.
	BlockCanvas BlockKind = "canvas"
)

// BlockKinds is the whitelist, in the order the picker offers them.
var BlockKinds = []BlockKind{
	BlockCanvas, BlockHeading, BlockText, BlockImage, BlockBanner,
	BlockGallery, BlockQuote, BlockButton, BlockDivider, BlockSpacer,
}

// ValidBlockKind reports whether a submitted kind is known.
func ValidBlockKind(k string) bool {
	for _, v := range BlockKinds {
		if string(v) == k {
			return true
		}
	}
	return false
}

// BlockStyle carries every presentation option for a single block. Values are
// validated in Normalize before they reach a template.
type BlockStyle struct {
	Size      string  `json:"size"`      // xs | sm | md | lg | xl | xxl
	Align     string  `json:"align"`     // start | center | end
	Weight    int     `json:"weight"`    // 300..900
	Color     string  `json:"color"`     // text | muted | primary | accent | price | custom
	CustomCol string  `json:"customCol"` // used when Color == "custom"
	Bg        string  `json:"bg"`        // none | surface | primary | accent | gradient | glass | image
	Width     string  `json:"width"`     // narrow | normal | wide | full
	Padding   string  `json:"padding"`   // none | sm | md | lg | xl
	Radius    int     `json:"radius"`    // 0..48 px
	Shadow    int     `json:"shadow"`    // 0..5
	Effect    string  `json:"effect"`    // none | fade | slide | zoom | float | glow | shine
	Ratio     string  `json:"ratio"`     // auto | square | 4:3 | 16:9 | 21:9 | 3:4
	Overlay   float64 `json:"overlay"`   // 0..1, for banner/image text legibility
	Height    int     `json:"height"`    // spacer height, or banner min-height, px
	Columns   int     `json:"columns"`   // gallery columns, 1..5
	Italic    bool    `json:"italic"`
	Uppercase bool    `json:"uppercase"`
	Border    bool    `json:"border"`
	Divider   string  `json:"divider"` // line | dots | dashed | gradient | ornament
}

// DefaultBlockStyle returns sensible presentation defaults for a new block.
func DefaultBlockStyle(kind BlockKind) BlockStyle {
	s := BlockStyle{
		Size: "md", Align: "center", Weight: 400, Color: "text",
		CustomCol: "#ffffff", Bg: "none", Width: "normal", Padding: "md",
		Radius: 18, Shadow: 0, Effect: "fade", Ratio: "16:9",
		Overlay: 0.45, Height: 40, Columns: 3, Divider: "line",
	}

	switch kind {
	case BlockHeading:
		s.Size, s.Weight, s.Color = "lg", 800, "primary"
	case BlockBanner:
		s.Bg, s.Padding, s.Height, s.Color = "image", "xl", 260, "text"
	case BlockQuote:
		s.Italic, s.Bg, s.Color = true, "surface", "muted"
	case BlockButton:
		s.Bg, s.Size, s.Padding = "primary", "md", "sm"
	case BlockSpacer:
		s.Padding, s.Height = "none", 48
	case BlockDivider:
		s.Padding = "sm"
	case BlockImage:
		s.Padding = "sm"
	}
	return s
}

// Normalize clamps numbers and rejects any value outside its whitelist, so the
// style is safe to turn into CSS.
func (s *BlockStyle) Normalize(kind BlockKind) {
	d := DefaultBlockStyle(kind)

	s.Size = oneOfStr(s.Size, []string{"xs", "sm", "md", "lg", "xl", "xxl"}, d.Size)
	s.Align = oneOfStr(s.Align, []string{"start", "center", "end"}, d.Align)
	s.Color = oneOfStr(s.Color, []string{"text", "muted", "primary", "accent", "price", "custom"}, d.Color)
	s.Bg = oneOfStr(s.Bg, []string{"none", "surface", "primary", "accent", "gradient", "glass", "image"}, d.Bg)
	s.Width = oneOfStr(s.Width, []string{"narrow", "normal", "wide", "full"}, d.Width)
	s.Padding = oneOfStr(s.Padding, []string{"none", "sm", "md", "lg", "xl"}, d.Padding)
	s.Effect = oneOfStr(s.Effect, []string{"none", "fade", "slide", "zoom", "float", "glow", "shine"}, d.Effect)
	s.Ratio = oneOfStr(s.Ratio, []string{"auto", "square", "4:3", "16:9", "21:9", "3:4"}, d.Ratio)
	s.Divider = oneOfStr(s.Divider, []string{"line", "dots", "dashed", "gradient", "ornament"}, d.Divider)

	s.CustomCol = hexOr(s.CustomCol, d.CustomCol)
	s.Weight = clampI(s.Weight, 300, 900, d.Weight)
	s.Radius = clampI(s.Radius, 0, 48, d.Radius)
	s.Shadow = clampI(s.Shadow, 0, 5, d.Shadow)
	s.Height = clampI(s.Height, 8, 700, d.Height)
	s.Columns = clampI(s.Columns, 1, 5, d.Columns)

	if s.Overlay < 0 {
		s.Overlay = 0
	}
	if s.Overlay > 1 {
		s.Overlay = 1
	}
}

// Block is one placed piece of custom content.
type Block struct {
	ID           int64
	RestaurantID int64
	Kind         BlockKind

	TitleFa string
	TitleEn string
	BodyFa  string
	BodyEn  string

	ImagePath string
	Images    []string

	LinkURL     string
	LinkLabelFa string
	LinkLabelEn string

	Style     BlockStyle
	SortOrder int
	IsVisible bool

	// Canvas is populated only for BlockCanvas.
	Canvas CanvasData

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Block) Title(locale string) string { return pick(locale, b.TitleFa, b.TitleEn) }
func (b *Block) Body(locale string) string  { return pick(locale, b.BodyFa, b.BodyEn) }
func (b *Block) LinkLabel(locale string) string {
	return pick(locale, b.LinkLabelFa, b.LinkLabelEn)
}

// Is reports whether the block is of a given kind, for template branching.
func (b *Block) Is(kind string) bool { return string(b.Kind) == kind }

// SectionKey is how a block is referenced inside the theme's section order.
func (b *Block) SectionKey() string { return "block:" + strconv.FormatInt(b.ID, 10) }

// StyleJSON serialises the presentation options for storage.
func (b *Block) StyleJSON() string {
	s := b.Style
	s.Normalize(b.Kind)
	out, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}
	return string(out)
}

// SetStyleJSON parses stored options on top of the kind's defaults.
func (b *Block) SetStyleJSON(raw string) {
	b.Style = DefaultBlockStyle(b.Kind)
	if strings.TrimSpace(raw) != "" && raw != "{}" {
		_ = json.Unmarshal([]byte(raw), &b.Style)
	}
	b.Style.Normalize(b.Kind)
}

// ImagesJSON serialises the gallery image list.
func (b *Block) ImagesJSON() string {
	if len(b.Images) == 0 {
		return "[]"
	}
	out, err := json.Marshal(b.Images)
	if err != nil {
		return "[]"
	}
	return string(out)
}

// SetImagesJSON parses the gallery image list, tolerating malformed data.
func (b *Block) SetImagesJSON(raw string) {
	b.Images = nil
	if strings.TrimSpace(raw) == "" {
		return
	}
	_ = json.Unmarshal([]byte(raw), &b.Images)
}

// --- small validation helpers ---------------------------------------------

func oneOfStr(v string, allowed []string, fallback string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	for _, a := range allowed {
		if v == a {
			return a
		}
	}
	return fallback
}

func clampI(v, min, max, fallback int) int {
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

const hexChars = "0123456789abcdefABCDEF"

func hexOr(v, fallback string) string {
	v = strings.TrimSpace(v)
	if len(v) != 4 && len(v) != 7 && len(v) != 9 {
		return fallback
	}
	if v[0] != '#' {
		return fallback
	}
	for _, r := range v[1:] {
		if !strings.ContainsRune(hexChars, r) {
			return fallback
		}
	}
	return v
}
