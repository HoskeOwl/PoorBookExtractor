package app

import (
	"context"
	"iter"
	"path/filepath"
	"strings"

	"github.com/HoskeOwl/PoorBookExtractor/internal/entities"
	"github.com/HoskeOwl/PoorBookExtractor/internal/logs"
	"github.com/HoskeOwl/PoorBookExtractor/internal/sources/inpx"
	"github.com/HoskeOwl/PoorBookExtractor/internal/storage/memory"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type App struct {
	storage memory.MemoryStorage
	inpx    inpx.InpxParser

	log *zap.Logger
}

func (a *App) ExportBooks(bookIDsByAuthor map[string][]string, directory string) error {
	for author, book_ids := range bookIDsByAuthor {
		books := a.storage.GetBooks(book_ids)
		if len(books) == 0 {
			a.log.Info("no books to export")
			return nil
		}
		err := inpx.ExportBooks(filepath.Join(directory, author), books)
		if err != nil {
			a.log.Error("error exporting books", zap.Error(err))
			return err
		}
		a.log.Info("exported books", zap.String("author", author), zap.Int("count", len(books)))
	}
	return nil
}

func NewApp(log *zap.Logger) *App {
	storage := memory.NewMemoryStorage(nil)
	return &App{
		log:     log,
		storage: *storage,
		inpx:    *inpx.NewInpxParser(""),
	}
}

func (a *App) ParseInpx(path string) error {
	ctx := logs.WithLog(context.Background(), a.log, zap.String("action", "parse_inpx"))
	books, err := a.inpx.ParseBooks(ctx, path)
	if err != nil {
		return err
	}
	a.log.Debug("parsed books", zap.Int("count", len(books)))
	a.storage.AddBooks(books)
	return nil
}

func (a *App) GetAuthors(value string) []string {
	return a.storage.GetAuthors(value)
}

func (a *App) GetAuthorBooks(author string) []entities.Book {
	return a.storage.GetAuthorBooks(author)
}

func (a *App) GetProgress() int {
	return a.inpx.GetProgress()
}

func (a *App) ClearStorage() {
	a.storage.Clear()
}

func (a *App) IterBooksByAuthor() iter.Seq2[string, []entities.Book] {
	return a.storage.IterBooksByAuthor()
}

func (a *App) SortBooks(books []entities.Book) {
	slices.SortFunc(books, func(a, b entities.Book) int {
		return strings.Compare(a.FullName(), b.FullName())
	})
}

func (a *App) Export(path string, books []entities.Book) error {
	if len(books) == 0 {
		a.log.Debug("no books to export")
		return nil
	}

	return inpx.ExportBooks(path, books)
}

func (a *App) AuthorsLen() int {
	return a.storage.AuthorsLen()
}

func (a *App) BooksLen() int {
	return a.storage.BooksLen()
}

func (a *App) GetBook(libID string) (entities.Book, bool) {
	return a.storage.GetBook(libID)
}
