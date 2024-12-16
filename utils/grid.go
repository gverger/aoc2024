package utils

import (
	"fmt"
	"iter"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

func (g *Grid[T]) SetAll(value T) {
	for i := range g.cells {
		g.cells[i] = value
	}
}

func (g *Grid[T]) SetAllFunc(init func() T) {
	for i := range g.cells {
		g.cells[i] = init()
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

func (g Grid[T]) Count(filter func(Cell[T]) bool) int {
	count := 0
	for c := range g.AllCells() {
		if filter(c) {
			count++
		}
	}
	return count
}

func MapGrid[T any, U any](g Grid[T], mapper func(t T) U) *Grid[U] {
	dst := NewGridEx[U](g.Width, g.Height, g.MinX, g.MinY)

	for i := 0; i < len(g.cells); i++ {
		dst.cells[i] = mapper(g.cells[i])
	}

	return dst
}

func (g Grid[T]) String() string {
	return g.Stringf(func(t T) string { return fmt.Sprint(t) })
}

func (g Grid[T]) Stringf(format func(T) string) string {
	maxCellWidth := 0
	for _, v := range g.cells {
		maxCellWidth = max(maxCellWidth, lipgloss.Width(format(v)))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%dx%d] grid:\n", g.Width, g.Height))
	for i, v := range g.cells {
		x := i % int(g.Width)
		y := i / int(g.Height)

		if x == 0 && y > 0 {
			sb.WriteString("\n")
		}

		text := format(v)
		sb.WriteString(text)
		sb.WriteString(strings.Repeat(" ", maxCellWidth-lipgloss.Width(text)))
	}

	return sb.String()
}

func (g Grid[T]) StringDots(filter func(T) bool) string {
	var sb strings.Builder

	for y := 0; y < int(g.Height)+4; y += 4 {
		for x := 0; x < int(g.Width)+2; x += 2 {
			current := Braille{}
			for dy := 0; dy < 4; dy++ {
				if y+dy >= int(g.Height) {
					break
				}
				for dx := 0; dx < 2; dx++ {
					if x+dx >= int(g.Width) {
						break
					}
					current[dx][dy] = filter(g.At(x+dx, y+dy))
				}
			}
			sb.WriteRune(current.Rune())
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (g Grid[T]) coordsToIdx(x, y int) int {
	Assert(g.IsCoordValid(x, y), "coord not valid: %d,%d", x, y)

	return (y-g.MinY)*int(g.Width) + (x - g.MinX)
}

func (g Grid[T]) AllCells() iter.Seq[Cell[T]] {
	return func(yield func(Cell[T]) bool) {
		for i, v := range g.cells {
			x := g.MinX + i%int(g.Width)
			y := g.MinY + i/int(g.Width)
			if !yield(Cell[T]{X: x, Y: y, Value: v}) {
				return
			}
		}
	}
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

func (g Grid[T]) Clone() *Grid[T] {

	cells := make([]T, len(g.cells))
	copy(cells, g.cells)

	return &Grid[T]{
		Width:  g.Width,
		Height: g.Height,
		cells:  cells,
		MinX:   g.MinX,
		MinY:   g.MinY,
		MaxX:   g.MaxX,
		MaxY:   g.MaxY,
	}
}

var Dirs4 = []Direction{DirUp, DirRight, DirDown, DirLeft}
var Dirs8 = []Direction{DirUp, DirUR, DirRight, DirDR, DirDown, DirDL, DirLeft, DirUL}

func NewNeighbors4[T any]() GridNeighbor[T] {
	return GridNeighbor[T]{
		Dirs: Dirs4,
	}
}

func NewNeighbors8[T any]() GridNeighbor[T] {
	return GridNeighbor[T]{
		Dirs: Dirs8,
	}
}
