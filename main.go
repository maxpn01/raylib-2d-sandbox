package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

/* Window */
type Window struct {
	width   uint16
	height  uint16
	title   string
	bgColor color.RGBA
}

var window = &Window{
	width:   1440,
	height:  820,
	title:   "2D sandbox",
	bgColor: rl.Black,
}

var windowCenter = rl.NewVector2(float32(window.width)/2, float32(window.height)/2)

/* Game entities */
var gameMap = NewMap(rl.NewVector2(200, 200), rl.NewVector2(10, 10), rl.Black, rl.Gray)
var player = NewPlayer(rl.NewVector2(300, 300), rl.NewVector2(30, 30), rl.Red, 400, 100, 100)
var fruitSpawner = NewFruitSpawner(rl.NewVector2(15, 15), rl.Yellow, 1, 20)
var creature = NewCreature(rl.NewVector2(windowCenter.X, windowCenter.Y), rl.NewVector2(30, 30), rl.Green, 200, 100, 100)

var entities = []GameObject{gameMap, player, fruitSpawner, creature}

/* HUD entities */
const hudTopPadding = 20
const hudTextSize = 22

var playerHpText = NewHUDText(rl.NewVector2(30, hudTopPadding), hudTextSize, rl.RayWhite, "hp:", func() float32 { return player.hp })
var playerLvlText = NewHUDText(rl.NewVector2(130, hudTopPadding), hudTextSize, rl.RayWhite, "lvl:", func() float32 { return float32(player.lvl) })
var playerExpText = NewHUDText(rl.NewVector2(230, hudTopPadding), hudTextSize, rl.RayWhite, "exp:", func() float32 { return player.exp })
var playerSpeedText = NewHUDText(rl.NewVector2(350, hudTopPadding), hudTextSize, rl.RayWhite, "speed:", func() float32 { return player.speed })

var creatureHpText = NewHUDText(rl.NewVector2(float32(window.width-210), hudTopPadding), hudTextSize, rl.RayWhite, "creature hp:", func() float32 { return creature.hp })
var creatureLvlText = NewHUDText(rl.NewVector2(float32(window.width-420), hudTopPadding), hudTextSize, rl.RayWhite, "creature lvl:", func() float32 { return float32(creature.lvl) })
var creatureExpText = NewHUDText(rl.NewVector2(float32(window.width-650), hudTopPadding), hudTextSize, rl.RayWhite, "creature exp:", func() float32 { return creature.exp })

var hud = []GameObject{
	playerHpText,
	playerLvlText,
	playerExpText,
	playerSpeedText,
	creatureHpText,
	creatureLvlText,
	creatureExpText,
}

/* Camera */
var playerTarget = rl.NewVector2(player.pos.X+player.size.X/2.0, player.pos.Y+player.size.Y/2.0)
var camera = rl.NewCamera2D(windowCenter, playerTarget, 0, 1)

func main() {
	rl.InitWindow(int32(window.width), int32(window.height), window.title)
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		for _, e := range entities {
			e.update(dt)
		}

		camera.Target = rl.NewVector2(player.pos.X+player.size.X/2.0, player.pos.Y+player.size.Y/2.0)

		for _, e := range hud {
			e.update(dt)
		}

		rl.BeginDrawing()
		rl.ClearBackground(window.bgColor)

		rl.BeginMode2D(camera)

		for _, e := range entities {
			e.draw()
		}

		rl.EndMode2D()

		for _, e := range hud {
			e.draw()
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
