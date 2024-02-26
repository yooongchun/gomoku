package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

const (
	boardSize = 15
	gridSize  = 30
)

func main() {
	a := app.New()
	w := a.NewWindow("Gomoku")

	content := container.NewWithoutLayout()

	// Draw the grid lines
	for i := 1; i < boardSize; i++ {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 2
		line.Move(fyne.NewPos(float32(i*gridSize), 0))
		line.Resize(fyne.NewSize(2, float32(boardSize*gridSize)))
		content.Add(line)

		line = canvas.NewLine(color.Black)
		line.StrokeWidth = 2
		line.Move(fyne.NewPos(0, float32(i*gridSize)))
		line.Resize(fyne.NewSize(float32(boardSize*gridSize), 2))
		content.Add(line)
	}

	w.SetContent(content)
	w.Resize(fyne.NewSize(float32(boardSize*gridSize), float32(boardSize*gridSize)))
	w.ShowAndRun()
}
