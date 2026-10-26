package repository

import (
	"ResturantMenue/internal/models"
)

// --- categories ------------------------------------------------------------

type CategoryRepo struct{ db dbQuery }

const categoryColumns = `
	id, restaurant_id, name_fa, name_en, description_fa, description_en,
	icon, image_path, sort_order, is_visible, created_at, updated_at`

func (r *CategoryRepo) Create(c *models.Category) error {
	// New categories land at the end of the list.
	if c.SortOrder == 0 {
		_ = r.db.QueryRow(`SELECT COALESCE(MAX(sort_order), 0) + 1 FROM categories WHERE restaurant_id = ?`,
			c.RestaurantID).Scan(&c.SortOrder)
	}
	res, err := r.db.Exec(`
		INSERT INTO categories (restaurant_id, name_fa, name_en, description_fa, description_en,
			icon, image_path, sort_order, is_visible)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.RestaurantID, c.NameFa, c.NameEn, c.DescriptionFa, c.DescriptionEn,
		c.Icon, c.ImagePath, c.SortOrder, boolInt(c.IsVisible))
	if err != nil {
		return translateErr(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *CategoryRepo) Update(c *models.Category) error {
	_, err := r.db.Exec(`
		UPDATE categories SET
			name_fa = ?, name_en = ?, description_fa = ?, description_en = ?,
			icon = ?, image_path = ?, is_visible = ?, updated_at = datetime('now')
		WHERE id = ? AND restaurant_id = ?`,
		c.NameFa, c.NameEn, c.DescriptionFa, c.DescriptionEn,
		c.Icon, c.ImagePath, boolInt(c.IsVisible), c.ID, c.RestaurantID)
	return translateErr(err)
}

// Delete scopes by restaurant so one owner can never remove another's data.
func (r *CategoryRepo) Delete(id, restaurantID int64) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = ? AND restaurant_id = ?`, id, restaurantID)
	return translateErr(err)
}

func (r *CategoryRepo) ByID(id, restaurantID int64) (*models.Category, error) {
	var c models.Category
	var visible int
	var created, updated string
	err := r.db.QueryRow(`SELECT `+categoryColumns+` FROM categories WHERE id = ? AND restaurant_id = ?`,
		id, restaurantID).Scan(
		&c.ID, &c.RestaurantID, &c.NameFa, &c.NameEn, &c.DescriptionFa, &c.DescriptionEn,
		&c.Icon, &c.ImagePath, &c.SortOrder, &visible, &created, &updated)
	if err != nil {
		return nil, translateErr(err)
	}
	c.IsVisible = visible == 1
	c.CreatedAt = parseTime(created)
	c.UpdatedAt = parseTime(updated)
	return &c, nil
}

