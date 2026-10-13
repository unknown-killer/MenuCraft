package theme

import "fmt"

// Generated templates: one archetype crossed with one palette, per venue type.
//
// Each venue type gets its own hue set — the colours that read as "café" are
// not the ones that read as "grill house" — and every entry pairs a different
// hue with a different structure, so the ten templates an owner sees differ in
// layout as well as colour.

type businessLook struct {
	business string
	labelFa  string
	labelEn  string
	// hues are the brand colours that suit this venue type.
	hues []float64
	// order is which archetype each of the ten templates uses.
	order []string
}

var businessLooks = []businessLook{
	{
		business: "restaurant", labelFa: "رستوران", labelEn: "Restaurant",
		hues:  []float64{18, 38, 355, 150, 200},
		order: []string{"classic", "editorial", "showcase", "bistro", "listing", "spotlight", "gallery", "poster", "modern", "catalog"},
	},
	{
		business: "fastfood", labelFa: "فست‌فود", labelEn: "Fast food",
		hues:  []float64{6, 32, 48, 280, 195},
		order: []string{"gallery", "tiles", "showcase", "spotlight", "modern", "classic", "listing", "editorial", "catalog", "poster"},
	},
	{
		business: "cafe", labelFa: "کافه", labelEn: "Café",
		hues:  []float64{28, 150, 12, 205, 340},
		order: []string{"editorial", "minimal", "bistro", "gallery", "modern", "listing", "classic", "poster", "tiles", "showcase"},
	},
	{
		business: "jigaraki", labelFa: "جگرکی", labelEn: "Grill house",
		hues:  []float64{14, 2, 30, 45, 350},
		order: []string{"showcase", "classic", "listing", "spotlight", "tiles", "bistro", "poster", "gallery", "catalog", "modern"},
	},
	{
		business: "teahouse", labelFa: "سفره‌خانه", labelEn: "Teahouse",
		hues:  []float64{175, 35, 90, 15, 260},
		order: []string{"bistro", "editorial", "poster", "classic", "listing", "minimal", "showcase", "catalog", "gallery", "spotlight"},
	},
	{
		business: "bakery", labelFa: "نانوایی و شیرینی", labelEn: "Bakery",
		hues:  []float64{335, 30, 45, 190, 15},
		order: []string{"gallery", "editorial", "tiles", "minimal", "classic", "showcase", "bistro", "listing", "modern", "poster"},
	},
	{
		business: "icecream", labelFa: "بستنی", labelEn: "Ice cream",
		hues:  []float64{330, 195, 45, 265, 155},
		order: []string{"tiles", "gallery", "modern", "showcase", "spotlight", "classic", "editorial", "listing", "minimal", "catalog"},
	},
	{
		business: "juicebar", labelFa: "آبمیوه", labelEn: "Juice bar",
		hues:  []float64{95, 25, 340, 185, 55},
		order: []string{"gallery", "tiles", "modern", "classic", "showcase", "listing", "spotlight", "minimal", "editorial", "catalog"},
	},
}

// archetypeByKey looks up a structural archetype.
func archetypeByKey(key string) archetype {
	for _, a := range archetypes {
		if a.key == key {
			return a
		}
	}
	return archetypes[0]
}

// buildGenerated produces the full generated library, ten templates per venue
// type. It runs once at package initialisation.
func buildGenerated() []Preset {
	out := make([]Preset, 0, len(businessLooks)*10)

	for _, look := range businessLooks {
		for i, archKey := range look.order {
			hue := look.hues[i%len(look.hues)]
			style := paletteStyles[i%len(paletteStyles)]
			arch := archetypeByKey(archKey)

			// These are declared inside the loop body, so each closure below
			// captures its own copy.
			palette := style.build(hue)

			out = append(out, Preset{
				Key:       fmt.Sprintf("%s-%s-%s", look.business, arch.key, style.key),
				LabelFa:   fmt.Sprintf("%s · %s", look.labelFa, arch.labelFa),
				LabelEn:   fmt.Sprintf("%s · %s", look.labelEn, arch.labelEn),
				NoteFa:    paletteNoteFa(style.key),
				NoteEn:    paletteNoteEn(style.key),
				For:       []string{look.business},
				Generated: true,
				Swatch:    []string{palette.Background, palette.Primary, palette.Accent},
				Build: func() Config {
					c := Default()
					c.Colors = palette
					arch.apply(&c)
					return c
				},
			})
		}
	}
	return out
}

func paletteNoteFa(style string) string {
	switch style {
	case "airy":
		return "روشن و تمیز"
	case "cream":
		return "کرم و گرم"
	case "deep":
		return "تیره و عمیق"
	case "ink":
		return "مشکی و شیک"
	case "vivid":
		return "تیره و پرانرژی"
	}
	return ""
}

func paletteNoteEn(style string) string {
	switch style {
	case "airy":
		return "Light and clean"
	case "cream":
		return "Warm cream"
	case "deep":
		return "Deep and dark"
	case "ink":
		return "Sleek black"
	case "vivid":
		return "Dark and vivid"
	}
	return ""
}
