package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/HoskeOwl/PoorBookExtractor/internal/entities"
	"go.uber.org/zap"
	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func (impl *MainForm) updateStatus(text string) {
	impl.Statusbar.Configure(Txt(text))
}

func (impl *MainForm) clearLists() {
	impl.ResultList.Delete(impl.ResultList.Children(""))
	impl.AuthorList.Delete(impl.AuthorList.Children(""))
}

func (impl *MainForm) refreshAuthorList() {
	impl.AuthorList.Delete(impl.AuthorList.Children(""))
	for author, books := range impl.app.IterBooksByAuthor() {
		impl.app.SortBooks(books)
		impl.AuthorList.Insert("", "end", Id(author), Txt(author))
		for _, book := range books {
			impl.AuthorList.Insert(author, "end", Id(book.ExtendId(author)), Txt(book.FullName()))
		}
	}
}

func (impl *MainForm) clearFind() {
	impl.FindInput.Configure(Textvariable(""))
	impl.refreshAuthorList()
}

func (impl *MainForm) findAuthor() {
	author := impl.FindInput.Textvariable()
	if author == "" {
		impl.refreshAuthorList()
		return
	}
	impl.AuthorList.Delete(impl.AuthorList.Children(""))
	for _, author := range impl.app.GetAuthors(author) {
		books := impl.app.GetAuthorBooks(author)
		impl.app.SortBooks(books)
		impl.AuthorList.Insert("", "end", Id(author), Txt(author))
		for _, book := range books {
			impl.AuthorList.Insert(author, "end", Id(book.ExtendId(author)), Txt(book.FullName()))
		}
	}
}

func (impl *MainForm) removeFromResultList() {
	lst := impl.ResultList.Selection("")
	if len(lst) == 0 {
		return
	}
	selected := lst[0]
	parent := impl.ResultList.Parent(selected)
	if parent == "" {
		children := impl.ResultList.Children(selected)
		if len(children) != 0 {
			impl.ResultList.Delete(children)
		}
		impl.ResultList.Delete(selected)
		return
	}
	impl.ResultList.Delete(selected)
	children := impl.ResultList.Children(parent)
	if len(children) == 0 {
		impl.ResultList.Delete(parent)
	}

}

func (impl *MainForm) checkAuthorExistsInResult(item string) bool {
	for _, parent := range impl.ResultList.Children("") {
		if parent == item {
			return true
		}
	}
	return false
}

func (impl *MainForm) checkBookExistsInResult(item string) bool {
	bookId := entities.GetBookIdFromExtended(item)
	for _, parent := range impl.ResultList.Children("") {
		for _, child := range impl.ResultList.Children(parent) {
			childBookId := entities.GetBookIdFromExtended(child)
			if childBookId == bookId {
				book, ok := impl.app.GetBook(bookId)
				if !ok {
					impl.log.Error("book from list not found", zap.String("book", item))
					return true
				}
				impl.updateStatus(fmt.Sprintf("Author %s already contains book %s", parent, book.FullName()))
				originalColor := impl.Statusbar.Background()
				impl.Statusbar.Configure(Background("red"))
				Update()
				time.Sleep(50 * time.Millisecond)
				impl.Statusbar.Configure(Background(originalColor))
				Update()
				return true
			}
		}
	}
	return false
}

func (impl *MainForm) addToResultList() {
	lst := impl.AuthorList.Selection("")
	if len(lst) == 0 {
		return
	}
	selected := lst[0]
	parent := impl.AuthorList.Parent(selected)
	impl.log.Debug("addToResultList", zap.String("selected", selected), zap.String("parent", parent))
	if parent == "" {
		if !impl.checkAuthorExistsInResult(selected) {
			// add author to result
			impl.ResultList.Insert("", "end", Id(selected), Txt(selected))
		}
		books := impl.app.GetAuthorBooks(selected)
		impl.app.SortBooks(books)
		for _, book := range books {
			if impl.checkBookExistsInResult(book.ExtendId(selected)) {
				continue
			}
			impl.ResultList.Insert(selected, "end", Id(book.ExtendId(selected)), Txt(book.FullName()))
		}
		impl.ResultList.Item(selected, Open(true))
		return
	}
	if !impl.checkAuthorExistsInResult(parent) {
		// add author to result
		impl.ResultList.Insert("", "end", Id(parent), Txt(parent))
	}
	// book selected. So we have an book id
	book, ok := impl.app.GetBook(entities.GetBookIdFromExtended(selected))
	if !ok {
		return
	}
	if impl.checkBookExistsInResult(book.ExtendId(parent)) {
		return
	}
	impl.ResultList.Insert(parent, "end", Id(book.ExtendId(parent)), Txt(book.FullName()))
	impl.ResultList.Item(parent, Open(true))
}

func (impl *MainForm) clearResultList() {
	impl.ResultList.Delete(impl.ResultList.Children(""))
}

func (impl *MainForm) openFile() {
	home := os.Getenv("HOME")
	files := GetOpenFile(
		Initialdir(home),
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
	impl.clearLists()
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
	impl.refreshAuthorList()
	impl.updateStatus(fmt.Sprintf("Imported %d authors, %d books.", impl.app.AuthorsLen(), impl.app.BooksLen()))
}

func (impl *MainForm) exportFiles() {
	authors := impl.ResultList.Children("")
	if len(authors) == 0 {
		return
	}
	bookIdsByAuthor := make(map[string][]string)
	totalBooks := 0
	for _, author := range authors {
		books := impl.ResultList.Children(author)
		if len(books) == 0 {
			continue
		}
		author_book_ids := impl.ResultList.Children(author)
		if len(author_book_ids) == 0 {
			continue
		}
		for _, book := range author_book_ids {
			bookIdsByAuthor[author] = append(bookIdsByAuthor[author], entities.GetBookIdFromExtended(book))
			totalBooks++
		}
	}
	impl.log.Debug("export files", zap.Int("count", totalBooks))
	home := os.Getenv("HOME")
	directory := ChooseDirectory(
		Initialdir(home),
		Title("Select directory to export files"),
	)
	impl.log.Info("export files", zap.String("directory", directory))

	if directory == "" {
		return
	}
	err := impl.app.ExportBooks(bookIdsByAuthor, directory)
	if err != nil {
		impl.updateStatus(fmt.Sprintf("Error exporting books: %s", err.Error()))
		return
	}
	impl.log.Info("exported books", zap.Int("count", totalBooks))
	impl.updateStatus(fmt.Sprintf("Exported %d books to %s", totalBooks, directory))
}
