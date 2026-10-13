package theme

// Layout archetypes: the structural half of a template.
//
// Keeping structure separate from colour is what lets the library offer ten
// genuinely different designs per venue type without hand-writing hundreds of
// configs — a template is one archetype crossed with one palette, and the two
// vary independently.

type archetype struct {
	key     string
	labelFa string
	labelEn string
	apply   func(*Config)
}

var archetypes = []archetype{
	{
		key: "classic", labelFa: "کلاسیک", labelEn: "Classic",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.Columns = "grid", 2
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "elevated", 18, 2
			c.Cards.ImageRatio, c.Cards.ImagePos = "4:3", "top"
			c.Header.Style = "hero"
			c.Layout.CategoryNav = "pills"
		},
	},
	{
		key: "gallery", labelFa: "گالری", labelEn: "Gallery",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.Columns, c.Layout.Gap = "grid", 3, 14
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "flat", 24, 1
			c.Cards.ImageRatio, c.Cards.ImagePos = "square", "top"
			c.Cards.HoverEffect = "zoom"
			c.Header.Style = "hero"
			c.Layout.Density = "compact"
		},
	},
	{
		key: "listing", labelFa: "لیستی", labelEn: "Listing",
		apply: func(c *Config) {
			c.Layout.Mode = "list"
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "flat", 12, 0
			c.Cards.ImagePos, c.Cards.ImageRatio = "side", "square"
			c.Cards.BorderWidth = 1
			c.Header.Style = "centered"
			c.Layout.CategoryNav = "tabs"
			c.Layout.MaxWidth = 900
		},
	},
	{
		key: "editorial", labelFa: "مجله‌ای", labelEn: "Editorial",
		apply: func(c *Config) {
			c.Layout.Mode = "magazine"
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "flat", 6, 0
			c.Cards.ImageRatio = "16:9"
			c.Typography.FontEn, c.Typography.HeadingScale = "playfair", 1.4
			c.Typography.LineHeight = 1.9
			c.Layout.Density = "relaxed"
			c.Header.Style = "centered"
		},
	},
	{
		key: "poster", labelFa: "پوستری", labelEn: "Poster",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.MaxWidth = "list", 760
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "flat", 0, 0
			c.Cards.BorderWidth, c.Cards.ImagePos = 0, "none"
			c.Display.ShowImages = false
			c.Typography.HeadingScale, c.Typography.HeadingWeight = 1.45, 900
			c.Header.Style = "centered"
			c.Layout.CategoryNav = "none"
			c.Layout.Density = "relaxed"
		},
	},
	{
		key: "tiles", labelFa: "کاشی", labelEn: "Tiles",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.Gap = "compact", 10
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "flat", 14, 1
			c.Cards.ImageRatio = "square"
			c.Display.ShowDescriptions = false
			c.Layout.Density = "compact"
			c.Header.Style = "minimal"
		},
	},
	{
		key: "showcase", labelFa: "ویترین", labelEn: "Showcase",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.Columns = "grid", 2
			c.Cards.ImagePos, c.Cards.Radius = "background", 22
			c.Cards.Shadow, c.Cards.Style = 3, "flat"
			c.Cards.HoverEffect = "zoom"
			c.Header.Style = "hero"
			c.Header.CoverHeight = 420
		},
	},
	{
		key: "minimal", labelFa: "خیلی ساده", labelEn: "Minimal",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.MaxWidth = "list", 820
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "flat", 0, 0
			c.Cards.BorderWidth, c.Cards.ImagePos = 0, "none"
			c.Cards.HoverEffect = "none"
			c.Layout.CategoryNav = "tabs"
			c.Typography.HeadingWeight = 600
			c.Effects.RoundedBtns, c.Effects.Animations = false, false
			c.Header.Style = "minimal"
		},
	},
	{
		key: "bistro", labelFa: "بیسترو", labelEn: "Bistro",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.Columns = "grid", 2
			c.Cards.Style, c.Cards.Radius = "outline", 8
			c.Cards.BorderWidth, c.Cards.Shadow = 1, 0
			c.Typography.FontFa, c.Typography.FontEn = "nazanin", "georgia"
			c.Typography.HeadingScale = 1.34
			c.Header.Style = "centered"
			c.Layout.Density = "relaxed"
		},
	},
	{
		key: "modern", labelFa: "مدرن", labelEn: "Modern",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.Columns = "grid", 3
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "glass", 26, 2
			c.Cards.HoverEffect = "lift"
			c.Effects.Glass, c.Effects.SmoothScroll = true, true
			c.Typography.FontEn = "poppins"
			c.Layout.StickyNav = true
			c.Header.Style = "hero"
		},
	},
	{
		key: "catalog", labelFa: "کاتالوگ", labelEn: "Catalog",
		apply: func(c *Config) {
			c.Layout.Mode, c.Layout.CategoryNav = "list", "sidebar"
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "flat", 10, 0
			c.Cards.ImagePos, c.Cards.BorderWidth = "side", 1
			c.Layout.Density, c.Layout.MaxWidth = "compact", 1240
			c.Header.Style = "minimal"
		},
	},
	{
		key: "spotlight", labelFa: "ویژه", labelEn: "Spotlight",
		apply: func(c *Config) {
			c.Layout.Mode = "magazine"
			c.Cards.Style, c.Cards.Radius, c.Cards.Shadow = "gradient", 20, 3
			c.Cards.ImageRatio, c.Cards.HoverEffect = "4:3", "glow"
			c.Typography.HeadingWeight = 900
			c.Header.Style = "hero"
			c.Header.CoverHeight = 380
		},
	},
}

