// Package menutemplate holds the starter category sets offered per venue type,
// so a new owner can populate a realistic menu skeleton in one click.
package menutemplate

// Category is one starter category.
type Category struct {
	NameFa string
	NameEn string
	Icon   string
}

// Template is a named group of starter categories.
type Template struct {
	Key        string
	LabelFa    string
	LabelEn    string
	NoteFa     string
	NoteEn     string
	Icon       string
	Categories []Category
}

var templates = []Template{
	{
		Key: "restaurant", LabelFa: "رستوران ایرانی", LabelEn: "Persian restaurant",
		NoteFa: "کباب، خورش و غذاهای اصیل ایرانی", NoteEn: "Kebabs, stews and Persian classics",
		Icon: "🍽️",
		Categories: []Category{
			{"پیش‌غذا", "Starters", "🥗"},
			{"کباب‌ها", "Kebabs", "🍢"},
			{"خورش‌ها", "Stews", "🍲"},
			{"چلو و پلو", "Rice dishes", "🍚"},
			{"غذای دریایی", "Seafood", "🐟"},
			{"سالاد و مخلفات", "Salads & sides", "🥙"},
			{"نوشیدنی‌ها", "Drinks", "🥤"},
			{"دسر", "Desserts", "🍮"},
		},
	},
	{
		Key: "western", LabelFa: "غذای فرنگی", LabelEn: "Western food",
		NoteFa: "استیک، پاستا و غذاهای بین‌المللی", NoteEn: "Steaks, pasta and international dishes",
		Icon: "🍝",
		Categories: []Category{
			{"پیش‌غذا", "Appetisers", "🥖"},
			{"سوپ و سالاد", "Soups & salads", "🥣"},
			{"استیک و گریل", "Steak & grill", "🥩"},
			{"پاستا و ریزوتو", "Pasta & risotto", "🍝"},
			{"غذای دریایی", "Seafood", "🦐"},
			{"دسر", "Desserts", "🍰"},
			{"نوشیدنی‌ها", "Drinks", "🍹"},
		},
	},
	{
		Key: "fastfood", LabelFa: "فست‌فود", LabelEn: "Fast food",
		NoteFa: "برگر، پیتزا، ساندویچ و سوخاری", NoteEn: "Burgers, pizza, sandwiches and fried chicken",
		Icon: "🍔",
		Categories: []Category{
			{"برگرها", "Burgers", "🍔"},
			{"پیتزاها", "Pizzas", "🍕"},
			{"ساندویچ‌ها", "Sandwiches", "🥪"},
			{"سوخاری", "Fried chicken", "🍗"},
			{"هات‌داگ", "Hot dogs", "🌭"},
			{"سیب‌زمینی و پیش‌غذا", "Fries & sides", "🍟"},
			{"سالاد", "Salads", "🥗"},
			{"نوشیدنی‌ها", "Drinks", "🥤"},
			{"دسر", "Desserts", "🍦"},
		},
	},
	{
		Key: "cafe", LabelFa: "کافه", LabelEn: "Café",
		NoteFa: "قهوه، دمنوش، کیک و صبحانه", NoteEn: "Coffee, tea, cake and breakfast",
		Icon: "☕",
		Categories: []Category{
			{"قهوه گرم", "Hot coffee", "☕"},
			{"قهوه سرد", "Iced coffee", "🧊"},
			{"چای و دمنوش", "Tea & herbal", "🫖"},
			{"شیک و اسموتی", "Shakes & smoothies", "🥤"},
			{"ماکتیل", "Mocktails", "🍹"},
			{"صبحانه", "Breakfast", "🍳"},
			{"کیک و دسر", "Cakes & desserts", "🍰"},
			{"ساندویچ و سالاد", "Sandwiches & salads", "🥪"},
		},
	},
	{
		Key: "jigaraki", LabelFa: "جگرکی", LabelEn: "Grill house",
		NoteFa: "جگر، دل و قلوه و کباب", NoteEn: "Liver, offal and kebab",
		Icon: "🔥",
		Categories: []Category{
			{"جگر و دل و قلوه", "Liver & offal", "🔥"},
			{"کباب‌ها", "Kebabs", "🍢"},
			{"ساندویچ", "Sandwiches", "🥙"},
			{"مخلفات", "Sides", "🧅"},
			{"سالاد و ترشی", "Salad & pickles", "🥗"},
			{"نوشیدنی‌ها", "Drinks", "🥤"},
		},
	},
	{
		Key: "teahouse", LabelFa: "سفره‌خانه و قهوه‌خانه", LabelEn: "Teahouse",
		NoteFa: "دیزی، قلیان، چای و دمنوش", NoteEn: "Dizi, hookah, tea and herbal infusions",
		Icon: "🫖",
		Categories: []Category{
			{"چای و دمنوش", "Tea & herbal", "🫖"},
			{"قلیان", "Hookah", "💨"},
			{"دیزی و آبگوشت", "Dizi & abgoosht", "🍲"},
			{"غذای سنتی", "Traditional dishes", "🍚"},
			{"پیش‌غذا و مخلفات", "Starters & sides", "🥗"},
			{"صبحانه سنتی", "Traditional breakfast", "🍳"},
			{"نوشیدنی سرد", "Cold drinks", "🧊"},
			{"دسر سنتی", "Traditional desserts", "🍮"},
		},
	},
	{
		Key: "bakery", LabelFa: "نانوایی و شیرینی", LabelEn: "Bakery & pastry",
		NoteFa: "نان، کیک و شیرینی", NoteEn: "Bread, cakes and pastries",
		Icon: "🥐",
		Categories: []Category{
			{"نان‌ها", "Breads", "🥖"},
			{"شیرینی تر", "Cream pastries", "🧁"},
			{"شیرینی خشک", "Dry pastries", "🍪"},
			{"کیک", "Cakes", "🎂"},
			{"دسر", "Desserts", "🍮"},
			{"نوشیدنی‌ها", "Drinks", "☕"},
		},
	},
	{
		Key: "icecream", LabelFa: "بستنی و آبمیوه", LabelEn: "Ice cream & juice",
		NoteFa: "بستنی، فالوده، معجون و آبمیوه", NoteEn: "Ice cream, faloodeh, juices and blends",
		Icon: "🍨",
		Categories: []Category{
			{"بستنی", "Ice cream", "🍨"},
			{"فالوده", "Faloodeh", "🍧"},
			{"معجون", "Majoon blends", "🥛"},
			{"آبمیوه طبیعی", "Fresh juice", "🧃"},
			{"شیک و اسموتی", "Shakes & smoothies", "🥤"},
			{"دسر سرد", "Cold desserts", "🍮"},
		},
	},

	// --- add-on packs: appended to an existing menu ---
	{
		Key: "pack_drinks", LabelFa: "بسته نوشیدنی‌ها", LabelEn: "Drinks pack",
		NoteFa: "فقط دسته‌بندی‌های نوشیدنی", NoteEn: "Drink categories only",
		Icon: "🥤",
		Categories: []Category{
			{"نوشابه و دلستر", "Soft drinks", "🥤"},
			{"آب‌معدنی", "Water", "💧"},
			{"دوغ و عرقیجات", "Doogh & distillates", "🥛"},
			{"آبمیوه طبیعی", "Fresh juice", "🧃"},
			{"نوشیدنی گرم", "Hot drinks", "☕"},
		},
	},
	{
		Key: "pack_hookah", LabelFa: "بسته قلیان", LabelEn: "Hookah pack",
		NoteFa: "طعم‌ها و سرویس قلیان", NoteEn: "Hookah flavours and service",
		Icon: "💨",
		Categories: []Category{
			{"قلیان میوه‌ای", "Fruit hookah", "🍓"},
			{"قلیان ویژه", "Premium hookah", "💨"},
			{"سرویس و شارژ", "Service & refill", "♨️"},
		},
	},
	{
		Key: "pack_herbal", LabelFa: "بسته دمنوش", LabelEn: "Herbal tea pack",
		NoteFa: "چای و دمنوش‌های گیاهی", NoteEn: "Teas and herbal infusions",
		Icon: "🌿",
		Categories: []Category{
			{"چای", "Tea", "🫖"},
			{"دمنوش گیاهی", "Herbal infusions", "🌿"},
			{"دمنوش سرد", "Iced infusions", "🧊"},
		},
	},
	{
		Key: "pack_breakfast", LabelFa: "بسته صبحانه", LabelEn: "Breakfast pack",
		NoteFa: "صبحانه ایرانی و فرنگی", NoteEn: "Persian and continental breakfast",
		Icon: "🍳",
		Categories: []Category{
			{"صبحانه ایرانی", "Persian breakfast", "🫓"},
			{"املت و تخم‌مرغ", "Eggs & omelettes", "🍳"},
			{"صبحانه فرنگی", "Continental breakfast", "🥐"},
		},
	},
}

// All returns every template, venue sets first then add-on packs.
func All() []Template { return templates }

// ByKey looks up one template.
func ByKey(key string) (Template, bool) {
	for _, t := range templates {
		if t.Key == key {
			return t, true
		}
	}
	return Template{}, false
}

// ForBusiness returns the starter set matching a venue type, used to seed a
// menu right after registration.
func ForBusiness(businessType string) Template {
	if t, ok := ByKey(businessType); ok {
		return t
	}
	t, _ := ByKey("restaurant")
	return t
}
