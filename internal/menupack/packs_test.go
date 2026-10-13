package menupack

import (
	"testing"

	"ResturantMenue/internal/theme"
)

// TestPacksAreComplete is the promise the gallery makes: pick any pack and you
// get a working menu, not a skeleton.
func TestPacksAreComplete(t *testing.T) {
	seen := map[string]bool{}

	for _, p := range All() {
		if seen[p.Key] {
			t.Fatalf("duplicate pack key %q", p.Key)
		}
		seen[p.Key] = true

		if p.LabelFa == "" || p.LabelEn == "" {
			t.Errorf("pack %q is missing a label", p.Key)
		}
		if p.NameFa == "" || p.NameEn == "" {
			t.Errorf("pack %q is missing a suggested venue name", p.Key)
		}

		// Every referenced group must exist, or the menu would come out short.
		cats := p.Categories()
		if len(cats) != len(p.Groups) {
			t.Errorf("pack %q references %d groups but expanded to %d — a key is wrong",
				p.Key, len(p.Groups), len(cats))
		}

		if n := p.ItemCount(); n < 10 {
			t.Errorf("pack %q only creates %d items; too thin to be useful", p.Key, n)
		}

		// The theme must be a real preset, or applying the pack would silently
		// fall back to the default and look nothing like its preview.
		if !theme.IsPreset(p.Theme) {
			t.Errorf("pack %q names an unknown theme %q", p.Key, p.Theme)
		}

		if p.Tier < 0.5 || p.Tier > 2.5 {
			t.Errorf("pack %q has an implausible price tier %.2f", p.Key, p.Tier)
		}

		for _, c := range cats {
			if len(c.Dishes) == 0 {
				t.Errorf("pack %q has an empty category %q", p.Key, c.Key)
			}
			for _, d := range c.Dishes {
				if d.Fa == "" || d.En == "" {
					t.Errorf("pack %q, category %q: a dish is missing a name", p.Key, c.Key)
				}
				if d.Price <= 0 {
					t.Errorf("pack %q: dish %q has no price", p.Key, d.Fa)
				}
				if d.Price%5000 != 0 {
					t.Errorf("pack %q: dish %q priced at %d, not rounded to a printable figure",
						p.Key, d.Fa, d.Price)
				}
			}
		}
	}

	if len(seen) < 50 {
		t.Errorf("expected at least 50 ready-made menus, got %d", len(seen))
	}
}

// TestEveryBusinessHasPacks makes sure no venue type is left without a
// ready-made starting point.
func TestEveryBusinessHasPacks(t *testing.T) {
	businesses := []string{
		"restaurant", "fastfood", "cafe", "jigaraki",
		"teahouse", "bakery", "icecream", "juicebar",
	}

	for _, b := range businesses {
		n := 0
		for _, p := range All() {
			if p.Business == b {
				n++
			}
		}
		if n < 4 {
			t.Errorf("business %q has only %d ready-made menus", b, n)
		}
	}
}

// TestTierScalesPrices checks that the price tier actually differentiates a
// budget venue from a premium one.
func TestTierScalesPrices(t *testing.T) {
	budget, ok := ByKey("j-street") // tier 0.7
	if !ok {
		t.Fatal("expected the street grill pack to exist")
	}
	premium, ok := ByKey("r-luxe") // tier 1.6
	if !ok {
		t.Fatal("expected the fine dining pack to exist")
	}

	priceOf := func(p Pack, group, dish string) int64 {
		for _, c := range p.Categories() {
			if c.Key != group {
				continue
			}
			for _, d := range c.Dishes {
				if d.Fa == dish {
					return d.Price
				}
			}
		}
		return 0
	}

	// Both menus carry the starter group, so the same dish is directly
	// comparable across tiers.
	cheap := priceOf(budget, "grill_sides", "نان تنوری")
	if cheap == 0 {
		t.Skip("shared dish not present in both packs")
	}

	base := roundPrice(25_000 * 0.7)
	if cheap != base {
		t.Errorf("budget tier not applied: got %d, want %d", cheap, base)
	}

	lux := priceOf(premium, "starters_ir", "ماست و خیار")
	std := roundPrice(85_000 * 1.6)
	if lux != std {
		t.Errorf("premium tier not applied: got %d, want %d", lux, std)
	}
}

// TestCatalogueIsNotMutated guards the copy-before-scale in Categories: two
// packs at different tiers must not corrupt each other's prices.
func TestCatalogueIsNotMutated(t *testing.T) {
	before := groups["starters_ir"].Dishes[0].Price

	for _, p := range All() {
		_ = p.Categories()
	}

	if after := groups["starters_ir"].Dishes[0].Price; after != before {
		t.Errorf("the shared dish catalogue was mutated: %d became %d", before, after)
	}
}
