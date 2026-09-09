package game

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	width   int32
	height  int32
	title   string
	bgColor color.RGBA
}

func NewWindow(width, height int32, title string, bgColor color.RGBA) *Window {
	return &Window{
		width:   width,
		height:  height,
		title:   title,
		bgColor: bgColor,
	}
}

func (w *Window) initWindow() {
	rl.InitWindow(w.width, w.height, w.title)

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)
}

func (w *Window) close() {
	rl.CloseWindow()
}

func (w *Window) getWindowCenter() rl.Vector2 {
	return rl.NewVector2(float32(w.width)/2, float32(w.height)/2)
}
