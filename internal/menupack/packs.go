package menupack

// Pack is a complete ready-made menu: which categories to create, which dishes
// go in them, a matching theme, and suggested venue wording.
//
// Tier scales every price, which is what makes one dish catalogue serve a
// street kebab shop and a fine-dining room from the same data.
type Pack struct {
	Key     string
	LabelFa string
	LabelEn string
	NoteFa  string
	NoteEn  string

	Business string
	Theme    string  // theme preset key
	Tier     float64 // price multiplier

	Groups []string

	NameFa    string
	NameEn    string
	TaglineFa string
	TaglineEn string
}

// packs is the whole library, ordered by venue type.
var packs = []Pack{

	// ============================================================ restaurant
	{
		Key: "r-sonnati", LabelFa: "رستوران سنتی", LabelEn: "Traditional Persian",
		NoteFa: "کباب، خورش و پیش‌غذای اصیل", NoteEn: "Kebab, stew and classic starters",
		Business: "restaurant", Theme: "saffron", Tier: 1.0,
		Groups: []string{"starters_ir", "kebabs", "stews", "rice_ir", "drinks_cold", "desserts_ir"},
		NameFa: "سفره‌خانه سنتی", NameEn: "Traditional House",
		TaglineFa: "طعم اصیل ایرانی", TaglineEn: "Authentic Persian flavour",
	},
	{
		Key: "r-kababi", LabelFa: "کبابی", LabelEn: "Kebab House",
		NoteFa: "تمرکز روی کباب و مخلفات", NoteEn: "Focused on kebab and sides",
		Business: "restaurant", Theme: "ember", Tier: 0.9,
		Groups: []string{"kebabs", "grill_sides", "starters_ir", "drinks_cold"},
		NameFa: "کبابسرا", NameEn: "Kebab House",
		TaglineFa: "روی آتش، تازه و داغ", TaglineEn: "Fresh off the fire",
	},
	{
		Key: "r-luxe", LabelFa: "رستوران لوکس", LabelEn: "Fine Dining",
		NoteFa: "منوی کامل با قیمت بالا", NoteEn: "Full menu, premium pricing",
		Business: "restaurant", Theme: "royal", Tier: 1.6,
		Groups: []string{"starters_ir", "kebabs", "stews", "seafood", "steaks", "desserts_ir", "drinks_cold"},
		NameFa: "رستوران شهرزاد", NameEn: "Shahrzad Restaurant",
		TaglineFa: "تجربه‌ای متفاوت از مهمان‌نوازی ایرانی", TaglineEn: "Persian hospitality, elevated",
	},
	{
		Key: "r-daryaee", LabelFa: "رستوران دریایی", LabelEn: "Seafood Restaurant",
		NoteFa: "غذای دریایی و سالاد", NoteEn: "Seafood and salads",
		Business: "restaurant", Theme: "coral", Tier: 1.25,
		Groups: []string{"seafood", "salads", "starters_ir", "rice_ir", "drinks_cold"},
		NameFa: "رستوران ساحل", NameEn: "Sahel Seafood",
		TaglineFa: "تازه از دریای جنوب", TaglineEn: "Fresh from the southern sea",
	},
	{
		Key: "r-farangi", LabelFa: "رستوران فرنگی", LabelEn: "Western Restaurant",
		NoteFa: "استیک، پاستا و سالاد", NoteEn: "Steak, pasta and salad",
		Business: "restaurant", Theme: "mocha", Tier: 1.3,
		Groups: []string{"salads", "pasta", "steaks", "cakes", "drinks_cold"},
		NameFa: "بیسترو وست", NameEn: "West Bistro",
		TaglineFa: "آشپزی اروپایی در قلب شهر", TaglineEn: "European cooking downtown",
	},
	{
		Key: "r-khanegi", LabelFa: "غذای خانگی", LabelEn: "Home Cooking",
		NoteFa: "خورش و پلو، قیمت اقتصادی", NoteEn: "Stews and rice, budget pricing",
		Business: "restaurant", Theme: "clay", Tier: 0.75,
		Groups: []string{"stews", "rice_ir", "starters_ir", "drinks_cold"},
		NameFa: "غذای خانگی مادر", NameEn: "Mother's Kitchen",
		TaglineFa: "مثل دست‌پخت خانه", TaglineEn: "Just like home",
	},
	{
		Key: "r-modern", LabelFa: "رستوران مدرن", LabelEn: "Modern Restaurant",
		NoteFa: "ترکیب ایرانی و بین‌المللی", NoteEn: "Persian meets international",
		Business: "restaurant", Theme: "aurora", Tier: 1.35,
		Groups: []string{"starters_ir", "steaks", "pasta", "kebabs", "cakes", "drinks_cold"},
		NameFa: "رستوران نُوا", NameEn: "Nova Kitchen",
		TaglineFa: "آشپزی امروز، ریشه‌های دیروز", TaglineEn: "Today's cooking, yesterday's roots",
	},
	{
		Key: "r-emerald", LabelFa: "رستوران باغی", LabelEn: "Garden Restaurant",
		NoteFa: "فضای سبز، منوی کامل", NoteEn: "Garden setting, full menu",
		Business: "restaurant", Theme: "emerald", Tier: 1.15,
		Groups: []string{"starters_ir", "kebabs", "stews", "salads", "desserts_ir", "tea_herbal"},
		NameFa: "باغ رستوران چنار", NameEn: "Chenar Garden",
		TaglineFa: "زیر سایه چنارها", TaglineEn: "Under the plane trees",
	},
	{
		Key: "r-paper", LabelFa: "منوی چاپی", LabelEn: "Printed Menu",
		NoteFa: "بدون عکس، شبیه منوی کاغذی", NoteEn: "No photos, like a paper menu",
		Business: "restaurant", Theme: "paper", Tier: 1.1,
		Groups: []string{"starters_ir", "kebabs", "stews", "rice_ir", "desserts_ir", "drinks_cold"},
		NameFa: "رستوران البرز", NameEn: "Alborz Restaurant",
		TaglineFa: "از سال ۱۳۶۵", TaglineEn: "Since 1986",
	},
	{
		Key: "r-velvet", LabelFa: "رستوران شبانه", LabelEn: "Night Restaurant",
		NoteFa: "فضای تاریک و شیک", NoteEn: "Dark and elegant",
		Business: "restaurant", Theme: "velvet", Tier: 1.45,
		Groups: []string{"starters_ir", "steaks", "seafood", "pasta", "cakes", "drinks_cold"},
		NameFa: "رستوران مخمل", NameEn: "Velvet Room",
		TaglineFa: "شام‌های به‌یادماندنی", TaglineEn: "Dinners worth remembering",
	},

	// ============================================================= fast food
	{
		Key: "f-burger", LabelFa: "برگرخانه", LabelEn: "Burger Joint",
		NoteFa: "برگر، سیب‌زمینی و نوشیدنی", NoteEn: "Burgers, fries and drinks",
		Business: "fastfood", Theme: "burger", Tier: 1.0,
		Groups: []string{"burgers", "sides_ff", "sandwiches", "drinks_cold"},
		NameFa: "برگرخانه", NameEn: "The Burger Joint",
		TaglineFa: "دست‌ساز، هر روز", TaglineEn: "Handmade, every day",
	},
	{
		Key: "f-pizza", LabelFa: "پیتزا فروشی", LabelEn: "Pizzeria",
		NoteFa: "پیتزا با خمیر تازه", NoteEn: "Pizza with fresh dough",
		Business: "fastfood", Theme: "pizzeria", Tier: 1.0,
		Groups: []string{"pizzas", "sides_ff", "salads", "drinks_cold"},
		NameFa: "پیتزا ناپولی", NameEn: "Napoli Pizza",
		TaglineFa: "خمیر تازه، هر روز صبح", TaglineEn: "Dough made fresh each morning",
	},
	{
		Key: "f-sookhari", LabelFa: "سوخاری", LabelEn: "Fried Chicken",
		NoteFa: "مرغ سوخاری و استریپس", NoteEn: "Fried chicken and strips",
		Business: "fastfood", Theme: "sunset", Tier: 1.0,
		Groups: []string{"fried", "sides_ff", "burgers", "drinks_cold"},
		NameFa: "سوخاری طلایی", NameEn: "Golden Crisp",
		TaglineFa: "ترد و طلایی", TaglineEn: "Crispy and golden",
	},
	{
		Key: "f-sandwich", LabelFa: "ساندویچی", LabelEn: "Sandwich Shop",
		NoteFa: "ساندویچ سرد و گرم، اقتصادی", NoteEn: "Hot and cold sandwiches, budget",
		Business: "fastfood", Theme: "charcoal", Tier: 0.8,
		Groups: []string{"sandwiches", "sides_ff", "drinks_cold"},
		NameFa: "ساندویچ شهر", NameEn: "City Sandwich",
		TaglineFa: "سریع، ساده، خوشمزه", TaglineEn: "Quick, simple, tasty",
	},
	{
		Key: "f-arcade", LabelFa: "فست‌فود رترو", LabelEn: "Retro Fast Food",
		NoteFa: "رنگارنگ و پرانرژی", NoteEn: "Colourful and loud",
		Business: "fastfood", Theme: "arcade", Tier: 1.1,
		Groups: []string{"burgers", "pizzas", "fried", "sides_ff", "shakes", "drinks_cold"},
		NameFa: "فست‌فود آرکید", NameEn: "Arcade Eats",
		TaglineFa: "بازی، برگر، شیک", TaglineEn: "Play, burger, shake",
	},
	{
		Key: "f-neon", LabelFa: "فست‌فود شبانه", LabelEn: "Late Night",
		NoteFa: "تیره و نئونی، باز تا دیروقت", NoteEn: "Dark and neon, open late",
		Business: "fastfood", Theme: "neon", Tier: 1.05,
		Groups: []string{"burgers", "sandwiches", "fried", "sides_ff", "drinks_cold"},
		NameFa: "نایت بایت", NameEn: "Night Bite",
		TaglineFa: "تا آخرین مشتری", TaglineEn: "Open till the last guest",
	},
	{
		Key: "f-full", LabelFa: "فست‌فود کامل", LabelEn: "Full Fast Food",
		NoteFa: "همه‌چیز در یک منو", NoteEn: "Everything on one menu",
		Business: "fastfood", Theme: "fastfood-showcase-deep", Tier: 1.0,
		Groups: []string{"burgers", "pizzas", "fried", "sandwiches", "sides_ff", "salads", "shakes", "drinks_cold"},
		NameFa: "فست‌فود مرکزی", NameEn: "Central Fast Food",
		TaglineFa: "هر چی بخوای، اینجاست", TaglineEn: "Whatever you fancy",
	},
	{
		Key: "f-veggie", LabelFa: "فست‌فود گیاهی", LabelEn: "Veggie Fast Food",
		NoteFa: "گزینه‌های گیاهی و سالم", NoteEn: "Plant-based and lighter options",
		Business: "fastfood", Theme: "pistachio", Tier: 1.05,
		Groups: []string{"burgers", "sandwiches", "salads", "sides_ff", "juices"},
		NameFa: "سبزخوان", NameEn: "Green Bite",
		TaglineFa: "گیاهی، ولی سیرکننده", TaglineEn: "Plant-based and filling",
	},
	{
		Key: "f-tiles", LabelFa: "فست‌فود کاشی", LabelEn: "Tile Grid",
		NoteFa: "چیدمان کاشی‌وار با عکس مربعی", NoteEn: "Square-tile layout",
		Business: "fastfood", Theme: "fastfood-tiles-cream", Tier: 0.95,
		Groups: []string{"burgers", "pizzas", "fried", "sides_ff", "drinks_cold"},
		NameFa: "فست‌فود کاشی", NameEn: "Tile Kitchen",
		TaglineFa: "انتخاب سریع از روی عکس", TaglineEn: "Pick it straight off the picture",
	},
	{
		Key: "f-brutal", LabelFa: "فست‌فود بروتال", LabelEn: "Brutalist",
		NoteFa: "سیاه و سفید با تاکید نارنجی", NoteEn: "Black and white with orange",
		Business: "fastfood", Theme: "brutalist", Tier: 1.0,
		Groups: []string{"burgers", "sandwiches", "sides_ff", "drinks_cold"},
		NameFa: "برگر بلاک", NameEn: "Burger Block",
		TaglineFa: "بدون حاشیه، فقط برگر", TaglineEn: "No fuss, just burgers",
	},

	// ================================================================== café
	{
		Key: "c-specialty", LabelFa: "کافه تخصصی", LabelEn: "Specialty Coffee",
		NoteFa: "قهوه تخصصی، منوی کوتاه", NoteEn: "Specialty coffee, short menu",
		Business: "cafe", Theme: "espresso", Tier: 1.15,
		Groups: []string{"coffee_hot", "coffee_cold", "cakes"},
		NameFa: "کافه دم", NameEn: "Dam Coffee",
		TaglineFa: "قهوه، درست دم‌شده", TaglineEn: "Coffee, properly brewed",
	},
	{
		Key: "c-brunch", LabelFa: "کافه صبحانه", LabelEn: "Brunch Café",
		NoteFa: "صبحانه و قهوه", NoteEn: "Breakfast and coffee",
		Business: "cafe", Theme: "latte", Tier: 1.1,
		Groups: []string{"breakfast", "coffee_hot", "coffee_cold", "cakes", "juices"},
		NameFa: "کافه صبح", NameEn: "Morning Café",
		TaglineFa: "روزت را اینجا شروع کن", TaglineEn: "Start your day here",
	},
	{
		Key: "c-dessert", LabelFa: "کافه دسر", LabelEn: "Dessert Café",
		NoteFa: "کیک، شیک و قهوه", NoteEn: "Cake, shakes and coffee",
		Business: "cafe", Theme: "vanilla", Tier: 1.1,
		Groups: []string{"cakes", "shakes", "coffee_hot", "coffee_cold"},
		NameFa: "کافه شیرین", NameEn: "Sweet Café",
		TaglineFa: "بهانه‌ای برای شیرینی", TaglineEn: "An excuse for something sweet",
	},
	{
		Key: "c-matcha", LabelFa: "کافه سبز", LabelEn: "Green Café",
		NoteFa: "دمنوش، ماچا و گزینه سالم", NoteEn: "Infusions, matcha and lighter options",
		Business: "cafe", Theme: "matcha", Tier: 1.1,
		Groups: []string{"tea_herbal", "coffee_hot", "shakes", "salads", "cakes"},
		NameFa: "کافه برگ", NameEn: "Leaf Café",
		TaglineFa: "آرام، سبز، تازه", TaglineEn: "Calm, green, fresh",
	},
	{
		Key: "c-book", LabelFa: "کافه کتاب", LabelEn: "Book Café",
		NoteFa: "منوی آرام با فضای مطالعه", NoteEn: "Quiet menu for a reading room",
		Business: "cafe", Theme: "linen", Tier: 1.05,
		Groups: []string{"coffee_hot", "tea_herbal", "cakes", "breakfast"},
		NameFa: "کافه کتاب", NameEn: "The Book Café",
		TaglineFa: "یک فنجان و یک فصل", TaglineEn: "A cup and a chapter",
	},
	{
		Key: "c-night", LabelFa: "کافه شبانه", LabelEn: "Night Café",
		NoteFa: "تیره و گرم، باز تا دیروقت", NoteEn: "Dark and warm, open late",
		Business: "cafe", Theme: "nightgarden", Tier: 1.15,
		Groups: []string{"coffee_hot", "coffee_cold", "cakes", "shakes", "tea_herbal"},
		NameFa: "کافه شب", NameEn: "Night Café",
		TaglineFa: "شب‌های طولانی، قهوه‌های گرم", TaglineEn: "Long nights, warm coffee",
	},
	{
		Key: "c-minimal", LabelFa: "کافه مینیمال", LabelEn: "Minimal Café",
		NoteFa: "منوی خیلی ساده و تمیز", NoteEn: "Very clean, very simple",
		Business: "cafe", Theme: "monoline", Tier: 1.1,
		Groups: []string{"coffee_hot", "coffee_cold", "cakes"},
		NameFa: "کافه خط", NameEn: "Line Café",
		TaglineFa: "کمتر، بهتر", TaglineEn: "Less, better",
	},
	{
		Key: "c-full", LabelFa: "کافه رستوران", LabelEn: "Café Restaurant",
		NoteFa: "قهوه به‌علاوه غذای کامل", NoteEn: "Coffee plus a full kitchen",
		Business: "cafe", Theme: "mocha", Tier: 1.2,
		Groups: []string{"breakfast", "salads", "pasta", "sandwiches", "coffee_hot", "coffee_cold", "cakes"},
		NameFa: "کافه رستوران وینتیج", NameEn: "Vintage Café",
		TaglineFa: "از صبحانه تا شام", TaglineEn: "Breakfast through dinner",
	},
	{
		Key: "c-sorbet", LabelFa: "کافه پاستلی", LabelEn: "Pastel Café",
		NoteFa: "رنگ‌های شاد و شیک", NoteEn: "Bright pastel look",
		Business: "cafe", Theme: "sorbet", Tier: 1.05,
		Groups: []string{"shakes", "coffee_cold", "cakes", "icecream", "juices"},
		NameFa: "کافه قندی", NameEn: "Candy Café",
		TaglineFa: "شیرین و خنک", TaglineEn: "Sweet and cool",
	},
	{
		Key: "c-editorial", LabelFa: "کافه مجله‌ای", LabelEn: "Editorial Café",
		NoteFa: "چیدمان مجله‌ای با قلم سریف", NoteEn: "Magazine layout with serif type",
		Business: "cafe", Theme: "cafe-editorial-airy", Tier: 1.1,
		Groups: []string{"coffee_hot", "breakfast", "cakes", "tea_herbal"},
		NameFa: "کافه نشریه", NameEn: "The Editorial",
		TaglineFa: "قهوه با یک داستان", TaglineEn: "Coffee with a story",
	},

	// ============================================================= jigaraki
	{
		Key: "j-classic", LabelFa: "جگرکی سنتی", LabelEn: "Classic Grill House",
		NoteFa: "جگر، دل و قلوه و مخلفات", NoteEn: "Liver, offal and sides",
		Business: "jigaraki", Theme: "ember", Tier: 1.0,
		Groups: []string{"liver_grill", "grill_sides", "drinks_cold"},
		NameFa: "جگرکی برادران", NameEn: "Brothers Grill",
		TaglineFa: "روی ذغال، لحظه‌ای", TaglineEn: "Straight off the coals",
	},
	{
		Key: "j-kebab", LabelFa: "جگرکی و کبابی", LabelEn: "Grill & Kebab",
		NoteFa: "جگر به‌علاوه کباب کامل", NoteEn: "Offal plus a full kebab menu",
		Business: "jigaraki", Theme: "jigaraki-classic-cream", Tier: 1.0,
		Groups: []string{"liver_grill", "kebabs", "grill_sides", "starters_ir", "drinks_cold"},
		NameFa: "جگر و کباب شعله", NameEn: "Shole Grill",
		TaglineFa: "از جگر تا شیشلیک", TaglineEn: "From liver to shishlik",
	},
	{
		Key: "j-street", LabelFa: "جگرکی خیابانی", LabelEn: "Street Grill",
		NoteFa: "منوی کوتاه و ارزان", NoteEn: "Short, cheap menu",
		Business: "jigaraki", Theme: "obsidian", Tier: 0.7,
		Groups: []string{"liver_grill", "grill_sides", "drinks_cold"},
		NameFa: "جگرکی گوشه میدان", NameEn: "Corner Grill",
		TaglineFa: "سیخی، ایستاده، عالی", TaglineEn: "By the skewer, standing up",
	},
	{
		Key: "j-modern", LabelFa: "جگرکی مدرن", LabelEn: "Modern Grill",
		NoteFa: "ظاهر امروزی برای غذای سنتی", NoteEn: "A modern face on traditional food",
		Business: "jigaraki", Theme: "charcoal", Tier: 1.1,
		Groups: []string{"liver_grill", "kebabs", "grill_sides", "salads", "drinks_cold"},
		NameFa: "جگرکی زغال", NameEn: "Charcoal",
		TaglineFa: "سنتی، با نگاه امروز", TaglineEn: "Traditional, seen today",
	},
	{
		Key: "j-family", LabelFa: "جگرکی خانوادگی", LabelEn: "Family Grill",
		NoteFa: "منوی کامل برای خانواده", NoteEn: "Full family menu",
		Business: "jigaraki", Theme: "clay", Tier: 0.95,
		Groups: []string{"liver_grill", "kebabs", "rice_ir", "grill_sides", "starters_ir", "drinks_cold", "desserts_ir"},
		NameFa: "جگرکی خانواده", NameEn: "Family Grill",
		TaglineFa: "جای همه هست", TaglineEn: "Room for everyone",
	},

	// ============================================================= teahouse
	{
		Key: "t-sofrekhane", LabelFa: "سفره‌خانه سنتی", LabelEn: "Traditional Teahouse",
		NoteFa: "دیزی، قلیان و چای", NoteEn: "Dizi, hookah and tea",
		Business: "teahouse", Theme: "clay", Tier: 1.0,
		Groups: []string{"dizi", "tea_herbal", "hookah", "starters_ir", "drinks_cold"},
		NameFa: "سفره‌خانه سنتی", NameEn: "Traditional Teahouse",
		TaglineFa: "دیزی، چای، قلیان", TaglineEn: "Dizi, tea and hookah",
	},
	{
		Key: "t-ghahvekhane", LabelFa: "قهوه‌خانه", LabelEn: "Old Coffee House",
		NoteFa: "چای، صبحانه سنتی و قلیان", NoteEn: "Tea, traditional breakfast and hookah",
		Business: "teahouse", Theme: "turquoise", Tier: 0.9,
		Groups: []string{"tea_herbal", "breakfast", "hookah", "dizi"},
		NameFa: "قهوه‌خانه قدیم", NameEn: "The Old Coffee House",
		TaglineFa: "از صبح تا شب، چای دم است", TaglineEn: "The tea is always brewing",
	},
	{
		Key: "t-garden", LabelFa: "سفره‌خانه باغی", LabelEn: "Garden Teahouse",
		NoteFa: "فضای باز با منوی کامل", NoteEn: "Open air with a full menu",
		Business: "teahouse", Theme: "nightgarden", Tier: 1.15,
		Groups: []string{"dizi", "kebabs", "tea_herbal", "hookah", "starters_ir", "desserts_ir"},
		NameFa: "باغ سفره‌خانه", NameEn: "Garden Teahouse",
		TaglineFa: "زیر آسمان، کنار حوض", TaglineEn: "Under the sky, by the pool",
	},
	{
		Key: "t-hookah", LabelFa: "قلیان‌سرا", LabelEn: "Hookah Lounge",
		NoteFa: "تمرکز روی قلیان و نوشیدنی", NoteEn: "Hookah and drinks focused",
		Business: "teahouse", Theme: "royal", Tier: 1.2,
		Groups: []string{"hookah", "tea_herbal", "drinks_cold", "shakes", "cakes"},
		NameFa: "لانژ دود", NameEn: "Smoke Lounge",
		TaglineFa: "شب‌های آرام", TaglineEn: "Easy evenings",
	},
	{
		Key: "t-poster", LabelFa: "سفره‌خانه پوستری", LabelEn: "Poster Teahouse",
		NoteFa: "بدون عکس، تایپوگرافی درشت", NoteEn: "No photos, big type",
		Business: "teahouse", Theme: "teahouse-poster-deep", Tier: 1.0,
		Groups: []string{"dizi", "tea_herbal", "hookah", "breakfast"},
		NameFa: "سفره‌خانه یاران", NameEn: "Yaran Teahouse",
		TaglineFa: "ساده و صمیمی", TaglineEn: "Simple and warm",
	},

	// =============================================================== bakery
	{
		Key: "b-nanvaei", LabelFa: "نانوایی", LabelEn: "Bakery",
		NoteFa: "نان تازه و صبحانه", NoteEn: "Fresh bread and breakfast",
		Business: "bakery", Theme: "sandstone", Tier: 0.8,
		Groups: []string{"breads", "pastries", "tea_herbal"},
		NameFa: "نانوایی سنگک", NameEn: "Sangak Bakery",
		TaglineFa: "هر روز، از تنور", TaglineEn: "From the oven, every day",
	},
	{
		Key: "b-shirini", LabelFa: "شیرینی‌فروشی", LabelEn: "Pastry Shop",
		NoteFa: "شیرینی تر و خشک", NoteEn: "Cream and dry pastries",
		Business: "bakery", Theme: "vanilla", Tier: 1.05,
		Groups: []string{"pastries", "cakes", "coffee_hot", "tea_herbal"},
		NameFa: "شیرینی‌سرای گل", NameEn: "Gol Pastry",
		TaglineFa: "شیرینی خانگی، هر روز تازه", TaglineEn: "Homemade, baked daily",
	},
	{
		Key: "b-cafebakery", LabelFa: "کافه نانوایی", LabelEn: "Café Bakery",
		NoteFa: "نان، قهوه و صبحانه", NoteEn: "Bread, coffee and breakfast",
		Business: "bakery", Theme: "latte", Tier: 1.1,
		Groups: []string{"breads", "pastries", "breakfast", "coffee_hot", "coffee_cold"},
		NameFa: "نان و قهوه", NameEn: "Bread & Coffee",
		TaglineFa: "بوی نان تازه و قهوه", TaglineEn: "Fresh bread and coffee",
	},
	{
		Key: "b-cake", LabelFa: "کیک‌فروشی", LabelEn: "Cake Shop",
		NoteFa: "کیک سفارشی و دسر", NoteEn: "Custom cakes and desserts",
		Business: "bakery", Theme: "sorbet", Tier: 1.15,
		Groups: []string{"cakes", "pastries", "shakes", "coffee_cold"},
		NameFa: "کیک‌سرا", NameEn: "The Cake House",
		TaglineFa: "برای هر جشنی", TaglineEn: "For every celebration",
	},
	{
		Key: "b-gallery", LabelFa: "شیرینی گالری", LabelEn: "Gallery Bakery",
		NoteFa: "چیدمان عکس‌محور", NoteEn: "Photo-led layout",
		Business: "bakery", Theme: "bakery-gallery-airy", Tier: 1.0,
		Groups: []string{"pastries", "cakes", "breads", "coffee_hot"},
		NameFa: "شیرینی گالری", NameEn: "Pastry Gallery",
		TaglineFa: "با چشم انتخاب کن", TaglineEn: "Choose with your eyes",
	},

	// ============================================================= icecream
	{
		Key: "i-sonnati", LabelFa: "بستنی سنتی", LabelEn: "Traditional Ice Cream",
		NoteFa: "بستنی زعفرانی و فالوده", NoteEn: "Saffron ice cream and faloodeh",
		Business: "icecream", Theme: "sorbet", Tier: 1.0,
		Groups: []string{"icecream", "juices", "shakes"},
		NameFa: "بستنی سنتی اکبر", NameEn: "Akbar Ice Cream",
		TaglineFa: "زعفران و پسته، از سال ۱۳۵۰", TaglineEn: "Saffron and pistachio since 1971",
	},
	{
		Key: "i-modern", LabelFa: "بستنی مدرن", LabelEn: "Modern Gelato",
		NoteFa: "ژلاتو و دسر سرد", NoteEn: "Gelato and cold desserts",
		Business: "icecream", Theme: "vanilla", Tier: 1.15,
		Groups: []string{"icecream", "shakes", "cakes", "coffee_cold"},
		NameFa: "ژلاتو", NameEn: "Gelato Bar",
		TaglineFa: "هر روز تازه، ایتالیایی", TaglineEn: "Made fresh, Italian style",
	},
	{
		Key: "i-summer", LabelFa: "بستنی تابستانی", LabelEn: "Summer Stand",
		NoteFa: "منوی خنک و کوتاه", NoteEn: "Short, cool menu",
		Business: "icecream", Theme: "citrus", Tier: 0.9,
		Groups: []string{"icecream", "juices", "drinks_cold"},
		NameFa: "بستنی تابستان", NameEn: "Summer Scoop",
		TaglineFa: "برای روزهای گرم", TaglineEn: "Built for hot days",
	},
	{
		Key: "i-tiles", LabelFa: "بستنی کاشی", LabelEn: "Tile Gelato",
		NoteFa: "چیدمان کاشی با عکس مربعی", NoteEn: "Square-tile photo layout",
		Business: "icecream", Theme: "icecream-tiles-airy", Tier: 1.05,
		Groups: []string{"icecream", "shakes", "juices", "cakes"},
		NameFa: "بستنی رنگی", NameEn: "Colour Scoop",
		TaglineFa: "از روی عکس انتخاب کن", TaglineEn: "Pick it from the picture",
	},
	{
		Key: "i-arcade", LabelFa: "بستنی رترو", LabelEn: "Retro Parlour",
		NoteFa: "رنگ‌های نئونی و شاد", NoteEn: "Neon and playful",
		Business: "icecream", Theme: "arcade", Tier: 1.05,
		Groups: []string{"icecream", "shakes", "cakes", "drinks_cold"},
		NameFa: "بستنی رترو", NameEn: "Retro Parlour",
		TaglineFa: "طعم‌های قدیمی، حال‌وهوای تازه", TaglineEn: "Old flavours, new mood",
	},

	// ============================================================= juice bar
	{
		Key: "jb-classic", LabelFa: "آبمیوه بستنی", LabelEn: "Juice & Ice Cream",
		NoteFa: "آبمیوه طبیعی و معجون", NoteEn: "Fresh juice and blends",
		Business: "juicebar", Theme: "citrus", Tier: 1.0,
		Groups: []string{"juices", "icecream", "shakes"},
		NameFa: "آبمیوه بستنی شهر", NameEn: "City Juice",
		TaglineFa: "تازه، همین‌جا گرفته می‌شود", TaglineEn: "Squeezed right here",
	},
	{
		Key: "jb-health", LabelFa: "بار سالم", LabelEn: "Health Bar",
		NoteFa: "اسموتی، سالاد و آبمیوه", NoteEn: "Smoothies, salads and juice",
		Business: "juicebar", Theme: "matcha", Tier: 1.15,
		Groups: []string{"juices", "shakes", "salads", "tea_herbal"},
		NameFa: "بار سلامت", NameEn: "Wellness Bar",
		TaglineFa: "سبک، تازه، سالم", TaglineEn: "Light, fresh, good for you",
	},
	{
		Key: "jb-street", LabelFa: "آبمیوه خیابانی", LabelEn: "Street Juice",
		NoteFa: "منوی کوتاه و ارزان", NoteEn: "Short and cheap",
		Business: "juicebar", Theme: "juicebar-tiles-cream", Tier: 0.8,
		Groups: []string{"juices", "shakes", "drinks_cold"},
		NameFa: "آبمیوه گوشه خیابان", NameEn: "Corner Juice",
		TaglineFa: "یک لیوان، همین حالا", TaglineEn: "One glass, right now",
	},
	{
		Key: "jb-cafe", LabelFa: "آبمیوه کافه", LabelEn: "Juice Café",
		NoteFa: "آبمیوه به‌علاوه قهوه و کیک", NoteEn: "Juice plus coffee and cake",
		Business: "juicebar", Theme: "pistachio", Tier: 1.1,
		Groups: []string{"juices", "shakes", "coffee_cold", "cakes", "icecream"},
		NameFa: "کافه آبمیوه سبز", NameEn: "Green Juice Café",
		TaglineFa: "میوه، قهوه، آرامش", TaglineEn: "Fruit, coffee, calm",
	},
}

