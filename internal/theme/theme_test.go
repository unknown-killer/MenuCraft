package theme

import "testing"

// TestPresetsAreValid guards the preset library: every entry needs a unique
// key and must produce a config that survives sanitising unchanged. A preset
// that shifts under Sanitize contains a value outside its allowed range.
func TestPresetsAreValid(t *testing.T) {
	seen := map[string]bool{}

	for _, p := range Presets() {
		if p.Key == "" {
			t.Fatal("preset with an empty key")
		}
		if seen[p.Key] {
			t.Fatalf("duplicate preset key %q", p.Key)
		}
		seen[p.Key] = true

		if p.LabelFa == "" || p.LabelEn == "" {
			t.Errorf("preset %q is missing a label", p.Key)
		}
		if len(p.Swatch) < 2 {
			t.Errorf("preset %q needs at least two swatch colours", p.Key)
		}

		cfg := ApplyPreset(p.Key)

		before := cfg.JSON()
		cfg.Sanitize()
		if after := cfg.JSON(); before != after {
			t.Errorf("preset %q changes under Sanitize\nbefore: %s\nafter:  %s", p.Key, before, after)
		}

		if len(cfg.Layout.Sections) == 0 {
			t.Errorf("preset %q produced no page sections", p.Key)
		}
	}

	if len(seen) < 25 {
		t.Errorf("expected a broad preset library, got %d", len(seen))
	}
}

// TestEveryBusinessHasTenTemplates is the promise the picker makes: whatever
// kind of venue an owner runs, they get at least ten designs built for it.
func TestEveryBusinessHasTenTemplates(t *testing.T) {
	businesses := []string{
		"restaurant", "fastfood", "cafe", "jigaraki",
		"teahouse", "bakery", "icecream", "juicebar",
	}

	for _, b := range businesses {
		var specific []Preset
		for _, p := range Presets() {
			if len(p.For) > 0 && p.Matches(b) {
				specific = append(specific, p)
			}
		}

		if len(specific) < 10 {
			t.Errorf("business %q has only %d dedicated templates, want at least 10", b, len(specific))
		}

		// Variety check: the templates must not all share one layout, or the
		// "ten designs" would only be ten colour schemes.
		layouts := map[string]bool{}
		for _, p := range specific {
			layouts[ApplyPreset(p.Key).Layout.Mode] = true
		}
		if len(layouts) < 3 {
			t.Errorf("business %q templates use only %d distinct layouts, want 3+", b, len(layouts))
		}
	}
}

// TestPresetsHaveReadableBothModes checks that every preset works in light and
// dark: the generated counterpart palette must keep text legible against its
// own background, or switching modes would produce an unreadable menu.
func TestPresetsHaveReadableBothModes(t *testing.T) {
	for _, p := range Presets() {
		cfg := ApplyPreset(p.Key)

		for _, pal := range []struct {
			name   string
			colors Colors
			// wantDark says the background should be the darker of the pair.
			wantDark bool
		}{
			{"light", cfg.Colors, false},
			{"dark", cfg.DarkColors, true},
		} {
			bg := luminance(pal.colors.Background)
			text := luminance(pal.colors.Text)

			if pal.wantDark && bg > 0.5 {
				t.Errorf("preset %q: dark palette background is light (lum %.2f)", p.Key, bg)
			}
			if !pal.wantDark && bg < 0.5 {
				t.Errorf("preset %q: light palette background is dark (lum %.2f)", p.Key, bg)
			}

			// Text and background must be clearly separated either way.
			if diff := text - bg; diff < 0.35 && diff > -0.35 {
				t.Errorf("preset %q %s palette: text/background contrast too low (%.2f vs %.2f)",
					p.Key, pal.name, text, bg)
			}

			if luminance(pal.colors.Primary) == 0.5 && pal.colors.Primary == "" {
				t.Errorf("preset %q %s palette has no primary colour", p.Key, pal.name)
			}
		}

		if cfg.Mode != "light" && cfg.Mode != "dark" && cfg.Mode != "auto" {
			t.Errorf("preset %q has an invalid mode %q", p.Key, cfg.Mode)
		}
	}
}

// TestLegacyThemeRepair covers themes saved before dark mode existed: the dark
// design sits in the light slot and must be moved across, not shown as-is.
func TestLegacyThemeRepair(t *testing.T) {
	legacy := Default()
	legacy.Colors = Colors{
		Primary: "#4ade80", Secondary: "#0d3b30", Accent: "#d4af37",
		Background: "#07231c", Surface: "#0d3227", SurfaceAlt: "#123c30",
		Text: "#eafaf2", TextMuted: "#94b8a9", Border: "#1c5343",
		Price: "#d4af37", Badge: "#4ade80",
	}
	legacy.Sanitize()

	if luminance(legacy.DarkColors.Background) > 0.45 {
		t.Errorf("the dark design was not moved into the dark palette: %s", legacy.DarkColors.Background)
	}
	if legacy.DarkColors.Background != "#07231c" {
		t.Errorf("dark palette should be the original design, got %s", legacy.DarkColors.Background)
	}
	if luminance(legacy.Colors.Background) < 0.5 {
		t.Errorf("light palette is still dark: %s", legacy.Colors.Background)
	}

	// A customised dark palette must survive untouched.
	custom := Default()
	custom.Colors = legacy.DarkColors
	custom.DarkColors = Colors{
		Primary: "#ff0000", Secondary: "#111111", Accent: "#00ff00",
		Background: "#010101", Surface: "#121212", SurfaceAlt: "#1a1a1a",
		Text: "#ffffff", TextMuted: "#888888", Border: "#333333",
		Price: "#ff0000", Badge: "#00ff00",
	}
	custom.Sanitize()

	if custom.DarkColors.Background != "#010101" {
		t.Errorf("a customised dark palette was overwritten: %s", custom.DarkColors.Background)
	}
}

