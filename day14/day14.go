package day14

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"regexp"
	"strconv"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type Pos struct {
	X int
	Y int
}

type Robot struct {
	Position  Pos
	Direction Direction
}

func (r Robot) PositionInTurn(n int) Pos {
	d := Direction{
		Dx: r.Direction.Dx * n,
		Dy: r.Direction.Dy * n,
	}
	x, y := d.Apply(r.Position.X, r.Position.Y)
	return Pos{
		X: x,
		Y: y,
	}
}

type Input struct {
	Robots []Robot
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	r := regexp.MustCompile(`-?\d+`)

	robots := make([]Robot, 0)
	for scanner.Scan() {
		numbers := r.FindAllString(scanner.Text(), -1)
		Assert(len(numbers) == 4, "need 4 numbers per line: %q -> %v", scanner.Text(), numbers)

		r := Robot{
			Position: Pos{
				X: Must(strconv.Atoi(numbers[0])),
				Y: Must(strconv.Atoi(numbers[1])),
			},
			Direction: Direction{
				Dx: Must(strconv.Atoi(numbers[2])),
				Dy: Must(strconv.Atoi(numbers[3])),
			},
		}

		robots = append(robots, r)
	}
	AssertNoErr(scanner.Err(), "reading input file")

	return Input{
		Robots: robots,
	}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

func PositionsAtTurn(input Input, turn int) []Pos {
	positions := make([]Pos, 0, len(input.Robots))
	for _, r := range input.Robots {
		p := r.PositionInTurn(turn)
		positions = append(positions, p)
	}

	return positions
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	realInput := true

	filename := "sample.txt"
	width := 11
	height := 7

	if realInput {
		filename = "input.txt"
		width = 101
		height = 103
	}

	input := ReadInput(filename)
	// input := ReadInput("sample.txt")
	// width := 11
	// height := 7

	callback(ctx, InputLoaded{Input: input})


	g := NewGrid[int](uint(width), uint(height))
	positions := PositionsAtTurn(input, 100)

	quadrants := [4]int{0, 0, 0, 0}
	for _, p := range positions {
		x := Mod(p.X, width)
		y := Mod(p.Y, height)
		g.Set(x, y, g.At(x, y)+1)

		if x == width/2 || y == height/2 {
			continue
		}

		quadrant := 0
		if x > width/2 {
			quadrant += 1
		}

		if y > height/2 {
			quadrant += 2
		}

		quadrants[quadrant]++
	}

	fmt.Println(g.Stringf(func(i int) string {
		if i == 0 {
			return "."
		} else {
			return strconv.Itoa(i)
		}
	}))

	solution1 := 1
	for _, v := range quadrants {
		solution1 *= v
	}

	log.Info().Interface("quadrants", quadrants).Msg("solution")
	callback(ctx, SolutionFound{Part: 1, Solution: solution1})
}
