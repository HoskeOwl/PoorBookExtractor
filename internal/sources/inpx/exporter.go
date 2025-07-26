package inpx

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/HoskeOwl/PoorBookExtractor/internal/entities"
)

var re = regexp.MustCompile(`[^\w\s\-а-яА-Я]+`)

func createName(book entities.Book) string {
	name := book.Title
	if len(name) > 100 {
		name = name[:100]
	}
	name = re.ReplaceAllLiteralString(name, "")
	if len(name) == 0 {
		name = book.Filename
	}
	return name + "." + book.Ext
}

func exportBook(zipPath, outFileName string, books []entities.Book) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	for _, bookData := range books {
		book, err := zipReader.Open(bookData.Filename + "." + bookData.Ext)
		if err != nil {
			return err
		}
		outFile, err := os.Create(filepath.Join(outFileName, createName(bookData)))
		if err != nil {
			book.Close()
			return err
		}
		_, err = io.Copy(outFile, book)
		if err != nil {
			book.Close()
			outFile.Close()
			return err
		}
		book.Close()
		outFile.Close()
	}
	return nil
}

func ExportBooks(path string, books []entities.Book) error {
	if len(books) == 0 {
		return nil
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	bookByArchive := make(map[string][]entities.Book)
	for _, book := range books {
		bookByArchive[book.Metadata.ArchiveName] = append(bookByArchive[book.Metadata.ArchiveName], book)
	}
	archivePath := books[0].Metadata.Filepath
	for archive, books := range bookByArchive {
		p := filepath.Join(archivePath, archive)
		err = exportBook(p, path, books)
		if err != nil {
			return err
		}
	}
	return nil
}
