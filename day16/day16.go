package day16

import (
	"bufio"
	"context"
	"embed"
	"slices"

	. "github.com/gverger/aoc2024/utils"
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
	Grid     *Grid[CellType]
	Start    Pos
	StartDir Direction
	End      Pos
}

type Pos struct {
	X int
	Y int
}

type Reindeer struct {
	Pos Pos
	Dir Direction
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

	g := NewGrid[CellType](uint(len(lines[0])), uint(len(lines)))
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
		Grid:     g,
		Start:    start,
		StartDir: DirRight,
		End:      end,
	}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
	Path     []Reindeer
	Grid     Grid[CellType]
}

func countParents(parents map[Reindeer]Set[Reindeer], current Reindeer, counted Set[Reindeer]) int {
	if counted.Exists(current) {
		return 0
	}
	if len(parents[current]) == 0 {
		return 1
	}
	counted.Add(current)
	sum := 1
	for p := range parents[current] {
		sum += countParents(parents, p, counted)
	}
	return sum
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	filename := "input.txt"
	input := ReadInput(filename)

	callback(ctx, InputLoaded{Input: input})

	start := WithCost[Reindeer, int]{Value: Reindeer{Pos: input.Start, Dir: input.StartDir}, Cost: 0}
	isDone := func(p Reindeer) bool { return p.Pos == input.End }
	neighbors := func(r WithCost[Reindeer, int]) []WithCost[Reindeer, int] {
		turned := []Direction{
			{Dx: r.Value.Dir.Dy, Dy: r.Value.Dir.Dx},
			{Dx: -r.Value.Dir.Dy, Dy: -r.Value.Dir.Dx},
		}

		neighbors := make([]WithCost[Reindeer, int], 0)
		x, y := r.Value.Dir.Apply(r.Value.Pos.X, r.Value.Pos.Y)
		if input.Grid.IsCoordValid(x, y) && input.Grid.At(x, y) != Wall {
			reindeer := Reindeer{Pos: Pos{X: x, Y: y}, Dir: r.Value.Dir}
			neighbors = append(neighbors, WithCost[Reindeer, int]{Value: reindeer, Cost: r.Cost + 1})
		}
		neighbors = append(neighbors, WithCost[Reindeer, int]{Value: Reindeer{Pos: r.Value.Pos, Dir: turned[0]}, Cost: r.Cost + 1000})
		neighbors = append(neighbors, WithCost[Reindeer, int]{Value: Reindeer{Pos: r.Value.Pos, Dir: turned[1]}, Cost: r.Cost + 1000})

		return neighbors
	}

	p, cost, ok := Dijkstra(start, isDone, neighbors)

	Assert(ok, "dijkstra")
	g := input.Grid.Clone()
	for _, r := range p {
		g.Set(r.Pos.X, r.Pos.Y, Footprints)
	}

	callback(ctx, SolutionFound{Part: 1, Solution: cost, Grid: *g})

	parents, costs := DijkstraAll(start, neighbors)
	parentsIn := NewSet[Reindeer]()

	possibleEnds := []Reindeer{
		{Pos: input.End, Dir: DirUp},
		{Pos: input.End, Dir: DirDown},
		{Pos: input.End, Dir: DirLeft},
		{Pos: input.End, Dir: DirRight},
	}

	endCost := slices.Min(MapTo(possibleEnds, func(r Reindeer) int { return costs[r] }))

	for _, r := range possibleEnds {
		if costs[r] == endCost {
			countParents(parents, r, parentsIn)
		}
	}

	positions := NewSet[Pos]()
	for r := range parentsIn {
		positions.Add(r.Pos)
	}

	g = input.Grid.Clone()
	for p := range positions {
		g.Set(p.X, p.Y, Footprints)
	}

	callback(ctx, SolutionFound{Part: 2, Solution: len(positions), Grid: *g})
}
