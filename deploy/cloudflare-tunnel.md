# اجرا روی Cloudflare — با Tunnel

> **چرا Workers نه؟**
> Cloudflare Workers کد جاوااسکریپت/WASM را روی V8 اجرا می‌کند: فایل‌سیستم ندارد،
> پروسه‌ی دائمی نگه نمی‌دارد و کتابخانه‌های Go که به دیسک نیاز دارند روی آن اجرا
> نمی‌شوند. این پروژه یک باینری Go با دیتابیس SQLite روی دیسک و آپلود فایل است،
> پس روی Workers بالا نمی‌آید — برای آن باید کل بک‌اند با JS و D1 بازنویسی شود.
>
> **Cloudflare Tunnel** اما دقیقاً همان چیزی را می‌دهد که می‌خواهی: یک آدرس
> عمومی HTTPS روی شبکه‌ی Cloudflare، بدون باز کردن هیچ پورتی و بدون سرور.

---

## راه ۱ — تست فوری (بدون حساب، ۲ دقیقه)

برای وقتی که فقط می‌خواهی یک لینک موقت بدهی دست بازاریاب.

**۱. برنامه را اجرا کن:**

```bash
go run .
```

**۲. در یک ترمینال دیگر، `cloudflared` را بگیر و تونل بزن:**

```bash
cloudflared tunnel --url http://localhost:8080
```

`cloudflared` یک آدرس مثل `https://random-words-1234.trycloudflare.com` چاپ
می‌کند. همان را بفرست.

> نکته: این آدرس تا وقتی ترمینال باز است کار می‌کند. لپ‌تاپ که خاموش شد، لینک
> هم می‌میرد. برای تست کوتاه عالی است، برای تحویل واقعی نه.

### نصب cloudflared روی ویندوز

```bash
winget install --id Cloudflare.cloudflared
```

---

## راه ۲ — تونل دائمی با دامنه‌ی خودت

برای وقتی که می‌خواهی آدرس ثابت بماند (مثلاً `menu.yourdomain.com`).

**پیش‌نیاز:** یک دامنه که در Cloudflare ثبت شده باشد.

```bash
# یک بار وارد شو
cloudflared tunnel login

# تونل بساز
cloudflared tunnel create menucraft

# دامنه را به تونل وصل کن
cloudflared tunnel route dns menucraft menu.yourdomain.com
```

بعد فایل `deploy/cloudflared-config.yml` را بردار، `<TUNNEL-ID>` و دامنه را
عوض کن و بگذارش در `~/.cloudflared/config.yml` (روی ویندوز:
`%USERPROFILE%\.cloudflared\config.yml`)، سپس:

```bash
cloudflared tunnel run menucraft
```

برای اینکه با روشن شدن سیستم خودکار بالا بیاید:

```bash
cloudflared service install
```

---

## نکته‌ی مهم درباره‌ی کوکی‌ها

وقتی پشت تونل هستی، ترافیک HTTPS است، پس حالت production را روشن کن تا
کوکی‌ها `Secure` شوند:

```bash
MENU_ENV=production go run .
```

روی ویندوز در PowerShell:

```powershell
$env:MENU_ENV="production"; go run .
```

اگر بدون HTTPS (مستقیم روی `http://localhost`) با `MENU_ENV=production` اجرا
کنی، مرورگر کوکی `Secure` را نمی‌پذیرد و ورود کار نمی‌کند. برای تست محلی
`MENU_ENV` را ست نکن.
