package utils_test

import (
	"testing"

	"github.com/gverger/aoc2024/utils"
	"github.com/matryer/is"
)

func TestGrid(t *testing.T) {
	is := is.New(t)

	g := utils.NewGrid[int](10, 5)

	value := 0
	for y := 0; y < 5; y++ {
		for x := 0; x < 10; x++ {
			g.Set(x, y, value)
			value++
		}
	}

	is.True(g.IsCoordValid(7, 2))
	is.True(g.IsCoordValid(0, 4))
	is.True(!g.IsCoordValid(-1, 4))
	is.True(!g.IsCoordValid(7, 5))

	is.Equal(g.At(0, 0), 0)  // First value is 0
	is.Equal(g.At(0, 1), 10) // Second line value is 10
	is.Equal(g.At(9, 4), 49) // Last item is 10
}

func TestNeighbors(t *testing.T) {
	is := is.New(t)

	g := utils.NewGrid[bool](10, 10)

	g.Set(2, 9, true)
	g.Set(2, 7, true)
	g.Set(2, 8, true)

	n := utils.GridNeighbor[bool]{
		Dirs: []utils.Direction{
			{Dx: 0, Dy: -2},
			{Dx: 0, Dy: 2},
			{Dx: 2, Dy: 0},
			{Dx: -2, Dy: 0},
		},
	}

	is.Equal(n.NeighborCells(g, 2, 9), []utils.Cell[bool]{
		{X: 2, Y: 7, Value: true},
		{X: 4, Y: 9, Value: false},
		{X: 0, Y: 9, Value: false},
	})

	is.Equal(utils.NewNeighbors4[bool]().NeighborCells(g, 2, 9), []utils.Cell[bool]{
		{X: 2, Y: 8, Value: true},
		{X: 3, Y: 9, Value: false},
		{X: 1, Y: 9, Value: false},
	})

	is.Equal(utils.NewNeighbors8[bool]().NeighborCells(g, 2, 9), []utils.Cell[bool]{
		{X: 2, Y: 8, Value: true},
		{X: 3, Y: 8, Value: false},
		{X: 3, Y: 9, Value: false},
		{X: 1, Y: 9, Value: false},
		{X: 1, Y: 8, Value: false},
	})
}
