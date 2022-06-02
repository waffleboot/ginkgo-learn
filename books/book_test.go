package books_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/waffleboot/ginkgo-learn/books"
)

var _ = Describe("Books", func() {
	var foxInSocks, lesMis *books.Book

	BeforeEach(func() {
		lesMis = &books.Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
			Pages:  2783,
		}

		foxInSocks = &books.Book{
			Title:  "Fox In Socks",
			Author: "Dr. Seuss",
			Pages:  24,
		}
	})

	Describe("Categorizing books", func() {
		Context("with more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(lesMis.Category()).To(Equal(books.CategoryNovel))
			})
		})
		Context("with fewer than 300 pages", func() {
			It("should be a short story", func() {
				Expect(foxInSocks.Category()).To(Equal(books.CategoryShortStory))
			})
		})
	})
})

var _ = Describe("Books", func() {
	var book *books.Book
	It("can extract the author's last name", func() {
		book = &books.Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
			Pages:  2783,
		}
		Expect(book.AuthorLastName()).To(Equal("Hugo"))
	})
})

var _ = Describe("Books", func() {
	var book *books.Book

	BeforeEach(func() {
		book = &books.Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
			Pages:  2783,
		}
		Expect(book.IsValid()).To(BeTrue())
	})

	It("can extract the author's last name", func() {
		Expect(book.AuthorLastName()).To(Equal("Hugo"))
	})

	It("interprets a single author name as a last name", func() {
		book.Author = "Hugo"
		Expect(book.AuthorLastName()).To(Equal("Hugo"))
	})

	It("can extract the author's first name", func() {
		Expect(book.AuthorFirstName()).To(Equal("Victor"))
	})

	It("returns no first name when there is a single author name", func() {
		book.Author = "Hugo"
		Expect(book.AuthorFirstName()).To(BeZero()) // BeZero asserts the value is the zero-value for its type.  In this case: ""
	})
})
