// Package menupack provides complete ready-made menus: real dishes with real
// prices, grouped into categories, paired with a matching theme.
//
// The library is split in two so fifty packs do not mean fifty hand-written
// menus: a shared catalogue of dishes grouped by course, and packs that pick
// groups, a theme and a price tier. The same dish can appear in a cheap
// sandwich shop and an upmarket restaurant at different prices, which is how
// real menus work anyway.
package menupack

// Dish is one menu item. Price is in Toman at the "standard" tier; a pack's
// price tier scales it.
type Dish struct {
	Fa     string
	En     string
	DescFa string
	DescEn string
	Price  int64
	Tags   []string
	Cal    int
	Min    int
}

// Group is a category's worth of dishes.
type Group struct {
	Key    string
	NameFa string
	NameEn string
	Icon   string
	Dishes []Dish
}

// groups is the whole catalogue, keyed for packs to reference.
var groups = map[string]Group{

	// ------------------------------------------------ Persian restaurant
	"starters_ir": {"starters_ir", "پیش‌غذا", "Starters", "🥗", []Dish{
		{"ماست و خیار", "Mast-o-Khiar", "ماست محلی با خیار، نعنا و گردو", "Yoghurt with cucumber, mint and walnut", 85_000, []string{"cold", "vegetarian"}, 120, 5},
		{"ماست موسیر", "Mast-o-Musir", "ماست چکیده با موسیر کوهی", "Strained yoghurt with wild shallot", 78_000, []string{"cold", "vegetarian"}, 110, 5},
		{"کشک بادمجان", "Kashk-e Bademjan", "بادمجان سرخ‌شده با کشک، نعنا داغ و پیاز", "Fried aubergine with whey, mint and onion", 145_000, []string{"popular", "hot"}, 310, 12},
		{"میرزا قاسمی", "Mirza Ghasemi", "بادمجان دودی با گوجه، سیر و تخم‌مرغ", "Smoked aubergine with tomato, garlic and egg", 165_000, []string{"vegetarian"}, 280, 15},
		{"سالاد شیرازی", "Shirazi Salad", "گوجه، خیار و پیاز با آبغوره تازه", "Tomato, cucumber and onion with verjuice", 95_000, []string{"vegetarian", "cold"}, 90, 7},
		{"زیتون پرورده", "Zeytoon Parvardeh", "زیتون شمالی با گردو، انار و سبزی", "Northern olives with walnut and pomegranate", 135_000, []string{"cold", "chef"}, 210, 6},
		{"سوپ جو", "Barley Soup", "سوپ جو خامه‌ای با هویج و جعفری", "Creamy barley soup with carrot and parsley", 110_000, []string{"hot"}, 190, 10},
		{"آش رشته", "Ash Reshteh", "آش سنتی با رشته، حبوبات و کشک", "Herb and noodle soup with whey", 155_000, []string{"popular", "hot"}, 340, 12},
	}},

	"kebabs": {"kebabs", "کباب‌ها", "Kebabs", "🍢", []Dish{
		{"کباب کوبیده", "Koobideh Kebab", "دو سیخ کوبیده گوشت گوسفندی با برنج ایرانی", "Two skewers of minced lamb with Persian rice", 385_000, []string{"popular"}, 720, 25},
		{"جوجه کباب زعفرانی", "Saffron Chicken Kebab", "سینه مرغ زعفرانی با کره و لیمو", "Saffron chicken breast with butter and lime", 340_000, []string{"chef"}, 610, 22},
		{"کباب برگ", "Barg Kebab", "فیله گوسفندی زعفرانی، نرم و آبدار", "Saffron lamb fillet, tender and juicy", 520_000, []string{"popular", "chef"}, 680, 30},
		{"کباب بختیاری", "Bakhtiari Kebab", "ترکیب فیله گوسفند و جوجه زعفرانی", "Lamb fillet combined with saffron chicken", 495_000, nil, 700, 28},
		{"شیشلیک", "Shishlik", "راسته گوسفندی با استخوان، مخصوص خانه", "Bone-in lamb chops, house speciality", 680_000, []string{"chef"}, 820, 35},
		{"کباب لقمه", "Loghmeh Kebab", "کوبیده پهن با چاشنی مخصوص", "Wide minced kebab with house seasoning", 410_000, nil, 740, 26},
		{"جوجه با استخوان", "Bone-in Chicken Kebab", "ران مرغ مزه‌دار شده با ادویه مخصوص", "Marinated chicken thigh with house spices", 320_000, nil, 640, 24},
		{"کباب ترش", "Kabab Torsh", "کباب گیلانی با رب انار و گردو", "Gilan kebab with pomegranate and walnut", 540_000, []string{"chef"}, 700, 32},
	}},

	"stews": {"stews", "خورش‌ها", "Stews", "🍲", []Dish{
		{"قورمه سبزی", "Ghormeh Sabzi", "خورش سنتی با سبزی معطر، لوبیا و گوشت", "Herb stew with beans and lamb", 295_000, []string{"popular"}, 540, 20},
		{"قیمه بادمجان", "Gheimeh Bademjan", "خورش لپه با بادمجان سرخ‌شده", "Split-pea stew with fried aubergine", 285_000, nil, 520, 20},
		{"فسنجان", "Fesenjan", "خورش گردو و انار با مرغ", "Walnut and pomegranate stew with chicken", 340_000, []string{"chef"}, 620, 25},
		{"باقالی پلو با ماهیچه", "Baghali Polo Mahicheh", "پلوی شوید و باقالی با ماهیچه گوسفندی", "Dill and broad bean rice with lamb shank", 620_000, []string{"chef", "popular"}, 780, 30},
		{"خورش کرفس", "Celery Stew", "خورش کرفس با نعنا و جعفری", "Celery stew with mint and parsley", 275_000, nil, 480, 20},
		{"آبگوشت", "Abgoosht", "دیزی سنتی با گوشت، نخود و لوبیا", "Traditional lamb and chickpea stew", 265_000, []string{"popular"}, 690, 25},
	}},

	"rice_ir": {"rice_ir", "چلو و پلو", "Rice dishes", "🍚", []Dish{
		{"ته‌چین مرغ", "Tahchin", "ته‌چین زعفرانی لایه‌ای با مرغ و زرشک", "Layered saffron rice cake with chicken and barberry", 320_000, []string{"popular", "chef"}, 690, 25},
		{"لوبیا پلو", "Loobia Polo", "پلوی لوبیا سبز با گوشت چرخ‌کرده", "Green bean rice with minced meat", 265_000, nil, 610, 20},
		{"عدس پلو", "Adas Polo", "پلوی عدس با کشمش و خرما", "Lentil rice with raisin and date", 235_000, []string{"vegetarian"}, 570, 18},
		{"استانبولی پلو", "Estamboli Polo", "پلوی گوجه با گوشت و سیب‌زمینی", "Tomato rice with meat and potato", 245_000, nil, 590, 20},
		{"زرشک پلو با مرغ", "Zereshk Polo", "پلوی زرشک زعفرانی با ران مرغ", "Barberry saffron rice with chicken", 295_000, []string{"popular"}, 640, 22},
	}},

	"seafood": {"seafood", "غذای دریایی", "Seafood", "🐟", []Dish{
		{"قزل‌آلای کبابی", "Grilled Trout", "قزل‌آلای تازه با لیمو و سبزیجات", "Fresh trout with lemon and vegetables", 420_000, nil, 450, 25},
		{"میگو سوخاری", "Crispy Prawns", "میگوی درشت سوخاری با سس مخصوص", "Large breaded prawns with house sauce", 480_000, []string{"new"}, 560, 18},
		{"ماهی شکم‌پر", "Stuffed Fish", "ماهی سفید شکم‌پر با سبزی و رب انار", "Whitefish stuffed with herbs and pomegranate", 560_000, []string{"chef"}, 520, 35},
		{"میگو پلو", "Shrimp Rice", "پلوی سبزی جنوبی با میگو", "Southern herb rice with shrimp", 490_000, []string{"spicy"}, 610, 28},
	}},

	// ------------------------------------------------ grill house
	"liver_grill": {"liver_grill", "جگر و دل و قلوه", "Liver & offal", "🔥", []Dish{
		{"جگر گوسفندی", "Lamb Liver", "جگر تازه گوسفندی، سیخی", "Fresh lamb liver, per skewer", 145_000, []string{"popular"}, 290, 12},
		{"دل و قلوه", "Heart & Kidney", "دل و قلوه مخلوط با ادویه مخصوص", "Mixed heart and kidney with house spice", 155_000, nil, 320, 14},
		{"قلوه‌گاه", "Lamb Ribs", "قلوه‌گاه چرب و آبدار", "Rich, juicy lamb ribs", 190_000, []string{"chef"}, 480, 18},
		{"دنبلان", "Donbalan", "دنبلان تازه کبابی", "Grilled lamb sweetbread", 240_000, nil, 350, 16},
		{"سنگدان", "Gizzard", "سنگدان مرغ کبابی با لیمو", "Grilled chicken gizzard with lime", 120_000, nil, 210, 12},
	}},

	"grill_sides": {"grill_sides", "مخلفات", "Sides", "🧅", []Dish{
		{"گوجه کبابی", "Grilled Tomato", "دو عدد گوجه کبابی", "Two grilled tomatoes", 45_000, []string{"vegetarian"}, 60, 8},
		{"پیاز و سبزی", "Onion & Herbs", "پیاز تازه با سبزی خوردن", "Fresh onion with herb platter", 40_000, []string{"vegetarian", "cold"}, 45, 3},
		{"نان تنوری", "Tandoor Bread", "نان تازه از تنور", "Fresh bread from the tandoor", 25_000, []string{"vegetarian"}, 180, 5},
		{"ترشی مخلوط", "Mixed Pickles", "ترشی خانگی مخلوط", "House mixed pickles", 45_000, []string{"vegetarian", "cold"}, 40, 2},
		{"کره و زعفران", "Butter & Saffron", "کره محلی با زعفران دم‌کرده", "Local butter with brewed saffron", 35_000, nil, 120, 2},
	}},

	// ------------------------------------------------ fast food
	"burgers": {"burgers", "برگرها", "Burgers", "🍔", []Dish{
		{"همبرگر کلاسیک", "Classic Burger", "برگر دست‌ساز گوشت با کاهو، گوجه و سس مخصوص", "Handmade beef burger with lettuce, tomato and house sauce", 235_000, []string{"popular"}, 680, 15},
		{"چیزبرگر", "Cheeseburger", "برگر گوشت با پنیر چدار دوبل", "Beef burger with double cheddar", 265_000, []string{"popular"}, 760, 15},
		{"برگر قارچ و پنیر", "Mushroom Swiss", "برگر با قارچ سوته و پنیر سوئیسی", "Burger with sautéed mushroom and Swiss cheese", 285_000, nil, 790, 17},
		{"دوبل برگر", "Double Burger", "دو لایه گوشت با پنیر و سس دودی", "Two beef patties with cheese and smoky sauce", 345_000, []string{"chef"}, 980, 18},
		{"برگر تند مکزیکی", "Mexican Spicy", "برگر با فلفل هالوپینو و سس تند", "Burger with jalapeño and hot sauce", 295_000, []string{"spicy"}, 820, 17},
		{"چیکن برگر", "Chicken Burger", "فیله مرغ سوخاری با سس سزار", "Crispy chicken fillet with Caesar sauce", 245_000, nil, 640, 15},
		{"وگی برگر", "Veggie Burger", "برگر گیاهی با نخود و سبزیجات", "Chickpea and vegetable patty", 205_000, []string{"vegetarian", "new"}, 480, 14},
	}},

	"pizzas": {"pizzas", "پیتزاها", "Pizzas", "🍕", []Dish{
		{"پیتزا مخصوص", "House Special", "ژامبون، قارچ، فلفل دلمه و پنیر پیتزا", "Ham, mushroom, pepper and mozzarella", 385_000, []string{"popular"}, 980, 20},
		{"پیتزا پپرونی", "Pepperoni", "پپرونی تند با پنیر دوبل", "Spicy pepperoni with double cheese", 365_000, []string{"spicy", "popular"}, 940, 20},
		{"پیتزا مارگاریتا", "Margherita", "گوجه، ریحان تازه و موزارلا", "Tomato, fresh basil and mozzarella", 295_000, []string{"vegetarian"}, 780, 18},
		{"پیتزا سبزیجات", "Vegetarian", "قارچ، ذرت، زیتون و فلفل دلمه", "Mushroom, corn, olive and pepper", 315_000, []string{"vegetarian"}, 760, 18},
		{"پیتزا چهار فصل", "Four Seasons", "چهار طعم در یک پیتزا", "Four flavours on one base", 425_000, []string{"chef"}, 1020, 22},
		{"پیتزا مرغ و قارچ", "Chicken & Mushroom", "فیله مرغ گریل با قارچ و سس آلفردو", "Grilled chicken with mushroom and Alfredo", 395_000, nil, 900, 20},
	}},

	"sandwiches": {"sandwiches", "ساندویچ‌ها", "Sandwiches", "🥪", []Dish{
		{"ساندویچ فلافل", "Falafel", "فلافل تازه با سبزیجات و سس سیر", "Fresh falafel with vegetables and garlic sauce", 125_000, []string{"vegetarian"}, 520, 10},
		{"هات‌داگ ویژه", "Special Hot Dog", "سوسیس آلمانی با پنیر و سس", "German sausage with cheese and sauce", 185_000, nil, 610, 12},
		{"ساندویچ مرغ گریل", "Grilled Chicken", "فیله مرغ گریل با سس هزار جزیره", "Grilled chicken fillet with thousand island", 215_000, []string{"popular"}, 580, 14},
		{"ساندویچ ژامبون", "Ham & Cheese", "ژامبون گوشت با پنیر گودا", "Beef ham with gouda", 195_000, nil, 620, 10},
		{"ساندویچ بندری", "Bandari", "سوسیس بندری تند با نان باگت", "Spicy bandari sausage in baguette", 165_000, []string{"spicy"}, 640, 12},
	}},

	"fried": {"fried", "سوخاری", "Fried chicken", "🍗", []Dish{
		{"مرغ سوخاری دو تکه", "2pc Fried Chicken", "دو تکه مرغ سوخاری با سیب‌زمینی", "Two pieces with fries", 285_000, []string{"popular"}, 780, 18},
		{"مرغ سوخاری چهار تکه", "4pc Fried Chicken", "چهار تکه با سیب‌زمینی و سالاد کلم", "Four pieces with fries and coleslaw", 465_000, nil, 1240, 22},
		{"استریپس", "Chicken Strips", "پنج عدد استریپس با سس", "Five strips with dipping sauce", 265_000, nil, 690, 16},
		{"بال سوخاری تند", "Spicy Wings", "شش عدد بال با سس تند", "Six wings in hot sauce", 245_000, []string{"spicy"}, 640, 16},
		{"قارچ سوخاری", "Fried Mushroom", "قارچ سوخاری با سس رنچ", "Breaded mushroom with ranch", 175_000, []string{"vegetarian"}, 420, 12},
	}},

	"sides_ff": {"sides_ff", "سیب‌زمینی و پیش‌غذا", "Fries & sides", "🍟", []Dish{
		{"سیب‌زمینی سرخ‌کرده", "French Fries", "سیب‌زمینی تازه با نمک دریا", "Fresh cut fries with sea salt", 115_000, []string{"vegetarian", "popular"}, 380, 8},
		{"سیب‌زمینی با پنیر", "Cheese Fries", "سیب‌زمینی با سس پنیر چدار", "Fries with cheddar sauce", 165_000, []string{"vegetarian"}, 520, 10},
		{"پیاز حلقه‌ای", "Onion Rings", "هشت عدد پیاز حلقه‌ای سوخاری", "Eight breaded onion rings", 145_000, []string{"vegetarian"}, 410, 10},
		{"ناچو", "Nachos", "چیپس ذرت با سس پنیر و هالوپینو", "Corn chips with cheese sauce and jalapeño", 185_000, []string{"spicy", "vegetarian"}, 490, 8},
		{"سالاد کلم", "Coleslaw", "سالاد کلم با سس مایونز", "Cabbage salad with mayo dressing", 75_000, []string{"vegetarian", "cold"}, 180, 5},
	}},

	// ------------------------------------------------ western
	"pasta": {"pasta", "پاستا و ریزوتو", "Pasta & risotto", "🍝", []Dish{
		{"پاستا آلفردو", "Alfredo Pasta", "فتوچینی با سس خامه و مرغ گریل", "Fettuccine with cream sauce and grilled chicken", 335_000, []string{"popular"}, 720, 20},
		{"اسپاگتی بولونز", "Spaghetti Bolognese", "اسپاگتی با سس گوشت ایتالیایی", "Spaghetti with Italian meat sauce", 315_000, nil, 690, 20},
		{"پنه آرابیاتا", "Penne Arrabbiata", "پنه با سس گوجه تند", "Penne in spicy tomato sauce", 275_000, []string{"spicy", "vegetarian"}, 580, 18},
		{"لازانیا", "Lasagna", "لازانیای خانگی با گوشت و بشامل", "Homemade lasagna with beef and béchamel", 365_000, []string{"chef"}, 810, 25},
		{"ریزوتو قارچ", "Mushroom Risotto", "ریزوتو با قارچ پورچینی و پارمزان", "Risotto with porcini and parmesan", 345_000, []string{"vegetarian"}, 640, 25},
	}},

	"steaks": {"steaks", "استیک و گریل", "Steak & grill", "🥩", []Dish{
		{"استیک راسته", "Sirloin Steak", "راسته گوساله با سس قارچ و سیب‌زمینی", "Beef sirloin with mushroom sauce and potato", 780_000, []string{"chef", "popular"}, 890, 30},
		{"استیک فیله", "Fillet Steak", "فیله گوساله نرم با سس فلفل", "Tender beef fillet with pepper sauce", 890_000, []string{"chef"}, 820, 32},
		{"مرغ گریل", "Grilled Chicken", "سینه مرغ گریل با سبزیجات بخارپز", "Grilled chicken breast with steamed vegetables", 385_000, nil, 560, 24},
		{"کباب مخلوط گریل", "Mixed Grill", "ترکیب گوشت، مرغ و سوسیس گریل", "Beef, chicken and sausage on one plate", 690_000, nil, 940, 30},
	}},

	"salads": {"salads", "سالادها", "Salads", "🥙", []Dish{
		{"سالاد سزار", "Caesar Salad", "کاهو رومی، مرغ گریل، پارمزان و کروتون", "Romaine, grilled chicken, parmesan and croutons", 245_000, []string{"popular"}, 380, 12},
		{"سالاد یونانی", "Greek Salad", "گوجه، خیار، زیتون و پنیر فتا", "Tomato, cucumber, olive and feta", 215_000, []string{"vegetarian", "cold"}, 290, 10},
		{"سالاد سبز", "Garden Salad", "سبزیجات تازه فصل با سس بالزامیک", "Seasonal greens with balsamic", 165_000, []string{"vegetarian", "cold"}, 180, 8},
		{"سالاد اندونزی", "Indonesian Salad", "سالاد ماکارونی با سس مخصوص", "Pasta salad with house dressing", 175_000, []string{"cold"}, 340, 8},
	}},

	// ------------------------------------------------ café
	"coffee_hot": {"coffee_hot", "قهوه گرم", "Hot coffee", "☕", []Dish{
		{"اسپرسو", "Espresso", "تک شات قهوه عربیکا", "Single shot arabica", 85_000, nil, 5, 3},
		{"آمریکانو", "Americano", "اسپرسو با آب داغ", "Espresso with hot water", 95_000, nil, 10, 4},
		{"کاپوچینو", "Cappuccino", "اسپرسو با شیر بخارپز و فوم", "Espresso with steamed milk and foam", 125_000, []string{"popular"}, 140, 5},
		{"لاته", "Café Latte", "اسپرسو با شیر فراوان", "Espresso with plenty of milk", 135_000, []string{"popular"}, 180, 5},
		{"موکا", "Café Mocha", "لاته با شکلات تلخ", "Latte with dark chocolate", 155_000, nil, 260, 6},
		{"فلت وایت", "Flat White", "دو شات ریسترتو با شیر مخملی", "Double ristretto with velvet milk", 145_000, []string{"chef"}, 150, 5},
		{"هات چاکلت", "Hot Chocolate", "شکلات داغ غلیظ با خامه", "Thick hot chocolate with cream", 145_000, nil, 320, 6},
	}},

	"coffee_cold": {"coffee_cold", "قهوه سرد", "Iced coffee", "🧊", []Dish{
		{"آیس آمریکانو", "Iced Americano", "آمریکانو با یخ فراوان", "Americano over ice", 115_000, []string{"cold"}, 15, 4},
		{"آیس لاته", "Iced Latte", "لاته سرد با شیر تازه", "Cold latte with fresh milk", 145_000, []string{"cold", "popular"}, 190, 5},
		{"فراپه", "Frappé", "قهوه بلند شده با یخ و خامه", "Blended coffee with ice and cream", 175_000, []string{"cold"}, 340, 7},
		{"کولد برو", "Cold Brew", "دم سرد ۱۲ ساعته", "Twelve-hour cold brew", 165_000, []string{"cold", "chef"}, 12, 3},
		{"آفوگاتو", "Affogato", "بستنی وانیلی با شات اسپرسو", "Vanilla ice cream with an espresso shot", 185_000, []string{"cold", "popular"}, 290, 4},
	}},

	"tea_herbal": {"tea_herbal", "چای و دمنوش", "Tea & herbal", "🫖", []Dish{
		{"چای ایرانی", "Persian Tea", "چای دم‌کشیده با نبات", "Brewed black tea with rock candy", 55_000, []string{"hot", "popular"}, 5, 5},
		{"چای زعفرانی", "Saffron Tea", "چای تازه‌دم با زعفران و هل", "Fresh tea with saffron and cardamom", 85_000, []string{"hot"}, 15, 6},
		{"دمنوش به‌لیمو", "Lemon Verbena", "دمنوش آرام‌بخش به‌لیمو", "Calming lemon verbena infusion", 75_000, []string{"hot", "vegetarian"}, 5, 6},
		{"دمنوش بهار نارنج", "Orange Blossom", "دمنوش معطر بهار نارنج", "Fragrant orange blossom infusion", 85_000, []string{"hot"}, 5, 6},
		{"چای سبز", "Green Tea", "چای سبز با نعنا تازه", "Green tea with fresh mint", 75_000, []string{"hot"}, 5, 5},
		{"دمنوش ترش", "Hibiscus", "دمنوش چای ترش با عسل", "Hibiscus infusion with honey", 80_000, []string{"hot"}, 40, 6},
	}},

	"shakes": {"shakes", "شیک و اسموتی", "Shakes & smoothies", "🥤", []Dish{
		{"میلک شیک شکلات", "Chocolate Shake", "شیک غلیظ شکلاتی با بستنی", "Thick chocolate shake with ice cream", 185_000, []string{"cold", "popular"}, 480, 6},
		{"میلک شیک وانیل", "Vanilla Shake", "شیک وانیلی با خامه", "Vanilla shake with cream", 175_000, []string{"cold"}, 450, 6},
		{"شیک نوتلا", "Nutella Shake", "شیک نوتلا با فندق", "Nutella shake with hazelnut", 215_000, []string{"cold", "chef"}, 620, 7},
		{"اسموتی توت‌فرنگی", "Strawberry Smoothie", "اسموتی توت‌فرنگی تازه با ماست", "Fresh strawberry smoothie with yoghurt", 165_000, []string{"cold", "vegetarian"}, 260, 6},
		{"اسموتی انبه", "Mango Smoothie", "اسموتی انبه با شیر نارگیل", "Mango smoothie with coconut milk", 175_000, []string{"cold", "new"}, 280, 6},
	}},

	"cakes": {"cakes", "کیک و دسر", "Cakes & desserts", "🍰", []Dish{
		{"چیز کیک", "Cheesecake", "چیز کیک نیویورکی با سس توت", "New York cheesecake with berry sauce", 195_000, []string{"popular"}, 420, 5},
		{"کیک شکلاتی", "Chocolate Cake", "کیک شکلاتی خیس با گاناش", "Moist chocolate cake with ganache", 185_000, []string{"popular"}, 480, 5},
		{"تیرامیسو", "Tiramisu", "تیرامیسوی خانگی با قهوه", "Homemade tiramisu with coffee", 205_000, []string{"chef"}, 390, 5},
		{"کیک هویج", "Carrot Cake", "کیک هویج با فراستینگ پنیری", "Carrot cake with cream cheese frosting", 175_000, []string{"vegetarian"}, 410, 5},
		{"براونی", "Brownie", "براونی گردویی با بستنی", "Walnut brownie with ice cream", 195_000, nil, 520, 6},
	}},

	"breakfast": {"breakfast", "صبحانه", "Breakfast", "🍳", []Dish{
		{"املت گوجه", "Tomato Omelette", "املت گوجه با نان تازه", "Tomato omelette with fresh bread", 145_000, []string{"vegetarian", "hot"}, 320, 12},
		{"نیمرو", "Fried Eggs", "دو عدد تخم‌مرغ نیمرو با کره", "Two fried eggs in butter", 125_000, []string{"vegetarian"}, 280, 8},
		{"عدسی", "Adasi", "عدسی گرم با نان بربری", "Warm lentil soup with barbari bread", 115_000, []string{"vegetarian", "hot"}, 340, 10},
		{"حلیم", "Haleem", "حلیم گندم با دارچین و شکر", "Wheat porridge with cinnamon and sugar", 135_000, []string{"hot", "popular"}, 420, 10},
		{"پنیر و سبزی و گردو", "Cheese, Herbs & Walnut", "پنیر تبریز با سبزی و گردو", "Tabriz cheese with herbs and walnut", 155_000, []string{"vegetarian", "cold"}, 380, 6},
		{"صبحانه انگلیسی", "English Breakfast", "تخم‌مرغ، سوسیس، قارچ و لوبیا", "Eggs, sausage, mushroom and beans", 285_000, nil, 680, 18},
	}},

	// ------------------------------------------------ teahouse
	"hookah": {"hookah", "قلیان", "Hookah", "💨", []Dish{
		{"قلیان دو سیب", "Double Apple", "تنباکو دو سیب با ذغال طبیعی", "Double apple tobacco with natural charcoal", 320_000, []string{"popular"}, 0, 10},
		{"قلیان هلو", "Peach", "تنباکو هلو با یخ", "Peach tobacco with ice", 340_000, nil, 0, 10},
		{"قلیان نعنا لیمو", "Mint Lemon", "ترکیب خنک نعنا و لیمو", "Cool mint and lemon blend", 340_000, []string{"cold"}, 0, 10},
		{"قلیان ویژه خانه", "House Special", "ترکیب مخصوص سفره‌خانه", "The house's own blend", 420_000, []string{"chef"}, 0, 12},
		{"شارژ ذغال", "Charcoal Refill", "تعویض ذغال", "Charcoal change", 60_000, nil, 0, 3},
	}},

	"dizi": {"dizi", "دیزی و سنتی", "Dizi & traditional", "🍲", []Dish{
		{"دیزی سنگی", "Stone Dizi", "دیزی سنتی در ظرف سنگی با نان و سبزی", "Traditional dizi in a stone pot with bread and herbs", 320_000, []string{"popular", "chef"}, 780, 25},
		{"کله پاچه", "Kalleh Pacheh", "سرویس کامل با نان تازه", "Full service with fresh bread", 380_000, nil, 720, 20},
		{"آش دوغ", "Ash-e Doogh", "آش سنتی با دوغ و سبزی", "Traditional yoghurt and herb soup", 145_000, []string{"hot"}, 290, 12},
		{"حلیم بادمجان", "Halim Bademjan", "حلیم بادمجان با کشک و نعنا", "Aubergine halim with whey and mint", 175_000, []string{"vegetarian"}, 340, 15},
	}},

	// ------------------------------------------------ bakery
	"breads": {"breads", "نان‌ها", "Breads", "🥖", []Dish{
		{"نان بربری", "Barbari", "نان تازه بربری از تنور", "Fresh barbari from the oven", 25_000, []string{"vegetarian"}, 280, 2},
		{"نان سنگک", "Sangak", "سنگک سنتی با کنجد", "Traditional sangak with sesame", 30_000, []string{"vegetarian", "popular"}, 250, 2},
		{"باگت فرانسوی", "French Baguette", "باگت تازه با پوسته ترد", "Fresh baguette with a crisp crust", 55_000, []string{"vegetarian"}, 320, 2},
		{"نان جو", "Barley Bread", "نان جو سبوس‌دار", "Wholegrain barley bread", 45_000, []string{"vegetarian"}, 240, 2},
		{"کروسان کره", "Butter Croissant", "کروسان ورقه‌ای با کره تازه", "Flaky croissant with fresh butter", 95_000, []string{"vegetarian", "popular"}, 380, 3},
	}},

	"pastries": {"pastries", "شیرینی", "Pastries", "🧁", []Dish{
		{"شیرینی خامه‌ای", "Cream Puff", "نان خامه‌ای تازه", "Fresh cream puff", 65_000, []string{"popular"}, 290, 2},
		{"دانمارکی", "Danish", "شیرینی دانمارکی با سیب و دارچین", "Danish pastry with apple and cinnamon", 85_000, nil, 340, 2},
		{"نان خشک", "Dry Sweets", "شیرینی خشک مخلوط، صد گرم", "Assorted dry sweets, per 100g", 95_000, []string{"vegetarian"}, 420, 2},
		{"باقلوا", "Baklava", "باقلوای پسته‌ای، صد گرم", "Pistachio baklava, per 100g", 165_000, []string{"chef"}, 480, 2},
		{"زولبیا بامیه", "Zoolbia Bamieh", "زولبیا و بامیه تازه، صد گرم", "Fresh zoolbia and bamieh, per 100g", 85_000, nil, 460, 2},
	}},

	// ------------------------------------------------ ice cream & juice
	"icecream": {"icecream", "بستنی", "Ice cream", "🍨", []Dish{
		{"بستنی سنتی زعفرانی", "Saffron Ice Cream", "بستنی سنتی با خامه و پسته", "Traditional ice cream with cream and pistachio", 120_000, []string{"cold", "popular"}, 320, 3},
		{"فالوده شیرازی", "Faloodeh Shirazi", "رشته فالوده با شربت لیمو", "Rice noodle dessert with lime syrup", 110_000, []string{"cold"}, 210, 3},
		{"بستنی وانیلی", "Vanilla Ice Cream", "دو اسکوپ بستنی وانیل", "Two scoops of vanilla", 95_000, []string{"cold", "vegetarian"}, 280, 2},
		{"بستنی شکلاتی", "Chocolate Ice Cream", "دو اسکوپ بستنی شکلات بلژیکی", "Two scoops of Belgian chocolate", 105_000, []string{"cold"}, 310, 2},
		{"معجون ویژه", "Special Majoon", "معجون با بستنی، خرما و مغزیجات", "Blend with ice cream, dates and nuts", 195_000, []string{"cold", "chef"}, 540, 6},
	}},

	"juices": {"juices", "آبمیوه طبیعی", "Fresh juice", "🧃", []Dish{
		{"آب پرتقال", "Orange Juice", "آب پرتقال تازه، بدون شکر", "Freshly squeezed, no sugar", 125_000, []string{"cold", "vegetarian"}, 180, 4},
		{"آب هویج", "Carrot Juice", "آب هویج تازه با بستنی", "Fresh carrot juice with ice cream", 145_000, []string{"cold", "popular"}, 220, 5},
		{"آب انار", "Pomegranate Juice", "آب انار تازه فصل", "Seasonal fresh pomegranate", 175_000, []string{"cold"}, 190, 5},
		{"آب طالبی", "Melon Juice", "آب طالبی خنک", "Chilled melon juice", 135_000, []string{"cold"}, 160, 4},
		{"شیک موز", "Banana Shake", "شیک موز با شیر و عسل", "Banana shake with milk and honey", 155_000, []string{"cold"}, 320, 5},
	}},

	// ------------------------------------------------ drinks (universal)
	"drinks_cold": {"drinks_cold", "نوشیدنی‌ها", "Drinks", "🥤", []Dish{
		{"نوشابه قوطی", "Canned Soda", "کوکا، فانتا یا اسپرایت", "Coke, Fanta or Sprite", 45_000, []string{"cold"}, 140, 1},
		{"دلستر", "Malt Drink", "دلستر با طعم‌های مختلف", "Malt beverage, assorted flavours", 55_000, []string{"cold"}, 120, 1},
		{"دوغ محلی", "Local Doogh", "دوغ گازدار با نعنا خشک", "Sparkling yoghurt drink with dried mint", 55_000, []string{"cold", "popular"}, 80, 2},
		{"آب معدنی", "Mineral Water", "بطری نیم لیتری", "Half-litre bottle", 25_000, []string{"cold"}, 0, 1},
		{"شربت سکنجبین", "Sekanjabin", "شربت سنتی نعنا و سرکه", "Traditional mint and vinegar cordial", 65_000, []string{"cold", "vegetarian"}, 140, 3},
		{"لیموناد", "Lemonade", "لیموناد تازه با نعنا", "Fresh lemonade with mint", 95_000, []string{"cold"}, 150, 4},
	}},

	"desserts_ir": {"desserts_ir", "دسر", "Desserts", "🍮", []Dish{
		{"شله زرد", "Sholeh Zard", "دسر برنج زعفرانی با دارچین", "Saffron rice pudding with cinnamon", 95_000, []string{"vegetarian"}, 260, 4},
		{"فرنی", "Ferni", "فرنی سنتی با گلاب", "Rice pudding with rosewater", 85_000, []string{"vegetarian"}, 240, 4},
		{"بستنی سنتی", "Traditional Ice Cream", "بستنی زعفرانی با پسته", "Saffron ice cream with pistachio", 120_000, []string{"cold", "popular"}, 320, 3},
		{"حلوا", "Halva", "حلوای زعفرانی خانگی", "Homemade saffron halva", 95_000, nil, 380, 4},
	}},
}

// Group returns one catalogue group.
func Group_(key string) (Group, bool) {
	g, ok := groups[key]
	return g, ok
}

// AllGroups exposes the catalogue for tests and tooling.
func AllGroups() map[string]Group { return groups }
