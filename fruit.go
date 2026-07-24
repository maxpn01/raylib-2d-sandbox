package main

import (
	"image/color"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Fruit struct {
	pos   rl.Vector2
	size  rl.Vector2
	color color.RGBA
}

type FruitSpawner struct {
	fruits             []Fruit
	fruitSize          rl.Vector2
	fruitColor         color.RGBA
	fruitSpawnTimer    float32
	fruitSpawnInterval int
	maxFruits          int
}

var fruitSpawnIntervalMaxSeconds int

func NewFruitSpawner(size rl.Vector2, color color.RGBA, spawnIntervalMaxSeconds int, maxFruits int) *FruitSpawner {
	fruitSpawnIntervalMaxSeconds = spawnIntervalMaxSeconds

	return &FruitSpawner{
		fruits:             []Fruit{},
		fruitSize:          size,
		fruitColor:         color,
		fruitSpawnTimer:    0,
		fruitSpawnInterval: rand.Intn(spawnIntervalMaxSeconds),
		maxFruits:          maxFruits,
	}
}

func (fs *FruitSpawner) update(dt float32) {
	spawnFruit(fs, dt)
}

func (fs *FruitSpawner) draw() {
	for _, fruit := range fs.fruits {
		rl.DrawRectangle(int32(fruit.pos.X), int32(fruit.pos.Y), int32(fruit.size.X), int32(fruit.size.Y), fruit.color)
	}
}

func spawnFruit(fs *FruitSpawner, dt float32) {
	fs.fruitSpawnTimer += dt

	if fs.fruitSpawnTimer >= float32(fs.fruitSpawnInterval) && len(fs.fruits) < fs.maxFruits {
		fruitRandX := rand.Intn(int(gameMap.size.X) - int(fs.fruitSize.X))
		fruitRandY := rand.Intn(int(gameMap.size.Y) - int(fs.fruitSize.Y))

		fs.fruits = append(fs.fruits, Fruit{
			pos:   rl.NewVector2(float32(fruitRandX), float32(fruitRandY)),
			size:  rl.NewVector2(fs.fruitSize.X, fs.fruitSize.X),
			color: fs.fruitColor,
		})

		fs.fruitSpawnTimer = 0
		fs.fruitSpawnInterval = 1 + rand.Intn(fruitSpawnIntervalMaxSeconds)
	}
}

func (fs *FruitSpawner) despawnFruit(fruitIndex int) {
	// when order matters, but much slower as it shifts the elements to the left
	// fs.fruits = append(fs.fruits[:fruitIndex], fs.fruits[fruitIndex+1:]...)
	fs.fruits[fruitIndex] = fs.fruits[len(fs.fruits)-1]
	fs.fruits = fs.fruits[:len(fs.fruits)-1]
}
