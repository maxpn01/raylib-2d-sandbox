package main

import (
	"image/color"
	"math"

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

func (m *Map) update(dt float32) {}

func (m *Map) draw() {
	colStart, rowStart, colEnd, rowEnd := m.visibleCellRange()
	for row := rowStart; row < rowEnd; row++ {
		for col := colStart; col < colEnd; col++ {
			cell := m.cells[row][col]
			x := float32(col) * cell.size.X
			y := float32(row) * cell.size.Y
			rl.DrawRectangle(int32(x), int32(y), int32(cell.size.X), int32(cell.size.Y), cell.color)
		}
	}
}

/* Implemented with AI, too mathy for me */
/*
	Raylib's Camera2D world→screen transform (with rotation=0) is:
	screen = (world - target) * zoom + offset
	Inverting it gives the world-space corners of the visible rectangle:
	worldMin = target - offset / zoom              // top-left
	worldMax = target + (screenSize - offset)/zoom // bottom-right
	Then divide by cell size → floor for the inclusive start index, ceil for the exclusive end index →
 	clamp to map bounds → loop only that sub-rectangle.
*/
func (m *Map) visibleCellRange() (colStart, rowStart, colEnd, rowEnd int) {
	if camera.Rotation != 0 {
		return 0, 0, int(m.sizeInCells.X), int(m.sizeInCells.Y)
	}

	zoom := camera.Zoom
	if zoom == 0 {
		zoom = 1
	}

	screenW := float32(window.width)
	screenH := float32(window.height)

	worldMinX := camera.Target.X - camera.Offset.X/zoom
	worldMinY := camera.Target.Y - camera.Offset.Y/zoom
	worldMaxX := camera.Target.X + (screenW-camera.Offset.X)/zoom
	worldMaxY := camera.Target.Y + (screenH-camera.Offset.Y)/zoom

	colStart = clampInt(int(math.Floor(float64(worldMinX/m.cellSize.X))), 0, int(m.sizeInCells.X))
	rowStart = clampInt(int(math.Floor(float64(worldMinY/m.cellSize.Y))), 0, int(m.sizeInCells.Y))
	colEnd = clampInt(int(math.Ceil(float64(worldMaxX/m.cellSize.X))), 0, int(m.sizeInCells.X))
	rowEnd = clampInt(int(math.Ceil(float64(worldMaxY/m.cellSize.Y))), 0, int(m.sizeInCells.Y))

	return colStart, rowStart, colEnd, rowEnd
}
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
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
