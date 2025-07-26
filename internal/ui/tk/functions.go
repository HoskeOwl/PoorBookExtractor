package ui

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func (impl *MainForm) updateStatus(text string) {
	impl.Statusbar.Configure(Txt(text))
}

func (impl *MainForm) refreshAuthorList(authors []string) {
	impl.AuthorList.Delete(impl.AuthorList.Children(""))
	for _, author := range authors {
		impl.AuthorList.Insert("", "end", Id(author), Txt(author), Value(author))
	}
}

func (impl *MainForm) openFile() {
	files := GetOpenFile(
		Title("Open library file"),
		Multiple(false),
		Filetypes(
			[]FileType{
				{"INPX files", []string{".inpx"}, ""},
			},
		),
	)

	filename := strings.Join(files, " ")

	if len(files) == 0 {
		return
	}
	impl.log.Debug("open file", zap.String("file", filename))
	impl.app.ClearStorage()
	var err error
	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} }()
		err = impl.app.ParseInpx(filename)
	}()
LOOP:
	for {
		select {
		case <-done:
			close(done)
			break LOOP
		default:
			progress := impl.app.GetProgress()
			impl.updateStatus(fmt.Sprintf("Parsing file %s - %d%%", filename, progress))
			Update()
			time.Sleep(50 * time.Millisecond)
		}
	}
	if err != nil {
		impl.updateStatus(fmt.Sprintf("Error parsing file: %s", err.Error()))
		return
	}
	authors := impl.app.GetAuthors("")
	impl.refreshAuthorList(authors)
	impl.updateStatus(fmt.Sprintf("File %s. Loaded %d authors.", files[0], len(authors)))
}

func (impl *MainForm) exportFiles() {
	book_ids := impl.ResultList.Selection("")
	if len(book_ids) == 0 {
		return
	}
	impl.log.Debug("export files", zap.Int("count", len(book_ids)))
	directory := ChooseDirectory(
		Title("Select directory to export files"),
	)
	impl.log.Info("export files", zap.String("directory", directory))

	if directory == "" {
		return
	}
	err := impl.app.ExportBooks(book_ids, directory)
	if err != nil {
		impl.updateStatus(fmt.Sprintf("Error exporting books: %s", err.Error()))
		return
	}
	impl.log.Info("exported books", zap.Strings("count", book_ids))
	impl.updateStatus(fmt.Sprintf("Exported %d books.", len(book_ids)))
}
