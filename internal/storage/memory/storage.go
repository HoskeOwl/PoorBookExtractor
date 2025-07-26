package memory

import (
	"maps"
	"slices"
	"strings"

	"github.com/HoskeOwl/PoorBoockExtractor/internal/entities"
)

type MemoryStorage struct {
	books    map[string]entities.Book
	ByAuthor map[string][]entities.Book
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
	return &MemoryStorage{books: booksMap, ByAuthor: byAuthor}
}

func (ms *MemoryStorage) AddBooks(books []entities.Book) {
	for _, book := range books {
		ms.books[book.LibID] = book
		for _, author := range book.Authors {
			ms.ByAuthor[author] = append(ms.ByAuthor[author], book)
		}
	}
}

func (ms *MemoryStorage) GetAuthors(value string) []string {
	if value == "" {
		keys := slices.Collect(maps.Keys(ms.ByAuthor))
		slices.Sort(keys)
		return keys
	}
	value = strings.ToLower(value)
	tokens := strings.Split(value, " ")
	authors := make([]string, 0)
	for author := range ms.ByAuthor {
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
	return ms.ByAuthor[author]
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
	ms.ByAuthor = make(map[string][]entities.Book)
}
