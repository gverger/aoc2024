package day15

import (
	"bufio"
	"context"
	"embed"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed input.txt
//go:embed sample.txt
//go:embed small.txt
var f embed.FS

type CellType uint8

const (
	Empty  CellType = 0
	Wall   CellType = 1
	Box    CellType = 2
	Player CellType = 3
)

type Pos struct {
	X int
	Y int
}

type Input struct {
	Grid   *Grid[CellType]
	Player Pos
	Moves  []Direction
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)

	for scanner.Scan() {
		log.Debug().Str("line", scanner.Text()).Msg("reading")
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}
	AssertNoErr(scanner.Err(), "reading input file")

	g := NewGrid[CellType](uint(len(lines[0])), uint(len(lines)))

	var player Pos

	for j := 0; j < len(lines); j++ {
		for i := 0; i < len(lines[j]); i++ {
			value := Empty
			switch lines[j][i] {
			case '.':
				value = Empty
			case '#':
				value = Wall
			case 'O':
				value = Box
			case '@':
				player = Pos{X: i, Y: j}
				value = Player
			}
			g.Set(i, j, value)
		}
	}

	moves := make([]Direction, 0)
	for scanner.Scan() {
		log.Debug().Str("line", scanner.Text()).Msg("reading")
		for _, c := range scanner.Text() {
			switch c {
			case 'v':
				moves = append(moves, DirDown)
			case '^':
				moves = append(moves, DirUp)
			case '>':
				moves = append(moves, DirRight)
			case '<':
				moves = append(moves, DirLeft)
			}
		}
	}
	AssertNoErr(scanner.Err(), "reading input file")

	return Input{
		Grid:   g,
		Player: player,
		Moves:  moves,
	}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

type Moved struct {
	Grid Grid[CellType]
}

func move(input *Input, dir Direction) {
	g := input.Grid
	px, py := dir.Apply(input.Player.X, input.Player.Y)
	x, y := px, py

	Assert(g.IsCoordValid(x, y), "out of bounds")

	for g.At(x, y) == Box {
		x, y = dir.Apply(x, y)
	}

	if g.At(x, y) == Wall {
		return
	}

	Assert(g.At(x, y) == Empty, "should be an empty cell behind")

	g.Set(input.Player.X, input.Player.Y, Empty)
	input.Player = Pos{px, py}
	g.Set(px, py, Player)
	if x != px || y != py {
		g.Set(x, y, Box)
	}
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	log.DefaultLogger.SetLevel(log.InfoLevel)

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})
	callback(ctx, Moved{Grid: *input.Grid})

	for _, m := range input.Moves {
		move(&input, m)

		callback(ctx, Moved{Grid: *input.Grid})
	}

	score := 0
	for cell := range input.Grid.AllCells() {
		if cell.Value == Box {
			score += cell.Y*100 + cell.X
		}
	}

	callback(ctx, SolutionFound{Part: 1, Solution: score})
}
