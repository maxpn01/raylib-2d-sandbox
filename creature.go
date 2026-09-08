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

type Creature struct {
	pos   rl.Vector2
	size  rl.Vector2
	color color.RGBA

	lvl    int
	maxLvl int
	exp    float32

	speed float32
	hp    float32
	maxHP float32

	actionState CreatureActionState

	awakeSeconds       float32
	asleepSeconds      float32
	awakeTimerSeconds  float32
	asleepTimerSeconds float32
}

func NewCreature(pos, size rl.Vector2, color color.RGBA, speed, maxHP float32, maxLvl int, awake, asleep float32) *Creature {
	return &Creature{
		pos:                pos,
		size:               size,
		color:              color,
		lvl:                1,
		maxLvl:             maxLvl,
		speed:              speed,
		hp:                 1,
		maxHP:              maxHP,
		awakeSeconds:       awake,
		asleepSeconds:      asleep,
		awakeTimerSeconds:  awake,
		asleepTimerSeconds: asleep,
	}
}

const creatureExpIncrement = 0.25
const creatureHpIncrement = 0.25
const creatureSpeedIncrement = 0.25

const creatureAsleepExpIncrement = 0.001
const creatureAsleepHpIncrement = 0.001

func (c *Creature) update(dt float32) {
	c.updateAwakenessStatus(dt)

	switch c.actionState {
	case ActionSearchingFood:
		c.searchFood(fruitSpawner.fruits, dt)
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

func (c *Creature) searchFood(fruits []Fruit, dt float32) {
	if len(fruits) == 0 {
		return
	}

	var closestFruitDistance float32 = findSquaredEuclideanDistance(c.pos, fruits[0].pos)
	var closestFruitPos rl.Vector2 = fruits[0].pos
	var closestFruitSize rl.Vector2 = fruits[0].size

	for i, fruit := range fruits {
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

	hasCreatureCollidedWithFruit, fruitIndex := c.checkFruitCollision(fruitSpawner)

	if hasCreatureCollidedWithFruit && fruitIndex >= 0 && fruitIndex < len(fruits) {
		fruitSpawner.despawnFruit(fruitIndex)
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
	c.pos = rl.Vector2MoveTowards(c.pos, target, c.speed*dt)

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
	c.exp += creatureExpIncrement
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
	c.exp += creatureAsleepExpIncrement * dt

	if c.hp < c.maxHP {
		c.hp = min(c.hp+creatureAsleepHpIncrement*dt, c.maxHP)
	}
}

func (c *Creature) handleLevelUp() {
	for c.lvl < c.maxLvl && c.exp >= calcExpForNextLvl(c.lvl) {
		c.exp -= calcExpForNextLvl(c.lvl)
		c.lvl++

		c.hp = min(c.hp+creatureHpIncrement, c.maxHP)
		c.speed += creatureSpeedIncrement
	}
}