// --- palette styles --------------------------------------------------------

// paletteStyle turns a hue into a complete, coherent palette. The generated
// colours all derive from one hue, which is what keeps a template looking
// designed rather than assembled from unrelated swatches.
type paletteStyle struct {
	key   string
	dark  bool
	build func(hue float64) Colors
}

var paletteStyles = []paletteStyle{
	{
		key: "airy", dark: false,
		build: func(h float64) Colors {
			return Colors{
				Primary:    hslToHex(h, 0.62, 0.38),
				Secondary:  hslToHex(h+24, 0.30, 0.34),
				Accent:     hslToHex(h+180, 0.55, 0.45),
				Background: hslToHex(h, 0.30, 0.975),
				Surface:    "#ffffff",
				SurfaceAlt: hslToHex(h, 0.26, 0.94),
				Text:       hslToHex(h, 0.25, 0.13),
				TextMuted:  hslToHex(h, 0.12, 0.45),
				Border:     hslToHex(h, 0.20, 0.88),
				Price:      hslToHex(h, 0.66, 0.36),
				Badge:      hslToHex(h+180, 0.55, 0.45),
			}
		},
	},
	{
		key: "cream", dark: false,
		build: func(h float64) Colors {
			return Colors{
				Primary:    hslToHex(h, 0.58, 0.36),
				Secondary:  hslToHex(h+150, 0.28, 0.32),
				Accent:     hslToHex(h+32, 0.62, 0.48),
				Background: hslToHex(40, 0.42, 0.955),
				Surface:    hslToHex(42, 0.55, 0.985),
				SurfaceAlt: hslToHex(40, 0.36, 0.91),
				Text:       hslToHex(h, 0.30, 0.14),
				TextMuted:  hslToHex(h, 0.14, 0.44),
				Border:     hslToHex(38, 0.28, 0.85),
				Price:      hslToHex(h, 0.62, 0.34),
				Badge:      hslToHex(h+32, 0.62, 0.46),
			}
		},
	},
	{
		key: "deep", dark: true,
		build: func(h float64) Colors {
			return Colors{
				Primary:    hslToHex(h, 0.70, 0.62),
				Secondary:  hslToHex(h, 0.30, 0.22),
				Accent:     hslToHex(h+165, 0.62, 0.58),
				Background: hslToHex(h, 0.24, 0.08),
				Surface:    hslToHex(h, 0.22, 0.12),
				SurfaceAlt: hslToHex(h, 0.20, 0.16),
				Text:       hslToHex(h, 0.18, 0.96),
				TextMuted:  hslToHex(h, 0.12, 0.64),
				Border:     hslToHex(h, 0.18, 0.24),
				Price:      hslToHex(h, 0.76, 0.64),
				Badge:      hslToHex(h+165, 0.62, 0.58),
			}
		},
	},
	{
		key: "ink", dark: true,
		build: func(h float64) Colors {
			return Colors{
				Primary:    hslToHex(h, 0.80, 0.60),
				Secondary:  hslToHex(h, 0.16, 0.18),
				Accent:     hslToHex(h+40, 0.72, 0.58),
				Background: hslToHex(h, 0.14, 0.045),
				Surface:    hslToHex(h, 0.14, 0.085),
				SurfaceAlt: hslToHex(h, 0.13, 0.125),
				Text:       "#f7f7fa",
				TextMuted:  hslToHex(h, 0.08, 0.58),
				Border:     hslToHex(h, 0.14, 0.19),
				Price:      hslToHex(h, 0.85, 0.62),
				Badge:      hslToHex(h+40, 0.72, 0.56),
			}
		},
	},
	{
		key: "vivid", dark: true,
		build: func(h float64) Colors {
			return Colors{
				Primary:    hslToHex(h, 0.88, 0.62),
				Secondary:  hslToHex(h+200, 0.50, 0.28),
				Accent:     hslToHex(h+140, 0.85, 0.60),
				Background: hslToHex(h+205, 0.42, 0.09),
				Surface:    hslToHex(h+205, 0.38, 0.14),
				SurfaceAlt: hslToHex(h+205, 0.34, 0.19),
				Text:       "#fdf9ff",
				TextMuted:  hslToHex(h+205, 0.20, 0.66),
				Border:     hslToHex(h+205, 0.32, 0.27),
				Price:      hslToHex(h+140, 0.88, 0.62),
				Badge:      hslToHex(h, 0.88, 0.60),
			}
		},
	},
}
