package theme

import "testing"

// TestLibrarySize reports the shape of the template library. It is written as a
// test so the numbers stay honest as the library grows.
func TestLibrarySize(t *testing.T) {
	all := Presets()
	t.Logf("total templates: %d", len(all))

	byBusiness := map[string]int{}
	universal := 0
	for _, p := range all {
		if len(p.For) == 0 {
			universal++
			continue
		}
		for _, b := range p.For {
			byBusiness[b]++
		}
	}

	t.Logf("universal: %d", universal)
	for _, b := range []string{
		"restaurant", "fastfood", "cafe", "jigaraki",
		"teahouse", "bakery", "icecream", "juicebar",
	} {
		t.Logf("  %-11s %d", b, byBusiness[b])
	}

	t.Logf("shown to a cafe owner: %d", len(PresetsFor("cafe")))
}
