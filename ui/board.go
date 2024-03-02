package ui

import (
	"fyne.io/fyne/v2/app"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type BoardUI struct {
	win      fyne.Window
	size     int     //棋盘大小
	cSize    float32 // 棋子大小
	turnover bool    // 用于交替出现黑白棋子
	bgColor  color.NRGBA
}

func NewBoardUI(size int, cSize float32, window fyne.Window) *BoardUI {
	return &BoardUI{
		win:      window,
		size:     size,
		cSize:    cSize,
		turnover: false,
		bgColor:  color.NRGBA{R: 0, G: 0, B: 0, A: 0},
	}

}
func (b *BoardUI) draw() {
	lblMsg := widget.NewLabel("label")
	lblMsg.Alignment = fyne.TextAlignCenter
	n1 := 0
	var ss = make([][col]int, col)
	btnList := make([]*ClickImageWidget, col*col)

	// 刷新，初始化相关数据
	btnRefresh := widget.NewButtonWithIcon("",
		theme.ViewRefreshIcon(), func() {
			n1 = 0
			for i := 0; i < b.size-2; i++ {
				for j := 0; j < b.size-2; j++ {
					i0 := i
					j0 := j
					ss[i0][j0] = 0
					btnList[n1].Image.FillColor = b.bgColor
					btnList[n1].Refresh()
					n1++
				}
			}
		})
	bg := canvas.NewRectangle(b.bgColor)
	bg.Resize(fyne.NewSquareSize(b.cSize * 14))

	// 棋盘容器
	cc := container.NewWithoutLayout(bg)

	for i := 0; i < b.size-2; i++ {
		for j := 0; j < b.size-2; j++ {
			i0 := i
			j0 := j
			ss[i0][j0] = 0
			var btn0 *ClickImageWidget
			btn0 = NewClickImage(canvas.NewCircle(b.bgColor), func(ci *ClickImageWidget) {
				if btn0.Image.FillColor == b.bgColor {
					if b.turnover {
						btn0.Image.FillColor = color.Black
						lblMsg.SetText("BLACK")
						ss[i0][j0] = -1
					} else {
						btn0.Image.FillColor = color.White
						// btn0.Icon = whites
						lblMsg.SetText("WHITE")
						ss[i0][j0] = 1
					}
					b.turnover = !b.turnover
					btn0.Refresh()
					// 判断获胜算法，请各位高手各显神通
					// result := ""
					// if ss[i0][j0] == 1 {
					// 	result = "White"
					// } else {
					// 	result = "Black"
					// }
					// dialog.ShowInformation("Win", result+"  win! Good Job!", w) 作者：bl4cyy https://www.bilibili.com/read/cv26941387/ 出处：bilibili
				}
			})
			btn0.Resize(fyne.NewSquareSize(b.cSize - 2))
			btn0.Move(fyne.NewPos(float32(j)*b.cSize+b.cSize/2+1, float32(i)*b.cSize+b.cSize/2+1))
			btnList[n1] = btn0
			n1++
			cc.Add(btn0)
		}
	}
	for i := 0; i < b.size; i++ {
		// 15条横线
		cc.Add(b.newLine(nil, b.cSize*(float32)(b.size-1), 0, 0, float32(i)*b.cSize))

		// 15条竖线
		cc.Add(b.newLine(nil, 0, b.cSize*(float32)(b.size-1), float32(i)*b.cSize, 0))
	}
	c := container.NewBorder(container.NewHBox(lblMsg, btnRefresh),
		nil, nil, nil, cc)

	b.win.SetContent(c)
}

func (b *BoardUI) newLine(c color.Color, x2, y2, x3, y3 float32) *canvas.Line {
	if c == nil {
		c = theme.ForegroundColor()
	}
	l := canvas.NewLine(c)
	l.Position1 = fyne.NewPos(0, 0)
	l.Position2 = fyne.NewPos(x2, y2)
	l.StrokeWidth = 2
	l.Move(fyne.NewPos(x3, y3))
	return l
}

func Show() {
	a := app.NewWithID("gomoku")
	//a.Settings().SetTheme(&theme1{})
	w := a.NewWindow("Gomoku AI")
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	NewBoardUI(15, 30, w).draw()

	w.ShowAndRun()
}
