ALTER TABLE items ADD COLUMN price_pending INTEGER NOT NULL DEFAULT 0 CHECK(price_pending IN (0,1));
