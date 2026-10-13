-- Core schema for the restaurant menu service.

CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT    NOT NULL UNIQUE,
    email           TEXT    NOT NULL UNIQUE,
    password_hash   TEXT    NOT NULL,
    display_name    TEXT    NOT NULL DEFAULT '',
    locale          TEXT    NOT NULL DEFAULT 'fa',
    -- Bumping this invalidates every previously issued session cookie.
    session_version INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS restaurants (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug            TEXT    NOT NULL UNIQUE,
    business_type   TEXT    NOT NULL DEFAULT 'restaurant',

    name_fa         TEXT    NOT NULL DEFAULT '',
    name_en         TEXT    NOT NULL DEFAULT '',
    tagline_fa      TEXT    NOT NULL DEFAULT '',
    tagline_en      TEXT    NOT NULL DEFAULT '',
    description_fa  TEXT    NOT NULL DEFAULT '',
    description_en  TEXT    NOT NULL DEFAULT '',

    logo_path       TEXT    NOT NULL DEFAULT '',
    cover_path      TEXT    NOT NULL DEFAULT '',

    phone           TEXT    NOT NULL DEFAULT '',
    address_fa      TEXT    NOT NULL DEFAULT '',
    address_en      TEXT    NOT NULL DEFAULT '',
    instagram       TEXT    NOT NULL DEFAULT '',
    telegram        TEXT    NOT NULL DEFAULT '',
    whatsapp        TEXT    NOT NULL DEFAULT '',
    website         TEXT    NOT NULL DEFAULT '',

    currency        TEXT    NOT NULL DEFAULT 'IRT',
    is_published    INTEGER NOT NULL DEFAULT 1,

    created_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_restaurants_user ON restaurants(user_id);

CREATE TABLE IF NOT EXISTS working_hours (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    day_of_week   INTEGER NOT NULL,             -- 0 = Saturday .. 6 = Friday
    opens_at      TEXT    NOT NULL DEFAULT '',
    closes_at     TEXT    NOT NULL DEFAULT '',
    is_closed     INTEGER NOT NULL DEFAULT 0,
    UNIQUE(restaurant_id, day_of_week)
);

CREATE TABLE IF NOT EXISTS categories (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    restaurant_id  INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name_fa        TEXT    NOT NULL DEFAULT '',
    name_en        TEXT    NOT NULL DEFAULT '',
    description_fa TEXT    NOT NULL DEFAULT '',
    description_en TEXT    NOT NULL DEFAULT '',
    icon           TEXT    NOT NULL DEFAULT '',   -- emoji or short glyph
    image_path     TEXT    NOT NULL DEFAULT '',
    sort_order     INTEGER NOT NULL DEFAULT 0,
    is_visible     INTEGER NOT NULL DEFAULT 1,
    created_at     TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at     TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_categories_restaurant ON categories(restaurant_id, sort_order);

CREATE TABLE IF NOT EXISTS items (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    restaurant_id  INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    category_id    INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,

    name_fa        TEXT    NOT NULL DEFAULT '',
    name_en        TEXT    NOT NULL DEFAULT '',
    description_fa TEXT    NOT NULL DEFAULT '',
    description_en TEXT    NOT NULL DEFAULT '',

    price          INTEGER NOT NULL DEFAULT 0,   -- minor-unit free: stored as whole currency units
    discount_price INTEGER NOT NULL DEFAULT 0,   -- 0 means "no discount"
    image_path     TEXT    NOT NULL DEFAULT '',

    tags           TEXT    NOT NULL DEFAULT '[]', -- JSON array: spicy, vegan, new, popular...
    is_available   INTEGER NOT NULL DEFAULT 1,
    is_featured    INTEGER NOT NULL DEFAULT 0,
    calories       INTEGER NOT NULL DEFAULT 0,
    prep_minutes   INTEGER NOT NULL DEFAULT 0,
    sort_order     INTEGER NOT NULL DEFAULT 0,

    created_at     TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at     TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_items_category ON items(category_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_items_restaurant ON items(restaurant_id);

-- The full visual configuration lives as JSON so the editor can evolve without
-- a migration for every new knob.
CREATE TABLE IF NOT EXISTS themes (
    restaurant_id INTEGER PRIMARY KEY REFERENCES restaurants(id) ON DELETE CASCADE,
    preset        TEXT NOT NULL DEFAULT 'midnight',
    config        TEXT NOT NULL DEFAULT '{}',
    updated_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS media (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    path          TEXT    NOT NULL,
    original_name TEXT    NOT NULL DEFAULT '',
    size_bytes    INTEGER NOT NULL DEFAULT 0,
    mime_type     TEXT    NOT NULL DEFAULT '',
    created_at    TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_media_restaurant ON media(restaurant_id, created_at DESC);
