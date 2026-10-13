package repository

import (
	"ResturantMenue/internal/models"
)

type RestaurantRepo struct{ db dbQuery }

const restaurantColumns = `
	id, user_id, slug, business_type,
	name_fa, name_en, tagline_fa, tagline_en, description_fa, description_en,
	logo_path, cover_path,
	phone, address_fa, address_en, instagram, telegram, whatsapp, website,
	currency, is_published, created_at, updated_at`

func (r *RestaurantRepo) Create(v *models.Restaurant) error {
	res, err := r.db.Exec(`
		INSERT INTO restaurants (user_id, slug, business_type, name_fa, name_en, currency, is_published)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		v.UserID, v.Slug, v.BusinessType, v.NameFa, v.NameEn, v.Currency, boolInt(v.IsPublished))
	if err != nil {
		return translateErr(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	v.ID = id
	return nil
}

func (r *RestaurantRepo) ByID(id int64) (*models.Restaurant, error) {
	return r.scanOne(`SELECT `+restaurantColumns+` FROM restaurants WHERE id = ?`, id)
}

func (r *RestaurantRepo) BySlug(slug string) (*models.Restaurant, error) {
	return r.scanOne(`SELECT `+restaurantColumns+` FROM restaurants WHERE slug = ?`, slug)
}

// ByUser returns the venue owned by a user. The current product is one venue
// per account; this is the seam to relax that later.
func (r *RestaurantRepo) ByUser(userID int64) (*models.Restaurant, error) {
	return r.scanOne(`SELECT `+restaurantColumns+` FROM restaurants WHERE user_id = ? ORDER BY id LIMIT 1`, userID)
}

func (r *RestaurantRepo) Update(v *models.Restaurant) error {
	_, err := r.db.Exec(`
		UPDATE restaurants SET
			slug = ?, business_type = ?,
			name_fa = ?, name_en = ?, tagline_fa = ?, tagline_en = ?,
			description_fa = ?, description_en = ?,
			logo_path = ?, cover_path = ?,
			phone = ?, address_fa = ?, address_en = ?,
			instagram = ?, telegram = ?, whatsapp = ?, website = ?,
			currency = ?, is_published = ?, updated_at = datetime('now')
		WHERE id = ?`,
		v.Slug, v.BusinessType,
		v.NameFa, v.NameEn, v.TaglineFa, v.TaglineEn,
		v.DescriptionFa, v.DescriptionEn,
		v.LogoPath, v.CoverPath,
		v.Phone, v.AddressFa, v.AddressEn,
		v.Instagram, v.Telegram, v.WhatsApp, v.Website,
		v.Currency, boolInt(v.IsPublished), v.ID)
	return translateErr(err)
}

// SlugTaken reports whether a slug is used by a venue other than excludeID.
func (r *RestaurantRepo) SlugTaken(slug string, excludeID int64) (bool, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM restaurants WHERE slug = ? AND id <> ?`, slug, excludeID).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *RestaurantRepo) scanOne(query string, args ...any) (*models.Restaurant, error) {
	var v models.Restaurant
	var published int
	var created, updated string
	err := r.db.QueryRow(query, args...).Scan(
		&v.ID, &v.UserID, &v.Slug, &v.BusinessType,
		&v.NameFa, &v.NameEn, &v.TaglineFa, &v.TaglineEn, &v.DescriptionFa, &v.DescriptionEn,
		&v.LogoPath, &v.CoverPath,
		&v.Phone, &v.AddressFa, &v.AddressEn, &v.Instagram, &v.Telegram, &v.WhatsApp, &v.Website,
		&v.Currency, &published, &created, &updated)
	if err != nil {
		return nil, translateErr(err)
	}
	v.IsPublished = published == 1
	v.CreatedAt = parseTime(created)
	v.UpdatedAt = parseTime(updated)
	return &v, nil
}

// --- working hours ---------------------------------------------------------

func (r *RestaurantRepo) WorkingHours(restaurantID int64) ([]*models.WorkingHour, error) {
	rows, err := r.db.Query(`
		SELECT id, restaurant_id, day_of_week, opens_at, closes_at, is_closed
		FROM working_hours WHERE restaurant_id = ? ORDER BY day_of_week`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*models.WorkingHour, 0, 7)
	for rows.Next() {
		var h models.WorkingHour
		var closed int
		if err := rows.Scan(&h.ID, &h.RestaurantID, &h.DayOfWeek, &h.OpensAt, &h.ClosesAt, &closed); err != nil {
			return nil, err
		}
		h.IsClosed = closed == 1
		out = append(out, &h)
	}
	return out, rows.Err()
}

// SaveWorkingHours upserts the seven-day schedule in one transaction.
func (r *RestaurantRepo) SaveWorkingHours(restaurantID int64, hours []*models.WorkingHour) error {
	return inTransaction(r.db, func(tx dbQuery) error {
		for _, h := range hours {
			if _, err := tx.Exec(`
			INSERT INTO working_hours (restaurant_id, day_of_week, opens_at, closes_at, is_closed)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(restaurant_id, day_of_week) DO UPDATE SET
				opens_at = excluded.opens_at,
				closes_at = excluded.closes_at,
				is_closed = excluded.is_closed`,
				restaurantID, h.DayOfWeek, h.OpensAt, h.ClosesAt, boolInt(h.IsClosed)); err != nil {
				return err
			}
		}
		return nil
	})
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
