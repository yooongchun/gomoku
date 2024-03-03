package ui

import (
	_ "embed"
	"fyne.io/fyne/v2/app"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/img/black.png
var blackChessSource []byte

//go:embed assets/img/white.png
var whiteChessSource []byte

type BoardUI struct {
	win        fyne.Window
	size       int     //棋盘大小
	cSize      float32 // 棋子大小
	turnover   bool    // 用于交替出现黑白棋子
	bgColor    color.NRGBA
	blackChess *fyne.StaticResource
	whiteChess *fyne.StaticResource
	btnList    []*widget.Button
	tip        *widget.Label
}

func NewBoardUI(size int, cSize float32, window fyne.Window) *BoardUI {
	return &BoardUI{
		win:        window,
		size:       size,
		cSize:      cSize,
		turnover:   false,
		bgColor:    color.NRGBA{R: 245, G: 245, B: 220, A: 255},
		blackChess: fyne.NewStaticResource("black.png", blackChessSource),
		whiteChess: fyne.NewStaticResource("white.png", whiteChessSource),
		btnList:    make([]*widget.Button, (size-2)*(size-2)),
		tip:        widget.NewLabel("who's turn?"),
	}

}
func (b *BoardUI) resetBoard() {
	for i := 0; i < b.size-2; i++ {
		for j := 0; j < b.size-2; j++ {
			n := i*(b.size-2) + j
			b.btnList[n].Icon = nil
			b.btnList[n].Refresh()
			b.turnover = false
		}
	}
}

func (b *BoardUI) onClicked(btn *widget.Button) {
	if btn.Icon != nil {
		return
	}
	if b.turnover {
		btn.Icon = b.blackChess
		b.tip.SetText("BLACK")
	} else {
		btn.Icon = b.whiteChess
		b.tip.SetText("WHITE")
	}
	b.turnover = !b.turnover
	btn.Refresh()
	btn.Show()

}
func (b *BoardUI) draw() {
	b.tip.Alignment = fyne.TextAlignCenter
	// 刷新，初始化相关数据
	btnRefresh := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), b.resetBoard)
	bg := canvas.NewRectangle(b.bgColor)
	bg.Resize(fyne.NewSquareSize(b.cSize * (float32)(b.size-1)))

	// 棋盘容器
	cc := container.NewWithoutLayout(bg)
	for i := 0; i < b.size-2; i++ {
		for j := 0; j < b.size-2; j++ {
			btn := widget.NewButtonWithIcon("", nil, func() { b.onClicked(b.btnList[i*(b.size-2)+j]) })
			btn.Resize(fyne.NewSquareSize(b.cSize))
			btn.Move(fyne.NewPos(float32(j)*b.cSize+b.cSize/2, float32(i)*b.cSize+b.cSize/2))
			b.btnList[i*(b.size-2)+j] = btn
			cc.Add(btn)
		}
	}
	for i := 0; i < b.size; i++ {
		// 15条横线
		cc.Add(b.newLine(nil, b.cSize*(float32)(b.size-1), 0, 0, float32(i)*b.cSize))
		// 15条竖线
		cc.Add(b.newLine(nil, 0, b.cSize*(float32)(b.size-1), float32(i)*b.cSize, 0))
	}
	c := container.NewBorder(container.NewHBox(b.tip, btnRefresh),
		nil, nil, nil, cc)

	b.win.SetContent(c)
}

func (b *BoardUI) newLine(c color.Color, x2, y2, x3, y3 float32) *canvas.Line {
	if c == nil {
		c = color.NRGBA{R: 50, G: 50, B: 50, A: 100}
	}
	l := canvas.NewLine(c)
	l.Position1 = fyne.NewPos(0, 0)
	l.Position2 = fyne.NewPos(x2, y2)
	l.StrokeWidth = 1
	l.Move(fyne.NewPos(x3, y3))
	return l
}

const (
	SIZE   = 15
	CELL   = 30
	WIDTH  = SIZE * CELL * 0.95
	HEIGHT = SIZE * CELL * 1.05
)

func Show() {
	a := app.NewWithID("gomoku")
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Gomoku AI")
	w.Resize(fyne.NewSize(WIDTH, HEIGHT))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	NewBoardUI(SIZE, CELL, w).draw()

	w.ShowAndRun()
}
