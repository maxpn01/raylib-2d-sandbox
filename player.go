package main

import (
	"image/color"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerStats struct {
	lvl int
	exp float32

	hp               float32
	speed            float32
	sprintMultiplier float32

	maxLvl int
	maxHP  float32
}

type PlayerGrowthStats struct {
	ExpPerFruit   float32
	HPPerLevel    float32
	SpeedPerLevel float32
}

type Player struct {
	pos   rl.Vector2
	size  rl.Vector2
	color color.RGBA

	stats PlayerStats

	growth PlayerGrowthStats
}

func NewPlayer(pos, size rl.Vector2, color color.RGBA, stats PlayerStats, growth PlayerGrowthStats) *Player {
	return &Player{
		pos:    pos,
		size:   size,
		color:  color,
		stats:  stats,
		growth: growth,
	}
}

func (p *Player) update(dt float32) {
	p.move(dt)
	p.handleFruitCollision(fruitSpawner)
	p.handleLevelUp()
}

func (p *Player) draw() {
	rl.DrawRectangleV(p.pos, p.size, p.color)
}

func (p *Player) move(dt float32) {
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

	speed := p.stats.speed
	if rl.IsKeyDown(rl.KeyLeftShift) {
		speed *= p.stats.sprintMultiplier
	}

	// normalize the diagonal speed
	if move.X != 0 || move.Y != 0 {
		move = rl.Vector2Normalize(move)

		p.pos.X += move.X * speed * dt
		p.pos.Y += move.Y * speed * dt
	}

	clamp(gameMap.edgePosX.start, &p.pos.X, &p.size.X, gameMap.edgePosX.end)
	clamp(gameMap.edgePosY.start, &p.pos.Y, &p.size.Y, gameMap.edgePosY.end)
}

func (p *Player) handleFruitCollision(fs *FruitSpawner) {
	for i, v := range slices.Backward(fs.fruits) {
		hasPlayerCollidedWithFruit := checkCollisions(p.pos, p.size, v.pos, v.size)

		if hasPlayerCollidedWithFruit {
			p.stats.exp += p.growth.ExpPerFruit
			fs.despawnFruit(i)
		}
	}
}

func (p *Player) handleLevelUp() {
	for p.stats.lvl < p.stats.maxLvl && p.stats.exp >= calcExpForNextLvl(p.stats.lvl) {
		p.stats.exp -= calcExpForNextLvl(p.stats.lvl)
		p.stats.lvl++

		p.stats.hp = min(p.stats.hp+p.growth.HPPerLevel, p.stats.maxHP)
		p.stats.speed += p.growth.SpeedPerLevel
	}
}
