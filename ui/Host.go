package ui

import (
	"fmt"
	"image"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
)

var (
    ScreenNum int = -1
    img *canvas.Image = nil
    w fyne.Window
)

func Host() {
	w = App.NewWindow("Screem :: Hosting")
	w.Resize(fyne.NewSize(800, 600))

	w.SetContent(
		container.NewGridWithColumns(2,
			container.NewCenter(container.NewVBox(screenSelector())),
			container.NewCenter(screenPreview()),
		),
	)
	w.ShowAndRun()
}

func screenPreview() *canvas.Image {
	img = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 200, 200)))
	return img
}

func screenSelector() *widget.Select {
	numOfScreens := screenshot.NumActiveDisplays()
	screenSlectArray := make([]string, numOfScreens)
	for i := 0; i < numOfScreens; i++ {
		screenSlectArray[i] = fmt.Sprintf("Screen %d", i)
	}
    var screenSelect *widget.Select
	screenSelect = widget.NewSelect(screenSlectArray, func(value string) {
		log.Println("Selected", value)
        ScreenNum = screenSelect.SelectedIndex()
        bounds := screenshot.GetDisplayBounds(ScreenNum)
        img.SetMinSize(fyne.NewSize(float32(bounds.Dx())/5, float32(bounds.Dy())/5))
        w.Content().Refresh()
	})
	return screenSelect
}

func UpdateScreenPreview(newImg *image.RGBA) {
    img.Image = newImg
    img.Refresh()
}
