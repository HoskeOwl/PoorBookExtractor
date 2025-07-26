package inp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HoskeOwl/PoorBoockExtractor/internal/entities"
	"github.com/HoskeOwl/PoorBoockExtractor/internal/logs"
	"go.uber.org/zap"
)

const (
	separator   = "\x04"
	endOfRecord = "\r\n"
	columnCount = 14

	authorIndex       = 0
	genreIndex        = 1
	titleIndex        = 2
	seriesIndex       = 3
	seriesNumberIndex = 4
	fileIndex         = 5
	sizeIndex         = 6
	libIDIndex        = 7
	delIndex          = 8
	extIndex          = 9
	dateIndex         = 10
	langIndex         = 11
	librateIndex      = 12
	keywordsIndex     = 13
)

/*
inp structure:
AUTHOR;GENRE;TITLE;SERIES;SERNO;FILE_TYPE;SIZE;LIBID;DEL;EXT;DATE;LANG;KEYWORDS;<CR><LF>
separator of fields (instead of ';') - <0x04>
end of record - <CR><LF> - <0x0D,0x0A>
*/

/*
#  DEFAULTSTRUCTURE = 'AUTHOR;GENRE;TITLE;SERIES;SERNO;FILE;SIZE;LIBID;DEL;EXT;DATE;LANG;LIBRATE;KEYWORDS';
class Book(namedtuple("Book", "authors genres title series ser_no filename size lib_id deleted ext date lang librate keywords archive_filename")):

*/

func parseAuthors(field string) []string {
	raw := strings.Split(field, ":")
	authors := make([]string, 0, len(raw))
	for _, author := range raw {
		author = strings.ReplaceAll(author, ",", " ")
		author = strings.TrimSpace(author)
		if author == "" {
			continue
		}
		authors = append(authors, author)
	}
	return authors
}

func parseGenres(field string) []string {
	raw := strings.Split(field, ":")
	genres := make([]string, 0, len(raw))
	for _, genre := range raw {
		genre = strings.TrimSpace(genre)
		if genre == "" {
			continue
		}
		genres = append(genres, genre)
	}
	return genres
}

func parseKeywords(field string) []string {
	raw := strings.Split(field, ":")
	keywords := make([]string, 0, len(raw))
	for _, keyword := range raw {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}
		keywords = append(keywords, strings.TrimSpace(keyword))
	}
	return keywords
}

func parseBook(ctx context.Context, fields []string) (entities.Book, error) {
	log := logs.GetFromContext(ctx)

	if len(fields) < columnCount {
		return entities.Book{}, fmt.Errorf("not enough fields in book: %v", fields)
	}
	if len(fields) > columnCount {
		fields = fields[:columnCount]
	}

	size, err := strconv.ParseInt(fields[sizeIndex], 10, 64)
	if err != nil {
		log.Error("error parsing size", zap.String("error", err.Error()), zap.String("size", fields[sizeIndex]))
		return entities.Book{}, ErrInvalidSize
	}

	date, err := time.Parse("2006-01-02", fields[dateIndex])
	if err != nil {
		log.Error("error parsing date", zap.String("error", err.Error()), zap.String("date", fields[dateIndex]))
		return entities.Book{}, ErrInvalidDate
	}

	return entities.Book{
		Authors:      parseAuthors(fields[authorIndex]),
		Genres:       parseGenres(fields[genreIndex]),
		Title:        fields[titleIndex],
		Series:       fields[seriesIndex],
		SeriesNumber: fields[seriesNumberIndex],
		Filename:     fields[fileIndex],
		Size:         size,
		LibID:        fields[libIDIndex],
		Ext:          fields[extIndex],
		Date:         date,
		Lang:         fields[langIndex],
		Keywords:     parseKeywords(fields[keywordsIndex]),
	}, nil
}

func ParseBooks(ctx context.Context, inp []byte) ([]entities.Book, error) {
	log := logs.GetFromContext(ctx).With(zap.String("action", "parse_books"))
	books := make([]entities.Book, 0)

	scanner := bufio.NewScanner(bytes.NewReader(inp))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\x04")
		book, err := parseBook(ctx, fields)
		if err != nil {
			log.Error("error parsing book", zap.String("error", err.Error()))
			continue
		}
		books = append(books, book)
	}
	return books, nil
}

func ParseBooksWithMetadata(ctx context.Context, inp []byte, metadata entities.BookMetadata) ([]entities.Book, error) {
	log := logs.GetFromContext(ctx).With(zap.String("action", "parse_books_with_metadata"))
	ctx = logs.WithLog(ctx, log)

	books := make([]entities.Book, 0)
	scanner := bufio.NewScanner(bytes.NewReader(inp))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\x04")
		book, err := parseBook(ctx, fields)
		if err != nil {
			log.Error("error parsing book", zap.String("error", err.Error()))
			continue
		}
		book.Metadata = metadata
		books = append(books, book)
	}
	return books, nil
}
