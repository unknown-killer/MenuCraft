# Third-party notices

## Vazirmatn

Copyright 2015 The Vazirmatn Project Authors. Distributed under the SIL Open Font License, Version 1.1.

Source: https://github.com/rastikerdar/vazirmatn
The bundled variable WOFF2 was obtained from `fonts/webfonts/Vazirmatn[wght].woff2`, Git blob `a501289a85595158570b0c2badcb4608b042e748`.

The full license is included in `web/static/fonts/OFL.txt` and `wordpress/menucraft/assets/OFL.txt`. The font and its license are embedded into the Go application's static filesystem and included in the WordPress ZIP.

## Go dependencies

Dependency versions and integrity hashes are recorded in `go.mod` and `go.sum`. Their own licenses apply, including:

- modernc.org/sqlite and its transitive runtime dependencies: https://pkg.go.dev/modernc.org/sqlite
- github.com/skip2/go-qrcode: https://github.com/skip2/go-qrcode
- golang.org/x/*: https://go.dev/LICENSE

Run `go list -m all` for the complete dependency inventory. No license has been selected here for the original application's own source.

## JavaScript and PHP

The public-menu JavaScript and WordPress integration in this repository use browser/PHP/WordPress APIs directly, without vendored JavaScript frameworks. WordPress itself is a separate dependency for the plugin and is not bundled in releases.
