package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*ClickImageWidget)(nil)

// ClickImageWidget is a widget for displaying a separator with themeable color.
//
// Since: 1.4
type ClickImageWidget struct {
	widget.BaseWidget
	hovered  bool
	Data1    string
	Image    *canvas.Circle
	OnTapped func(*ClickImageWidget) `json:"-"`
}

// NewClickImage creates a new separator.
//
// Since: 1.4
func NewClickImage(img *canvas.Circle, fn func(*ClickImageWidget)) *ClickImageWidget {
	s := &ClickImageWidget{}
	s.Image = img
	s.OnTapped = fn
	s.ExtendBaseWidget(s)
	return s
}

func (b *ClickImageWidget) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped(b)
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (b *ClickImageWidget) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *ClickImageWidget) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (b *ClickImageWidget) MouseOut() {
	b.hovered = false
	b.Refresh()
}
func (b *ClickImageWidget) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// CreateRenderer returns a new renderer for the separator.
//
// Implements: fyne.Widget
func (b *ClickImageWidget) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	return &separatorRenderer{
		WidgetRenderer: widget.NewSimpleRenderer(
			container.NewStack(b.Image)),
		d: b,
	}
}

// MinSize returns the minimal size of the separator.
//
// Implements: fyne.Widget
func (b *ClickImageWidget) MinSize() fyne.Size {
	return b.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*separatorRenderer)(nil)

type separatorRenderer struct {
	fyne.WidgetRenderer
	d *ClickImageWidget
}

func (r *separatorRenderer) MinSize() fyne.Size {
	// t := theme.ClickImageThicknessSize()
	return fyne.NewSquareSize(10)
}

func (r *separatorRenderer) Refresh() {
	if r.d.hovered {
		r.d.Image.StrokeColor = color.Black
		r.d.Image.StrokeWidth = 2
	} else {
		r.d.Image.StrokeWidth = 0
	}
	r.d.Image.Refresh()
}
