# Build stage ---------------------------------------------------------------
FROM golang:1.26.10-alpine AS build

ENV CGO_ENABLED=0

WORKDIR /src

# Dependencies first so a code change does not re-download the module cache.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# The SQLite driver is pure Go, so a fully static binary needs no C toolchain
# and the final image needs no libc.
RUN go build -trimpath -ldflags="-s -w" -o /out/menu . \
 && go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

# Runtime stage -------------------------------------------------------------
FROM alpine:3.23

# tzdata so opening-hours display in the venue's own timezone;
# ca-certificates for any outbound HTTPS a future integration might need.
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -u 10001 menu

COPY --from=build /out/menu /usr/local/bin/menu
COPY --from=build /out/seed /usr/local/bin/seed
COPY deploy/entrypoint.sh /usr/local/bin/entrypoint.sh
# Strip CR before making it executable. The script is authored on Windows, and
# a CRLF shebang fails in the container with "/bin/sh^M: bad interpreter" —
# normalising here means the build does not depend on the host's git settings.
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh \
 && chmod +x /usr/local/bin/entrypoint.sh

# Everything mutable lives here, so this is the only path that needs a volume.
RUN mkdir -p /data/db /data/uploads && chown -R menu:menu /data
VOLUME ["/data"]

ENV MENU_ENV=production \
    MENU_ADDR=:8080 \
    MENU_DB=/data/db/menu.db \
    MENU_UPLOADS=/data/uploads \
    MENU_LOCALE=fa \
    TZ=Asia/Tehran

USER menu
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
