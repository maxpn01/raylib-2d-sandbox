package main

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

var window = &Window{
	width:   1440,
	height:  820,
	title:   "2D sandbox",
	bgColor: rl.Black,
}

var windowCenter = rl.NewVector2(float32(window.width)/2, float32(window.height)/2)

var world *World

var camera = rl.NewCamera2D(windowCenter, rl.Vector2{}, 0, 1)

var menu = NewMenu([]*MenuButton{
	NewMenuButton("play", ActionPlay),
	NewMenuButton("exit", ActionExit),
}, 30, rl.Black)

func restartGame() {
	world = NewWorld()
}

func main() {
	rl.InitWindow(int32(window.width), int32(window.height), window.title)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	restartGame()

	camera.Target = rl.NewVector2(
		world.player.pos.X+world.player.size.X/2,
		world.player.pos.Y+world.player.size.Y/2,
	)

	menuOpen := true
	exitRequested := false

	for !rl.WindowShouldClose() && !exitRequested {
		dt := rl.GetFrameTime()

		if rl.IsKeyPressed(rl.KeyEscape) {
			menuOpen = !menuOpen
		}

		if !menuOpen {
			world.update(dt)
			if world.player.isDeadFlag {
				restartGame()
			}

			camera.Target = rl.NewVector2(
				world.player.pos.X+world.player.size.X/2,
				world.player.pos.Y+world.player.size.Y/2,
			)

			world.hud.update(dt)
		}

		rl.BeginDrawing()
		rl.ClearBackground(window.bgColor)

		if menuOpen {
			switch menu.draw() {
			case ActionPlay:
				menuOpen = false
			case ActionExit:
				exitRequested = true
			}
		} else {
			rl.BeginMode2D(camera)

			world.draw()

			rl.EndMode2D()

			world.hud.draw()
		}

		rl.EndDrawing()
	}
}
