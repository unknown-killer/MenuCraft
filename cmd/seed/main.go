// Command seed creates a demo restaurant with a full Persian menu, so the app
// can be explored with realistic content straight after a fresh install.
//
//	go run ./cmd/seed
//
// Re-running replaces the demo account's menu; other accounts are untouched.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"ResturantMenue/internal/auth"
	"ResturantMenue/internal/config"
	"ResturantMenue/internal/database"
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
	"ResturantMenue/internal/theme"
)

type seedItem struct {
	fa, en, descFa, descEn string
	price, discount        int64
	tags                   []string
	calories, minutes      int
	featured               bool
}

type seedCategory struct {
	fa, en, icon string
	items        []seedItem
}

var demoMenu = []seedCategory{
	{"پیش‌غذا", "Starters", "🥗", []seedItem{
		{"ماست و خیار", "Mast-o-Khiar", "ماست محلی با خیار، نعنا و گردو", "Local yoghurt with cucumber, mint and walnut", 85_000, 0, []string{"cold", "vegetarian"}, 120, 5, false},
		{"سالاد شیرازی", "Shirazi Salad", "گوجه، خیار و پیاز با آبغوره تازه", "Tomato, cucumber and onion with verjuice", 95_000, 75_000, []string{"vegetarian", "cold"}, 90, 7, false},
		{"کشک بادمجان", "Kashk-e Bademjan", "بادمجان سرخ‌شده با کشک، نعنا داغ و پیاز", "Fried aubergine with whey, mint and onion", 145_000, 0, []string{"popular", "hot"}, 310, 12, true},
		{"میرزا قاسمی", "Mirza Ghasemi", "بادمجان دودی با گوجه، سیر و تخم‌مرغ", "Smoked aubergine with tomato, garlic and egg", 165_000, 0, []string{"vegetarian"}, 280, 15, false},
	}},
	{"کباب‌ها", "Kebabs", "🍢", []seedItem{
		{"کباب کوبیده", "Koobideh Kebab", "دو سیخ کوبیده گوشت گوسفندی با برنج ایرانی", "Two skewers of minced lamb with Persian rice", 385_000, 0, []string{"popular"}, 720, 25, true},
		{"جوجه کباب زعفرانی", "Saffron Chicken Kebab", "سینه مرغ زعفرانی با کره و لیمو", "Saffron chicken breast with butter and lime", 340_000, 295_000, []string{"chef"}, 610, 22, false},
		{"کباب برگ", "Barg Kebab", "فیله گوسفندی زعفرانی، نرم و آبدار", "Saffron lamb fillet, tender and juicy", 520_000, 0, []string{"popular", "chef"}, 680, 30, true},
		{"کباب بختیاری", "Bakhtiari Kebab", "ترکیب فیله گوسفند و جوجه زعفرانی", "Lamb fillet combined with saffron chicken", 495_000, 0, nil, 700, 28, false},
		{"شیشلیک شاندیز", "Shishlik Shandiz", "راسته گوسفندی با استخوان، مخصوص خانه", "Bone-in lamb chops, house speciality", 680_000, 610_000, []string{"chef", "spicy"}, 820, 35, true},
	}},
	{"خورش‌ها", "Stews", "🍲", []seedItem{
		{"قورمه سبزی", "Ghormeh Sabzi", "خورش سنتی با سبزی معطر، لوبیا و گوشت", "Herb stew with beans and lamb", 295_000, 0, []string{"popular"}, 540, 20, false},
		{"قیمه بادمجان", "Gheimeh Bademjan", "خورش لپه با بادمجان سرخ‌شده", "Split-pea stew with fried aubergine", 285_000, 0, nil, 520, 20, false},
		{"فسنجان", "Fesenjan", "خورش گردو و انار با مرغ", "Walnut and pomegranate stew with chicken", 340_000, 0, []string{"chef"}, 620, 25, false},
	}},
	{"غذای دریایی", "Seafood", "🐟", []seedItem{
		{"قزل‌آلای کبابی", "Grilled Trout", "قزل‌آلای تازه با لیمو و سبزیجات", "Fresh trout with lemon and vegetables", 420_000, 0, nil, 450, 25, false},
		{"میگو سوخاری", "Crispy Prawns", "میگوی درشت سوخاری با سس مخصوص", "Large breaded prawns with house sauce", 480_000, 430_000, []string{"new"}, 560, 18, false},
	}},
	{"نوشیدنی‌ها", "Drinks", "🥤", []seedItem{
		{"دوغ محلی", "Local Doogh", "دوغ گازدار با نعنا خشک", "Sparkling yoghurt drink with dried mint", 55_000, 0, []string{"cold"}, 80, 2, false},
		{"شربت سکنجبین", "Sekanjabin Sharbat", "شربت سنتی نعنا و سرکه", "Traditional mint and vinegar cordial", 65_000, 0, []string{"cold", "vegetarian"}, 140, 3, false},
		{"چای زعفرانی", "Saffron Tea", "چای تازه‌دم با زعفران و هل", "Fresh-brewed tea with saffron and cardamom", 45_000, 0, []string{"hot"}, 15, 5, false},
		{"آب‌معدنی", "Mineral Water", "", "", 25_000, 0, nil, 0, 1, false},
	}},
	{"دسر", "Desserts", "🍮", []seedItem{
		{"بستنی سنتی زعفرانی", "Saffron Ice Cream", "بستنی سنتی با خامه و پسته", "Traditional ice cream with cream and pistachio", 120_000, 0, []string{"cold", "popular"}, 320, 3, false},
		{"شله زرد", "Sholeh Zard", "دسر برنج زعفرانی با دارچین", "Saffron rice pudding with cinnamon", 95_000, 0, []string{"vegetarian"}, 260, 4, false},
		{"فالوده شیرازی", "Faloodeh Shirazi", "رشته فالوده با شربت لیمو", "Rice noodle dessert with lime syrup", 110_000, 90_000, []string{"cold", "new"}, 210, 3, false},
	}},
}

