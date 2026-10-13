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
	query       dbQuery
	Users       *UserRepo
	Restaurants *RestaurantRepo
	Categories  *CategoryRepo
	Items       *ItemRepo
	Themes      *ThemeRepo
	Media       *MediaRepo
	Blocks      *BlockRepo
}

func New(db *database.DB) *Store {
	return newStore(db, db)
}

// dbQuery is implemented by both the connection pool and a transaction.
type dbQuery interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

func newStore(db *database.DB, query dbQuery) *Store {
	return &Store{
		db:          db,
		query:       query,
		Users:       &UserRepo{query},
		Restaurants: &RestaurantRepo{query},
		Categories:  &CategoryRepo{query},
		Items:       &ItemRepo{query},
		Themes:      &ThemeRepo{query},
		Media:       &MediaRepo{query},
		Blocks:      &BlockRepo{query},
	}
}

func inTransaction(db dbQuery, fn func(dbQuery) error) error {
	if _, ok := db.(*sql.Tx); ok {
		return fn(db)
	}
	pool, ok := db.(interface{ Begin() (*sql.Tx, error) })
	if !ok {
		return errors.New("repository: connection cannot start transaction")
	}
	tx, err := pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// Transaction commits all repository changes together, or rolls them all back.
// A separate Store is used so concurrent requests never share the transaction.
func (s *Store) Transaction(fn func(*Store) error) error {
	return inTransaction(s.query, func(query dbQuery) error { return fn(newStore(s.db, query)) })
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
