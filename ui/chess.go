package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"image/jpeg"
	"os"
)

const (
	SIZE      = 15
	GRID_SIZE = 30
	MARGIN_X  = 10
	MARGIN_Y  = 10
	BSIZE     = GRID_SIZE * SIZE
)

func Show() {
	a := app.New()
	win := a.NewWindow("Gomoku")

	file, err := os.Open("ui/assets/img/board.jpg")
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer file.Close()

	img1, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Error: File could not be decoded")
		os.Exit(1)
	}

	fmt.Println(img1.Bounds())

	img := canvas.NewImageFromFile("ui/assets/img/board.jpg")
	img.FillMode = canvas.ImageFillOriginal

	//content := container.NewWithoutLayout()
	//content.Add(img)

	win.Resize(fyne.NewSize(800, 600))
	win.SetContent(img)
	win.ShowAndRun()
}
