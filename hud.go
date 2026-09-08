package main

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Text struct {
	pos   rl.Vector2
	color color.RGBA

	text     string
	fontSize int32
}

func (t *Text) update(dt float32) {}

func (t *Text) draw() {
	rl.DrawText(t.text, int32(t.pos.X), int32(t.pos.Y), t.fontSize, t.color)
}

type HUDText struct {
	Text
	label        string
	getStatValue func() float32
}

func NewHUDText(pos rl.Vector2, fontSize int32, color color.RGBA, text string, getValue func() float32) *HUDText {
	return &HUDText{
		Text: Text{
			pos:      pos,
			fontSize: fontSize,
			color:    color,
		},
		label:        text,
		getStatValue: getValue,
	}
}

func (h *HUDText) update(dt float32) {
	h.text = fmt.Sprintf("%s %.2f", h.label, h.getStatValue())
}
