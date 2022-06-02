package books

type Category int

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
	return b.Author
}

func (b *Book) AuthorFirstName() string {
	return b.Author
}

func (b *Book) IsValid() bool {
	return b.Title != "" && b.Author != "" && b.Pages > 0
}
