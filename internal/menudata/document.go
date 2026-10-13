// Package menudata defines the portable, versioned menu exchange format.
package menudata

import (
	"ResturantMenue/internal/models"
	"errors"
	"strings"
)

const Version = 1
const MaxPrice int64 = 1_000_000_000_000

type Document struct {
	Version    int        `json:"version"`
	SourceURL  string     `json:"source_url"`
	Restaurant Restaurant `json:"restaurant"`
	Categories []Category `json:"categories"`
}
type Restaurant struct {
	NameFa        string `json:"name_fa"`
	NameEn        string `json:"name_en"`
	DescriptionFa string `json:"description_fa"`
	DescriptionEn string `json:"description_en"`
	Currency      string `json:"currency"`
}
type Category struct {
	NameFa        string `json:"name_fa"`
	NameEn        string `json:"name_en"`
	DescriptionFa string `json:"description_fa"`
	DescriptionEn string `json:"description_en"`
	Icon          string `json:"icon"`
	Visible       bool   `json:"visible"`
	Items         []Item `json:"items"`
}
type Item struct {
	NameFa        string   `json:"name_fa"`
	NameEn        string   `json:"name_en"`
	DescriptionFa string   `json:"description_fa"`
	DescriptionEn string   `json:"description_en"`
	Price         int64    `json:"price"`
	DiscountPrice int64    `json:"discount_price"`
	ImagePath     string   `json:"image_path"`
	Tags          []string `json:"tags"`
	Available     bool     `json:"available"`
	Featured      bool     `json:"featured"`
	Calories      int      `json:"calories"`
	PrepMinutes   int      `json:"prep_minutes"`
}

func Export(venue *models.Restaurant, cats []*models.Category, source string) Document {
	d := Document{Version: Version, SourceURL: source, Restaurant: Restaurant{
		NameFa: venue.NameFa, NameEn: venue.NameEn, DescriptionFa: venue.DescriptionFa,
		DescriptionEn: venue.DescriptionEn, Currency: venue.Currency,
	}, Categories: []Category{}}
	for _, c := range cats {
		out := Category{NameFa: c.NameFa, NameEn: c.NameEn, DescriptionFa: c.DescriptionFa,
			DescriptionEn: c.DescriptionEn, Icon: c.Icon, Visible: c.IsVisible, Items: []Item{}}
		for _, it := range c.Items {
			out.Items = append(out.Items, Item{NameFa: it.NameFa, NameEn: it.NameEn,
				DescriptionFa: it.DescriptionFa, DescriptionEn: it.DescriptionEn, Price: it.Price,
				DiscountPrice: it.DiscountPrice, ImagePath: it.ImagePath, Tags: it.Tags,
				Available: it.IsAvailable, Featured: it.IsFeatured, Calories: it.Calories, PrepMinutes: it.PrepMinutes})
		}
		d.Categories = append(d.Categories, out)
	}
	return d
}

// Validate rejects malformed imports before any existing row can be touched.
func (d Document) Validate() error {
	bad := errors.New("invalid menu: version, names, currency, counts or amounts")
	if d.Version != Version || len(d.Categories) > 200 {
		return bad
	}
	switch d.Restaurant.Currency {
	case "IRT", "IRR", "USD", "EUR", "AED", "TRY", "IQD":
	default:
		return bad
	}
	total := 0
	for _, c := range d.Categories {
		if !namesValid(c.NameFa, c.NameEn) || len(c.DescriptionFa) > 8000 || len(c.DescriptionEn) > 8000 || len(c.Icon) > 64 {
			return bad
		}
		total += len(c.Items)
		if total > 5000 {
			return bad
		}
		for _, it := range c.Items {
			if !namesValid(it.NameFa, it.NameEn) || len(it.DescriptionFa) > 8000 || len(it.DescriptionEn) > 8000 ||
				it.Price < 0 || it.Price > MaxPrice || it.DiscountPrice < 0 ||
				(it.DiscountPrice > 0 && it.DiscountPrice >= it.Price) ||
				it.Calories < 0 || it.Calories > 100000 || it.PrepMinutes < 0 || it.PrepMinutes > 10080 ||
				len(it.Tags) > 20 || len(it.ImagePath) > 500 {
				return bad
			}
			for _, tag := range it.Tags {
				if len(tag) > 40 {
					return bad
				}
			}
		}
	}
	return nil
}
func namesValid(fa, en string) bool {
	return (strings.TrimSpace(fa) != "" || strings.TrimSpace(en) != "") && len(fa) <= 500 && len(en) <= 500
}
