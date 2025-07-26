package app

import (
	"context"

	"github.com/HoskeOwl/PoorBoockExtractor/internal/entities"
	"github.com/HoskeOwl/PoorBoockExtractor/internal/logs"
	"github.com/HoskeOwl/PoorBoockExtractor/internal/sources/inpx"
	"github.com/HoskeOwl/PoorBoockExtractor/internal/storage/memory"
	"go.uber.org/zap"
)

type App struct {
	storage memory.MemoryStorage
	inpx    inpx.InpxParser

	log *zap.Logger
}

func (a *App) ExportBooks(book_ids []string, directory string) error {
	books := a.storage.GetBooks(book_ids)
	if len(books) == 0 {
		a.log.Info("no books to export")
		return nil
	}
	err := inpx.ExportBooks(directory, books)
	if err != nil {
		a.log.Error("error exporting books", zap.Error(err))
		return err
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

func (a *App) Export(path string, books []entities.Book) error {
	if len(books) == 0 {
		a.log.Debug("no books to export")
		return nil
	}

	return inpx.ExportBooks(path, books)
}
