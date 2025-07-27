package ui

import (
	_ "embed"
	"fmt"
	"runtime"

	"github.com/HoskeOwl/PoorBookExtractor/internal/app"
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
	Bind(findInput, "<Return>", Command(impl.findAuthor))
	findLabel := fr.Label(Txt("Find"))
	clearBtn := fr.Button(Txt("❌"), Width(1), Height(1), Command(impl.clearFind))
	findBtn := fr.Button(Txt("🔍"), Width(1), Height(1), Command(impl.findAuthor))
	Pack(findLabel, Side("left"))
	Pack(findInput, Side("left"), Expand(true), Fill("x"))
	Pack(findBtn, Side("right"), Expand(false), Fill("x"))
	Pack(clearBtn, Side("right"), Expand(false), Fill("x"))
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

	lv.Heading("#0", Txt("Books by authors"), Anchor("center"))
	lv.Column("#0", Width(500), Stretch(true), Separator(false))

	Pack(lv, Expand(true), Fill("both"))
	sb.Configure(Command(func(e *Event) { e.Yview(lv) }))
	impl.AuthorList = lv

	Bind(lv, "<<TreeviewOpen>>", Command(impl.authorListOpen))

	return fr
}

func (impl *MainForm) CreateResultList() *TFrameWidget {
	// Results
	fr := TFrame()

	// Additional frame for buttons
	buttonFrame := fr.TFrame()
	Pack(buttonFrame, Side("left"), Padx("2p"), Anchor("center"))

	// Two small buttons
	addBtn := buttonFrame.Button(Txt("➡"), Width(1), Height(1), Command(impl.addToResultList))
	removeBtn := buttonFrame.Button(Txt("⬅"), Width(1), Height(1), Command(impl.removeFromResultList))
	clearBtn := buttonFrame.Button(Txt("❌"), Width(1), Height(1), Command(impl.clearResultList))

	// Pack buttons in the button frame, centered vertically
	Pack(addBtn, Side("top"), Pady("1p"))
	Pack(removeBtn, Side("top"), Pady("1p"))
	Pack(clearBtn, Side("top"), Pady("1p"))

	// Scrollbar
	sb := fr.TScrollbar()
	Pack(sb, Side("right"), Fill("both"))
	// Listview
	lv := fr.TTreeview(Selectmode("browse"), Height(30),
		Yscrollcommand(func(e *Event) { e.ScrollSet(sb) }))

	lv.Heading("#0", Txt("Books to export"), Anchor("center"))
	lv.Column("#0", Width(600), Stretch(true), Separator(false))
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
	App.Center().Wait()
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
	App.WmTitle("PoorBookExtractor")

	App.IconPhoto(NewPhoto(Data(ico)))

	return impl
}