// All returns every pack.
func All() []Pack { return packs }

// ByKey finds one pack.
func ByKey(key string) (Pack, bool) {
	for _, p := range packs {
		if p.Key == key {
			return p, true
		}
	}
	return Pack{}, false
}

// ForBusiness returns the packs suited to a venue type, that type's own first.
func ForBusiness(business string) []Pack {
	if business == "" {
		return packs
	}

	matched := make([]Pack, 0, len(packs))
	rest := make([]Pack, 0, len(packs))
	for _, p := range packs {
		if p.Business == business {
			matched = append(matched, p)
		} else {
			rest = append(rest, p)
		}
	}
	return append(matched, rest...)
}

// Categories expands a pack into the categories and dishes it will create,
// with prices already scaled by the pack's tier.
func (p Pack) Categories() []Group {
	out := make([]Group, 0, len(p.Groups))

	for _, key := range p.Groups {
		g, ok := groups[key]
		if !ok {
			continue
		}

		// Copy before scaling so the shared catalogue is never mutated.
		scaled := Group{Key: g.Key, NameFa: g.NameFa, NameEn: g.NameEn, Icon: g.Icon}
		scaled.Dishes = make([]Dish, len(g.Dishes))
		for i, d := range g.Dishes {
			d.Price = roundPrice(float64(d.Price) * p.Tier)
			scaled.Dishes[i] = d
		}
		out = append(out, scaled)
	}
	return out
}

// ItemCount reports how many dishes a pack will create.
func (p Pack) ItemCount() int {
	n := 0
	for _, key := range p.Groups {
		if g, ok := groups[key]; ok {
			n += len(g.Dishes)
		}
	}
	return n
}

// roundPrice keeps generated prices looking like prices — to the nearest five
// thousand Toman rather than an odd number nobody would print on a menu.
func roundPrice(v float64) int64 {
	const step = 5_000
	n := int64((v + step/2) / step)
	if n < 1 {
		n = 1
	}
	return n * step
}
