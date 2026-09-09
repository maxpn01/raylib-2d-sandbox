package game

import (
	"image/color"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AIActionState int

const (
	ActionSearchingFood AIActionState = iota
	ActionEating
	ActionSleeping
)

type AIStats struct {
	lvl int
	exp float32

	hp    float32
	speed float32

	maxLvl int
	maxHP  float32
}

type AIGrowthStats struct {
	HPPerLevel    float32
	SpeedPerLevel float32

	ExpAsleep float32
	HPAsleep  float32
}

type AI struct {
	pos   rl.Vector2
	size  rl.Vector2
	color color.RGBA

	stats AIStats

	growth AIGrowthStats

	actionState AIActionState

	awakeSeconds       float32
	asleepSeconds      float32
	awakeTimerSeconds  float32
	asleepTimerSeconds float32

	nutritionalFactor float32
}

func NewAI(
	pos,
	size rl.Vector2,
	color color.RGBA,
	stats AIStats,
	growth AIGrowthStats,
	nutritionalFactor float32,
	awake, asleep float32,
) *AI {
	return &AI{
		pos:                pos,
		size:               size,
		color:              color,
		stats:              stats,
		growth:             growth,
		nutritionalFactor:  nutritionalFactor,
		awakeSeconds:       awake,
		asleepSeconds:      asleep,
		awakeTimerSeconds:  awake,
		asleepTimerSeconds: asleep,
	}
}

func (ai *AI) update(w *World, dt float32) {
	ai.updateAwakenessStatus(dt)

	switch ai.actionState {
	case ActionSearchingFood:
		ai.seekFood(w.fruitSpawner, w.player, w.gameMap, dt)
	case ActionEating:
		ai.actionState = ActionSearchingFood
	case ActionSleeping:
		ai.sleep(dt)
	}

	ai.handleLevelUp()
}

func (ai *AI) draw() {
	rl.DrawRectangleV(ai.pos, ai.size, ai.color)
}

func (ai *AI) seekFood(fs *FruitSpawner, p *Player, gameMap *Map, dt float32) {
	targetPos, targetSize, targetFound := ai.chooseFoodTarget(fs, p)

	if !targetFound {
		return
	}

	ai.move(targetPos, targetSize, gameMap, dt)

	ai.handleFoodInteractions(fs, p)
}

func (ai *AI) chooseFoodTarget(fs *FruitSpawner, p *Player) (targetPos rl.Vector2, targetSize rl.Vector2, targetFound bool) {
	closestFruit, fruitFound := ai.findClosestFruit(fs)

	if ai.canEatPlayer(p) {
		if !fruitFound {
			return p.pos, p.size, true
		}

		playerDistance := ai.distanceTo(p.pos, p.size)
		closestFruitDistance := ai.distanceTo(closestFruit.pos, closestFruit.size)
		playerIsCloser := playerDistance <= closestFruitDistance
		if playerIsCloser {
			return p.pos, p.size, true
		}
	}

	if fruitFound {
		return closestFruit.pos, closestFruit.size, true
	}

	return rl.Vector2{}, rl.Vector2{}, false
}

func (ai *AI) findClosestFruit(fs *FruitSpawner) (Fruit, bool) {
	if len(fs.fruits) == 0 {
		return Fruit{}, false
	}

	closestFruit := fs.fruits[0]
	closestFruitDistance := ai.distanceTo(closestFruit.pos, closestFruit.size)

	for _, fruit := range fs.fruits[1:] {
		distance := ai.distanceTo(fruit.pos, fruit.size)
		if distance <= closestFruitDistance {
			closestFruitDistance = distance
			closestFruit = fruit
		}
	}

	return closestFruit, true
}

func (ai *AI) distanceTo(pos, size rl.Vector2) float32 {
	return findSquaredEuclideanDistance(rectCenter(ai.pos, ai.size), rectCenter(pos, size))
}

func (ai *AI) canEatPlayer(p *Player) bool {
	return ai.stats.hp > p.stats.hp
}

func (ai *AI) move(targetPos, targetSize rl.Vector2, gameMap *Map, dt float32) {
	// Align our center with the target's center, accounting for their different sizes.
	target := rl.NewVector2(
		targetPos.X+targetSize.X/2-ai.size.X/2,
		targetPos.Y+targetSize.Y/2-ai.size.Y/2,
	)
	// Travel at most speed * dt, stopping at the target rather than overshooting it.
	ai.pos = rl.Vector2MoveTowards(ai.pos, target, ai.stats.speed*dt)

	clamp(gameMap.edgePosX.start, &ai.pos.X, &ai.size.X, gameMap.edgePosX.end)
	clamp(gameMap.edgePosY.start, &ai.pos.Y, &ai.size.Y, gameMap.edgePosY.end)
}

func (ai *AI) handleFoodInteractions(fs *FruitSpawner, p *Player) {
	hasAICollidedWithFruit, fruitIndex := ai.checkFruitCollision(fs)

	if hasAICollidedWithFruit && fruitIndex >= 0 && fruitIndex < len(fs.fruits) {
		ai.eatFood(fs.fruits[fruitIndex].nutritionalValue)
		fs.despawnFruit(fruitIndex)
	}

	hasAICollidedWithPlayer := checkCollisions(ai.pos, ai.size, p.pos, p.size)

	if hasAICollidedWithPlayer {
		if ai.canEatPlayer(p) {
			ai.eatFood(p.nutritionalValue())
			p.markPlayerDeath()
		} else {
			// kill ai
		}
	}
}

func (ai *AI) checkFruitCollision(fs *FruitSpawner) (bool, int) {
	for i, v := range slices.Backward(fs.fruits) {
		hasAICollidedWithFruit := checkCollisions(ai.pos, ai.size, v.pos, v.size)

		if hasAICollidedWithFruit {
			return true, i
		}
	}

	return false, -1
}

func (ai *AI) nutritionalValue() float32 {
	return ai.nutritionalFactor * float32(ai.stats.lvl)
}

func (ai *AI) eatFood(nutrition float32) {
	ai.stats.exp += nutrition
	ai.actionState = ActionEating
}

func (ai *AI) updateAwakenessStatus(dt float32) {
	switch ai.actionState {
	case ActionSearchingFood, ActionEating:
		ai.awakeTimerSeconds -= dt
		if ai.awakeTimerSeconds <= 0 {
			ai.actionState = ActionSleeping
			ai.asleepTimerSeconds = ai.asleepSeconds
		}
	case ActionSleeping:
		ai.asleepTimerSeconds -= dt
		if ai.asleepTimerSeconds <= 0 {
			ai.actionState = ActionSearchingFood
			ai.awakeTimerSeconds = ai.awakeSeconds
		}
	}
}

func (ai *AI) sleep(dt float32) {
	ai.stats.exp += ai.growth.ExpAsleep * dt

	if ai.stats.hp < ai.stats.maxHP {
		ai.stats.hp = min(ai.stats.hp+ai.growth.HPAsleep*dt, ai.stats.maxHP)
	}
}

func (ai *AI) handleLevelUp() {
	for ai.stats.lvl < ai.stats.maxLvl && ai.stats.exp >= calcExpForNextLvl(ai.stats.lvl) {
		ai.stats.exp -= calcExpForNextLvl(ai.stats.lvl)
		ai.stats.lvl++

		ai.stats.hp = min(ai.stats.hp+ai.growth.HPPerLevel, ai.stats.maxHP)
		ai.stats.speed += ai.growth.SpeedPerLevel
	}
}
