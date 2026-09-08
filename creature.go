package main

import (
	"image/color"

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
}

func NewCreature(pos, size rl.Vector2, color color.RGBA, speed, maxHP float32, maxLvl int) *Creature {
	return &Creature{
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

const creatureExpIncrement = 0.25
const creatureHpIncrement = 0.25
const creatureSpeedIncrement = 0.25

const creatureAsleepExpIncrement = 0.001
const creatureAsleepHpIncrement = 0.001

const creatureAwakeSeconds = 120
const creatureAsleepSeconds = 60

var awakeTimerSeconds float32 = creatureAwakeSeconds
var asleepTimerSeconds float32 = creatureAsleepSeconds

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

	c.moveCreature(closestFruitPos, closestFruitSize, dt)

	hasCreatureCollidedWithFruit, fruitIndex := c.checkCreatureFruitCollision(fruitSpawner)

	if hasCreatureCollidedWithFruit && fruitIndex >= 0 && fruitIndex < len(fruits) {
		fruitSpawner.despawnFruit(fruitIndex)
		c.actionState = ActionEating
	}
}

func (c *Creature) moveCreature(targetPos, targetSize rl.Vector2, dt float32) {
	move := rl.Vector2{}

	creatureCenterX := c.pos.X + c.size.X/2
	creatureCenterY := c.pos.Y + c.size.Y/2
	targetCenterX := targetPos.X + targetSize.X/2
	targetCenterY := targetPos.Y + targetSize.Y/2

	step := c.speed * dt

	diffX := creatureCenterX - targetCenterX
	diffY := creatureCenterY - targetCenterY

	if diffX > step {
		move.X = -1
	} else if diffX < -step {
		move.X = 1
	}
	if diffY > step {
		move.Y = -1
	} else if diffY < -step {
		move.Y = 1
	}

	// normalize the diagonal speed
	if move.X != 0 || move.Y != 0 {
		move = rl.Vector2Normalize(move)

		c.pos.X += move.X * c.speed * dt
		c.pos.Y += move.Y * c.speed * dt
	}

	clamp(gameMap.edgePosX.start, &c.pos.X, &c.size.X, gameMap.edgePosX.end)
	clamp(gameMap.edgePosX.start, &c.pos.Y, &c.size.Y, gameMap.edgePosY.end)
}

func (c *Creature) checkCreatureFruitCollision(fs *FruitSpawner) (bool, int) {
	for i := len(fs.fruits) - 1; i >= 0; i-- {
		hasCreatureCollidedWithFruit := checkCollisions(c.pos, c.size, fs.fruits[i].pos, fs.fruits[i].size)

		if hasCreatureCollidedWithFruit && c.hp < c.maxHP {
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
		if awakeTimerSeconds > 0 {
			awakeTimerSeconds -= 1 * dt
		} else {
			c.actionState = ActionSleeping
			asleepTimerSeconds = creatureAsleepSeconds
		}
	case ActionSleeping:
		if asleepTimerSeconds > 0 {
			asleepTimerSeconds -= 1 * dt
		} else {
			c.actionState = ActionSearchingFood
			awakeTimerSeconds = creatureAwakeSeconds
		}
	}
}

func (c *Creature) sleep(dt float32) {
	c.exp += creatureAsleepExpIncrement * dt

	if c.hp < c.maxHP {
		c.hp += creatureAsleepHpIncrement * dt
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
