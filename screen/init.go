package screen
import (
	"log"
	"os"
	"time"

	"github.com/kbinani/screenshot"
	"screem.frankmayer.io/ui"
)

var (
	StopHosting chan bool
)

func InitHosting() {
	StopHosting = make(chan bool, 1)
	go backgroundLoop()
}

func backgroundLoop() {
	for {
		select {
		case <-StopHosting:
			log.Println("Stopping hosting")
			return
		default:
			captureScreen()
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func captureScreen() {
	if ui.ScreenNum < 0 {
		return
	}

	bounds := screenshot.GetDisplayBounds(ui.ScreenNum)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Println("Failed to capture screen:", err)
		return
	}

	ui.UpdateScreenPreview(img)

	fileName := "screenshot.avif"
	file, _ := os.Create(fileName)
	defer file.Close()
	avif.Encode(file, img)
}
