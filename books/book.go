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
