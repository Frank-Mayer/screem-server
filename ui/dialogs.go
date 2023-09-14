package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Guest() {
	w := App.NewWindow("Screem :: Guest")
	label := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(label))
	w.ShowAndRun()
}
