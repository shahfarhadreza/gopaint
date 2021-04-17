package main

import (
	"log"
	"runtime"

	"gopaint/reza"
)

type Gopaint struct {
	reza.Application
	Title string
}

var app *Gopaint

func logError(msg string) {
	log.Printf("error: %s\n", msg)
}

func logInfo(msg string) {
	//log.Println(msg)
}

func init() {
	// Make sure everything runs under the same (main) thread (otherwise program hangs randomly)
	logInfo("Lock thread...")
	runtime.LockOSThread()
}

func main() {
	app = &Gopaint{
		Application: reza.NewApplication(),
		Title:       "GoPaint",
	}
	app.SetMainWindow(NewMainWindow())
	if app.GetMainWindow().Show() {
		app.Run()
	}
}
