package day20

import (
	"bufio"
	"context"
	"embed"

	"github.com/gverger/aoc2024/utils"
)

//go:embed *.txt
var f embed.FS

type CellType int

const (
	Empty CellType = iota
	Wall
	Footprints
)

type Input struct {
	Grid  *utils.Grid[CellType]
	Start Pos
	End   Pos
}

type Pos struct {
	X int
	Y int
}

func ReadInput(filename string) Input {
	file := utils.Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	utils.AssertNoErr(scanner.Err(), "reading input file")

	g := utils.NewGrid[CellType](uint(len(lines[0])), uint(len(lines)))
	var start, end Pos

	for j, l := range lines {
		for i, r := range l {
			switch r {
			case '#':
				g.Set(i, j, Wall)
			case '.':
				g.Set(i, j, Empty)
			case 'S':
				g.Set(i, j, Empty)
				start = Pos{i, j}
			case 'E':
				g.Set(i, j, Empty)
				end = Pos{i, j}
			}
		}
	}

	return Input{
		Grid:  g,
		Start: start,
		End:   end,
	}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	filename := "input.txt"
	input := ReadInput(filename)
	callback(ctx, InputLoaded{Input: input})

	g := input.Grid

	start := utils.WithCost[Pos, int]{Value: input.Start, Cost: 0}
	isDone := func(p Pos) bool { return p == input.End }

	neighbors := func(r utils.WithCost[Pos, int]) []utils.WithCost[Pos, int] {
		neighbors := make([]utils.WithCost[Pos, int], 0)

		for _, d := range utils.Dirs4 {
			x, y := d.Apply(r.Value.X, r.Value.Y)
			if g.IsCoordValid(x, y) && g.At(x, y) == 0 {
				neighbors = append(neighbors, utils.WithCost[Pos, int]{Value: Pos{X: x, Y: y}, Cost: r.Cost + 1})
			}
		}

		return neighbors
	}

	path, _, ok := utils.Dijkstra(start, isDone, neighbors)
	utils.Assert(ok, "path not found")

	step := utils.NewGrid[int](g.Width, g.Height)
	for i, p := range path {
		step.Set(p.X, p.Y, i)
	}

	parts := []struct {
		dist    int
		minGain int
	}{
		{dist: 2, minGain: 100},
		{dist: 20, minGain: 100},
	}

	for i, config := range parts {
		dist := config.dist
		minGain := config.minGain
		sum := 0

		for _, p := range path {
			for i := -dist; i <= dist; i++ {
				for j := -dist + utils.Abs(i); j <= dist-utils.Abs(i); j++ {
					dir := utils.Direction{Dx: i, Dy: j}
					x, y := dir.Apply(p.X, p.Y)
					if step.IsCoordValid(x, y) && step.At(x, y) >= step.At(p.X, p.Y)+utils.Abs(i)+utils.Abs(j)+minGain {
						g.Set(p.X, p.Y, Footprints)
						sum++
					}
				}
			}
		}
		callback(ctx, SolutionFound{Part: i + 1, Solution: sum})
	}

}
