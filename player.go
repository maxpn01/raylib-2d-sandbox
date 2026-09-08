package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	pos   rl.Vector2
	size  rl.Vector2
	color color.RGBA

	lvl    int
	maxLvl int
	exp    float32

	speed float32
	hp    float32
	maxHP float32
}

func NewPlayer(pos, size rl.Vector2, color color.RGBA, speed, maxHP float32, maxLvl int) *Player {
	return &Player{
		pos:    pos,
		size:   size,
		color:  color,
		lvl:    1,
		maxLvl: maxLvl,
		speed:  speed,
		hp:     1,
		maxHP:  maxHP,
	}
}

const playerExpIncrement = 0.25

var playerHpIncrement float32 = 0.25
var playerSpeedIncrement float32 = 0.25

func (p *Player) update(dt float32) {
	p.movePlayer(dt)
	p.handlePlayerFruitCollision(fruitSpawner)
	p.handleLevelUp()
}

func (p *Player) draw() {
	rl.DrawRectangleV(p.pos, p.size, p.color)
}

func (p *Player) movePlayer(dt float32) {
	move := rl.Vector2{}

	if rl.IsKeyDown(rl.KeyW) {
		move.Y -= 1
	}
	if rl.IsKeyDown(rl.KeyS) {
		move.Y += 1
	}
	if rl.IsKeyDown(rl.KeyA) {
		move.X -= 1
	}
	if rl.IsKeyDown(rl.KeyD) {
		move.X += 1
	}

	speed := p.speed
	if rl.IsKeyDown(rl.KeyLeftShift) {
		speed *= 2
	}

	// normalize the diagonal speed
	if move.X != 0 || move.Y != 0 {
		move = rl.Vector2Normalize(move)

		p.pos.X += move.X * speed * dt
		p.pos.Y += move.Y * speed * dt
	}

	clamp(gameMap.edgePosX.start, &p.pos.X, &p.size.X, gameMap.edgePosX.end)
	clamp(gameMap.edgePosX.start, &p.pos.Y, &p.size.Y, gameMap.edgePosY.end)
}

func (p *Player) handlePlayerFruitCollision(fs *FruitSpawner) {
	for i := len(fs.fruits) - 1; i >= 0; i-- {
		hasPlayerCollidedWithFruit := checkCollisions(p.pos, p.size, fs.fruits[i].pos, fs.fruits[i].size)

		if hasPlayerCollidedWithFruit {
			p.exp += playerExpIncrement
			fs.despawnFruit(i)
		}
	}
}

func (p *Player) handleLevelUp() {
	for p.lvl < p.maxLvl && p.exp >= calcExpForNextLvl(p.lvl) {
		p.exp -= calcExpForNextLvl(p.lvl)
		p.lvl++

		p.hp = min(p.hp+playerHpIncrement, p.maxHP)
		p.speed += playerSpeedIncrement
	}
}
