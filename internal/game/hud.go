package game

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

func NewHUD(world *World, window *Window) *HUD {
	const (
		fontSize  = 22
		rowGap    = 10
		rowHeight = fontSize + rowGap
		top       = 20
		playerX   = 30
	)
	aiX := float32(window.width - 200)

	playerHpText := NewHUDText(
		rl.NewVector2(playerX, top),
		fontSize,
		rl.RayWhite,
		"hp:",
		func() float32 { return world.player.stats.hp },
	)
	playerLvlText := NewHUDText(
		rl.NewVector2(playerX, top+1*rowHeight),
		fontSize,
		rl.RayWhite,
		"lvl:",
		func() float32 { return float32(world.player.stats.lvl) },
	)
	playerExpText := NewHUDText(
		rl.NewVector2(playerX, top+2*rowHeight),
		fontSize,
		rl.RayWhite,
		"exp:",
		func() float32 { return world.player.stats.exp },
	)
	playerSpeedText := NewHUDText(
		rl.NewVector2(playerX, top+3*rowHeight),
		fontSize,
		rl.RayWhite,
		"speed:",
		func() float32 { return world.player.stats.speed },
	)

	aiHpText := NewHUDText(
		rl.NewVector2(aiX, top),
		fontSize,
		rl.RayWhite,
		"ai hp:",
		func() float32 { return world.ai.stats.hp },
	)
	aiLvlText := NewHUDText(
		rl.NewVector2(aiX, top+1*rowHeight),
		fontSize,
		rl.RayWhite,
		"ai lvl:",
		func() float32 { return float32(world.ai.stats.lvl) },
	)
	aiExpText := NewHUDText(
		rl.NewVector2(aiX, top+2*rowHeight),
		fontSize,
		rl.RayWhite,
		"ai exp:",
		func() float32 { return world.ai.stats.exp },
	)
	aiSpeedText := NewHUDText(
		rl.NewVector2(aiX, top+3*rowHeight),
		fontSize,
		rl.RayWhite,
		"ai speed:",
		func() float32 { return world.ai.stats.speed },
	)

	h := &HUD{texts: []*HUDText{
		playerHpText,
		playerLvlText,
		playerExpText,
		playerSpeedText,
		aiHpText,
		aiLvlText,
		aiExpText,
		aiSpeedText,
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
