package game

import rl "github.com/gen2brain/raylib-go/raylib"

type RenderContext struct {
	Camera       rl.Camera2D
	WindowWidth  int32
	WindowHeight int32
}
