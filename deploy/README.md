# راه‌اندازی برای تست

هدف: یک لینک عمومی که بدهی دست بازاریاب و کار کند.

---

## اول یک نکته‌ی مهم

**Cloudflare Workers نمی‌تواند این پروژه را اجرا کند.**

Workers کد جاوااسکریپت/WASM را روی V8 اجرا می‌کند — فایل‌سیستم ندارد، پروسه‌ی
دائمی نگه نمی‌دارد، و SQLite روی دیسک یا آپلود فایل روی آن معنا ندارد. این
پروژه یک باینری Go است با دیتابیس فایلی و پوشه‌ی آپلود. برای Workers باید کل
بک‌اند از صفر با JS و D1 و R2 نوشته شود؛ یعنی پروژه‌ی دیگری.

سه راهی که واقعاً کار می‌کنند در ادامه آمده. هر سه HTTPS می‌دهند.

---

## راه ۱ — سریع‌ترین: Cloudflare Tunnel

**۲ دقیقه، بدون حساب، بدون هزینه.** روی شبکه‌ی Cloudflare هم هست.

```bash
go run .
```

در ترمینال دوم:

```bash
cloudflared tunnel --url http://localhost:8080
```

آدرسی مثل `https://xyz-1234.trycloudflare.com` می‌دهد — همان را بفرست.

**محدودیت:** تا وقتی لپ‌تاپت روشن و ترمینال باز است کار می‌کند.

جزئیات بیشتر و تونل دائمی با دامنه‌ی خودت: [`cloudflare-tunnel.md`](cloudflare-tunnel.md)

---

## راه ۲ — لینک دائمی: Fly.io

**۵ دقیقه، رایگان برای این حجم.** بعد از این دیگر به سیستم تو وابسته نیست.

```bash
# یک بار
fly auth signup

# fly.toml را به ریشه‌ی پروژه کپی کن
cp deploy/fly.toml .

fly launch --no-deploy --copy-config --name YOUR-APP-NAME
fly volumes create menu_data --size 1 --region fra
fly deploy
```

آدرس نهایی: `https://YOUR-APP-NAME.fly.dev`

> **حتماً volume بساز.** بدون آن، هر بار deploy دیتابیس از صفر شروع می‌شود و
> همه‌ی منوها می‌پرند.

منطقه‌ها: `fra` (فرانکفورت) و `ams` (آمستردام) از ایران کم‌تأخیرترند.

---

## راه ۳ — سرور شخصی: Docker

روی هر VPS با Docker:

```bash
docker compose -f deploy/docker-compose.yml up -d --build
```

اگر می‌خواهی بدون باز کردن پورت، از طریق Cloudflare منتشر شود: در داشبورد
Cloudflare یک تونل بساز، توکنش را بگذار در `.env` کنار compose، و بخش `ports`
را کامنت کن:

```bash
echo 'CLOUDFLARE_TUNNEL_TOKEN=your-token-here' > deploy/.env
docker compose -f deploy/docker-compose.yml --profile tunnel up -d --build
```

---

## آنچه بازاریاب می‌بیند

روی یک دیتابیس نو، رستوران نمونه خودکار ساخته می‌شود:

| | |
|---|---|
| منوی عمومی | `/m/demo` |
| ورود به پنل | `/login` |
| نام کاربری | `demo` |
| رمز عبور | `Demo1234` |

**قبل از تحویل حتماً رمز را عوض کن** (پنل ← تنظیمات ← حساب کاربری)، یا اگر
نمی‌خواهی داده‌ی نمونه ساخته شود `MENU_SEED_DEMO=0` بگذار.

پیشنهاد مسیر نمایش برای بازاریاب:

1. `/m/demo` — منوی نهایی که مشتری رستوران می‌بیند
2. `/panel/gallery` — ۵۴ منوی آماده؛ یکی را بزند و ببیند کل منو ساخته می‌شود
3. `/panel/theme` — ۱۱۴ قالب و پیش‌نمایش زنده
4. `/panel/design/canvas/new` — طراح آزاد با درگ‌اند‌دراپ

---

## تنظیمات محیطی

| متغیر | پیش‌فرض | توضیح |
|---|---|---|
| `MENU_ADDR` | `:8080` | آدرس و پورت |
| `MENU_ENV` | `development` | در `production` کوکی‌ها `Secure` و کش فعال می‌شود |
| `MENU_DB` | `data/menu.db` | مسیر دیتابیس |
| `MENU_UPLOADS` | `uploads/` | پوشه‌ی تصاویر |
| `MENU_SEED_DEMO` | `1` | ساخت داده‌ی نمونه فقط روی دیتابیس نو |
| `MENU_MAX_UPLOAD` | `8388608` | حداکثر حجم آپلود |
| `MENU_SESSION_TTL` | `720h` | مدت اعتبار نشست |
| `TZ` | `Asia/Tehran` | منطقه‌ی زمانی، برای ساعات کاری |

`MENU_ENV=production` را فقط پشت HTTPS بگذار. روی `http://localhost` مرورگر
کوکی `Secure` را رد می‌کند و ورود کار نمی‌کند.

---

## پشتیبان‌گیری

همه‌چیز در یک پوشه است: دیتابیس، تصاویر آپلودی، و کلید امضای نشست.

```bash
# Docker
docker run --rm -v menucraft_menu_data:/data -v "$PWD":/backup alpine \
  tar czf /backup/menu-backup.tar.gz -C /data .

# Fly.io
fly ssh console -C "tar czf - -C /data ." > menu-backup.tar.gz
```

فایل `/data/db/.secret` را نگه دار — اگر گم شود همه‌ی کاربران یک بار از حساب
خارج می‌شوند (داده‌ای از دست نمی‌رود).

---

## محدودیت‌هایی که باید بدانی

**یک ماشین، نه بیشتر.** SQLite یک نویسنده دارد. اگر دو نمونه بالا بیاوری، هر
کدام روی volume خودش می‌نویسد و داده‌ها از هم جدا می‌شوند. برای مقیاس واقعی باید
به Postgres مهاجرت کرد — لایه‌ی repository برای همین جدا نوشته شده.

**آپلودها روی دیسک محلی‌اند.** برای چند سرور باید به S3/R2 منتقل شوند.

**نرخ‌محدودسازی ندارد.** برای دموی محدود مشکلی نیست؛ قبل از انتشار عمومی باید
روی `/login` و `/register` محدودیت گذاشت (Cloudflare WAF هم می‌تواند).
