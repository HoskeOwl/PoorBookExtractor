package ui

import (
	_ "embed"
	"fmt"
	"runtime"
	"slices"
	"strings"

	"github.com/HoskeOwl/PoorBoockExtractor/internal/app"
	"github.com/HoskeOwl/PoorBoockExtractor/internal/entities"
	"go.uber.org/zap"
	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

//go:embed small_ico.png
var ico []byte

type MainForm struct {
	Menubar    *MenuWidget
	FindInput  *TEntryWidget
	AuthorList *TTreeviewWidget
	// ResultList *ListboxWidget
	ResultList *TTreeviewWidget
	Statusbar  *LabelWidget

	FindValue *Opt

	app *app.App
	log *zap.Logger
}

func (impl *MainForm) CreateMenubar() {
	menubar := Menu()
	menubar.AddCommand(Lbl("Open"), Underline(0), Accelerator("Ctrl+O"), Command(impl.openFile))
	Bind(App, "<Control-o>", Command(func() { menubar.Invoke(1) }))
	menubar.AddSeparator()
	menubar.AddCommand(Lbl("Export"), Underline(2), Accelerator("Ctrl+E"), Command(impl.exportFiles))
	Bind(App, "<Control-s>", Command(func() { menubar.Invoke(2) }))
	menubar.AddSeparator()
	menubar.AddCommand(Lbl("Exit"), Underline(1), Accelerator("Ctrl+Q"), ExitHandler())
	Bind(App, "<Control-q>", Command(func() { menubar.Invoke(3) }))

	impl.Menubar = menubar
}

func (impl *MainForm) CreateFind() *TFrameWidget {
	// find
	fr := TFrame()
	eVal := Textvariable("")
	findInput := fr.TEntry(eVal)
	findLabel := fr.Label(Txt("Find"))
	Pack(findLabel, Side("left"))
	Pack(findInput, Side("left"), Expand(true), Fill("x"))
	impl.FindInput = findInput
	impl.FindValue = &eVal

	return fr
}

func (impl *MainForm) CreateAuthorList() *TFrameWidget {
	// Lists
	// Authors
	fr := TFrame()
	// Scrollbar
	sb := fr.TScrollbar()
	Pack(sb, Side("right"), Fill("y"))

	lv := fr.TTreeview(Selectmode("browse"), Height(30),
		Yscrollcommand(func(e *Event) { e.ScrollSet(sb) }))

	lv.Heading("#0", Txt("Авторы"), Anchor("center"))
	lv.Column("#0", Width(320), Stretch(true), Separator(false))

	Pack(lv, Expand(true), Fill("both"))
	sb.Configure(Command(func(e *Event) { e.Yview(lv) }))
	impl.AuthorList = lv

	Bind(lv, "<<TreeviewSelect>>", Command(func() {
		author := lv.Selection("")
		impl.log.Debug("author selected", zap.Int("author", len(author)))
		if len(author) == 0 {
			return
		}
		impl.ResultList.Delete(impl.ResultList.Children(""))
		books := impl.app.GetAuthorBooks(author[0])
		slices.SortFunc(books, func(a, b entities.Book) int {
			return strings.Compare(a.FullName(), b.FullName())
		})
		impl.log.Debug("book got", zap.Int("books", len(books)))
		for _, book := range books {
			impl.ResultList.Insert("", "end", Id(book.LibID), Txt(book.FullName()), Value(book.LibID))
		}
	}))

	return fr
}

func (impl *MainForm) CreateResultList() *TFrameWidget {
	// Results
	fr := TFrame()
	// Scrollbar
	sb := fr.TScrollbar()
	Pack(sb, Side("right"), Fill("both"))
	// Listview
	lv := fr.TTreeview(Selectmode("extended"), Height(30),
		Yscrollcommand(func(e *Event) { e.ScrollSet(sb) }))

	lv.Heading("#0", Txt("Книги"), Anchor("center"))
	lv.Column("#0", Width(800), Stretch(true), Separator(false))
	Pack(lv, Expand(true), Fill("both"))
	sb.Configure(Command(func(e *Event) { e.Yview(lv) }))
	impl.ResultList = lv

	return fr
}

func (impl *MainForm) CreateStatusbar() *TFrameWidget {
	// Status bar
	fr := TFrame()
	lb := fr.Label(Txt("Status"), Borderwidth("1p"), Relief("sunken"), Justify("left"), Anchor("w"))
	impl.Statusbar = lb
	Pack(lb, Side("left"), Fill("x"), Expand(true), Pady("3p"))
	return fr
}

func (impl *MainForm) Wait() {
	App.Wait()
}

func NewForm(logger *zap.Logger, app *app.App) *MainForm {
	impl := &MainForm{app: app, log: logger}
	impl.CreateMenubar()

	find := impl.CreateFind()
	Grid(find, Row(0), Column(0), Sticky("we"), Padx("1m"), Pady("1m"))

	authorList := impl.CreateAuthorList()
	Grid(authorList, Row(1), Column(0), Sticky("nesw"))
	results := impl.CreateResultList()
	Grid(results, Row(1), Column(1), Sticky("nesw"))

	statusbar := impl.CreateStatusbar()
	Grid(statusbar, Row(2), Column(0), Sticky("we"), Columnspan(2))

	GridColumnConfigure(App, 1, Weight(1))
	GridRowConfigure(App, 1, Weight(1))

	App.WmTitle(fmt.Sprintf("%s on %s", App.WmTitle("FreeLibrary"), runtime.GOOS))
	ActivateTheme("azure light")
	App.Configure(Mnu(impl.Menubar), Width("150c"), Height("6c"))

	App.IconPhoto(NewPhoto(Data(ico)))

	return impl
}
