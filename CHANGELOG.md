# Changelog / تاریخچهٔ تغییرات

## 1.2.0 — 2026-10-26

### فارسی

- دفترچهٔ مشتری: انتخاب غذا، تعداد، یادداشت هر مورد و یادداشت کلی، جمع تقریبی، نمای گارسون، کپی و چاپ.
- حفظ انتخاب‌ها در مرورگر برای هر مجموعه و هنگام تعویض زبان؛ نمایش شمارهٔ میز از لینک اسکن‌شده.
- سازندهٔ استند در تنظیمات: پنج چیدمان، شمارهٔ تکی یا بازه، چهار کارت معمولی یا دو استند تاشو در هر برگهٔ آ۴.
- دفترچه در افزونهٔ وردپرس مستقل و نمایش زنده؛ راهنمای فارسی و انگلیسی.
- قیمت نامشخص با مهاجرت خودکار ۰۰۴؛ حفظ وضعیت در ورود/خروج داده و وردپرس.
- چاپ با رنگ و پس‌زمینهٔ قالب، صفحهٔ معرفی تک‌رستوران، زبان پیش‌فرض فارسی، نقشهٔ سایت و دادهٔ ساختاریافته.
- چند شمارهٔ تماس، حفظ لوگو هنگام ذخیرهٔ تنظیمات و مرتب‌سازی قیمت نامشخص در انتهای فهرست.
- دفترچه سفارش آنلاین ارسال نمی‌کند.

### English

- Go 1.26.10 toolchain fixes GO-2026-6809 in the standard library.

- Guest notebook with quantities, item/general notes, estimated totals, waiter presentation, copy fallback and isolated printing.
- Restaurant-scoped browser storage, language continuity and table prefill from QR links.
- Owner stand generator: five layouts, single or batched table numbers, four flat cards or two folded tents per A4 sheet.
- Native WordPress notebook support; preserved table numbers in language and share links.
- Bounded input validation, current-price totals with exact arithmetic, safe text rendering and invalid/unavailable item recovery.
- Pending prices with automatic migration 004; JSON/CSV and WordPress support.
- Theme-preserving printing with image/font readiness; single-venue landing page, Persian default, sitemap and search metadata.
- Multiple phone links, logo preservation on settings save and unknown-price sorting.
- No order submission; back up before upgrading.

## 1.1.0 — 2026-10-13

### Added / افزوده‌ها
- Native WordPress JSON menus and live iframe integration.
- JSON import/export, safe CSV export and a bilingual transfer page.
- Public dietary filters, browser favorites, within-category price sorting, share, print and result announcements.
- Self-hosted Vazirmatn font and reduced-motion support.
- Database readiness endpoint, canonical public URL, explicit embed origins and registration switch.
- Windows/Linux amd64/arm64 packaging, checksums and bilingual deployment documentation.

### Fixed / اصلاحات
- Background decoration no longer disables pointer events or fades the whole page.
- Revoked sessions no longer cause login redirects; private menus and QR codes respect ownership.
- Request limits apply before multipart/CSRF parsing; account endpoints have bounded peer rate limiting.
- Registration, menu-pack replacement, import, settings and migration operations are atomic.
- Numeric validation handles Persian digits and rejects negative, fractional or malformed prices.
- Windows SQLite paths containing spaces or reserved URI characters open correctly.
- Session key failures stop startup; unsafe redirects, query-string CSRF tokens and upload directory listings are rejected.
- Image picker rendering uses DOM text instead of string HTML; phone numbers display as text.
- Theme/canvas sanitizers reject non-finite numbers; fully transparent canvas elements remain transparent.
- Demo seeding is opt-in in production. Go toolchain updated to 1.26.9.

### Existing capabilities retained / امکانات قبلی حفظ‌شده
Bilingual owners/panel, 114 themes, 54 complete menu packs, category/item editing, media library, section/block editing, free-positioned canvas, QR, venue settings and opening hours.

### Scope
No order/payment backend or automatic WordPress synchronization was added. Menu JSON is not a full application backup.
