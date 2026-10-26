# 🧪 Live demo

The [public demo](https://unknown-killer.github.io/MenuCraft/) is generated from the real Go public-menu renderer. It includes all 114 presets in Persian and English, plus a static design studio for colors, layouts, columns, card radius and category ordering. Guests can use the actual public-menu search, dietary filters and favorites.

Studio settings are saved in browser storage. They do not modify a restaurant database. The full owner panel, uploads, menu-content editing and freeform canvas designer require a running Go application. GitHub Pages hosts static files and cannot run that backend.

## Build locally

Run from the repository root with Go 1.26.10 or newer:

```sh
go run ./cmd/demo -output _site
python -m http.server 8765 --bind 127.0.0.1 --directory _site
```

Open `http://127.0.0.1:8765/`. Python is only an optional static preview server. The exporter creates and seeds its own temporary database; it ignores inherited `MENU_*` settings, closes and removes the temporary database, and never reads the operator's database or uploaded files. Generated output in `_site/` is ignored by Git.

The exporter calls `theme.ApplyPreset` so each preview uses the preset's intended light/dark palette. Public styles, scripts and local fonts are copied from the application. Contact links from the sample seed are removed before export.

## Publish

Repository Settings → Pages → Build and deployment → Source must be **GitHub Actions**. The `Publish live menu demo` workflow validates Go and JavaScript, renders 228 menu pages, uploads a Pages artifact and deploys it. Updates to demo source and public menu templates rebuild the site. The workflow also refreshes the existing 1.1.0 release description from `docs/RELEASE_NOTES.md`; it does not replace release binaries or move the version tag.

## Manual checks

- Switch Persian/English and confirm reading direction changes.
- Select a light and a dark preset; confirm their colors and layout differ.
- Change colors and layout, reorder categories by dragging and by buttons, then reload to check persistence.
- Search and filter inside the menu preview; save a favorite.
- Check phone and desktop widths and keyboard focus.
- Confirm that the standalone sample links and local font/styles load under the `/MenuCraft/` project prefix.
