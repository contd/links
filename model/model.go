package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Link is the struct for link database objects
type Link struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Category string    `json:"category"`
	Created  time.Time `json:"created_on"`
	Done     bool      `json:"done"`
}

// GetLink gets a single link given its id
func (l *Link) GetLink(db *sqlx.DB) error {
	return db.QueryRow(
		"SELECT id, url, category, created_on, done FROM links WHERE id=?",
		l.ID).Scan(&l.ID, &l.URL, &l.Category, &l.Created, &l.Done)
}

// UpdateLink updates a link given its id
func (l *Link) UpdateLink(db *sqlx.DB) error {
	_, err := db.Exec(
		"UPDATE links SET url=?, category=?, done=? WHERE id=?",
		l.URL, l.Category, l.Done, l.ID)
	return err
}

// DeleteLink deletes a link given its id
func (l *Link) DeleteLink(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM links WHERE id=?", l.ID)
	return err
}

// CreateLink makes a new link
func (l *Link) CreateLink(db *sqlx.DB) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO links(url, category, created_on, done) VALUES(?, ?, ?, ?)",
		l.URL, l.Category, l.Created, l.Done)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetLinks gets all Links
func GetLinks(db *sqlx.DB) ([]Link, error) {
	rows, err := db.Query("SELECT id, url, category, created_on, done FROM links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	links := []Link{}

	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.URL, &l.Category, &l.Created, &l.Done); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}
