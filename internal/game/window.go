package game

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	width   uint16
	height  uint16
	title   string
	bgColor color.RGBA
}

func NewWindow(width, height uint16, title string, bgColor color.RGBA) *Window {
	return &Window{
		width:   width,
		height:  height,
		title:   title,
		bgColor: bgColor,
	}
}

func (w *Window) getWindowCenter() rl.Vector2 {
	return rl.NewVector2(float32(w.width)/2, float32(w.height)/2)
}
