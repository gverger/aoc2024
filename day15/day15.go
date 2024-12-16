package day15

import (
	"bufio"
	"context"
	"embed"
	"sort"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed input.txt
//go:embed sample.txt
//go:embed small.txt
var f embed.FS

type CellType uint8

const (
	Empty CellType = 1 << iota
	Wall
	Box
	Player
	Left
	Right

	Highlighted
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

func createPart2(input Input) Input {
	g := NewGrid[CellType](2*input.Grid.Width, input.Grid.Height)
	for j := 0; j < int(input.Grid.Height); j++ {
		for i := 0; i < int(input.Grid.Width); i++ {
			cell := input.Grid.At(i, j)
			switch cell {
			case Box:
				g.Set(2*i, j, Box|Left)
				g.Set(2*i+1, j, Box|Right)
			case Player:
				g.Set(2*i, j, Player)
				g.Set(2*i+1, j, Empty)
			default:
				g.Set(2*i, j, cell)
				g.Set(2*i+1, j, cell)
			}
		}
	}
	return Input{
		Grid:   g,
		Player: Pos{X: 2 * input.Player.X, Y: input.Player.Y},
		Moves:  input.Moves,
	}

}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)

	for scanner.Scan() {
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

func pushedBoxes(g Grid[CellType], x, y int, dir Direction) (Set[Pos], bool) {
	nx, ny := dir.Apply(x, y)

	positions := NewSet[Pos]()

	cell := g.At(nx, ny)

	if cell&Wall != 0 {
		return nil, false
	}

	if cell&Empty != 0 {
		return NewSet[Pos](), true
	}

	if cell&Right != 0 {
		nx--
	}

	positions.Add(Pos{nx, ny})
	positions.Add(Pos{nx + 1, ny})

	if dir.Dy != 0 {
		if pos, ok := pushedBoxes(g, nx, ny, dir); ok {
			positions.Union(pos)
		} else {
			return nil, false
		}
		if pos, ok := pushedBoxes(g, nx+1, ny, dir); ok {
			positions.Union(pos)
		} else {
			return nil, false
		}
	} else if dir.Dx == 1 {
		if pos, ok := pushedBoxes(g, nx+1, ny, dir); ok {
			positions.Union(pos)
		} else {
			return nil, false
		}
	} else {
		if pos, ok := pushedBoxes(g, nx, ny, dir); ok {
			positions.Union(pos)
		} else {
			return nil, false
		}
	}

	return positions, true
}

func move2(input *Input, dir Direction) {
	g := input.Grid
	px, py := dir.Apply(input.Player.X, input.Player.Y)

	Assert(g.IsCoordValid(px, py), "out of bounds")

	pushed, ok := pushedBoxes(*g, input.Player.X, input.Player.Y, dir)
	if !ok {
		return
	}

	pos := make([]Pos, 0, len(pushed))
	for p := range pushed {
		pos = append(pos, p)
	}

	sort.Slice(pos, func(i, j int) bool {
		di := pos[i].X*dir.Dx + pos[i].Y*dir.Dy
		dj := pos[j].X*dir.Dx + pos[j].Y*dir.Dy
		return di > dj
	})

	for _, p := range pos {
		x, y := dir.Apply(p.X, p.Y)
		g.Set(x, y, g.At(p.X, p.Y)|Highlighted)
		g.Set(p.X, p.Y, Empty)
	}
	g.Set(input.Player.X, input.Player.Y, Empty)
	input.Player = Pos{px, py}
	g.Set(px, py, Player)
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	log.DefaultLogger.SetLevel(log.InfoLevel)

	input := ReadInput("input.txt")
	// callback(ctx, InputLoaded{Input: input})
	//
	// for _, m := range input.Moves {
	// 	move(&input, m)
	//
	// 	callback(ctx, Moved{Grid: *input.Grid})
	// }
	//
	// score := 0
	// for cell := range input.Grid.AllCells() {
	// 	if cell.Value == Box {
	// 		score += cell.Y*100 + cell.X
	// 	}
	// }
	//
	// callback(ctx, SolutionFound{Part: 1, Solution: score})

	input2 := createPart2(input)
	callback(ctx, InputLoaded{Input: input2})
	for _, m := range input2.Moves {
		move2(&input2, m)
		// callback(ctx, Moved{Grid: *input2.Grid})

		for c := range input2.Grid.AllCells() {
			nohighlight := c.Value & ^Highlighted
			input2.Grid.Set(c.X, c.Y, CellType(nohighlight))
		}
	}
	score := 0
	for cell := range input2.Grid.AllCells() {
		if cell.Value == Box|Left {
			score += cell.Y*100 + cell.X
		}
	}
	callback(ctx, SolutionFound{Part: 2, Solution: score})
}
