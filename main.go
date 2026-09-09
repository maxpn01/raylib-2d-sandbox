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

var menu = NewMenu([]*MenuButton{
	NewMenuButton("play", ActionPlay),
	NewMenuButton("exit", ActionExit),
}, 30, rl.Black)

/* Game entities */
var gameMap = NewMap(
	rl.NewVector2(200, 200),
	rl.NewVector2(10, 10),
	rl.Black, rl.Gray,
)
var player = NewPlayer(
	rl.NewVector2(300, 300),
	rl.NewVector2(30, 30),
	rl.Red,
	PlayerStats{
		lvl:              1,
		hp:               1,
		speed:            300,
		sprintMultiplier: 2,
		maxLvl:           100,
		maxHP:            100,
	},
	PlayerGrowthStats{
		HPPerLevel:    0.25,
		SpeedPerLevel: 0.25,
	},
	1,
)
var creature = NewCreature(
	rl.NewVector2(windowCenter.X, windowCenter.Y),
	rl.NewVector2(30, 30),
	rl.Green,
	CreatureStats{
		lvl:    1,
		hp:     1,
		speed:  300,
		maxLvl: 100,
		maxHP:  100,
	},
	CreatureGrowthStats{
		HPPerLevel:    0.25,
		SpeedPerLevel: 0.25,
		ExpAsleep:     0.001,
		HPAsleep:      0.001,
	},
	1,
	120, 60,
)
var fruitSpawner = NewFruitSpawner(
	rl.NewVector2(15, 15),
	rl.Yellow,
	0.25,
	1, 20,
)

var entities = []GameObject{gameMap, player, fruitSpawner, creature}

/* HUD entities */
const hudTopPadding = 20
const hudTextSize = 22

var playerHpText = NewHUDText(
	rl.NewVector2(30, hudTopPadding),
	hudTextSize,
	rl.RayWhite,
	"hp:",
	func() float32 { return player.stats.hp },
)
var playerLvlText = NewHUDText(
	rl.NewVector2(130, hudTopPadding),
	hudTextSize,
	rl.RayWhite,
	"lvl:",
	func() float32 { return float32(player.stats.lvl) },
)
var playerExpText = NewHUDText(
	rl.NewVector2(230, hudTopPadding),
	hudTextSize,
	rl.RayWhite,
	"exp:",
	func() float32 { return player.stats.exp },
)
var playerSpeedText = NewHUDText(
	rl.NewVector2(350, hudTopPadding),
	hudTextSize,
	rl.RayWhite,
	"speed:",
	func() float32 { return player.stats.speed },
)

var creatureHpText = NewHUDText(
	rl.NewVector2(float32(window.width-210), hudTopPadding),
	hudTextSize,
	rl.RayWhite,
	"creature hp:",
	func() float32 { return creature.stats.hp },
)
var creatureLvlText = NewHUDText(
	rl.NewVector2(float32(window.width-420), hudTopPadding),
	hudTextSize,
	rl.RayWhite,
	"creature lvl:",
	func() float32 { return float32(creature.stats.lvl) },
)
var creatureExpText = NewHUDText(
	rl.NewVector2(float32(window.width-650), hudTopPadding),
	hudTextSize,
	rl.RayWhite,
	"creature exp:",
	func() float32 { return creature.stats.exp },
)

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
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	menuOpen := true
	exitRequested := false

	for !rl.WindowShouldClose() && !exitRequested {
		dt := rl.GetFrameTime()

		if rl.IsKeyPressed(rl.KeyEscape) {
			menuOpen = !menuOpen
		}

		if !menuOpen {
			for _, e := range entities {
				e.update(dt)
			}

			camera.Target = rl.NewVector2(
				player.pos.X+player.size.X/2.0,
				player.pos.Y+player.size.Y/2.0,
			)

			for _, e := range hud {
				e.update(dt)
			}
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

			for _, e := range entities {
				e.draw()
			}

			rl.EndMode2D()

			for _, e := range hud {
				e.draw()
			}
		}

		rl.EndDrawing()
	}
}
