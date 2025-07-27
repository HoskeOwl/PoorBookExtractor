package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/HoskeOwl/PoorBookExtractor/internal/app"
	ui "github.com/HoskeOwl/PoorBookExtractor/internal/ui/tk"
	"go.uber.org/zap"
	_ "modernc.org/tk9.0/themes/azure"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	lvl := os.Getenv("DEBUG")
	var logger *zap.Logger
	var err error
	if lvl == "true" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	var mainForm *ui.MainForm
	app := app.NewApp(logger)

	mainForm = ui.NewForm(logger, app)
	mainForm.Wait()
}
