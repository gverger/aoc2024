package day8

import (
	"bufio"
	"context"
	"embed"
	"fmt"

	. "github.com/gverger/aoc2024/utils"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type Antenna rune

func (a Antenna) String() string {
	return string(a)
}

type Input struct {
	Grid Grid[Antenna]
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
	g := NewGrid[Antenna](uint(len(lines[0])), uint(len(lines)))

	for j, line := range lines {
		for i, r := range line {
			g.Set(i, j, Antenna(r))
		}
	}

	return Input{Grid: *g}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

type Point struct {
	X int
	Y int
}

func antennas(g Grid[Antenna]) map[Antenna][]Point {
	antennas := make(map[Antenna][]Point)
	for y := g.MinY; y <= g.MaxY; y++ {
		for x := g.MinX; x <= g.MaxX; x++ {
			a := g.At(x, y)
			if a == '.' {
				continue
			}
			antennas[a] = append(antennas[a], Point{X: x, Y: y})
		}
	}

	return antennas
}

func antinodes1(g Grid[Antenna]) Grid[bool] {
	antinodes := NewGrid[bool](g.Width, g.Height)

	antennas := antennas(g)

	for _, points := range antennas {
		for i, p1 := range points {
			for _, p2 := range points[i+1:] {
				antinode1 := Point{X: 2*p1.X - p2.X, Y: 2*p1.Y - p2.Y}
				if antinodes.IsCoordValid(antinode1.X, antinode1.Y) {
					antinodes.Set(antinode1.X, antinode1.Y, true)
				}
				antinode2 := Point{X: 2*p2.X - p1.X, Y: 2*p2.Y - p1.Y}
				if antinodes.IsCoordValid(antinode2.X, antinode2.Y) {
					antinodes.Set(antinode2.X, antinode2.Y, true)
				}
			}
		}
	}

	return *antinodes
}

func antinodes2(g Grid[Antenna]) Grid[bool] {
	antinodes := NewGrid[bool](g.Width, g.Height)

	antennas := antennas(g)

	for _, points := range antennas {
		for i, p1 := range points {
			for _, p2 := range points[i+1:] {
				mult := 0
				for {
					node := Point{X: p1.X + mult*(p1.X-p2.X), Y: p1.Y + mult*(p1.Y-p2.Y)}
					if !antinodes.IsCoordValid(node.X, node.Y) {
						break
					}
					antinodes.Set(node.X, node.Y, true)
					mult++
				}

				mult = 0
				for {
					node := Point{X: p2.X + mult*(p2.X-p1.X), Y: p2.Y + mult*(p2.Y-p1.Y)}
					if !antinodes.IsCoordValid(node.X, node.Y) {
						break
					}
					antinodes.Set(node.X, node.Y, true)
					mult++
				}
			}
		}
	}

	return *antinodes
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})

	fmt.Println(input.Grid)
	antinodes1 := antinodes1(input.Grid)
	fmt.Println(antinodes1.Stringf(func(b bool) string {
		if b {
			return "@"
		} else {
			return " "
		}
	}))

	callback(ctx, SolutionFound{Part: 1, Solution: antinodes1.Count(func(b bool) bool { return b })})

	antinodes2 := antinodes2(input.Grid)
	fmt.Println(antinodes2.Stringf(func(b bool) string {
		if b {
			return "@"
		} else {
			return " "
		}
	}))

	callback(ctx, SolutionFound{Part: 2, Solution: antinodes2.Count(func(b bool) bool { return b })})
}
