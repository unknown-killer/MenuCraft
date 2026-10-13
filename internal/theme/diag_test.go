package theme

import (
	"strings"
	"testing"
)

// TestBackgroundIntegrity hunts for the failure the owner reported: some
// templates rendering as a visible checkerboard or mesh instead of their
// intended background.
//
// The pattern overlay is drawn with the text colour, so its visibility depends
// on the palette it lands on — a pattern authored against a dark palette turns
// into a harsh grid when the same template is viewed light, and vice versa.
func TestBackgroundIntegrity(t *testing.T) {
	for _, p := range Presets() {
		cfg := ApplyPreset(p.Key)

		if cfg.Background.Type != "pattern" {
			continue
		}

		// The pattern layer paints with --c-border at --bg-pattern-op. Against
		// its own background it must stay decorative, never structural.
		for _, pal := range []struct {
			name   string
			colors Colors
		}{
			{"light", cfg.Colors},
			{"dark", cfg.DarkColors},
		} {
			contrast := luminance(pal.colors.Border) - luminance(pal.colors.Background)
			if contrast < 0 {
				contrast = -contrast
			}

			effective := contrast * cfg.Background.PatternOpac
			if effective > 0.12 {
				t.Errorf("preset %q (%s): pattern overlay too strong — border contrast %.2f × opacity %.2f = %.3f",
					p.Key, pal.name, contrast, cfg.Background.PatternOpac, effective)
			}
		}
	}
}

// TestGradientMatchesPalette catches templates whose configured gradient was
// authored for the opposite palette, which is how a light template ends up
// painting a dark backdrop behind dark text.
func TestGradientMatchesPalette(t *testing.T) {
	for _, p := range Presets() {
		cfg := ApplyPreset(p.Key)
		if cfg.Background.Type != "gradient" {
			continue
		}

		css := string(cfg.Variables("fa"))

		// The :root block is the published mode; its --bg-base must agree with
		// the palette that block also sets.
		root := css[strings.Index(css, ":root{"):]
		if end := strings.Index(root, "}"); end > 0 {
			root = root[:end]
		}

		bgBase := extractVar(root, "--bg-base")
		cBg := extractVar(root, "--c-bg")
		if bgBase == "" || cBg == "" {
			t.Errorf("preset %q: missing background variables", p.Key)
			continue
		}

		// Pull the first colour stop out of the gradient.
		first := firstHexIn(bgBase)
		if first == "" {
			continue
		}

		if diff := luminance(first) - luminance(cBg); diff > 0.35 || diff < -0.35 {
			t.Errorf("preset %q: gradient (%s, lum %.2f) does not match palette background (%s, lum %.2f)",
				p.Key, first, luminance(first), cBg, luminance(cBg))
		}
	}
}

func extractVar(block, name string) string {
	i := strings.Index(block, name+":")
	if i < 0 {
		return ""
	}
	rest := block[i+len(name)+1:]
	end := strings.Index(rest, ";")
	if end < 0 {
		return ""
	}
	return strings.TrimSpace(rest[:end])
}

func firstHexIn(s string) string {
	i := strings.Index(s, "#")
	if i < 0 || len(s) < i+7 {
		return ""
	}
	return s[i : i+7]
}
