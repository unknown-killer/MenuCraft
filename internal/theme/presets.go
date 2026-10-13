package theme

// Preset is a named starting point the owner can apply and then fine-tune.
type Preset struct {
	Key     string
	LabelFa string
	LabelEn string
	NoteFa  string
	NoteEn  string
	// For is the list of business types this preset suits, used to filter the
	// picker. Empty means it works for any venue.
	For []string
	// Generated marks a template from the per-venue library. Those are only
	// offered to their own venue type; the curated presets are offered to
	// everyone regardless of their For tag.
	Generated bool
	// Swatch drives the little preview chips in the picker.
	Swatch []string
	Build  func() Config
}

// Matches reports whether a preset is a good fit for a business type.
func (p Preset) Matches(businessType string) bool {
	if len(p.For) == 0 {
		return true
	}
	for _, b := range p.For {
		if b == businessType {
			return true
		}
	}
	return false
}

// presets is ordered: the picker renders them in this sequence.
var presets = []Preset{
	{
		Key: "midnight", LabelFa: "نیمه‌شب", LabelEn: "Midnight",
		NoteFa: "تیره و لوکس با طلایی", NoteEn: "Dark luxe with gold",
		Swatch: []string{"#0f1117", "#e8b04b", "#ff7a59"},
		Build:  func() Config { return Default() },
	},
	{
		Key: "saffron", LabelFa: "زعفران", LabelEn: "Saffron",
		NoteFa: "روشن و گرم، ایرانی سنتی", NoteEn: "Warm traditional Persian",
		Swatch: []string{"#fdf6e9", "#c1440e", "#1f6f5c"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#c1440e", Secondary: "#1f6f5c", Accent: "#e0a204",
				Background: "#fdf6e9", Surface: "#fffdf7", SurfaceAlt: "#f7ecd8",
				Text: "#2e2013", TextMuted: "#7d6a52", Border: "#e6d5b8",
				Price: "#c1440e", Badge: "#1f6f5c",
			}
			c.Background.Type = "pattern"
			c.Background.Color = "#fdf6e9"
			c.Background.Pattern = "diagonal"
			c.Background.PatternOpac = 0.55
			c.Background.PatternScale = 26
			c.Typography.FontFa = "nazanin"
			c.Typography.HeadingWeight = 700
			c.Cards.Style = "outline"
			c.Cards.Radius = 14
			c.Cards.Shadow = 1
			c.Header.Style = "centered"
			return c
		},
	},
	{
		Key: "minimal", LabelFa: "مینیمال", LabelEn: "Minimal",
		NoteFa: "سفید، تمیز و مدرن", NoteEn: "Clean modern white",
		Swatch: []string{"#ffffff", "#111111", "#6b7280"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#111111", Secondary: "#6b7280", Accent: "#111111",
				Background: "#ffffff", Surface: "#ffffff", SurfaceAlt: "#f7f7f8",
				Text: "#111111", TextMuted: "#6b7280", Border: "#e6e6e8",
				Price: "#111111", Badge: "#111111",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#ffffff"
			c.Cards.Style = "flat"
			c.Cards.Radius = 8
			c.Cards.Shadow = 0
			c.Cards.BorderWidth = 1
			c.Cards.HoverEffect = "none"
			c.Layout.Mode = "list"
			c.Layout.CategoryNav = "tabs"
			c.Header.Style = "minimal"
			c.Typography.HeadingWeight = 600
			c.Effects.RoundedBtns = false
			return c
		},
	},
	{
		Key: "emerald", LabelFa: "زمرد", LabelEn: "Emerald",
		NoteFa: "سبز عمیق و شیک", NoteEn: "Deep elegant green",
		Swatch: []string{"#07231c", "#4ade80", "#d4af37"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#4ade80", Secondary: "#0d3b30", Accent: "#d4af37",
				Background: "#07231c", Surface: "#0d3227", SurfaceAlt: "#123c30",
				Text: "#eafaf2", TextMuted: "#94b8a9", Border: "#1c5343",
				Price: "#d4af37", Badge: "#4ade80",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#04180f"
			c.Background.GradientVia = "#07231c"
			c.Background.GradientTo = "#0c3a2b"
			c.Background.GradientDeg = 200
			c.Cards.Radius = 24
			c.Effects.Glass = true
			c.Cards.Style = "glass"
			return c
		},
	},
	{
		Key: "sunset", LabelFa: "غروب", LabelEn: "Sunset",
		NoteFa: "پرانرژی، مناسب فست‌فود", NoteEn: "High-energy, fast food",
		Swatch: []string{"#2b1055", "#ff5f6d", "#ffc371"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#ff5f6d", Secondary: "#7b2ff7", Accent: "#ffc371",
				Background: "#2b1055", Surface: "#3a1a68", SurfaceAlt: "#472178",
				Text: "#fff5f7", TextMuted: "#c9aee4", Border: "#5b2f93",
				Price: "#ffc371", Badge: "#ff5f6d",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#2b1055"
			c.Background.GradientVia = "#5b2472"
			c.Background.GradientTo = "#c94b6b"
			c.Background.GradientDeg = 150
			c.Cards.Style = "gradient"
			c.Cards.Radius = 28
			c.Cards.Shadow = 3
			c.Cards.HoverEffect = "zoom"
			c.Layout.Columns = 3
			c.Typography.HeadingWeight = 900
			return c
		},
	},
	{
		Key: "mocha", LabelFa: "موکا", LabelEn: "Mocha",
		NoteFa: "قهوه‌ای گرم، مناسب کافه", NoteEn: "Warm coffee tones for cafés",
		Swatch: []string{"#2c1e18", "#c8a07a", "#e8dcc8"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#c8a07a", Secondary: "#5b4034", Accent: "#e2b887",
				Background: "#2c1e18", Surface: "#3a2921", SurfaceAlt: "#46322a",
				Text: "#f4ebe1", TextMuted: "#b9a292", Border: "#543d33",
				Price: "#e2b887", Badge: "#c8a07a",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#211611"
			c.Background.GradientVia = "#2c1e18"
			c.Background.GradientTo = "#3b2a21"
			c.Background.GradientDeg = 170
			c.Typography.FontEn = "playfair"
			c.Typography.HeadingScale = 1.34
			c.Cards.Radius = 16
			c.Layout.Mode = "magazine"
			return c
		},
	},
	{
		Key: "neon", LabelFa: "نئون", LabelEn: "Neon",
		NoteFa: "سیاه با نئون، خیلی مدرن", NoteEn: "Black with neon accents",
		Swatch: []string{"#05060a", "#00e5ff", "#ff2d9b"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#00e5ff", Secondary: "#111827", Accent: "#ff2d9b",
				Background: "#05060a", Surface: "#0c0f18", SurfaceAlt: "#131828",
				Text: "#e8fbff", TextMuted: "#7f8ba3", Border: "#1e2740",
				Price: "#00e5ff", Badge: "#ff2d9b",
			}
			c.Background.Type = "pattern"
			c.Background.Color = "#05060a"
			c.Background.Pattern = "grid"
			c.Background.PatternOpac = 0.75
			c.Background.PatternScale = 34
			c.Cards.Style = "outline"
			c.Cards.BorderWidth = 1
			c.Cards.HoverEffect = "glow"
			c.Cards.Radius = 12
			c.Typography.FontEn = "montserrat"
			c.Typography.LetterSpacing = 0.4
			return c
		},
	},
	{
		Key: "royal", LabelFa: "سلطنتی", LabelEn: "Royal",
		NoteFa: "بنفش عمیق و طلایی", NoteEn: "Deep purple and gold",
		Swatch: []string{"#1a1033", "#d4af37", "#8b5cf6"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#d4af37", Secondary: "#3b2a63", Accent: "#8b5cf6",
				Background: "#1a1033", Surface: "#241645", SurfaceAlt: "#2e1d56",
				Text: "#f3eefe", TextMuted: "#a99cc9", Border: "#3d2a6f",
				Price: "#d4af37", Badge: "#8b5cf6",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#120a24"
			c.Background.GradientVia = "#1a1033"
			c.Background.GradientTo = "#2b1a52"
			c.Background.GradientDeg = 165
			c.Typography.FontEn = "playfair"
			c.Typography.HeadingScale = 1.36
			c.Cards.Radius = 22
			c.Cards.Shadow = 3
			c.Header.Style = "hero"
			return c
		},
	},
	{
		Key: "ocean", LabelFa: "اقیانوس", LabelEn: "Ocean",
		NoteFa: "آبی تازه و آرام", NoteEn: "Fresh calm blue",
		Swatch: []string{"#f0f7ff", "#0369a1", "#06b6d4"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#0369a1", Secondary: "#0e7490", Accent: "#06b6d4",
				Background: "#f0f7ff", Surface: "#ffffff", SurfaceAlt: "#e2eefb",
				Text: "#0b2540", TextMuted: "#5b7590", Border: "#cfe1f3",
				Price: "#0369a1", Badge: "#06b6d4",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#f0f7ff"
			c.Background.GradientVia = "#e6f2fd"
			c.Background.GradientTo = "#d8ecfb"
			c.Background.GradientDeg = 180
			c.Cards.Style = "elevated"
			c.Cards.Shadow = 2
			c.Cards.Radius = 18
			return c
		},
	},
	{
		Key: "rosewood", LabelFa: "شرابی", LabelEn: "Rosewood",
		NoteFa: "زرشکی کلاسیک", NoteEn: "Classic burgundy",
		Swatch: []string{"#2a0d14", "#e0607e", "#f2d8c2"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#e0607e", Secondary: "#5c1a2b", Accent: "#f2d8c2",
				Background: "#2a0d14", Surface: "#3a141f", SurfaceAlt: "#471a28",
				Text: "#fdeef1", TextMuted: "#c39aa5", Border: "#5a2233",
				Price: "#f2d8c2", Badge: "#e0607e",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#1d080e"
			c.Background.GradientVia = "#2a0d14"
			c.Background.GradientTo = "#43141f"
			c.Background.GradientDeg = 175
			c.Typography.FontEn = "georgia"
			c.Cards.Radius = 10
			c.Cards.Style = "flat"
			return c
		},
	},
	{
		Key: "pistachio", LabelFa: "پسته‌ای", LabelEn: "Pistachio",
		NoteFa: "سبز روشن و ملایم", NoteEn: "Soft light green",
		Swatch: []string{"#f4f8ee", "#5b8c3e", "#c2703d"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#5b8c3e", Secondary: "#3f6b2a", Accent: "#c2703d",
				Background: "#f4f8ee", Surface: "#ffffff", SurfaceAlt: "#e9f1de",
				Text: "#22301a", TextMuted: "#647159", Border: "#d5e3c5",
				Price: "#5b8c3e", Badge: "#c2703d",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#f4f8ee"
			c.Cards.Radius = 20
			c.Cards.Shadow = 1
			c.Layout.Mode = "grid"
			c.Layout.Columns = 3
			return c
		},
	},
	{
		Key: "charcoal", LabelFa: "زغالی", LabelEn: "Charcoal",
		NoteFa: "خاکستری مدرن با یک رنگ تاکیدی", NoteEn: "Modern grey, single accent",
		Swatch: []string{"#141517", "#f97316", "#a1a1aa"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#f97316", Secondary: "#27272a", Accent: "#fb923c",
				Background: "#141517", Surface: "#1c1d21", SurfaceAlt: "#25262b",
				Text: "#fafafa", TextMuted: "#a1a1aa", Border: "#303136",
				Price: "#f97316", Badge: "#f97316",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#141517"
			c.Cards.Style = "flat"
			c.Cards.Radius = 14
			c.Cards.Shadow = 1
			c.Layout.Mode = "compact"
			c.Layout.CategoryNav = "sidebar"
			return c
		},
	},
}

