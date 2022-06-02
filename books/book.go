package books

import (
	"encoding/json"
	"errors"
	"strings"
)

type Category int

var (
	ErrInvalidJSON    = errors.New("InvalidJSON")
	ErrIncompleteJSON = errors.New("IncompleteJSON")
)

const (
	CategoryNovel Category = iota + 1
	CategoryShortStory
)

type Book struct {
	Title  string
	Author string
	Pages  uint
}

func (b *Book) Category() Category {
	if b.Pages > 300 {
		return CategoryNovel
	}
	return CategoryShortStory
}

func (b *Book) AuthorLastName() string {
	fields := strings.Fields(b.Author)
	switch {
	case len(fields) == 0:
		return ""
	default:
		return fields[len(fields)-1]
	}
}

func (b *Book) AuthorFirstName() string {
	fields := strings.Fields(b.Author)
	switch {
	case len(fields) < 2:
		return ""
	default:
		return fields[0]
	}
}

func (b *Book) IsValid() bool {
	return b.Title != "" && b.Author != "" && b.Pages > 0
}

func NewBookFromJSON(encoded string) (*Book, error) {
	book := new(Book)
	if err := json.Unmarshal([]byte(encoded), book); err != nil {
		return nil, ErrInvalidJSON
	}
	if !book.IsValid() {
		return nil, ErrIncompleteJSON
	}
	return book, nil
}

func (b *Book) AsJSON() (string, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
