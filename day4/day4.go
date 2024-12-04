package day4

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

type Input struct {
	Grid Grid[string]
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	if len(lines) == 0 {
		log.Fatal().Str("file", filename).Msg("No line in the file")
	}

	g := NewGrid[string](uint(len(lines[0])), uint(len(lines)))

	for i, l := range lines {
		for j, r := range l {
			g.Set(i, j, string(r))
		}
	}

	return Input{
		Grid: g,
	}
}

func isXmas(g Grid[string], x, y int, d Direction) bool {
	word := "XMAS"

	for i := 0; i < len(word); i++ {
		nx := x + d.Dx*i
		ny := y + d.Dy*i
		if !g.IsCoordValid(nx, ny) || g.At(nx, ny) != string(word[i]) {
			return false
		}
	}
	return true
}

func isMaxInX(g Grid[string], x, y int) bool {
	if g.At(x, y) != "A" {
		return false
	}

	diag1 := g.At(DirUL.Apply(x, y)) + g.At(DirDR.Apply(x, y))
	if diag1 != "MS" && diag1 != "SM" {
		return false
	}

	diag2 := g.At(DirDL.Apply(x, y)) + g.At(DirUR.Apply(x, y))
	if diag2 != "MS" && diag2 != "SM" {
		return false
	}
	return true
}

func backOneCell(d Direction, x, y int) (int, int) {
	return x - d.Dx, y - d.Dy
}

type InputLoaded struct {
	Input Input
}

type XMasFound struct {
	X   int
	Y   int
	Dir Direction
}

type MasInXFound struct {
	X int
	Y int
}

type SolutionFound struct {
	Part     int
	Solution int
}

func Run(ctx context.Context, listener chan<- Event) {
	notify := func(event any) {
		log.Debug().Msg("Sending event")
		select {
		case listener <- event:
		case <-ctx.Done():
			log.Info().Msg("Terminated")
			return
		}
		log.Debug().Msg("event sent")
	}

	input := ReadInput("sample.txt")
	notify(InputLoaded{Input: input})

	// Part 1
	neighbour := NewNeighbors8[string]()
	nb := 0

	for _, d := range neighbour.Dirs {
		for y := 0; y < int(input.Grid.Height); y++ {
			for x := 0; x < int(input.Grid.Width); x++ {
				if isXmas(input.Grid, x, y, d) {
					notify(XMasFound{X: x, Y: y, Dir: d})
					nb++
				}
			}
		}
	}
	notify(SolutionFound{Part: 1, Solution: nb})
	log.Info().Int("nb of xmas", nb).Msg("Part 1")

	// Part 2
	part2Nb := 0

	for y := 1; y < int(input.Grid.Height)-1; y++ {
		for x := 1; x < int(input.Grid.Width)-1; x++ {
			if isMaxInX(input.Grid, x, y) {
				notify(MasInXFound{X: x, Y: y})
				part2Nb++
			}
		}
	}

	notify(SolutionFound{Part: 2, Solution: part2Nb})
	log.Info().Int("nb of mas in x", part2Nb).Msg("Part 2")
}
