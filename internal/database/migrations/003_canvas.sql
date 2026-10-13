-- Free-form canvas sections: absolutely positioned elements stored as JSON.
ALTER TABLE blocks ADD COLUMN canvas TEXT NOT NULL DEFAULT '{}';
