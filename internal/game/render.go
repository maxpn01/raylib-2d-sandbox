package game

import rl "github.com/gen2brain/raylib-go/raylib"

// RenderContext describes the camera and viewport used to draw the world.
type RenderContext struct {
	Camera       rl.Camera2D
	WindowWidth  int32
	WindowHeight int32
}
