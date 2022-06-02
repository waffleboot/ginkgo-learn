package books_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/waffleboot/ginkgo-learn/books"
)

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

	Describe("Extracting the author's first and last name", func() {
		Context("When the author has both names", func() {
			It("can extract the author's last name", func() {
				Expect(book.AuthorLastName()).To(Equal("Hugo"))
			})

			It("can extract the author's first name", func() {
				Expect(book.AuthorFirstName()).To(Equal("Victor"))
			})
		})

		Context("When the author only has one name", func() {
			BeforeEach(func() {
				book.Author = "Hugo"
			})

			It("interprets the single author name as a last name", func() {
				Expect(book.AuthorLastName()).To(Equal("Hugo"))
			})

			It("returns empty for the first name", func() {
				Expect(book.AuthorFirstName()).To(BeZero())
			})
		})

		Context("When the author has a middle name", func() {
			BeforeEach(func() {
				book.Author = "Victor Marie Hugo"
			})

			It("can extract the author's last name", func() {
				Expect(book.AuthorLastName()).To(Equal("Hugo"))
			})

			It("can extract the author's first name", func() {
				Expect(book.AuthorFirstName()).To(Equal("Victor"))
			})
		})

		Context("When the author has no name", func() {
			It("should not be a valid book and returns empty for first and last name", func() {
				book.Author = ""
				Expect(book.IsValid()).To(BeFalse())
				Expect(book.AuthorLastName()).To(BeZero())
				Expect(book.AuthorFirstName()).To(BeZero())
			})
		})
	})
})
