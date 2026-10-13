package repository

import (
	"strings"

	"ResturantMenue/internal/models"
)

type UserRepo struct{ db dbQuery }

const userColumns = `id, username, email, password_hash, display_name, locale, session_version, created_at, updated_at`

func (r *UserRepo) Create(u *models.User) error {
	// session_version is written explicitly rather than left to the column
	// default: the caller issues a session cookie straight from this struct, so
	// the in-memory value has to match what the row actually holds.
	if u.SessionVersion < 1 {
		u.SessionVersion = 1
	}

	res, err := r.db.Exec(`
		INSERT INTO users (username, email, password_hash, display_name, locale, session_version)
		VALUES (?, ?, ?, ?, ?, ?)`,
		strings.ToLower(u.Username), strings.ToLower(u.Email), u.PasswordHash,
		u.DisplayName, u.Locale, u.SessionVersion)
	if err != nil {
		return translateErr(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *UserRepo) ByID(id int64) (*models.User, error) {
	return r.scanOne(`SELECT `+userColumns+` FROM users WHERE id = ?`, id)
}

// ByLogin accepts either a username or an email address, so the login form has
// a single field.
func (r *UserRepo) ByLogin(login string) (*models.User, error) {
	login = strings.ToLower(strings.TrimSpace(login))
	return r.scanOne(`SELECT `+userColumns+` FROM users WHERE username = ? OR email = ?`, login, login)
}

func (r *UserRepo) UpdateProfile(u *models.User) error {
	_, err := r.db.Exec(`
		UPDATE users SET display_name = ?, locale = ?, updated_at = datetime('now')
		WHERE id = ?`, u.DisplayName, u.Locale, u.ID)
	return translateErr(err)
}

// UpdatePassword also bumps session_version, which signs every existing
// session cookie out — the expected behaviour after a password change.
func (r *UserRepo) UpdatePassword(userID int64, hash string) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET password_hash = ?, session_version = session_version + 1, updated_at = datetime('now')
		WHERE id = ?`, hash, userID)
	return translateErr(err)
}

// Exists reports whether a username or email is already taken.
func (r *UserRepo) Exists(username, email string) (bool, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`,
		strings.ToLower(username), strings.ToLower(email)).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *UserRepo) scanOne(query string, args ...any) (*models.User, error) {
	var u models.User
	var created, updated string
	err := r.db.QueryRow(query, args...).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName,
		&u.Locale, &u.SessionVersion, &created, &updated)
	if err != nil {
		return nil, translateErr(err)
	}
	u.CreatedAt = parseTime(created)
	u.UpdatedAt = parseTime(updated)
	return &u, nil
}
