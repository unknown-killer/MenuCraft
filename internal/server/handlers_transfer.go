package server

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"ResturantMenue/internal/menudata"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
)

func (s *Server) handleTransferPage(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	v.Title = v.T("transfer.title")
	s.render.Render(w, http.StatusOK, "panel", "transfer", v)
}

func (s *Server) handleMenuExport(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.NotFound(w, r)
		return
	}
	cats, err := s.store.FullMenu(v.Restaurant.ID, false)
	if err != nil {
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if r.URL.Query().Get("format") == "csv" {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=menu-"+v.Restaurant.Slug+".csv")
		_, _ = w.Write([]byte{0xef, 0xbb, 0xbf}) // Excel recognises Persian UTF-8.
		writer := csv.NewWriter(w)
		_ = writer.Write([]string{"category_fa", "category_en", "name_fa", "name_en", "description_fa", "description_en", "price", "discount_price", "currency", "available", "tags"})
		for _, c := range cats {
			for _, it := range c.Items {
				row := []string{c.NameFa, c.NameEn, it.NameFa, it.NameEn, it.DescriptionFa, it.DescriptionEn, strconv.FormatInt(it.Price, 10), strconv.FormatInt(it.DiscountPrice, 10), v.Restaurant.Currency, strconv.FormatBool(it.IsAvailable), strings.Join(it.Tags, ",")}
				for i := range row {
					row[i] = safeCSVCell(row[i])
				}
				_ = writer.Write(row)
			}
		}
		writer.Flush()
		return
	}
	doc := menudata.Export(v.Restaurant, cats, s.absoluteMenuURL(r, v.Restaurant.Slug))
	w.Header().Set("Content-Disposition", "attachment; filename=menu-"+v.Restaurant.Slug+".json")
	writeJSON(w, http.StatusOK, doc)
}

func safeCSVCell(s string) string {
	if trimmed := strings.TrimLeft(s, " \t\r\n"); trimmed != "" && strings.ContainsAny(trimmed[:1], "=+-@") {
		return "'" + s
	}
	if strings.HasPrefix(s, "\t") || strings.HasPrefix(s, "\r") || strings.HasPrefix(s, "\n") {
		return "'" + s
	}
	return s
}

func (s *Server) handleMenuImport(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.NotFound(w, r)
		return
	}
	redirect := func(kind, key string) {
		s.setFlash(w, kind, v.T(key))
		http.Redirect(w, r, "/panel/transfer", http.StatusSeeOther)
	}
	file, _, err := r.FormFile("menu_file")
	if err != nil {
		redirect("error", "transfer.invalid")
		return
	}
	defer file.Close()
	var doc menudata.Document
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&doc); err != nil {
		redirect("error", "transfer.invalid")
		return
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF || doc.Validate() != nil {
		redirect("error", "transfer.invalid")
		return
	}
	// Currency is intentionally not changed silently when appending menu data.
	if doc.Restaurant.Currency != v.Restaurant.Currency {
		redirect("error", "transfer.currency")
		return
	}
	if err := s.importMenu(v.Restaurant.ID, doc, formBool(r, "replace")); err != nil {
		logf("menu import: %v", err)
		redirect("error", "err.generic")
		return
	}
	redirect("success", "transfer.imported")
}

func (s *Server) importMenu(venueID int64, doc menudata.Document, replace bool) error {
	return s.store.Transaction(func(tx *repository.Store) error {
		if replace {
			cats, err := tx.Categories.List(venueID)
			if err != nil {
				return err
			}
			for _, c := range cats {
				if err := tx.Categories.Delete(c.ID, venueID); err != nil {
					return err
				}
			}
		}
		for _, c := range doc.Categories {
			cat := &models.Category{RestaurantID: venueID, NameFa: c.NameFa, NameEn: c.NameEn,
				DescriptionFa: c.DescriptionFa, DescriptionEn: c.DescriptionEn, Icon: c.Icon, IsVisible: c.Visible}
			if err := tx.Categories.Create(cat); err != nil {
				return err
			}
			for _, it := range c.Items {
				tags := []string{}
				for _, tag := range knownTags {
					for _, input := range it.Tags {
						if input == tag {
							tags = append(tags, tag)
							break
						}
					}
				}
				item := &models.Item{RestaurantID: venueID, CategoryID: cat.ID, NameFa: it.NameFa, NameEn: it.NameEn,
					DescriptionFa: it.DescriptionFa, DescriptionEn: it.DescriptionEn, Price: it.Price, DiscountPrice: it.DiscountPrice,
					ImagePath: s.sanitizeUploadPath(venueID, it.ImagePath), Tags: tags, IsAvailable: it.Available, IsFeatured: it.Featured,
					Calories: it.Calories, PrepMinutes: it.PrepMinutes}
				if err := tx.Items.Create(item); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
