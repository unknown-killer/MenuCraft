// Package repository is the only place that speaks SQL. Handlers depend on
// these types, which keeps the storage engine swappable.
package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"ResturantMenue/internal/database"
)

var (
	// ErrNotFound is returned instead of sql.ErrNoRows so callers do not need
	// to import database/sql.
	ErrNotFound = errors.New("not found")
	// ErrDuplicate signals a uniqueness conflict (username, email or slug).
	ErrDuplicate = errors.New("already exists")
)

// Store bundles every repository behind one value passed to the handlers.
type Store struct {
	db          *database.DB
	Users       *UserRepo
	Restaurants *RestaurantRepo
	Categories  *CategoryRepo
	Items       *ItemRepo
	Themes      *ThemeRepo
	Media       *MediaRepo
	Blocks      *BlockRepo
}

func New(db *database.DB) *Store {
	return &Store{
		db:          db,
		Users:       &UserRepo{db},
		Restaurants: &RestaurantRepo{db},
		Categories:  &CategoryRepo{db},
		Items:       &ItemRepo{db},
		Themes:      &ThemeRepo{db},
		Media:       &MediaRepo{db},
		Blocks:      &BlockRepo{db},
	}
}

// DB exposes the connection for the rare cross-repository transaction.
func (s *Store) DB() *database.DB { return s.db }

// parseTime tolerates the couple of layouts SQLite's datetime() and our own
// inserts can produce.
func parseTime(s string) time.Time {
	for _, layout := range []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// translateErr maps driver errors onto the package's sentinel errors.
func translateErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if isUniqueViolation(err) {
		return ErrDuplicate
	}
	return err
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
