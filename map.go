package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cell struct {
	size  rl.Vector2
	color color.RGBA
}

type Map struct {
	sizeInCells rl.Vector2
	cellSize    rl.Vector2
	size        rl.Vector2
	cells       [][]Cell
}

func NewMap(size, cellSize rl.Vector2, bgColor, accentColor color.RGBA) *Map {
	cells := make([][]Cell, int(size.Y))
	for row := range cells {
		cells[row] = make([]Cell, int(size.X))
	}

	for row := 0; row < int(size.Y); row++ {
		for col := 0; col < int(size.X); col++ {
			c := bgColor
			if isEdgeCell(row, col, int(size.X), int(size.Y)) || isCenterCell(row, col, int(size.X), int(size.Y)) {
				c = accentColor
			}
			cells[row][col] = Cell{
				size:  cellSize,
				color: c,
			}
		}
	}

	return &Map{
		sizeInCells: size,
		cellSize:    cells[0][0].size,
		size:        rl.NewVector2(size.X*cells[0][0].size.X, size.Y*cells[0][0].size.Y),
		cells:       cells,
	}
}

func isEdgeCell(row, col, width, height int) bool {
	return row == 0 || row == height-1 || col == 0 || col == width-1
}

func isCenterCell(row, col, width, height int) bool {
	centerRow := height / 2
	centerCol := width / 2
	radius := 2
	return row >= centerRow-radius && row <= centerRow+radius &&
		col >= centerCol-radius && col <= centerCol+radius
}

func (m *Map) update(dt float32) {}

func (m *Map) draw() {
	for row := 0; row < int(m.sizeInCells.Y); row++ {
		for col := 0; col < int(m.sizeInCells.X); col++ {
			cell := m.cells[row][col]
			x := float32(col) * cell.size.X
			y := float32(row) * cell.size.Y
			rl.DrawRectangle(int32(x), int32(y), int32(cell.size.X), int32(cell.size.Y), cell.color)
		}
	}
}
