package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	args := []any{title, content, expires}

	result, err := sm.DB.Exec(stmt, args...)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
	FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := sm.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &s, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
