package snippets

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yousifsabah0/snippets/internal/models"
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

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, expires, created)
						 VALUES
						 (?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), UTC_TIMESTAMP())
			`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	var snippet Snippet
	stmt := `SELECT id, title, content, expires, created FROM snippets WHERE id = ? AND expires > UTC_TIMESTAMP()`

	row := m.DB.QueryRow(stmt, id)
	if err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Expires, &snippet.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, models.ErrNoRecord
		}

		return Snippet{}, err
	}

	return snippet, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	var snippets []Snippet

	stmt := `SELECT id, title, content, expires, created FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snippet Snippet
		if err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Expires, &snippet.Created); err != nil {
			return nil, nil
		}

		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
