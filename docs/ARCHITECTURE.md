# Architecture and source review / معماری و بررسی سورس

## Application shape

A single Go HTTP process renders `html/template` pages, authenticates owners and reads/writes a SQLite database. `modernc.org/sqlite` keeps builds CGO-free. Public menus are server-rendered and progressively enhanced with vanilla JavaScript. Go embeds templates, static assets and migrations.

```text
main.go                  lifecycle, version flags and HTTP server
cmd/seed                 explicit demo fixture generator
cmd/release              four-platform archives, plugin ZIP and checksums
internal/config          validated environment configuration and signing key
internal/database        SQLite connection and transactional migrations
internal/auth            password hashes, signed sessions and CSRF
internal/repository      venue-scoped SQL and transaction-capable store
internal/models          domain models and canvas validation
internal/i18n            Persian/English strings, numbers and currency
internal/theme           presets, sanitization, CSS and page structure
internal/menupack         complete menu packs and example dishes
internal/menutemplate    category templates
internal/menudata         portable version-1 JSON schema
internal/server          routes, middleware, rendering and handlers
web/templates            public and panel templates
web/static               CSS, JavaScript and bundled font
wordpress/menucraft      native snapshot/live iframe WordPress plugin
```

## Request and data flow

The router wraps requests with security headers and bounded body parsing before handlers or CSRF middleware can parse them. Owner routes validate the signed session against the stored session version and select the venue through the authenticated user. Repository mutation queries include venue identifiers.

Passwords use salted PBKDF2 hashes with bounded parsing of stored work factors. Sessions are signed, HttpOnly, SameSite cookies, Secure in production. Changing a password increments `session_version`. CSRF values must be in the POST body.

A transaction-specific store wraps `sql.Tx`, so registration provisioning, menu-pack application, JSON import and settings updates can commit completely or roll back. Migrations apply both schema statements and migration tracking in a single transaction.

## Material issues corrected

| Finding | Correction |
|---|---|
| CSS `.bg-pattern` matched the whole body | Decorative layer rules now target the div; the menu stays opaque and interactive |
| Expired/revoked owner sessions could loop through login redirects | All login/public session resolution checks stored session version |
| Unpublished menu and QR visibility were inconsistent | Both enforce owner-only draft preview |
| Large multipart parsing could precede request limits | A body limit wraps routing and form parsing; temporary files are removed |
| Multi-step operations could leave partial changes | Transaction-scoped repositories and rollback integration tests |
| Host/forwarded-header assumptions affected links | Explicit canonical public origin for QR and JSON source URLs |
| Missing embedding policy for WordPress | Exact origin allowlist on public menu paths only |
| Loose number parsing accepted malformed values | Whole-unit, nonnegative bounded prices and validated optional metadata |
| SQLite DSN failed for certain Windows paths | Escaped file URI with drive-letter handling |
| Secret persistence failures could silently rotate sessions | Fail startup on invalid/unreadable/unwritable secrets |
| Phone formatting treated a string as a number | Preserve the escaped phone string |
| Unsafe HTML string construction in image picker | DOM element construction and text assignment |
| Non-finite canvas/theme input | Clamp or reset NaN/Inf; preserve valid opacity zero |
| Production demo credentials enabled by default | Demo population explicitly opt-in |
| Upload directory listing | Only generated image filename routes are served |

## Transfer contract

Version 1 contains source URL, export timestamp, venue descriptive fields/currency and category/dish content. Validation precedes mutation. Import appends by default; replacement requires an explicit selection. Go imports require matching currency and accept only image paths already owned by the receiving venue. Export includes hidden categories and unavailable dishes for editing/transfer, while WordPress native display omits them.

CSV is an export format only. Text cells beginning with formula-significant characters are prefixed to prevent spreadsheet formula execution. JSON contains no account credentials or session secrets. It does not include full settings, themes or image bytes.

## Review scope and remaining boundaries

This pass combined source inspection, integration tests, dependency vulnerability scanning and browser checks. It is not a claim that every possible security issue is excluded. Existing account, theme, designer and menu-pack features remain; advanced ordering, billing, recovery email and analytics services are outside this implementation.

SQLite and the in-memory limiter assume a single application instance. Favorites live only in each browser. WordPress native mode is a separate rendering layer with explicit limits described in its guide.

## توضیح فارسی

برنامه یک سرور Go با خروجی HTML سمت سرور، SQLite و جاوااسکریپت ساده است. منطق ذخیره‌سازی در repository، احراز هویت در auth و تولید ظاهر در theme قرار دارد. فایل‌های موردنیاز در نسخهٔ اجرایی جاسازی شده‌اند.

در بررسی، ایراد کلیک‌نشدن و نیمه‌شفاف‌شدن کل منو، نشست‌های باطل، نمایش پیش‌نویس، پردازش زودهنگام فایل بزرگ، عملیات چندمرحله‌ای غیرتراکنشی، مسیر SQLite در ویندوز و اعتبارسنجی قیمت اصلاح شدند. انتقال منو و افزونهٔ وردپرس به لایه‌های مستقل اضافه شده‌اند.

محدودیت‌ها صریح‌اند: JSON پشتیبان کامل نیست؛ افزونهٔ مستقل طراحی Go را بازسازی نمی‌کند؛ علاقه‌مندی‌ها محلی‌اند؛ سفارش و پرداخت وجود ندارد؛ استقرار چندنمونه‌ای روی یک فایل مشترک SQLite پشتیبانی نمی‌شود.
