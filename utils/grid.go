package utils

import (
	"fmt"
	"strings"
)

type Grid[T any] struct {
	Width  uint
	Height uint
	cells  []T

	MinX int
	MinY int
	MaxX int
	MaxY int
}

func NewGrid[T any](width, height uint) *Grid[T] {
	return NewGridEx[T](width, height, 0, 0)
}

func NewGridEx[T any](width, height uint, minX, minY int) *Grid[T] {
	return &Grid[T]{
		Width:  width,
		Height: height,
		cells:  make([]T, width*height),

		MinX: minX,
		MinY: minY,
		MaxX: minX + int(width) - 1,
		MaxY: minY + int(height) - 1,
	}
}

func (g Grid[T]) IsCoordValid(x, y int) bool {
	return x >= g.MinX && x <= g.MaxX && y >= g.MinY && y <= g.MaxY
}

func (g Grid[T]) At(x, y int) T {
	return g.cells[g.coordsToIdx(x, y)]
}

func (g *Grid[T]) Set(x, y int, value T) {
	g.cells[g.coordsToIdx(x, y)] = value
}

func (g Grid[T]) String() string {
	maxCellWidth := 0
	for _, v := range g.cells {
		maxCellWidth = max(maxCellWidth, len(fmt.Sprint(v)))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%dx%d] grid:\n", g.Width, g.Height))
	for i, v := range g.cells {
		x := i % int(g.Width)
		y := i / int(g.Height)

		if x == 0 && y > 0 {
			sb.WriteString("\n")
		}

		text := fmt.Sprint(v)
		sb.WriteString(text)
		sb.WriteString(strings.Repeat(" ", maxCellWidth-len(text)))
	}

	return sb.String()
}

func (g Grid[T]) coordsToIdx(x, y int) int {
	Assert(g.IsCoordValid(x, y))

	return (y-g.MinY)*int(g.Width) + (x - g.MinX)
}

type Cell[T any] struct {
	X     int
	Y     int
	Value T
}

type Direction struct {
	Dx int
	Dy int
}

func (d Direction) Apply(x, y int) (int, int) {
	return x + d.Dx, y + d.Dy
}

var DirUp = Direction{Dx: 0, Dy: -1}
var DirDown = Direction{Dx: 0, Dy: 1}
var DirRight = Direction{Dx: 1, Dy: 0}
var DirLeft = Direction{Dx: -1, Dy: 0}

var DirUR = Direction{Dx: 1, Dy: -1}
var DirUL = Direction{Dx: -1, Dy: -1}
var DirDR = Direction{Dx: 1, Dy: 1}
var DirDL = Direction{Dx: -1, Dy: 1}

type GridNeighbor[T any] struct {
	Dirs []Direction
}

func (n GridNeighbor[T]) NeighborCells(g Grid[T], x, y int) []Cell[T] {
	neighbors := make([]Cell[T], 0, len(n.Dirs))

	for _, d := range n.Dirs {
		nx, ny := x+d.Dx, y+d.Dy
		if !g.IsCoordValid(nx, ny) {
			continue
		}
		neighbors = append(neighbors, Cell[T]{
			X:     nx,
			Y:     ny,
			Value: g.At(nx, ny),
		})
	}

	return neighbors
}

func NewNeighbors4[T any]() GridNeighbor[T] {
	return GridNeighbor[T]{
		Dirs: []Direction{DirUp, DirRight, DirDown, DirLeft},
	}
}

func NewNeighbors8[T any]() GridNeighbor[T] {
	return GridNeighbor[T]{
		Dirs: []Direction{DirUp, DirUR, DirRight, DirDR, DirDown, DirDL, DirLeft, DirUL},
	}
}