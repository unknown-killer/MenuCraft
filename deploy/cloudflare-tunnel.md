# Cloudflare Tunnel / تونل کلادفلر

MenuCraft still runs on your computer or server with persistent SQLite storage. A tunnel forwards public HTTPS traffic to that running process; it does not host the application or keep an offline computer available.

برنامه همچنان روی رایانه یا سرور شما اجرا می‌شود. تونل فقط ترافیک HTTPS را به آن می‌رساند؛ خاموش‌شدن میزبان، منو را از دسترس خارج می‌کند.

For a temporary local preview, run MenuCraft on `127.0.0.1:8080`, then use the official cloudflared client:

```sh
cloudflared tunnel --url http://127.0.0.1:8080
```

The generated URL is temporary. To make QR/export links correct, restart MenuCraft with `MENU_PUBLIC_URL` set to that exact HTTPS origin. Never use this preview with the default demo password on an unattended public deployment.

For a stable deployment, configure a named tunnel in your Cloudflare account and route its hostname to your running MenuCraft service. The Compose `tunnel` profile accepts `CLOUDFLARE_TUNNEL_TOKEN`; its tunnel route should point to `http://menu:8080` inside the Compose network. Keep the token outside Git.

Set `MENU_ENV=production`, the real `MENU_PUBLIC_URL`, and, if required, `MENU_EMBED_ORIGINS`. Follow [the deployment guide](../docs/DEPLOYMENT.md) for persistent storage, account setup and backups. `cloudflared-config.yml` is a template for a locally managed tunnel; replace its placeholder ID, credential path and hostname.

برای استفادهٔ دائمی، تونل دارای نام و دامنهٔ خود را تنظیم کنید. کلید تونل را وارد گیت نکنید. تنظیمات production، دامنهٔ اصلی، داده‌های ماندگار و پشتیبان‌گیری همچنان لازم‌اند.
