# MenuCraft 1.1.0 · منوساز

Persian/English restaurant menus with Go, SQLite and a new WordPress plugin.

## Downloads / فایل‌ها

- Windows amd64 / arm64: extract ZIP and run `run.bat`.
- Linux amd64 / arm64: extract tar.gz and run `./run.sh`.
- WordPress: upload `menucraft-wordpress_1.1.0.zip` through Plugins → Add New.
- Verify downloads against `SHA256SUMS.txt`.

## Highlights

- Search both languages, dietary filters, favorites, price sorting, sharing and print.
- JSON import/export and formula-safe CSV exports.
- Native WordPress snapshot menus or live embedding of the complete Go design.
- Fixed menu click/opacity bug, draft privacy, revoked sessions, transactional operations and input validation.
- Patched Go 1.26.9 toolchain, bundled Vazirmatn, four-platform builds and bilingual documentation.

## فارسی

این نسخه شامل فایل اجرایی ویندوز و لینوکس برای دو معماری و افزونهٔ قابل‌نصب وردپرس است. جست‌وجوی دوزبانه، فیلتر غذایی، علاقه‌مندی‌ها، مرتب‌سازی قیمت، چاپ و انتقال JSON/CSV اضافه شده‌اند. ایراد تعامل منو، نشست باطل، پیش‌نویس، عملیات ناقص و اعتبارسنجی قیمت نیز برطرف شده‌اند.

برای اجرای محلی پس از استخراج از اسکریپت اجرا استفاده کنید. برای اینترنت، راهنمای HTTPS را بخوانید. افزونهٔ مستقل محتوای منو را وارد می‌کند؛ برای نمایش دقیق طراحی Go از حالت زنده استفاده کنید.

## Upgrade and scope / ارتقا و محدوده

Back up the database, signing secret and uploads together before replacing binaries. Do not run the demo seed against real data. JSON exports are content transfers, not full backups. Native WordPress imports need manual refresh and preserve external image dependencies unless images are rehosted.

پیش از ارتقا، از پایگاه‌داده، کلید و تصاویر پشتیبان بگیرید. JSON پشتیبان کامل نیست. ورود وردپرس مستقل، خودکار همگام نمی‌شود و تصاویر تا انتقال به میزبان جدید به نشانی قبلی وابسته‌اند.

See README.md, README.fa.md and docs/ for setup, configuration, architecture and verification.
