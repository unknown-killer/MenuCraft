package theme

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Small colour helpers used to derive one palette from another. Everything
// here takes and returns "#rrggbb" strings; malformed input falls back to a
// neutral rather than erroring, because these run during rendering.

// parseHex splits "#rgb" or "#rrggbb" into 0..1 components.
func parseHex(h string) (r, g, b float64, ok bool) {
	h = strings.TrimSpace(strings.TrimPrefix(h, "#"))

	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) == 8 {
		h = h[:6] // drop any alpha
	}
	if len(h) != 6 {
		return 0, 0, 0, false
	}

	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return float64((v>>16)&0xFF) / 255, float64((v>>8)&0xFF) / 255, float64(v&0xFF) / 255, true
}

func toHex(r, g, b float64) string {
	clamp := func(v float64) int {
		n := int(math.Round(v * 255))
		if n < 0 {
			return 0
		}
		if n > 255 {
			return 255
		}
		return n
	}
	return fmt.Sprintf("#%02x%02x%02x", clamp(r), clamp(g), clamp(b))
}

// luminance is the perceptual brightness used to decide whether a colour needs
// lifting or darkening for contrast.
func luminance(hex string) float64 {
	r, g, b, ok := parseHex(hex)
	if !ok {
		return 0.5
	}
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// hueOf returns the hue angle (0..360) of a colour, or 210 when it cannot be
// parsed — a neutral blue-grey that suits most palettes.
func hueOf(hex string) float64 {
	r, g, b, ok := parseHex(hex)
	if !ok {
		return 210
	}

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	d := max - min
	if d == 0 {
		return 210
	}

	var h float64
	switch max {
	case r:
		h = math.Mod((g-b)/d, 6)
	case g:
		h = (b-r)/d + 2
	default:
		h = (r-g)/d + 4
	}

	h *= 60
	if h < 0 {
		h += 360
	}
	return h
}

// hslToHex converts h (0..360), s and l (0..1) to a hex string.
func hslToHex(h, s, l float64) string {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return toHex(r+m, g+m, b+m)
}

// shadeFromHue builds a very desaturated dark neutral carrying a hint of the
// brand hue, which is what stops a generated dark theme looking generic.
func shadeFromHue(hue, lightness float64) string {
	return hslToHex(hue, 0.14, lightness)
}

// tintFromHue is the light-mode counterpart.
func tintFromHue(hue, lightness float64) string {
	return hslToHex(hue, 0.16, lightness)
}

// mixHex blends two colours; amount is how much of b to take.
func mixHex(a, b string, amount float64) string {
	ar, ag, ab, ok1 := parseHex(a)
	br, bg, bb, ok2 := parseHex(b)
	if !ok1 || !ok2 {
		return a
	}
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	return toHex(
		ar+(br-ar)*amount,
		ag+(bg-ag)*amount,
		ab+(bb-ab)*amount,
	)
}

// ensureBright lifts a colour toward white until it is readable on a dark
// background, leaving already-bright colours untouched.
func ensureBright(hex string, want float64) string {
	if luminance(hex) >= want {
		return hex
	}
	out := hex
	for i := 0; i < 12 && luminance(out) < want; i++ {
		out = mixHex(out, "#ffffff", 0.12)
	}
	return out
}

// ensureDark is the mirror image, for light backgrounds.
func ensureDark(hex string, want float64) string {
	if luminance(hex) <= want {
		return hex
	}
	out := hex
	for i := 0; i < 12 && luminance(out) > want; i++ {
		out = mixHex(out, "#000000", 0.12)
	}
	return out
}
