package memory

import (
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/HoskeOwl/PoorBookExtractor/internal/entities"
)

type MemoryStorage struct {
	books    map[string]entities.Book
	byAuthor map[string][]entities.Book
}

func NewMemoryStorage(books []entities.Book) *MemoryStorage {
	booksMap := make(map[string]entities.Book)
	byAuthor := make(map[string][]entities.Book)
	for _, book := range books {
		for _, author := range book.Authors {
			byAuthor[author] = append(byAuthor[author], book)
		}
		booksMap[book.LibID] = book
	}
	return &MemoryStorage{books: booksMap, byAuthor: byAuthor}
}

func (ms *MemoryStorage) AddBooks(books []entities.Book) {
	for _, book := range books {
		ms.books[book.LibID] = book
		for _, author := range book.Authors {
			ms.byAuthor[author] = append(ms.byAuthor[author], book)
		}
	}
}

func (ms *MemoryStorage) GetAuthors(value string) []string {
	if value == "" {
		keys := slices.Collect(maps.Keys(ms.byAuthor))
		slices.Sort(keys)
		return keys
	}
	value = strings.ToLower(value)
	tokens := strings.Split(value, " ")
	authors := make([]string, 0)
	for author := range ms.byAuthor {
		lowerAuthor := strings.ToLower(author)
		found := true
		for _, token := range tokens {
			found = found && strings.Contains(lowerAuthor, token)
		}
		if found {
			authors = append(authors, author)
		}
	}
	slices.Sort(authors)
	slices.Reverse(authors)
	return authors
}

func (ms *MemoryStorage) GetAuthorBooks(author string) []entities.Book {
	return ms.byAuthor[author]
}

func (ms *MemoryStorage) GetBooks(book_ids []string) []entities.Book {
	books := make([]entities.Book, 0, len(book_ids))
	for _, bid := range book_ids {
		book, ok := ms.books[bid]
		if !ok {
			continue
		}
		books = append(books, book)
	}
	return books
}

func (ms *MemoryStorage) Clear() {
	ms.books = make(map[string]entities.Book)
	ms.byAuthor = make(map[string][]entities.Book)
}

func (ms *MemoryStorage) IterBooksByAuthor() iter.Seq2[string, []entities.Book] {
	return func(yield func(author string, book []entities.Book) bool) {
		authors := slices.Collect(maps.Keys(ms.byAuthor))
		slices.Sort(authors)
		for _, author := range authors {
			books := ms.byAuthor[author]
			if !yield(author, books) {
				return
			}
		}
	}
}

func (ms *MemoryStorage) AuthorsLen() int {
	return len(ms.byAuthor)
}

func (ms *MemoryStorage) BooksLen() int {
	return len(ms.books)
}

func (ms *MemoryStorage) GetBook(libID string) (entities.Book, bool) {
	book, ok := ms.books[libID]
	return book, ok
}
