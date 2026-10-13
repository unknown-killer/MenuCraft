-- Free-form content blocks the owner can drop anywhere on the public page:
-- headings, paragraphs, images, banners, dividers and spacers.

CREATE TABLE IF NOT EXISTS blocks (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    -- text | heading | image | banner | gallery | divider | spacer | quote | button
    kind          TEXT    NOT NULL DEFAULT 'text',

    title_fa      TEXT    NOT NULL DEFAULT '',
    title_en      TEXT    NOT NULL DEFAULT '',
    body_fa       TEXT    NOT NULL DEFAULT '',
    body_en       TEXT    NOT NULL DEFAULT '',

    image_path    TEXT    NOT NULL DEFAULT '',
    -- Extra images for gallery blocks, as a JSON array of upload paths.
    images        TEXT    NOT NULL DEFAULT '[]',

    link_url      TEXT    NOT NULL DEFAULT '',
    link_label_fa TEXT    NOT NULL DEFAULT '',
    link_label_en TEXT    NOT NULL DEFAULT '',

    -- Presentation options as JSON so new knobs need no migration.
    style         TEXT    NOT NULL DEFAULT '{}',

    sort_order    INTEGER NOT NULL DEFAULT 0,
    is_visible    INTEGER NOT NULL DEFAULT 1,

    created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_blocks_restaurant ON blocks(restaurant_id, sort_order);
