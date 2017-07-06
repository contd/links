package main

import (
	"database/sql"
)

type link struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	Category string `json:"category"`
	Created  string `json:"created_on"`
	Done     int    `json:"done"`
}

func (l *link) getLink(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, url, category, created_on, done FROM links WHERE id=?",
		l.ID).Scan(&l.ID, &l.Url, &l.Category, &l.Created, &l.Done)
}

func (l *link) updateLink(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE links SET url=?, category=?, created_on=?, done=? WHERE id=?",
		l.Url, l.Category, l.Created, l.Done, l.ID)
	return err
}

func (l *link) deleteLink(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM links WHERE id=?", l.ID)
	return err
}

func (l *link) createLink(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO links(url, category, created_on, done) VALUES(?, ?, ?, ?)",
		l.Url, l.Category, l.Created, l.Done)
	if err != nil {
		return err
	}
	return nil
}

func getLinks(db *sql.DB) ([]link, error) {
	rows, err := db.Query("SELECT id, url, category, created_on, done FROM links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	links := []link{}

	for rows.Next() {
		var l link
		if err := rows.Scan(&l.ID, &l.Url, &l.Category, &l.Created, &l.Done); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}
