package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SnakeActionState int

const (
	ActionSearchingFood SnakeActionState = iota
	ActionEating
	ActionSleeping
)

// for now just one unit without the trailing units
type Snake struct {
	pos   rl.Vector2
	size  rl.Vector2
	color color.RGBA

	lvl    int
	maxLvl int
	exp    float32

	speed float32
	hp    float32
	maxHP float32

	targetIndex int

	actionState SnakeActionState
}

func NewSnake(pos, size rl.Vector2, color color.RGBA, speed, maxHP float32, maxLvl int) *Snake {
	return &Snake{
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

const snakeExpIncrement = 0.25
const snakeHpIncrement = 0.25
const snakeSpeedIncrement = 0.25

const snakeAsleepExpIncrement = 0.001
const snakeAsleepHpIncrement = 0.001

const snakeAwakeSeconds = 120
const snakeAsleepSeconds = 60

var awakeTimerSeconds float32 = snakeAwakeSeconds
var asleepTimerSeconds float32 = snakeAsleepSeconds

func (s *Snake) update(dt float32) {
	s.updateAwakenessStatus(dt)

	switch s.actionState {
	case ActionSearchingFood:
		s.searchFood(fruitSpawner.fruits, dt)
	case ActionEating:
		s.eatFood(fruitSpawner, s.targetIndex)
	case ActionSleeping:
		s.sleep(dt)
	}
}

func (s *Snake) draw() {
	rl.DrawRectangle(int32(s.pos.X), int32(s.pos.Y), int32(s.size.X), int32(s.size.Y), s.color)
}

func (s *Snake) searchFood(fruits []Fruit, dt float32) {
	if len(fruits) == 0 {
		return
	}

	var closestFruitDistance float32 = findSquaredEuclideanDistance(s.pos, fruits[0].pos)
	var closestFruitPos rl.Vector2 = fruits[0].pos
	var closestFruitSize rl.Vector2 = fruits[0].size

	for i, fruit := range fruits {
		if i == 0 {
			continue
		}

		distance := findSquaredEuclideanDistance(s.pos, fruit.pos)

		if distance <= closestFruitDistance {
			closestFruitDistance = distance
			closestFruitPos = fruit.pos
			closestFruitSize = fruit.size
		}
	}

	s.moveSnake(closestFruitPos, closestFruitSize, dt)

	hasSnakeCollidedWithFruit, fruitIndex := s.checkSnakeFruitCollision(fruitSpawner)

	if hasSnakeCollidedWithFruit && fruitIndex >= 0 && fruitIndex < len(fruits) {
		s.actionState = ActionEating
		s.targetIndex = fruitIndex
	}
}

func (s *Snake) moveSnake(targetPos, targetSize rl.Vector2, dt float32) {
	move := rl.Vector2{}

	snakeCenterX := s.pos.X + s.size.X/2
	snakeCenterY := s.pos.Y + s.size.Y/2
	targetCenterX := targetPos.X + targetSize.X/2
	targetCenterY := targetPos.Y + targetSize.Y/2

	step := s.speed * dt

	diffX := snakeCenterX - targetCenterX
	diffY := snakeCenterY - targetCenterY

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

		s.pos.X += move.X * s.speed * dt
		s.pos.Y += move.Y * s.speed * dt
	}

	clamp(0, &s.pos.X, &s.size.X, gameMap.size.X)
	clamp(0, &s.pos.Y, &s.size.Y, gameMap.size.Y)
}

func (s *Snake) checkSnakeFruitCollision(fs *FruitSpawner) (bool, int) {
	expForNextLvl := calcExpForNextLvl(s.lvl)

	for i := len(fs.fruits) - 1; i >= 0; i-- {
		hasSnakeCollidedWithFruit := checkCollisions(s.pos, s.size, fs.fruits[i].pos, fs.fruits[i].size)

		if hasSnakeCollidedWithFruit && s.hp < s.maxHP {
			if s.exp < expForNextLvl {
				s.exp += snakeExpIncrement
			}
			return true, i
		}
	}

	if s.exp == expForNextLvl {
		s.exp -= expForNextLvl
		s.lvl++

		if s.hp < s.maxHP {
			s.hp += playerHpIncrement
		}
		s.speed += snakeSpeedIncrement
	}

	return false, -1
}

func (s *Snake) eatFood(fs *FruitSpawner, fruitIndex int) {
	s.exp += snakeExpIncrement
	fs.despawnFruit(fruitIndex)
	s.actionState = ActionSearchingFood
	s.targetIndex = -1
}

func (s *Snake) updateAwakenessStatus(dt float32) {
	switch s.actionState {
	case ActionSearchingFood:
		fallthrough
	case ActionEating:
		if awakeTimerSeconds > 0 {
			awakeTimerSeconds -= 1 * dt
		} else {
			s.actionState = ActionSleeping
			asleepTimerSeconds = snakeAsleepSeconds
		}
	case ActionSleeping:
		if asleepTimerSeconds > 0 {
			asleepTimerSeconds -= 1 * dt
		} else {
			s.actionState = ActionSearchingFood
			awakeTimerSeconds = snakeAwakeSeconds
		}
	}
}

// initial crude version (might extend later)
func (s *Snake) sleep(dt float32) {
	s.exp += snakeAsleepExpIncrement * dt

	if s.hp < s.maxHP {
		s.hp += snakeAsleepHpIncrement * dt
	}
}
