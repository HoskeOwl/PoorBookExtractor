package entities

import (
	"fmt"
	"time"
)

type BookMetadata struct {
	ArchiveName string
	Filepath    string
}

type Book struct {
	Metadata     BookMetadata
	Authors      []string
	Genres       []string
	Title        string
	Series       string
	SeriesNumber string
	Filename     string
	Size         int64
	LibID        string
	Ext          string
	Date         time.Time
	Lang         string
	Keywords     []string

	fullName string
}

func (b *Book) FullName() string {
	if b.fullName == "" {
		b.fullName = fmt.Sprintf("%s %s %s %s %s %d", b.Lang, b.Date.Format(time.DateOnly), b.Title, b.Series, b.SeriesNumber, b.Size)
	}
	return b.fullName
}
