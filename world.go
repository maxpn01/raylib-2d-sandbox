package main

import rl "github.com/gen2brain/raylib-go/raylib"

type GameObject interface {
	update(w *World, dt float32)
	draw()
}

type World struct {
	gameMap      *Map
	player       *Player
	fruitSpawner *FruitSpawner
	creature     *Creature
	entities     []GameObject
	hud          *HUD
}

func NewWorld() *World {
	w := &World{}
	w.gameMap = NewMap(
		rl.NewVector2(200, 200),
		rl.NewVector2(10, 10),
		rl.Black, rl.Gray,
	)
	w.player = NewPlayer(
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
	w.creature = NewCreature(
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
	w.fruitSpawner = NewFruitSpawner(
		rl.NewVector2(15, 15),
		rl.Yellow,
		0.25,
		1, 20,
	)

	w.entities = []GameObject{w.gameMap, w.player, w.fruitSpawner, w.creature}
	w.hud = NewHUD(w)
	return w
}

func (w *World) update(dt float32) {
	for _, entity := range w.entities {
		if w.player.isDeadFlag {
			break
		}
		entity.update(w, dt)
	}
}

func (w *World) draw() {
	for _, entity := range w.entities {
		entity.draw()
	}
}
