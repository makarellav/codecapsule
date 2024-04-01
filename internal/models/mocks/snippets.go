package mocks

import (
	"github.com/makarellav/codecapsule/internal/models"
	"time"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return &mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{&mockSnippet}, nil
}
