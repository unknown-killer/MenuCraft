package repository

import (
	"ResturantMenue/internal/models"
	"ResturantMenue/internal/theme"
)

// --- themes ----------------------------------------------------------------

type ThemeRepo struct{ db dbQuery }

// Get returns the stored theme, or the default one when the venue has never
// been customised. It never fails just because a row is missing.
func (r *ThemeRepo) Get(restaurantID int64) (string, theme.Config, error) {
	var preset, raw string
	err := r.db.QueryRow(`SELECT preset, config FROM themes WHERE restaurant_id = ?`,
		restaurantID).Scan(&preset, &raw)
	if err != nil {
		if translateErr(err) == ErrNotFound {
			return "midnight", theme.Default(), nil
		}
		return "", theme.Config{}, err
	}
	return preset, theme.Parse(raw), nil
}

// Save upserts the whole theme document.
func (r *ThemeRepo) Save(restaurantID int64, preset string, cfg theme.Config) error {
	cfg.Sanitize()
	_, err := r.db.Exec(`
		INSERT INTO themes (restaurant_id, preset, config, updated_at)
		VALUES (?, ?, ?, datetime('now'))
		ON CONFLICT(restaurant_id) DO UPDATE SET
			preset = excluded.preset,
			config = excluded.config,
			updated_at = excluded.updated_at`,
		restaurantID, preset, cfg.JSON())
	return translateErr(err)
}

// --- media -----------------------------------------------------------------

type MediaRepo struct{ db dbQuery }

func (r *MediaRepo) Create(m *models.Media) error {
	res, err := r.db.Exec(`
		INSERT INTO media (restaurant_id, path, original_name, size_bytes, mime_type)
		VALUES (?, ?, ?, ?, ?)`,
		m.RestaurantID, m.Path, m.OriginalName, m.SizeBytes, m.MimeType)
	if err != nil {
		return translateErr(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

// List returns the venue's uploads, newest first, for the media picker.
func (r *MediaRepo) List(restaurantID int64, limit int) ([]*models.Media, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.db.Query(`
		SELECT id, restaurant_id, path, original_name, size_bytes, mime_type, created_at
		FROM media WHERE restaurant_id = ? ORDER BY created_at DESC, id DESC LIMIT ?`,
		restaurantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Media
	for rows.Next() {
		var m models.Media
		var created string
		if err := rows.Scan(&m.ID, &m.RestaurantID, &m.Path, &m.OriginalName,
			&m.SizeBytes, &m.MimeType, &created); err != nil {
			return nil, err
		}
		m.CreatedAt = parseTime(created)
		out = append(out, &m)
	}
	return out, rows.Err()
}

// ByPath looks a file up so deletes can verify ownership before touching disk.
func (r *MediaRepo) ByPath(restaurantID int64, path string) (*models.Media, error) {
	var m models.Media
	var created string
	err := r.db.QueryRow(`
		SELECT id, restaurant_id, path, original_name, size_bytes, mime_type, created_at
		FROM media WHERE restaurant_id = ? AND path = ?`, restaurantID, path).Scan(
		&m.ID, &m.RestaurantID, &m.Path, &m.OriginalName, &m.SizeBytes, &m.MimeType, &created)
	if err != nil {
		return nil, translateErr(err)
	}
	m.CreatedAt = parseTime(created)
	return &m, nil
}

func (r *MediaRepo) Delete(id, restaurantID int64) error {
	_, err := r.db.Exec(`DELETE FROM media WHERE id = ? AND restaurant_id = ?`, id, restaurantID)
	return translateErr(err)
}
