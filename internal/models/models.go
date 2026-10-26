// Package models holds the domain types shared by repositories and handlers.
package models

import (
	"encoding/json"
	"strings"
	"time"
)

type User struct {
	ID             int64
	Username       string
	Email          string
	PasswordHash   string
	DisplayName    string
	Locale         string
	SessionVersion int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// BusinessType drives which starter category set a new venue receives.
type BusinessType string

const (
	BusinessRestaurant BusinessType = "restaurant"
	BusinessFastFood   BusinessType = "fastfood"
	BusinessCafe       BusinessType = "cafe"
	BusinessJigaraki   BusinessType = "jigaraki"
	BusinessTeahouse   BusinessType = "teahouse"
	BusinessBakery     BusinessType = "bakery"
	BusinessIceCream   BusinessType = "icecream"
	BusinessJuiceBar   BusinessType = "juicebar"
)

type Restaurant struct {
	ID           int64
	UserID       int64
	Slug         string
	BusinessType string

	NameFa        string
	NameEn        string
	TaglineFa     string
	TaglineEn     string
	DescriptionFa string
	DescriptionEn string

	LogoPath  string
	CoverPath string

	Phone     string
	AddressFa string
	AddressEn string
	Instagram string
	Telegram  string
	WhatsApp  string
	Website   string

	Currency    string
	IsPublished bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Name returns the localised venue name, falling back to the other language
// when only one has been filled in.
func (r *Restaurant) Name(locale string) string {
	return pick(locale, r.NameFa, r.NameEn)
}

func (r *Restaurant) Tagline(locale string) string {
	return pick(locale, r.TaglineFa, r.TaglineEn)
}

func (r *Restaurant) Description(locale string) string {
	return pick(locale, r.DescriptionFa, r.DescriptionEn)
}

func (r *Restaurant) Address(locale string) string {
	return pick(locale, r.AddressFa, r.AddressEn)
}

// Phones separates multiple contact numbers without splitting formatting spaces.
func (r *Restaurant) Phones() []string {
	parts := strings.FieldsFunc(r.Phone, func(c rune) bool { return strings.ContainsRune("/,;|\n،", c) })
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

type Category struct {
	ID            int64
	RestaurantID  int64
	NameFa        string
	NameEn        string
	DescriptionFa string
	DescriptionEn string
	Icon          string
	ImagePath     string
	SortOrder     int
	IsVisible     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Populated by the repository when a full menu is loaded.
	Items []*Item
}

func (c *Category) Name(locale string) string { return pick(locale, c.NameFa, c.NameEn) }
func (c *Category) Description(locale string) string {
	return pick(locale, c.DescriptionFa, c.DescriptionEn)
}

// VisibleItems filters out items the owner has marked hidden/unavailable for
// the public view.
func (c *Category) VisibleItems() []*Item {
	out := make([]*Item, 0, len(c.Items))
	for _, it := range c.Items {
		if it.IsAvailable {
			out = append(out, it)
		}
	}
	return out
}

type Item struct {
	ID            int64
	RestaurantID  int64
	CategoryID    int64
	NameFa        string
	NameEn        string
	DescriptionFa string
	DescriptionEn string
	PricePending  bool
	Price         int64
	DiscountPrice int64
	ImagePath     string
	Tags          []string
	IsAvailable   bool
	IsFeatured    bool
	Calories      int
	PrepMinutes   int
	SortOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (i *Item) Name(locale string) string { return pick(locale, i.NameFa, i.NameEn) }
func (i *Item) Description(locale string) string {
	return pick(locale, i.DescriptionFa, i.DescriptionEn)
}

// HasDiscount reports whether a genuine reduced price is set.
func (i *Item) HasDiscount() bool {
	return !i.PricePending && i.DiscountPrice > 0 && i.DiscountPrice < i.Price
}

// EffectivePrice is the amount the customer actually pays.
func (i *Item) EffectivePrice() int64 {
	if i.HasDiscount() {
		return i.DiscountPrice
	}
	return i.Price
}

// DiscountPercent returns the rounded percentage saved, or 0 when not discounted.
func (i *Item) DiscountPercent() int {
	if !i.HasDiscount() || i.Price == 0 {
		return 0
	}
	return int((i.Price - i.DiscountPrice) * 100 / i.Price)
}

// TagsJSON serialises tags for storage.
func (i *Item) TagsJSON() string {
	if len(i.Tags) == 0 {
		return "[]"
	}
	b, err := json.Marshal(i.Tags)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// SetTagsJSON parses the stored tag list, tolerating malformed data.
func (i *Item) SetTagsJSON(raw string) {
	i.Tags = nil
	if strings.TrimSpace(raw) == "" {
		return
	}
	_ = json.Unmarshal([]byte(raw), &i.Tags)
}

type WorkingHour struct {
	ID           int64
	RestaurantID int64
	DayOfWeek    int
	OpensAt      string
	ClosesAt     string
	IsClosed     bool
}

type Media struct {
	ID           int64
	RestaurantID int64
	Path         string
	OriginalName string
	SizeBytes    int64
	MimeType     string
	CreatedAt    time.Time
}

// pick chooses the field matching the locale, falling back to whichever value
// is non-empty so a half-translated menu still renders something useful.
func pick(locale, fa, en string) string {
	if locale == "en" {
		if strings.TrimSpace(en) != "" {
			return en
		}
		return fa
	}
	if strings.TrimSpace(fa) != "" {
		return fa
	}
	return en
}
