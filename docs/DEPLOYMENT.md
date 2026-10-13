# Deployment / استقرار

## Local versus public

For local testing use `MENU_ADDR=127.0.0.1:8080` and the default development environment. For an internet-facing installation, terminate HTTPS at a reverse proxy, use `MENU_ENV=production` and set `MENU_PUBLIC_URL` to the real HTTPS origin. Production session cookies require HTTPS.

برای اجرای محلی، از نشانی لوکال استفاده کنید. برای دسترسی اینترنتی، HTTPS، حالت production و دامنهٔ اصلی را تنظیم کنید. کوکی ورود در حالت production روی HTTP معمولی ارسال نمی‌شود.

## Linux service

Extract the Linux archive under `/opt/menucraft`. Create a dedicated unprivileged `menucraft` service user and make `/var/lib/menucraft` writable by that user. Example systemd unit at `/etc/systemd/system/menucraft.service`:

```ini
[Unit]
Description=MenuCraft restaurant menus
After=network.target

[Service]
Type=simple
User=menucraft
Group=menucraft
WorkingDirectory=/var/lib/menucraft
ExecStart=/opt/menucraft/menucraft
Environment=MENU_ADDR=127.0.0.1:8080
Environment=MENU_ENV=production
Environment=MENU_DB=/var/lib/menucraft/data/menu.db
Environment=MENU_UPLOADS=/var/lib/menucraft/uploads
Environment=MENU_PUBLIC_URL=https://menu.example.com
Restart=on-failure
RestartSec=3
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=/var/lib/menucraft

[Install]
WantedBy=multi-user.target
```

Replace the domain, install your TLS proxy, then run `systemctl daemon-reload` and `systemctl enable --now menucraft`. Review logs with `journalctl -u menucraft`.

Example Caddy configuration (Caddy obtains HTTPS for a correctly configured domain):

```caddy
menu.example.com {
    reverse_proxy 127.0.0.1:8080
}
```

The service user and directories must exist before starting. Keep the database outside the executable directory so application upgrades do not replace it.

کاربر سرویس و پوشهٔ داده را پیش از اجرا بسازید. دامنه باید به سرور اشاره کند. فایل‌های داده را بیرون از پوشهٔ نسخه‌های اجرایی نگه دارید تا هنگام ارتقا جایگزین نشوند.

## Docker Compose

From the source repository:

```sh
export MENU_PUBLIC_URL=https://menu.example.com
docker compose -f deploy/docker-compose.yml up -d --build
```

Compose exposes only `127.0.0.1:8080`; configure an HTTPS reverse proxy on the host. A named `menu_data` volume preserves `/data/db/menu.db`, `/data/db/.secret` and `/data/uploads`. The process runs as UID 10001. The optional `tunnel` profile uses a user-supplied Cloudflare tunnel token and should route to `http://menu:8080`.

Set `MENU_DISABLE_REGISTRATION=true` after registering the required owners, then recreate the service. Demo seeding remains off unless explicitly enabled with `MENU_SEED_DEMO=1`; never enable the documented demo credentials on a public installation.

نسخهٔ Docker داده‌ها را در volume نگه می‌دارد و تنها پورت لوکال را باز می‌کند. پس از ساخت حساب مالک، در صورت نیاز ثبت‌نام را ببندید. حالت نمونه را روی سرور عمومی فعال نکنید.

## Windows

Use the release `run.bat` for local access. For an internet deployment, run `menucraft.exe` under your chosen Windows service manager with a writable working directory, explicit environment variables, persistent data and an HTTPS proxy.

Example PowerShell local configuration:

```powershell
$env:MENU_ADDR='127.0.0.1:8080'
$env:MENU_DB='C:\\MenuCraftData\\menu.db'
$env:MENU_UPLOADS='C:\\MenuCraftData\\uploads'
.\\menucraft.exe
```

## WordPress embedding

```sh
MENU_PUBLIC_URL=https://menu.example.com
MENU_EMBED_ORIGINS=https://www.example.com,https://example.com
```

Set these as actual service/container environment variables and restart. Origins must include scheme and optional port, with no path, wildcard, credentials or query. Only public `/m/` pages become embeddable; account and panel routes retain same-origin framing restrictions. Both sites need HTTPS in production.

## Persistence, backup and restore

A complete backup includes:

1. The SQLite database and any `-wal` / `-shm` files present.
2. The complete upload directory.
3. The `.secret` beside the database, or the separately managed `MENU_SECRET`.

For a straightforward consistent backup, **stop MenuCraft**, copy the data and uploads together to a protected backup location, and start it again. Do not copy only a live SQLite `.db` file while writes are occurring. Container users should stop the service and back up the named volume using their normal volume backup tooling.

Restore while the service is stopped, preserve owner permissions and the signing key, then start and check `/readyz`, an owner login and a published menu. Test restoration periodically. Changing the signing key invalidates all existing sessions.

برای نسخهٔ پشتیبان ساده و سازگار، برنامه را متوقف کنید و پایگاه‌داده، فایل‌های جانبی موجود، کلید امضا و تصاویر را با هم کپی کنید. خروجی JSON جایگزین این پشتیبان نیست. بازیابی را در حالت توقف سرویس انجام دهید.

## Updates and monitoring

Back up first, stop the service, replace only the executable/assets package, then restart. Schema migrations run in transactions at startup. Keep the previous executable and backup for recovery; restoring an older executable alone does not undo a schema change.

- `GET /healthz`: process responds.
- `GET /readyz`: database ping succeeds.
- `menucraft --version`: version, source commit and build timestamp.
- Logs include requests, startup and shutdown information. Store logs with appropriate access controls.

Login/registration limits use the network peer address: 20 attempts per minute shared across these endpoints, with a bounded in-memory client map. Forwarded headers are deliberately not used to establish identity. Behind a proxy, requests share its peer limit; also configure appropriate per-client rate limits at your trusted proxy. This single-process limiter resets on restart.

## Operating boundaries

Run one application instance per database with persistent local storage. Do not place the SQLite file on network/shared filesystems or expect multi-instance failover. Uploaded images are public assets, even when a menu is a draft; do not upload confidential files. The request limit for forms and JSON transfer is 2 MiB, with a separate configured image-upload allowance.

Money values are nonnegative whole units in the selected currency (for example whole toman), not fractional major units. Review example-pack prices and dietary tags before publishing.

مقادیر پولی در این نسخه عدد صحیح در واحد انتخاب‌شده‌اند. قیمت‌های نمونه و برچسب‌های غذایی را پیش از انتشار بررسی کنید. تصاویر بارگذاری‌شده فایل عمومی هستند، حتی اگر منو هنوز پیش‌نویس باشد.
