package main

import (
	"image/color"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CreatureActionState int

const (
	ActionSearchingFood CreatureActionState = iota
	ActionEating
	ActionSleeping
)

type CreatureStats struct {
	lvl int
	exp float32

	hp    float32
	speed float32

	maxLvl int
	maxHP  float32
}

type CreatureGrowthStats struct {
	ExpPerFruit float32

	HPPerLevel    float32
	SpeedPerLevel float32

	ExpAsleep float32
	HPAsleep  float32
}

type Creature struct {
	pos   rl.Vector2
	size  rl.Vector2
	color color.RGBA

	stats CreatureStats

	growth CreatureGrowthStats

	actionState CreatureActionState

	awakeSeconds       float32
	asleepSeconds      float32
	awakeTimerSeconds  float32
	asleepTimerSeconds float32
}

func NewCreature(pos, size rl.Vector2, color color.RGBA, stats CreatureStats, growth CreatureGrowthStats, awake, asleep float32) *Creature {
	return &Creature{
		pos:                pos,
		size:               size,
		color:              color,
		stats:              stats,
		growth:             growth,
		awakeSeconds:       awake,
		asleepSeconds:      asleep,
		awakeTimerSeconds:  awake,
		asleepTimerSeconds: asleep,
	}
}

func (c *Creature) update(dt float32) {
	c.updateAwakenessStatus(dt)

	switch c.actionState {
	case ActionSearchingFood:
		c.searchFood(fruitSpawner, dt)
	case ActionEating:
		c.eatFood()
	case ActionSleeping:
		c.sleep(dt)
	}

	c.handleLevelUp()
}

func (c *Creature) draw() {
	rl.DrawRectangleV(c.pos, c.size, c.color)
}

func (c *Creature) searchFood(fs *FruitSpawner, dt float32) {
	if len(fs.fruits) == 0 {
		return
	}

	var closestFruitDistance float32 = findSquaredEuclideanDistance(c.pos, fs.fruits[0].pos)
	var closestFruitPos rl.Vector2 = fs.fruits[0].pos
	var closestFruitSize rl.Vector2 = fs.fruits[0].size

	for i, fruit := range fs.fruits {
		if i == 0 {
			continue
		}

		distance := findSquaredEuclideanDistance(c.pos, fruit.pos)

		if distance <= closestFruitDistance {
			closestFruitDistance = distance
			closestFruitPos = fruit.pos
			closestFruitSize = fruit.size
		}
	}

	c.move(closestFruitPos, closestFruitSize, dt)

	hasCreatureCollidedWithFruit, fruitIndex := c.checkFruitCollision(fs)

	if hasCreatureCollidedWithFruit && fruitIndex >= 0 && fruitIndex < len(fs.fruits) {
		fs.despawnFruit(fruitIndex)
		c.actionState = ActionEating
	}
}

func (c *Creature) move(targetPos, targetSize rl.Vector2, dt float32) {
	// Align our center with the fruit's center, accounting for their different sizes.
	target := rl.NewVector2(
		targetPos.X+targetSize.X/2-c.size.X/2,
		targetPos.Y+targetSize.Y/2-c.size.Y/2,
	)
	// Travel at most speed * dt, stopping at the target rather than overshooting it.
	c.pos = rl.Vector2MoveTowards(c.pos, target, c.stats.speed*dt)

	clamp(gameMap.edgePosX.start, &c.pos.X, &c.size.X, gameMap.edgePosX.end)
	clamp(gameMap.edgePosY.start, &c.pos.Y, &c.size.Y, gameMap.edgePosY.end)
}

func (c *Creature) checkFruitCollision(fs *FruitSpawner) (bool, int) {
	for i, v := range slices.Backward(fs.fruits) {
		hasCreatureCollidedWithFruit := checkCollisions(c.pos, c.size, v.pos, v.size)

		if hasCreatureCollidedWithFruit {
			return true, i
		}
	}

	return false, -1
}

func (c *Creature) eatFood() {
	c.stats.exp += c.growth.ExpPerFruit
	c.actionState = ActionSearchingFood
}

func (c *Creature) updateAwakenessStatus(dt float32) {
	switch c.actionState {
	case ActionSearchingFood:
		if c.awakeTimerSeconds > 0 {
			c.awakeTimerSeconds -= 1 * dt
		} else {
			c.actionState = ActionSleeping
			c.asleepTimerSeconds = c.asleepSeconds
		}
	case ActionSleeping:
		if c.asleepTimerSeconds > 0 {
			c.asleepTimerSeconds -= 1 * dt
		} else {
			c.actionState = ActionSearchingFood
			c.awakeTimerSeconds = c.awakeSeconds
		}
	}
}

func (c *Creature) sleep(dt float32) {
	c.stats.exp += c.growth.ExpAsleep * dt

	if c.stats.hp < c.stats.maxHP {
		c.stats.hp = min(c.stats.hp+c.growth.HPAsleep*dt, c.stats.maxHP)
	}
}

func (c *Creature) handleLevelUp() {
	for c.stats.lvl < c.stats.maxLvl && c.stats.exp >= calcExpForNextLvl(c.stats.lvl) {
		c.stats.exp -= calcExpForNextLvl(c.stats.lvl)
		c.stats.lvl++

		c.stats.hp = min(c.stats.hp+c.growth.HPPerLevel, c.stats.maxHP)
		c.stats.speed += c.growth.SpeedPerLevel
	}
}
