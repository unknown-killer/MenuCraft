# Changelog / تاریخچهٔ تغییرات

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
