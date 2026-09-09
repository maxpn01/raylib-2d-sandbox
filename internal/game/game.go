package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var window = NewWindow(
	1440,
	820,
	"2D sandbox",
	rl.Black,
)

var camera = rl.NewCamera2D(window.getWindowCenter(), rl.Vector2{}, 0, 1)

var menu = NewMenu([]*MenuButton{
	NewMenuButton("play", ActionPlay),
	NewMenuButton("exit", ActionExit),
}, 30, rl.Black)

var world *World

func restartGame() {
	world = NewWorld()
}

func RunGame() {
	rl.InitWindow(int32(window.width), int32(window.height), window.title)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	restartGame()

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