func main() {
	username := flag.String("user", "demo", "demo account username")
	password := flag.String("pass", "Demo1234", "demo account password")
	flag.Parse()

	cfg := config.Load()

	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	store := repository.New(db)

	user, err := store.Users.ByLogin(*username)
	if err != nil {
		hash, err := auth.HashPassword(*password)
		if err != nil {
			log.Fatalf("hash password: %v", err)
		}
		user = &models.User{
			Username:     *username,
			Email:        *username + "@example.com",
			PasswordHash: hash,
			DisplayName:  "مدیر نمونه",
			Locale:       "fa",
		}
		if err := store.Users.Create(user); err != nil {
			log.Fatalf("create user: %v", err)
		}
		fmt.Printf("created account  %s / %s\n", *username, *password)
	} else {
		fmt.Printf("reusing account  %s\n", *username)
	}

	venue, err := store.Restaurants.ByUser(user.ID)
	if err != nil {
		venue = &models.Restaurant{
			UserID:       user.ID,
			Slug:         "demo",
			BusinessType: string(models.BusinessRestaurant),
			Currency:     "IRT",
			IsPublished:  true,
		}
		if err := store.Restaurants.Create(venue); err != nil {
			log.Fatalf("create venue: %v", err)
		}
	}

	venue.NameFa = "رستوران سنتی البرز"
	venue.NameEn = "Alborz Traditional Restaurant"
	venue.TaglineFa = "طعم اصیل ایرانی، از سال ۱۳۶۵"
	venue.TaglineEn = "Authentic Persian flavours since 1986"
	venue.DescriptionFa = "رستورانی خانوادگی با غذاهای سنتی ایرانی، کباب تازه و فضایی گرم."
	venue.DescriptionEn = "A family restaurant serving traditional Persian food and fresh kebab."
	venue.AddressFa = "تهران، خیابان ولیعصر، نبش کوچه بهار، پلاک ۱۲"
	venue.AddressEn = "12 Bahar Alley, Valiasr St, Tehran"
	venue.Phone = "02188776655"
	venue.Instagram = "alborz.restaurant"
	venue.IsPublished = true
	if err := store.Restaurants.Update(venue); err != nil {
		log.Fatalf("update venue: %v", err)
	}

	// Start from a clean menu so re-running does not pile up duplicates.
	existing, err := store.Categories.List(venue.ID)
	if err != nil {
		log.Fatalf("list categories: %v", err)
	}
	for _, c := range existing {
		if err := store.Categories.Delete(c.ID, venue.ID); err != nil {
			log.Fatalf("clear category %d: %v", c.ID, err)
		}
	}

	items := 0
	for _, sc := range demoMenu {
		cat := &models.Category{
			RestaurantID: venue.ID,
			NameFa:       sc.fa,
			NameEn:       sc.en,
			Icon:         sc.icon,
			IsVisible:    true,
		}
		if err := store.Categories.Create(cat); err != nil {
			log.Fatalf("create category %q: %v", sc.en, err)
		}

		for _, si := range sc.items {
			it := &models.Item{
				RestaurantID:  venue.ID,
				CategoryID:    cat.ID,
				NameFa:        si.fa,
				NameEn:        si.en,
				DescriptionFa: si.descFa,
				DescriptionEn: si.descEn,
				Price:         si.price,
				DiscountPrice: si.discount,
				Tags:          si.tags,
				IsAvailable:   true,
				IsFeatured:    si.featured,
				Calories:      si.calories,
				PrepMinutes:   si.minutes,
			}
			if err := store.Items.Create(it); err != nil {
				log.Fatalf("create item %q: %v", si.en, err)
			}
			items++
		}
	}

	// A warm preset suits the demo content better than the default dark theme.
	if err := store.Themes.Save(venue.ID, "saffron", theme.ApplyPreset("saffron")); err != nil {
		log.Fatalf("save theme: %v", err)
	}

	hours := make([]*models.WorkingHour, 0, 7)
	for d := 0; d < 7; d++ {
		hours = append(hours, &models.WorkingHour{
			RestaurantID: venue.ID, DayOfWeek: d,
			OpensAt: "12:00", ClosesAt: "23:30",
			IsClosed: d == 6, // closed on Fridays
		})
	}
	if err := store.Restaurants.SaveWorkingHours(venue.ID, hours); err != nil {
		log.Fatalf("save hours: %v", err)
	}

	fmt.Fprintf(os.Stdout,
		"seeded %d categories and %d items\npublic menu:  /m/%s\nsign in with: %s / %s\n",
		len(demoMenu), items, venue.Slug, *username, *password)
}
