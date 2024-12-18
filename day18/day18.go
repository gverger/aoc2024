package day18

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/gverger/aoc2024/utils"
)

//go:embed *.txt
var f embed.FS

type Fall struct {
	Pos
	At int
}

type Input struct {
	Falls []Fall
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

	falls := make([]Fall, 0)
	step := 0
	for scanner.Scan() {
		step++
		coords := strings.Split(scanner.Text(), ",")
		falls = append(falls, Fall{
			Pos: Pos{
				X: utils.Must(strconv.Atoi(coords[0])),
				Y: utils.Must(strconv.Atoi(coords[1])),
			},
			At: step,
		})
	}
	utils.AssertNoErr(scanner.Err(), "reading input file")

	start := Pos{X: 0, Y: 0}
	end := Pos{X: 6, Y: 6}
	if filename == "input.txt" {
		end = Pos{X: 70, Y: 70}
	}

	return Input{
		Falls: falls,
		Start: start,
		End:   end,
	}
}

func gridFromFalls(falls []Fall, mx, my int) *utils.Grid[int] {

	g := utils.NewGrid[int](uint(mx+1), uint(my+1))
	for _, f := range falls {
		g.Set(f.X, f.Y, f.At)
	}

	return g
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution string
}

type GridUpdated struct {
	Grid *utils.Grid[int]
}

func Part1(ctx context.Context, input Input, limit int, callback func(ctx context.Context, obj any)) {
	g := gridFromFalls(input.Falls[:limit], input.End.X, input.End.Y)

	callback(ctx, GridUpdated{Grid: g})

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

	p, cost, ok := utils.Dijkstra(start, isDone, neighbors)
	utils.Assert(ok, "path not found")

	for _, pos := range p {
		g.Set(pos.X, pos.Y, -1)
	}

	callback(ctx, GridUpdated{Grid: g})
	callback(ctx, SolutionFound{Part: 1, Solution: strconv.Itoa(cost)})

}

func Part2(ctx context.Context, input Input, limit int, callback func(ctx context.Context, obj any)) {
	g := utils.NewGrid[int](uint(input.End.X+1), uint(input.End.Y+1))

	callback(ctx, GridUpdated{Grid: g})
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

	i := 0
	for i < len(input.Falls) {
		f := input.Falls[i]
		g.Set(f.X, f.Y, i+1)

		p, _, ok := utils.Dijkstra(start, isDone, neighbors)

		if !ok {
			gWithPath := g.Clone()

			g.Set(f.X, f.Y, 0)
			p, _, ok := utils.Dijkstra(start, isDone, neighbors)
			utils.Assert(ok, "rerun n-1")
			g.Set(f.X, f.Y, i+1)

			for _, pos := range p {
				gWithPath.Set(pos.X, pos.Y, -1)
			}
			gWithPath.Set(f.X, f.Y, -2)

			callback(ctx, GridUpdated{Grid: gWithPath})
			callback(ctx, SolutionFound{Part: 2, Solution: fmt.Sprintf("%d,%d", f.X, f.Y)})
			return
		}

		gWithPath := g.Clone()

		for _, pos := range p {
			gWithPath.Set(pos.X, pos.Y, -1)
		}

		i++
		for i < len(input.Falls) && gWithPath.At(input.Falls[i].X, input.Falls[i].Y) == 0 {
			g.Set(input.Falls[i].X, input.Falls[i].Y, i+1)
			gWithPath.Set(input.Falls[i].X, input.Falls[i].Y, i+1)
			i++
			callback(ctx, GridUpdated{Grid: gWithPath})
		}

		callback(ctx, GridUpdated{Grid: gWithPath})
		// callback(ctx, SolutionFound{Part: 2, Solution: fmt.Sprintf("%d,%d", f.X, f.Y)})

	}
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	filename := "input.txt"
	input := ReadInput(filename)
	callback(ctx, InputLoaded{Input: input})

	limit := 12
	if filename == "input.txt" {
		limit = 1024
	}

	// Part1(ctx, input, limit, callback)
	Part2(ctx, input, limit, callback)
}
