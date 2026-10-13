package repository

import (
	"ResturantMenue/internal/database"
	"ResturantMenue/internal/models"
)

type BlockRepo struct{ db *database.DB }

const blockColumns = `
	id, restaurant_id, kind, title_fa, title_en, body_fa, body_en,
	image_path, images, link_url, link_label_fa, link_label_en,
	style, canvas, sort_order, is_visible, created_at, updated_at`

func (r *BlockRepo) Create(b *models.Block) error {
	if b.SortOrder == 0 {
		_ = r.db.QueryRow(`SELECT COALESCE(MAX(sort_order), 0) + 1 FROM blocks WHERE restaurant_id = ?`,
			b.RestaurantID).Scan(&b.SortOrder)
	}

	res, err := r.db.Exec(`
		INSERT INTO blocks (restaurant_id, kind, title_fa, title_en, body_fa, body_en,
			image_path, images, link_url, link_label_fa, link_label_en,
			style, canvas, sort_order, is_visible)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		b.RestaurantID, string(b.Kind), b.TitleFa, b.TitleEn, b.BodyFa, b.BodyEn,
		b.ImagePath, b.ImagesJSON(), b.LinkURL, b.LinkLabelFa, b.LinkLabelEn,
		b.StyleJSON(), b.Canvas.CanvasJSON(), b.SortOrder, boolInt(b.IsVisible))
	if err != nil {
		return translateErr(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

func (r *BlockRepo) Update(b *models.Block) error {
	_, err := r.db.Exec(`
		UPDATE blocks SET
			kind = ?, title_fa = ?, title_en = ?, body_fa = ?, body_en = ?,
			image_path = ?, images = ?, link_url = ?,
			link_label_fa = ?, link_label_en = ?,
			style = ?, canvas = ?, is_visible = ?, updated_at = datetime('now')
		WHERE id = ? AND restaurant_id = ?`,
		string(b.Kind), b.TitleFa, b.TitleEn, b.BodyFa, b.BodyEn,
		b.ImagePath, b.ImagesJSON(), b.LinkURL,
		b.LinkLabelFa, b.LinkLabelEn,
		b.StyleJSON(), b.Canvas.CanvasJSON(), boolInt(b.IsVisible), b.ID, b.RestaurantID)
	return translateErr(err)
}

func (r *BlockRepo) Delete(id, restaurantID int64) error {
	_, err := r.db.Exec(`DELETE FROM blocks WHERE id = ? AND restaurant_id = ?`, id, restaurantID)
	return translateErr(err)
}

func (r *BlockRepo) ByID(id, restaurantID int64) (*models.Block, error) {
	row := r.db.QueryRow(`SELECT `+blockColumns+` FROM blocks WHERE id = ? AND restaurant_id = ?`,
		id, restaurantID)
	return scanBlock(row.Scan)
}

// List returns every block for a venue. onlyVisible filters to what the public
// page should show.
func (r *BlockRepo) List(restaurantID int64, onlyVisible bool) ([]*models.Block, error) {
	query := `SELECT ` + blockColumns + ` FROM blocks WHERE restaurant_id = ?`
	if onlyVisible {
		query += ` AND is_visible = 1`
	}
	query += ` ORDER BY sort_order, id`

	rows, err := r.db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Block
	for rows.Next() {
		b, err := scanBlock(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func scanBlock(scan func(...any) error) (*models.Block, error) {
	var b models.Block
	var kind, images, style, canvas string
	var visible int
	var created, updated string

	err := scan(&b.ID, &b.RestaurantID, &kind, &b.TitleFa, &b.TitleEn, &b.BodyFa, &b.BodyEn,
		&b.ImagePath, &images, &b.LinkURL, &b.LinkLabelFa, &b.LinkLabelEn,
		&style, &canvas, &b.SortOrder, &visible, &created, &updated)
	if err != nil {
		return nil, translateErr(err)
	}

	b.Kind = models.BlockKind(kind)
	b.SetImagesJSON(images)
	// SetStyleJSON depends on Kind for its defaults, so it runs after.
	b.SetStyleJSON(style)
	b.Canvas = models.ParseCanvas(canvas)
	b.IsVisible = visible == 1
	b.CreatedAt = parseTime(created)
	b.UpdatedAt = parseTime(updated)
	return &b, nil
}