// TestPresetsForOrdersByBusiness pins the picker's contract: this venue's looks
// come first, every curated preset stays available whatever the venue type, and
// another venue type's generated templates are filtered out.
func TestPresetsForOrdersByBusiness(t *testing.T) {
	got := PresetsFor("cafe")

	if len(got[0].For) == 0 || !got[0].Matches("cafe") {
		t.Errorf("first preset %q is not cafe-specific", got[0].Key)
	}

	// Ordering: no general preset may appear before a cafe-specific one.
	seenGeneral := false
	for _, p := range got {
		specific := len(p.For) > 0 && p.Matches("cafe")
		if !specific {
			seenGeneral = true
		} else if seenGeneral {
			t.Errorf("cafe preset %q appears after a general one", p.Key)
			break
		}
	}

	inList := make(map[string]bool, len(got))
	for _, p := range got {
		inList[p.Key] = true
	}

	for _, p := range Presets() {
		switch {
		case p.Generated && !p.Matches("cafe"):
			if inList[p.Key] {
				t.Errorf("generated preset %q for another venue type leaked into the cafe list", p.Key)
			}
		default:
			if !inList[p.Key] {
				t.Errorf("preset %q is missing from the cafe list", p.Key)
			}
		}
	}

	// A cafe owner should still get at least ten dedicated designs.
	specific := 0
	for _, p := range got {
		if len(p.For) > 0 && p.Matches("cafe") {
			specific++
		}
	}
	if specific < 10 {
		t.Errorf("cafe owner sees only %d dedicated templates", specific)
	}
}

// TestSanitizeRejectsBadValues covers the CSS-injection guard: anything that is
// not a plain hex colour must fall back to the default.
func TestSanitizeRejectsBadValues(t *testing.T) {
	c := Default()
	c.Colors.Primary = "red; } body { display:none } .x{"
	c.Colors.Accent = "url(https://evil.example/x)"
	c.Background.ImagePath = "../../etc/passwd"
	c.Layout.Columns = 99
	c.Cards.Radius = -20
	c.Typography.BaseSize = 400

	c.Sanitize()

	if c.Colors.Primary != Default().Colors.Primary {
		t.Errorf("a non-hex primary colour was accepted: %q", c.Colors.Primary)
	}
	if c.Colors.Accent != Default().Colors.Accent {
		t.Errorf("a url() accent colour was accepted: %q", c.Colors.Accent)
	}
	if c.Background.ImagePath != "" {
		t.Errorf("a traversal image path was accepted: %q", c.Background.ImagePath)
	}
	if c.Layout.Columns != 4 {
		t.Errorf("columns not clamped: %d", c.Layout.Columns)
	}
	if c.Cards.Radius != 0 {
		t.Errorf("radius not clamped: %d", c.Cards.Radius)
	}
	if c.Typography.BaseSize != 22 {
		t.Errorf("base size not clamped: %d", c.Typography.BaseSize)
	}
}

// TestReconcileSections covers the drift the page designer has to survive.
func TestReconcileSections(t *testing.T) {
	t.Run("drops deleted blocks", func(t *testing.T) {
		got := ReconcileSections([]string{"header", "block:9", "menu", "footer"}, nil)
		for _, s := range got {
			if s == "block:9" {
				t.Fatalf("a reference to a deleted block survived: %v", got)
			}
		}
	})

	t.Run("adds new blocks before the footer", func(t *testing.T) {
		got := ReconcileSections([]string{"header", "menu", "footer"}, []int64{4})
		if got[len(got)-1] != "footer" {
			t.Fatalf("footer is no longer last: %v", got)
		}
		if got[len(got)-2] != "block:4" {
			t.Fatalf("new block was not placed before the footer: %v", got)
		}
	})

	t.Run("keeps the owner's order", func(t *testing.T) {
		got := ReconcileSections([]string{"footer", "menu", "header"}, nil)
		if got[0] != "footer" || got[2] != "header" {
			t.Fatalf("stored order was not preserved: %v", got)
		}
	})

	t.Run("restores missing built-in sections", func(t *testing.T) {
		got := ReconcileSections([]string{"menu"}, nil)
		if len(got) != len(FixedSections) {
			t.Fatalf("expected every built-in section back, got %v", got)
		}
	})

	t.Run("empty stored order does not panic", func(t *testing.T) {
		got := ReconcileSections(nil, []int64{2})
		if len(got) != len(FixedSections)+1 {
			t.Fatalf("expected every section plus the block, got %v", got)
		}
	})

	t.Run("empty falls back to the default", func(t *testing.T) {
		if got := SanitizeSections(nil); len(got) != len(FixedSections) {
			t.Fatalf("expected the default section list, got %v", got)
		}
	})

	t.Run("rejects unknown keys", func(t *testing.T) {
		got := SanitizeSections([]string{"header", "<script>", "block:abc", "menu"})
		for _, s := range got {
			if s != "header" && s != "menu" {
				t.Fatalf("an unknown section key survived: %v", got)
			}
		}
	})
}

// TestIsBlockSection checks the block-reference parsing both ways.
func TestIsBlockSection(t *testing.T) {
	if id, ok := IsBlockSection("block:12"); !ok || id != 12 {
		t.Errorf("block:12 did not parse: id=%d ok=%v", id, ok)
	}
	for _, bad := range []string{"header", "block:", "block:0", "block:-3", "block:x"} {
		if _, ok := IsBlockSection(bad); ok {
			t.Errorf("%q was accepted as a block reference", bad)
		}
	}
	if got := BlockSectionKey(7); got != "block:7" {
		t.Errorf("BlockSectionKey(7) = %q", got)
	}
}
