package day6

import (
	"bufio"
	"context"
	"embed"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type CellType uint8

const (
	EmptyCell      CellType = iota
	ObstacleCell   CellType = iota
	GuardCell      CellType = iota
	FootPrintsCell CellType = iota
)

type Guard struct {
	X   int
	Y   int
	Dir Direction
}

func (g Guard) PositionAfterStep() (int, int) {
	return g.Dir.Apply(g.X, g.Y)
}

func (g *Guard) StepForward() {
	g.X, g.Y = g.PositionAfterStep()
}

func (g *Guard) TurnRight() {
	switch g.Dir {
	case DirUp:
		g.Dir = DirRight
	case DirRight:
		g.Dir = DirDown
	case DirDown:
		g.Dir = DirLeft
	case DirLeft:
		g.Dir = DirUp
	}
}

type Input struct {
	Grid  Grid[CellType]
	Guard *Guard
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
	guard := Guard{}
	for j, line := range lines {
		for i, r := range line {
			switch r {
			case '#':
				g.Set(i, j, ObstacleCell)
			case '.':
				g.Set(i, j, EmptyCell)
			case '^':
				guard.X = i
				guard.Y = j
				guard.Dir = DirUp
			}
		}
	}

	return Input{
		Grid:  *g,
		Guard: &guard,
	}
}

type InputLoaded struct {
	Input Input
}

type GuardMoved struct {
	OldX int
	OldY int
	X    int
	Y    int
	Dir  Direction
}

type GuardTurned struct {
	OldDir Direction
	Dir    Direction
}

type VisitedGrid struct {
	Grid  Grid[CellType]
	Guard Guard
}
type SolutionFound struct {
	Part     int
	Solution int
}

type GuardResult int

const (
	GuardOut     GuardResult = 1
	GuardInCycle GuardResult = 2
)

func run(ctx context.Context, input Input, callback func(ctx context.Context, obj any)) (Grid[bool], GuardResult) {
	g := input.Grid
	guard := Guard{
		X:   input.Guard.X,
		Y:   input.Guard.Y,
		Dir: input.Guard.Dir,
	}
	turnedHere := NewGrid[bool](g.Width, g.Height) // No need to store the direction here
	visited := NewGrid[bool](g.Width, g.Height)

	for {
		oldx, oldy := guard.X, guard.Y

		for g.IsCoordValid(guard.PositionAfterStep()) && g.At(guard.PositionAfterStep()) != ObstacleCell {
			visited.Set(guard.X, guard.Y, true)
			if g.At(guard.X, guard.Y) != FootPrintsCell {
				g.Set(guard.X, guard.Y, FootPrintsCell)
			}
			guard.StepForward()
		}

		callback(
			ctx,
			GuardMoved{
				OldX: oldx,
				OldY: oldy,
				X:    guard.X,
				Y:    guard.Y,
				Dir:  guard.Dir,
			},
		)

		if !g.IsCoordValid(guard.PositionAfterStep()) {
			visited.Set(guard.X, guard.Y, true)
			if g.At(guard.X, guard.Y) != FootPrintsCell {
				g.Set(guard.X, guard.Y, FootPrintsCell)
			}
			break
		}

		if turnedHere.At(guard.X, guard.Y) {
			return *visited, GuardInCycle
		}
		turnedHere.Set(guard.X, guard.Y, true)

		oldDir := guard.Dir
		for g.At(guard.Dir.Apply(guard.X, guard.Y)) == ObstacleCell {
			guard.TurnRight()
		}
		callback(ctx, GuardTurned{
			OldDir: oldDir,
			Dir:    guard.Dir,
		})
	}

	return *visited, GuardOut
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})

	visited, result1 := run(ctx, Input{Grid: *input.Grid.Clone(), Guard: input.Guard}, callback)
	if result1 == GuardInCycle {
		log.Fatal().Msg("In cycle??")
	}
	callback(ctx, SolutionFound{Part: 1, Solution: visited.Count(func(b Cell[bool]) bool { return b.Value })})

	cycles := 0
	for y := 0; y < int(input.Grid.Height); y++ {
		for x := 0; x < int(input.Grid.Width); x++ {
			if !visited.At(x, y) {
				continue
			}
			if x == input.Guard.X && y == input.Guard.Y {
				continue
			}

			g := input.Grid.Clone()
			g.Set(x, y, ObstacleCell)
			_, result := run(ctx, Input{Grid: *g, Guard: input.Guard}, callback)
			if result == GuardInCycle {
				cycles++
				callback(ctx, VisitedGrid{Grid: *g, Guard: *input.Guard})
			}

		}
	}
	callback(ctx, SolutionFound{Part: 2, Solution: cycles})

}
