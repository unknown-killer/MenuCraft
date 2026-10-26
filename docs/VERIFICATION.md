# Verification / گزارش بررسی

## Version 1.2.0 — 2026-10-26

- Windows: `go test ./...` and `go vet ./...` passed with Go 1.26.10.
- `govulncheck@v1.9.0`: no vulnerabilities found after upgrading from Go 1.26.9 to 1.26.10 for GO-2026-6809.
- Node notebook tests (4) passed, including exact totals, input bounds, unavailable items and byte-identical WordPress assets.
- PHP plugin syntax and changed JavaScript syntax checks passed.
- Static demo generation produced 114 themes in both languages.
- New HTTP/SQLite tests cover table QR batches, draft access, pending-price save/export/import/manual pricing, canonical metadata, sitemap, default language, saved landing-page details and the single-venue WordPress embed boundary.
- Release automation runs native Windows/Linux tests, builds four platform targets and checks the packaged Linux executable readiness before publishing. Check the matching Actions run for final artifact evidence. ARM64 targets are cross-compiled, not hardware-tested.

<div dir="rtl">

آزمون‌های برنامه، دفترچه، نحو افزونه و اسکن امنیتی موفق بودند. نسخهٔ Go برای رفع آسیب‌پذیری کتابخانهٔ استاندارد ارتقا یافت. فایل‌های خصوصی استقرار، داده‌های مشتری، کلیدها و تصاویر اختصاصی رستوران در انتشار عمومی قرار ندارند. جزئیات اجرای ساخت بسته‌ها در گزارش همان ریلیز ثبت می‌شود.

</div>

## Historical checks: version 1.1.0 — 2026-10-13

The following records describe the previous release preparation, not a new WordPress runtime certification.

## Automated checks

- `go test ./...`: passed on Windows with Go 1.26.9, including real SQLite HTTP integration fixtures.
- `go vet ./...`: passed.
- `govulncheck v1.9.0` with Go 1.26.9: no known reachable vulnerabilities reported at scan time.
- PHP 8.4.26 syntax check: passed.
- Node syntax checks for changed Go-panel/public JavaScript and WordPress JavaScript: passed.

Integration tests exercise both languages, login/CSRF, revoked sessions, private menu/QR access, tenant boundaries, request-size enforcement, rollback on registration/menu-pack failures, JSON round trips, invalid import preservation, canonical URLs and database readiness.

## Browser checks

Using an isolated 21-dish, 6-category demo database:

- Persian and English public menus and panel templates render.
- Arabic `كباب` search returns five matching dishes; clearing returns 21.
- Vegetarian filter returns five; favorite selection persists after reloading.
- At 390 CSS pixels, the menu has no horizontal overflow.
- Body opacity is 1 and pointer events are active after the decorative-layer fix.
- WordPress 7.1.3 / PHP 8.4 in local WordPress Playground activated the real plugin, imported the Go JSON and rendered all 21 dishes.
- Native WordPress search returned five kebab dishes; dessert filtering returned three.
- Live iframe mode loaded the Go menu successfully with the configured WordPress origin.
- The checked public menu reported no JavaScript console errors.

## Release verification

The repository release workflow runs native Go tests on Windows and Linux before publishing. It builds all four platform/architecture combinations, runs a Linux amd64 executable readiness smoke check, packages the plugin and publishes SHA-256 checksums. The workflow log attached to each GitHub Actions run is the authoritative evidence for that run.

ARM64 files are cross-compiled; they have not been executed on physical ARM64 hardware in this review. Docker/systemd examples were reviewed but not deployed to a production server. No live customer WordPress installation was changed.

## فارسی

تست‌های Go، بررسی vet، اسکن وابستگی‌ها و بررسی نحو PHP/JavaScript موفق بودند. منوی آزمایشی در مرورگر، جست‌وجوی دو زبان، فیلتر، ذخیرهٔ علاقه‌مندی و عرض موبایل بررسی شد. افزونه نیز روی وردپرس واقعی، هم ورود مستقل JSON و هم iframe زنده را نمایش داد.

گردش‌کار انتشار، تست بومی ویندوز و لینوکس و بررسی اجرای فایل لینوکس را پیش از انتشار انجام می‌دهد؛ نتیجهٔ هر اجرا در GitHub Actions ثبت می‌شود. نسخه‌های ARM64 خروجی کامپایل هستند و روی سخت‌افزار ARM64 اجرا نشده‌اند. استقرار عملی روی سایت واقعی کاربر انجام نشده است.
