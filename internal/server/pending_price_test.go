package server

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"ResturantMenue/internal/menudata"
	"ResturantMenue/internal/models"
)

func TestPendingPriceRoundTripAndManualPricing(t *testing.T) {
	s, u, v := fixture(t)
	cats, _ := s.store.Categories.List(v.ID)
	form := url.Values{"category_id": {strconv.FormatInt(cats[0].ID, 10)}, "name_fa": {"Ø¨Ø±Ú¯Ø±"}, "name_en": {"Burger"}, "price_pending": {"1"}, "price": {"0"}, "is_available": {"1"}}
	if w := request(s, "POST", "/panel/menu/item", form, session(s, u)); w.Code != 303 {
		t.Fatalf("save: %d", w.Code)
	}
	items, _ := s.store.Items.ListByRestaurant(v.ID)
	var pending *models.Item
	for _, it := range items {
		if it.NameEn == "Burger" {
			pending = it
		}
	}
	if pending == nil || !pending.PricePending {
		t.Fatal("pending state not saved")
	}
	body := request(s, "GET", "/m/"+v.Slug+"?lang=en", nil, nil).Body.String()
	if !strings.Contains(body, `data-note-price="" data-note-pending="1"`) || !strings.Contains(body, "Price coming soon") {
		t.Fatal("public pending price must be blank for notebook and labelled")
	}
	export := request(s, "GET", "/panel/menu/export", nil, session(s, u))
	var doc menudata.Document
	if err := json.Unmarshal(export.Body.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, c := range doc.Categories {
		for _, it := range c.Items {
			if it.NameEn == "Burger" {
				found = it.PricePending
			}
		}
	}
	if !found || doc.Validate() != nil {
		t.Fatal("pending price lost in export")
	}
	if err := s.importMenu(v.ID, doc, false); err != nil {
		t.Fatal(err)
	}
	items, _ = s.store.Items.ListByRestaurant(v.ID)
	count := 0
	for _, it := range items {
		if it.PricePending {
			count++
		}
	}
	if count != 2 {
		t.Fatal("pending price lost on import")
	}
	form.Set("id", strconv.FormatInt(pending.ID, 10))
	form.Del("price_pending")
	form.Set("price", "139500")
	request(s, "POST", "/panel/menu/item", form, session(s, u))
	got, err := s.store.Items.ByID(pending.ID, v.ID)
	if err != nil || got.PricePending || got.EffectivePrice() != 139500 {
		t.Fatal("manual price must clear pending state")
	}
}