func (r *CategoryRepo) List(restaurantID int64) ([]*models.Category, error) {
	rows, err := r.db.Query(`SELECT `+categoryColumns+`
		FROM categories WHERE restaurant_id = ? ORDER BY sort_order, id`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Category
	for rows.Next() {
		var c models.Category
		var visible int
		var created, updated string
		if err := rows.Scan(&c.ID, &c.RestaurantID, &c.NameFa, &c.NameEn,
			&c.DescriptionFa, &c.DescriptionEn, &c.Icon, &c.ImagePath,
			&c.SortOrder, &visible, &created, &updated); err != nil {
			return nil, err
		}
		c.IsVisible = visible == 1
		c.CreatedAt = parseTime(created)
		c.UpdatedAt = parseTime(updated)
		out = append(out, &c)
	}
	return out, rows.Err()
}

// Reorder persists a new ordering in one transaction.
func (r *CategoryRepo) Reorder(restaurantID int64, orderedIDs []int64) error {
	return inTransaction(r.db, func(tx dbQuery) error {
		for i, id := range orderedIDs {
			if _, err := tx.Exec(`UPDATE categories SET sort_order = ? WHERE id = ? AND restaurant_id = ?`,
				i+1, id, restaurantID); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- items -----------------------------------------------------------------

type ItemRepo struct{ db dbQuery }

const itemColumns = `
	id, restaurant_id, category_id, name_fa, name_en, description_fa, description_en,
	price, discount_price, image_path, tags, is_available, is_featured,
	calories, prep_minutes, sort_order, created_at, updated_at, price_pending`

func (r *ItemRepo) Create(it *models.Item) error {
	if it.SortOrder == 0 {
		_ = r.db.QueryRow(`SELECT COALESCE(MAX(sort_order), 0) + 1 FROM items WHERE category_id = ?`,
			it.CategoryID).Scan(&it.SortOrder)
	}
	res, err := r.db.Exec(`
		INSERT INTO items (restaurant_id, category_id, name_fa, name_en,
			description_fa, description_en, price, discount_price, image_path,
			tags, is_available, is_featured, calories, prep_minutes, sort_order, price_pending)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		it.RestaurantID, it.CategoryID, it.NameFa, it.NameEn,
		it.DescriptionFa, it.DescriptionEn, it.Price, it.DiscountPrice, it.ImagePath,
		it.TagsJSON(), boolInt(it.IsAvailable), boolInt(it.IsFeatured),
		it.Calories, it.PrepMinutes, it.SortOrder, boolInt(it.PricePending))
	if err != nil {
		return translateErr(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	it.ID = id
	return nil
}

func (r *ItemRepo) Update(it *models.Item) error {
	_, err := r.db.Exec(`
		UPDATE items SET
			category_id = ?, name_fa = ?, name_en = ?,
			description_fa = ?, description_en = ?,
			price = ?, discount_price = ?, image_path = ?, tags = ?,
			is_available = ?, is_featured = ?, calories = ?, prep_minutes = ?, price_pending = ?,
			updated_at = datetime('now')
		WHERE id = ? AND restaurant_id = ?`,
		it.CategoryID, it.NameFa, it.NameEn,
		it.DescriptionFa, it.DescriptionEn,
		it.Price, it.DiscountPrice, it.ImagePath, it.TagsJSON(),
		boolInt(it.IsAvailable), boolInt(it.IsFeatured), it.Calories, it.PrepMinutes, boolInt(it.PricePending),
		it.ID, it.RestaurantID)
	return translateErr(err)
}

func (r *ItemRepo) Delete(id, restaurantID int64) error {
	_, err := r.db.Exec(`DELETE FROM items WHERE id = ? AND restaurant_id = ?`, id, restaurantID)
	return translateErr(err)
}

func (r *ItemRepo) ByID(id, restaurantID int64) (*models.Item, error) {
	row := r.db.QueryRow(`SELECT `+itemColumns+` FROM items WHERE id = ? AND restaurant_id = ?`, id, restaurantID)
	return scanItem(row.Scan)
}

func (r *ItemRepo) ListByRestaurant(restaurantID int64) ([]*models.Item, error) {
	rows, err := r.db.Query(`SELECT `+itemColumns+`
		FROM items WHERE restaurant_id = ? ORDER BY category_id, sort_order, id`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Item
	for rows.Next() {
		it, err := scanItem(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// ToggleAvailable flips the availability flag without loading the whole item.
func (r *ItemRepo) ToggleAvailable(id, restaurantID int64) error {
	_, err := r.db.Exec(`
		UPDATE items SET is_available = 1 - is_available, updated_at = datetime('now')
		WHERE id = ? AND restaurant_id = ?`, id, restaurantID)
	return translateErr(err)
}

func (r *ItemRepo) Reorder(restaurantID int64, orderedIDs []int64) error {
	return inTransaction(r.db, func(tx dbQuery) error {
		for i, id := range orderedIDs {
			if _, err := tx.Exec(`UPDATE items SET sort_order = ? WHERE id = ? AND restaurant_id = ?`,
				i+1, id, restaurantID); err != nil {
				return err
			}
		}
		return nil
	})
}

// scanItem works with both *sql.Row.Scan and *sql.Rows.Scan.
func scanItem(scan func(...any) error) (*models.Item, error) {
	var it models.Item
	var tags string
	var available, featured, pending int
	var created, updated string
	err := scan(&it.ID, &it.RestaurantID, &it.CategoryID, &it.NameFa, &it.NameEn,
		&it.DescriptionFa, &it.DescriptionEn, &it.Price, &it.DiscountPrice, &it.ImagePath,
		&tags, &available, &featured, &it.Calories, &it.PrepMinutes, &it.SortOrder,
		&created, &updated, &pending)
	if err != nil {
		return nil, translateErr(err)
	}
	it.SetTagsJSON(tags)
	it.IsAvailable = available == 1
	it.IsFeatured = featured == 1
	it.PricePending = pending == 1
	it.CreatedAt = parseTime(created)
	it.UpdatedAt = parseTime(updated)
	return &it, nil
}

// --- full menu -------------------------------------------------------------

// FullMenu loads categories with their items attached, in two queries rather
// than one per category.
func (s *Store) FullMenu(restaurantID int64, onlyVisible bool) ([]*models.Category, error) {
	cats, err := s.Categories.List(restaurantID)
	if err != nil {
		return nil, err
	}
	items, err := s.Items.ListByRestaurant(restaurantID)
	if err != nil {
		return nil, err
	}

	byCategory := make(map[int64][]*models.Item, len(cats))
	for _, it := range items {
		if onlyVisible && !it.IsAvailable {
			continue
		}
		byCategory[it.CategoryID] = append(byCategory[it.CategoryID], it)
	}

	out := make([]*models.Category, 0, len(cats))
	for _, c := range cats {
		if onlyVisible && !c.IsVisible {
			continue
		}
		c.Items = byCategory[c.ID]
		out = append(out, c)
	}
	return out, nil
}
