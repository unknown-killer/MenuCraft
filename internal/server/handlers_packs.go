package server

import (
	"net/http"
	"strings"

	"ResturantMenue/internal/menupack"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
	"ResturantMenue/internal/theme"
)

type galleryPage struct {
	Packs    []menupack.Pack
	Presets  []theme.Preset
	Business string
	Tab      string // "packs" | "themes"
}

// handleGallery shows the ready-made menus and themes side by side.
func (s *Server) handleGallery(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	tab := r.URL.Query().Get("tab")
	if tab != "themes" {
		tab = "packs"
	}

	v.Title = v.T("gal.title")
	v.Page = &galleryPage{
		Packs:    menupack.ForBusiness(v.Restaurant.BusinessType),
		Presets:  theme.PresetsFor(v.Restaurant.BusinessType),
		Business: v.Restaurant.BusinessType,
		Tab:      tab,
	}
	s.render.Render(w, http.StatusOK, "panel", "gallery", v)
}

// handleApplyPack installs a ready-made menu: categories, dishes, theme and —
// when the venue is still unnamed — the suggested wording.
func (s *Server) handleApplyPack(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := parseForm(r); err != nil {
		s.galleryRedirect(w, r, "error", v.T("err.generic"))
		return
	}

	pack, ok := menupack.ByKey(formStr(r, "pack"))
	if !ok {
		s.galleryRedirect(w, r, "error", v.T("err.not_found"))
		return
	}

	venue := v.Restaurant
	items := 0
	err := s.store.Transaction(func(tx *repository.Store) error {
		local := *s
		local.store = tx
		s := &local

		// "replace" wipes the current menu first. Categories cascade to items, so
		// deleting them is enough.
		if formBool(r, "replace") {
			existing, err := s.store.Categories.List(venue.ID)
			if err != nil {
				logf("pack: list categories for venue %d: %v", venue.ID, err)
				return err
			}
			for _, c := range existing {
				if err := s.store.Categories.Delete(c.ID, venue.ID); err != nil {
					logf("pack: clear category %d: %v", c.ID, err)
					return err
				}
			}
		}

		// Skip categories the venue already has, so appending twice is harmless.
		have := map[string]bool{}
		if current, err := s.store.Categories.List(venue.ID); err == nil {
			for _, c := range current {
				have[strings.ToLower(strings.TrimSpace(c.NameFa))] = true
				have[strings.ToLower(strings.TrimSpace(c.NameEn))] = true
			}
		} else {
			return err
		}

		for _, g := range pack.Categories() {
			if have[strings.ToLower(g.NameFa)] || have[strings.ToLower(g.NameEn)] {
				continue
			}

			cat := &models.Category{
				RestaurantID: venue.ID,
				NameFa:       g.NameFa,
				NameEn:       g.NameEn,
				Icon:         g.Icon,
				IsVisible:    true,
			}
			if err := s.store.Categories.Create(cat); err != nil {
				logf("pack: create category %q: %v", g.NameEn, err)
				return err
			}

			for _, d := range g.Dishes {
				it := &models.Item{
					RestaurantID:  venue.ID,
					CategoryID:    cat.ID,
					NameFa:        d.Fa,
					NameEn:        d.En,
					DescriptionFa: d.DescFa,
					DescriptionEn: d.DescEn,
					Price:         d.Price,
					Tags:          d.Tags,
					IsAvailable:   true,
					Calories:      d.Cal,
					PrepMinutes:   d.Min,
				}
				if err := s.store.Items.Create(it); err != nil {
					logf("pack: create item %q: %v", d.En, err)
					return err
				}
				items++
			}
		}

		// The pack's theme is half of what makes it look like its preview.
		if formBool(r, "with_theme") && theme.IsPreset(pack.Theme) {
			if err := s.store.Themes.Save(venue.ID, pack.Theme, theme.ApplyPreset(pack.Theme)); err != nil {
				logf("pack: apply theme %q: %v", pack.Theme, err)
				return err
			}
		}

		// Only fill in wording the owner has not written themselves.
		changed := false
		if strings.TrimSpace(venue.NameFa) == "" && strings.TrimSpace(venue.NameEn) == "" {
			venue.NameFa, venue.NameEn = pack.NameFa, pack.NameEn
			changed = true
		}
		if strings.TrimSpace(venue.TaglineFa) == "" && strings.TrimSpace(venue.TaglineEn) == "" {
			venue.TaglineFa, venue.TaglineEn = pack.TaglineFa, pack.TaglineEn
			changed = true
		}
		if changed {
			if err := s.store.Restaurants.Update(venue); err != nil {
				logf("pack: update venue wording: %v", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		s.galleryRedirect(w, r, "error", v.T("err.generic"))
		return
	}
	s.galleryRedirect(w, r, "success", v.T("gal.applied", v.N(items)))
}

func (s *Server) galleryRedirect(w http.ResponseWriter, r *http.Request, kind, msg string) {
	s.setFlash(w, kind, msg)
	http.Redirect(w, r, "/panel/gallery", http.StatusSeeOther)
}
