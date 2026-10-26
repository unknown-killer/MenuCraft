package i18n

func init() {
	entries := map[string][2]string{
		"menu.price_pending": {"قیمت هنوز مشخص نیست", "Price not yet set"},
		"pub.price_pending":  {"قیمت به‌زودی", "Price coming soon"},
		"qr.title":           {"🪧 استند منو و شمارهٔ میز", "🪧 Menu stands & table numbers"},
		"qr.intro":           {"شمارهٔ میز را داخل کد قرار دهید و استند آمادهٔ چاپ بسازید. پس از اسکن، شماره در دفترچهٔ مشتری نمایش داده می‌شود.", "Create printable stands with the table number inside each QR link. Scanning fills the guest notebook’s table number."},
		"qr.start":           {"از شمارهٔ میز", "First table number"},
		"qr.end":             {"تا شمارهٔ میز", "Last table number"},
		"qr.layout":          {"مدل استند", "Stand layout"},
		"qr.above":           {"شماره بالای کد", "Number above QR"},
		"qr.below":           {"شماره پایین کد", "Number below QR"},
		"qr.left":            {"شماره سمت چپ کد", "Number left of QR"},
		"qr.right":           {"شماره سمت راست کد", "Number right of QR"},
		"qr.tent":            {"استند تاشو دوطرفه", "Folded double-sided tent"},
		"qr.generate":        {"ساخت پیش‌نمایش", "Generate preview"},
		"qr.print":           {"🖨️ چاپ یا ذخیرهٔ PDF", "🖨️ Print or save PDF"},
		"qr.settings":        {"بازگشت به تنظیمات", "Back to settings"},
		"qr.table":           {"میز", "Table"},
		"qr.scan":            {"منو را اسکن کنید و انتخاب‌هایتان را یادداشت کنید", "Scan the menu and note your choices"},
		"qr.print_hint":      {"کاغذ A4، اندازهٔ واقعی ۱۰۰٪، بدون سربرگ و پابرگ مرورگر. کارت معمولی ۹۰×۱۲۵ میلی‌متر؛ مدل تاشو ۹۰×۱۸۰ میلی‌متر با خط تا در وسط. پیش از چاپ همهٔ میزها، یک نمونه را اسکن کنید. حداکثر ۱۰۰ میز در هر نوبت.", "Use A4 paper, 100% scale, browser headers/footers off. Flat cards: 90×125 mm; folded tents: 90×180 mm with a center fold. Scan one proof before printing the batch. Up to 100 tables per batch."},
		"qr.public_hint":     {"برای استفادهٔ مشتریان، برنامه باید با نشانی اینترنتی یا شبکهٔ قابل‌دسترسی اجرا شود؛ نشانی محلی رایانه روی گوشی مشتری باز نمی‌شود.", "Guests need a reachable public or network address. A computer’s localhost link will not open its menu on a guest’s phone."},
		"qr.draft":           {"منو هنوز منتشر نشده است؛ پیش از استفادهٔ مشتریان آن را در تنظیمات منتشر کنید.", "This menu is still a draft. Publish it in settings before guests scan these codes."},
		"qr.range_error":     {"شماره‌ها باید بین ۱ تا ۹۹۹۹ باشند؛ پایان از شروع کمتر نباشد و هر نوبت حداکثر ۱۰۰ میز باشد.", "Use table numbers 1–9999, last ≥ first, with at most 100 tables per batch."},
		"qr.layout_error":    {"مدل استند معتبر نیست.", "Invalid stand layout."},
	}
	for key, value := range entries {
		catalogFa[key] = value[0]
		catalogEn[key] = value[1]
	}
}
