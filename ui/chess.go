package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

const (
	boardSize = 15
	cellSize  = 30
)

func Show() {
	a := app.New()
	w := a.NewWindow("五子棋")

	board := make([][]*canvas.Rectangle, boardSize)
	for i := range board {
		board[i] = make([]*canvas.Rectangle, boardSize)
		for j := range board[i] {
			board[i][j] = canvas.NewRectangle(theme.BackgroundColor())
			board[i][j].SetMinSize(fyne.NewSize(cellSize, cellSize))
		}
	}

	grid := container.NewGridWithRows(boardSize)
	for _, row := range board {
		for _, cell := range row {
			grid.Add(cell)
		}
	}

	w.SetContent(grid)
	w.Resize(fyne.NewSize(boardSize*cellSize, boardSize*cellSize))
	w.ShowAndRun()
}
