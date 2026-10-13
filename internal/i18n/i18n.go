// Package i18n provides the Persian/English catalogues, locale detection and
// the locale-aware number and price formatting used across the app.
package i18n

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Locale describes one supported language.
type Locale struct {
	Code   string // "fa" | "en"
	Dir    string // "rtl" | "ltr"
	Name   string // endonym, shown in the switcher
	HTML   string // BCP-47 tag for the lang attribute
	IsRTL  bool
	Digits bool // render numbers with Persian digits
}

var locales = map[string]Locale{
	"fa": {Code: "fa", Dir: "rtl", Name: "فارسی", HTML: "fa-IR", IsRTL: true, Digits: true},
	"en": {Code: "en", Dir: "ltr", Name: "English", HTML: "en", IsRTL: false, Digits: false},
}

// Supported lists the locales in switcher order.
var Supported = []Locale{locales["fa"], locales["en"]}

// Get returns the locale for a code, defaulting to Persian.
func Get(code string) Locale {
	if l, ok := locales[strings.ToLower(strings.TrimSpace(code))]; ok {
		return l
	}
	return locales["fa"]
}

// Detect resolves the active locale from, in order: an explicit ?lang query,
// the locale cookie, then the Accept-Language header.
func Detect(r *http.Request, fallback string) Locale {
	if v := r.URL.Query().Get("lang"); v != "" {
		if l, ok := locales[strings.ToLower(v)]; ok {
			return l
		}
	}
	if c, err := r.Cookie("menu_locale"); err == nil {
		if l, ok := locales[strings.ToLower(c.Value)]; ok {
			return l
		}
	}
	if h := r.Header.Get("Accept-Language"); h != "" {
		for _, part := range strings.Split(h, ",") {
			tag := strings.ToLower(strings.TrimSpace(strings.SplitN(part, ";", 2)[0]))
			switch {
			case strings.HasPrefix(tag, "fa"):
				return locales["fa"]
			case strings.HasPrefix(tag, "en"):
				return locales["en"]
			}
		}
	}
	return Get(fallback)
}

// SetCookie persists the chosen locale for a year.
func SetCookie(w http.ResponseWriter, code string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "menu_locale",
		Value:    Get(code).Code,
		Path:     "/",
		MaxAge:   365 * 24 * 3600,
		HttpOnly: false,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// T looks a key up in the active catalogue. Missing keys fall back to English,
// then to the key itself, so a gap in translation is visible but never fatal.
func T(locale, key string, args ...any) string {
	cat := catalogFa
	if locale == "en" {
		cat = catalogEn
	}

	s, ok := cat[key]
	if !ok {
		if s, ok = catalogEn[key]; !ok {
			return key
		}
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}

// --- number formatting -----------------------------------------------------

var persianDigits = [...]rune{'۰', '۱', '۲', '۳', '۴', '۵', '۶', '۷', '۸', '۹'}

// ToPersianDigits swaps ASCII digits for their Persian forms.
func ToPersianDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s) * 2)
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(persianDigits[r-'0'])
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// GroupThousands inserts separators every three digits.
func GroupThousands(n int64, sep string) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.FormatInt(n, 10)

	if len(s) <= 3 {
		if neg {
			return "-" + s
		}
		return s
	}

	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)

	out := strings.Join(parts, sep)
	if neg {
		return "-" + out
	}
	return out
}

// Number formats an integer for display in the given locale.
func Number(locale string, n int64) string {
	s := GroupThousands(n, "٬")
	if locale != "en" {
		return ToPersianDigits(s)
	}
	return GroupThousands(n, ",")
}

// currencyLabels maps a currency code to its label in each language.
var currencyLabels = map[string][2]string{
	// code: {fa, en}
	"IRT": {"تومان", "Toman"},
	"IRR": {"ریال", "Rial"},
	"USD": {"دلار", "USD"},
	"EUR": {"یورو", "EUR"},
	"AED": {"درهم", "AED"},
	"TRY": {"لیر", "TRY"},
	"IQD": {"دینار", "IQD"},
}

// Currencies lists the selectable codes in a stable order.
var Currencies = []string{"IRT", "IRR", "USD", "EUR", "AED", "TRY", "IQD"}

// CurrencyLabel returns the localised name of a currency code.
func CurrencyLabel(locale, code string) string {
	pair, ok := currencyLabels[strings.ToUpper(code)]
	if !ok {
		return code
	}
	if locale == "en" {
		return pair[1]
	}
	return pair[0]
}

// Price renders an amount with its currency label, e.g. "۲۵۰٬۰۰۰ تومان".
func Price(locale, currency string, amount int64) string {
	return Number(locale, amount) + " " + CurrencyLabel(locale, currency)
}

// WeekdayNames returns the seven day labels starting at Saturday, matching the
// day_of_week column.
func WeekdayNames(locale string) []string {
	if locale == "en" {
		return []string{"Saturday", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	}
	return []string{"شنبه", "یکشنبه", "دوشنبه", "سه‌شنبه", "چهارشنبه", "پنجشنبه", "جمعه"}
}