// allPresets is the hand-made sets followed by the generated per-venue library,
// built once at start-up.
var allPresets = func() []Preset {
	gen := buildGenerated()
	out := make([]Preset, 0, len(presets)+len(extraPresets)+len(gen))
	out = append(out, presets...)
	out = append(out, extraPresets...)
	out = append(out, gen...)
	return out
}()

// Presets exposes the full registry to handlers and templates.
func Presets() []Preset { return allPresets }

// PresetsFor returns the presets suited to a business type first, followed by
// the universal ones, so the picker leads with the most relevant looks.
func PresetsFor(businessType string) []Preset {
	if businessType == "" {
		return allPresets
	}

	matched := make([]Preset, 0, len(allPresets))
	rest := make([]Preset, 0, len(allPresets))

	for _, p := range allPresets {
		fits := len(p.For) > 0 && p.Matches(businessType)

		// Another venue type's generated templates would only pad the picker
		// with looks that were not designed for this venue.
		if p.Generated && !fits {
			continue
		}
		if fits {
			matched = append(matched, p)
		} else {
			rest = append(rest, p)
		}
	}
	return append(matched, rest...)
}

// ApplyPreset returns the configuration for a preset key, falling back to the
// default theme when the key is unknown.
func ApplyPreset(key string) Config {
	for _, p := range allPresets {
		if p.Key == key {
			c := p.Build()
			c.finishPalettes()
			c.Sanitize()
			return c
		}
	}
	return Default()
}

// finishPalettes works out whether a preset was designed light or dark, then
// generates the counterpart palette.
//
// Each preset only declares the single look it was designed for; without this
// step, switching a dark-designed preset into light mode would fall back to
// the default palette and lose the preset's character entirely.
func (c *Config) finishPalettes() {
	if luminance(c.Colors.Background) < 0.45 {
		designed := c.Colors
		c.Mode = "dark"
		c.DarkColors = designed
		c.Colors = DeriveLight(designed)
		return
	}
	c.Mode = "light"
	c.DarkColors = DeriveDark(c.Colors)
}

// IsPreset reports whether key names a known preset.
func IsPreset(key string) bool {
	for _, p := range allPresets {
		if p.Key == key {
			return true
		}
	}
	return false
}
