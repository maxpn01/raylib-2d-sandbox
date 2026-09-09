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

type HUD struct {
	texts []*HUDText
}

func NewHUD(w *World) *HUD {
	playerHpText := NewHUDText(
		rl.NewVector2(30, 20),
		22,
		rl.RayWhite,
		"hp:",
		func() float32 { return w.player.stats.hp },
	)
	playerLvlText := NewHUDText(
		rl.NewVector2(130, 20),
		22,
		rl.RayWhite,
		"lvl:",
		func() float32 { return float32(w.player.stats.lvl) },
	)
	playerExpText := NewHUDText(
		rl.NewVector2(230, 20),
		22,
		rl.RayWhite,
		"exp:",
		func() float32 { return w.player.stats.exp },
	)
	playerSpeedText := NewHUDText(
		rl.NewVector2(350, 20),
		22,
		rl.RayWhite,
		"speed:",
		func() float32 { return w.player.stats.speed },
	)

	creatureHpText := NewHUDText(
		rl.NewVector2(float32(window.width-210), 20),
		22,
		rl.RayWhite,
		"creature hp:",
		func() float32 { return w.creature.stats.hp },
	)
	creatureLvlText := NewHUDText(
		rl.NewVector2(float32(window.width-420), 20),
		22,
		rl.RayWhite,
		"creature lvl:",
		func() float32 { return float32(w.creature.stats.lvl) },
	)
	creatureExpText := NewHUDText(
		rl.NewVector2(float32(window.width-650), 20),
		22,
		rl.RayWhite,
		"creature exp:",
		func() float32 { return w.creature.stats.exp },
	)

	h := &HUD{texts: []*HUDText{
		playerHpText,
		playerLvlText,
		playerExpText,
		playerSpeedText,
		creatureHpText,
		creatureLvlText,
		creatureExpText,
	}}
	h.update(0)
	return h
}

func (h *HUD) update(dt float32) {
	for _, text := range h.texts {
		text.update(dt)
	}
}

func (h *HUD) draw() {
	for _, text := range h.texts {
		text.draw()
	}
}
