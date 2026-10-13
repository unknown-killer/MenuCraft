# WordPress integration / راهنمای وردپرس

Requires WordPress 6.4+ and PHP 8.1+. The plugin is included at `wordpress/menucraft` and packaged as `menucraft-wordpress_1.1.0.zip`.

## Install

1. Download the plugin ZIP from the release.
2. In WordPress, open **Plugins → Add New → Upload Plugin**.
3. Upload, install and activate **MenuCraft Restaurant Menu**.
4. Open **Settings → MenuCraft**. Only administrators with `manage_options` may change settings or import menus.
5. Add a Shortcode block to a page with `[menucraft_menu lang="fa"]` or `[menucraft_menu lang="en"]`.

در پیشخوان وردپرس از بخش افزونه‌ها ZIP را نصب و فعال کنید. سپس تنظیمات MenuCraft را باز کنید و شورت‌کد زبان دلخواه را در یک برگه قرار دهید.

## Option A: Native WordPress menu

In the Go owner panel, open **Menu transfer** and download JSON. In WordPress choose **Native WordPress (JSON)**, save, then import the JSON file.

This stores a validated snapshot in WordPress options. It renders the venue name/description, visible categories and available dishes, bilingual content, prices/discounts, responsive cards, search and category filtering. Persian letters are normalized during search; Vazirmatn is served from the plugin itself.

Importing **replaces the previous WordPress snapshot**. Invalid files are rejected before updating it. Limits are 2 MiB, 200 categories and 5,000 dishes. The PHP/server upload limit must also permit the file.

حالت مستقل، متن و قیمت و دسته‌های منو را در خود وردپرس ذخیره می‌کند. ورود فایل جدید، نسخهٔ قبلی وردپرس را جایگزین می‌کند. فایل نامعتبر پیش از تغییر داده رد می‌شود. برای تغییرات بعدی از Go خروجی تازه بگیرید و دوباره وارد کنید.

### Native-mode boundaries

- No Go server is needed to serve imported text and prices.
- Photos are linked, not downloaded. Exported `/uploads/...` paths resolve against the JSON `source_url` origin. Set `MENU_PUBLIC_URL` to a reachable HTTPS origin before exporting. A `localhost` source will not work for real visitors.
- To remove the image dependency, upload photos to WordPress Media and replace the JSON items' `image_path` values with those HTTPS image URLs before importing.
- Go appearance presets, custom CSS, canvas, blocks, favorites, dietary controls and ordering are not replicated by native mode. WordPress provides its own presentation.
- There is no automatic sync, built-in dish editor, ecommerce integration or remote JSON fetching.
- Data remains in WordPress options when the plugin is deactivated/deleted. Remove `menucraft_menu_data` and `menucraft_settings` through your normal administration tools if permanent erasure is wanted.

برای متن و قیمت نیازی به روشن‌بودن Go نیست، اما تصاویر به میزبان اولیه وابسته‌اند. برای استقلال تصاویر، آن‌ها را در کتابخانهٔ وردپرس بارگذاری کنید و نشانی HTTPS را در `image_path` فایل JSON بگذارید. قالب، بوم و CSS اختصاصی Go منتقل نمی‌شوند. افزونه، داده‌های واردشده را هنگام حذف خودکار پاک نمی‌کند.

## Option B: Live Go menu

1. Publish the menu in Go and copy its full HTTPS URL, e.g. `https://menu.example.com/m/your-restaurant`.
2. On the Go server, set `MENU_PUBLIC_URL=https://menu.example.com` and `MENU_EMBED_ORIGINS=https://www.example.com`, using the actual WordPress origin.
3. Restart Go.
4. In WordPress select **Live Go embed**, enter the menu URL, a frame title and height (400–2400 pixels), then save.
5. Use the same shortcode on the page.

This keeps the complete current Go design and features, including updates made in the Go panel. It requires an online Go server. The iframe height is configured, not automatically synchronized across domains. A separate-tab link is included.

حالت زنده، منوی اصلی Go را داخل برگه نمایش می‌دهد؛ تغییرات پنل بلافاصله در نمایش بعدی دیده می‌شوند. سرور Go باید آنلاین باشد و دامنهٔ وردپرس در فهرست مجاز iframe قرار بگیرد. ارتفاع iframe دستی است.

## Troubleshooting

| Symptom | Check |
|---|---|
| Empty native menu | Import JSON; make sure categories are visible and dishes available |
| Frame refused | Check exact origin including www/port, restart Go, inspect proxy CSP/X-Frame-Options |
| Login fails on HTTP | Production uses Secure cookies; use HTTPS |
| Images missing | Set public source origin or host photos in WordPress; avoid HTTP images on HTTPS |
| Import rejected | Use version-1 MenuCraft JSON, correct types/whole prices and size limits |
| Old content | Native mode needs another import; also purge any WordPress page cache |
| Menu appears narrow | Choose a wide/full-width page layout in your WordPress theme |

Both integrations are display-only. No production WordPress account is required to build the ZIP. Installation on your live site is a separate administrator action.
