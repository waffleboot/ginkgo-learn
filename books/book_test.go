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

	DescribeTable("Extracting the author's first and last name",
		func(author string, isValid bool, firstName string, lastName string) {
			book.Author = author
			Expect(book.IsValid()).To(Equal(isValid))
			Expect(book.AuthorFirstName()).To(Equal(firstName))
			Expect(book.AuthorLastName()).To(Equal(lastName))
		},
		Entry("When the author has both names", "Victor Hugo", true, "Victor", "Hugo"),
		Entry("When the author only has one name", "Hugo", true, "", "Hugo"),
		Entry("When the author has a middle name", "Victor Marie Hugo", true, "Victor", "Hugo"),
		Entry("When the author has no name", "", false, "", ""))

	Describe("JSON encoding and decoding", func() {
		It("survives the round trip", func() {
			encoded, err := book.AsJSON()
			Expect(err).NotTo(HaveOccurred())

			decoded, err := books.NewBookFromJSON(encoded)
			Expect(err).NotTo(HaveOccurred())

			Expect(decoded).To(Equal(book))
		})

		Describe("some JSON decoding edge cases", func() {
			var err error

			When("the JSON fails to parse", func() {
				BeforeEach(func() {
					book, err = books.NewBookFromJSON(`{
				"title":"Les Miserables",
				"author":"Victor Hugo",
				"pages":2783oops
			  }`)
				})

				It("returns a nil book", func() {
					Expect(book).To(BeNil())
				})

				It("errors", func() {
					Expect(err).To(MatchError(books.ErrInvalidJSON))
				})
			})

			When("the JSON is incomplete", func() {
				BeforeEach(func() {
					book, err = books.NewBookFromJSON(`{
				"title":"Les Miserables",
				"author":"Victor Hugo"
			  }`)
				})

				It("returns a nil book", func() {
					Expect(book).To(BeNil())
				})

				It("errors", func() {
					Expect(err).To(MatchError(books.ErrIncompleteJSON))
				})
			})
		})
	})
})
