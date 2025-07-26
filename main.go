package main

import (
	"github.com/HoskeOwl/PoorBookExtractor/internal/app"
	ui "github.com/HoskeOwl/PoorBookExtractor/internal/ui/tk"
	"go.uber.org/zap"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	var mainForm *ui.MainForm
	app := app.NewApp(logger)

	mainForm = ui.NewForm(logger, app)
	mainForm.Wait()
}
