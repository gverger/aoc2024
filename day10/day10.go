package day10

import (
	"bufio"
	"context"
	"embed"
	"strconv"

	. "github.com/gverger/aoc2024/utils"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type Input struct {
	Grid *Grid[int]
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	AssertNoErr(scanner.Err(), "reading input file")

	g := NewGrid[int](uint(len(lines[0])), uint(len(lines)))
	for j, line := range lines {
		for i, r := range line {
			g.Set(i, j, Must(strconv.Atoi(string(r))))
		}
	}

	return Input{
		Grid: g,
	}
}

func access(g *Grid[int]) *Grid[Set[int]] {
	accessible := NewGrid[Set[int]](g.Width, g.Height)
	accessible.SetAllFunc(func() Set[int] { return NewSet[int]() })
	hillID := 1
	for c := range g.AllCells() {
		if c.Value == 9 {
			accessible.At(c.X, c.Y).Add(hillID)
			hillID++
		}
	}

	dirs := NewNeighbors4[Set[int]]()

	for i := 9; i > 0; i-- {
		for c := range g.AllCells() {
			if c.Value == i {
				for _, n := range dirs.NeighborCells(*accessible, c.X, c.Y) {
					if g.At(n.X, n.Y) == i-1 {
						n.Value.Union(accessible.At(c.X, c.Y))
					}
				}
			}
		}
	}

	return accessible
}

func nbPaths(g *Grid[int]) *Grid[int] {
	paths := NewGrid[int](g.Width, g.Height)
	for c := range g.AllCells() {
		if c.Value == 9 {
			paths.Set(c.X, c.Y, 1)
		}
	}

	dirs := NewNeighbors4[int]()

	for i := 9; i > 0; i-- {
		for c := range g.AllCells() {
			if c.Value == i {
				for _, n := range dirs.NeighborCells(*paths, c.X, c.Y) {
					if g.At(n.X, n.Y) == i-1 {
						paths.Set(n.X, n.Y, n.Value+paths.At(c.X, c.Y))
					}
				}
			}
		}
	}

	return paths
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})
	accessible := access(input.Grid)

	sum := 0
	for c := range input.Grid.AllCells() {
		if c.Value == 0 {
			sum += len(accessible.At(c.X, c.Y))
		}
	}

	callback(ctx, SolutionFound{Part: 1, Solution: sum})

	paths := nbPaths(input.Grid)
	sum = 0

	for c := range input.Grid.AllCells() {
		if c.Value == 0 {
			sum += paths.At(c.X, c.Y)
		}
	}
	callback(ctx, SolutionFound{Part: 2, Solution: sum})
}
