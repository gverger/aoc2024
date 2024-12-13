package day12

import (
	"bufio"
	"context"
	"embed"
	"fmt"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type Input struct {
	Farm Grid[rune]
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

	g := NewGrid[rune](uint(len(lines[0])), uint(len(lines)))

	for j, line := range lines {
		for i := 0; i < len(line); i++ {
			g.Set(i, j, rune(line[i]))
		}
	}

	return Input{Farm: *g}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

type Pos struct {
	X int
	Y int
}

func Region(farm Grid[rune], visited *Grid[bool], x int, y int) []Pos {
	positions := make([]Pos, 0)
	positions = append(positions, Pos{X: x, Y: y})
	visited.Set(x, y, true)

	idx := 0

	n := NewNeighbors4[rune]()
	region := farm.At(x, y)

	for idx < len(positions) {
		current := positions[idx]
		neighbors := n.NeighborCells(farm, current.X, current.Y)

		for _, c := range neighbors {
			if c.Value == region && !visited.At(c.X, c.Y) {
				positions = append(positions, Pos{c.X, c.Y})
				visited.Set(c.X, c.Y, true)
			}
		}

		idx++
	}

	return positions
}

func RegionPrice(farm Grid[rune], visited *Grid[bool], x int, y int) int {
	positions := Region(farm, visited, x, y)

	n := NewNeighbors4[rune]()
	fences := 0
	region := farm.At(x, y)

	for _, current := range positions {
		neighbors := n.NeighborCells(farm, current.X, current.Y)
		fences += 4 - len(neighbors)

		for _, c := range neighbors {
			if c.Value != region {
				fences++
			}
		}
	}

	return fences * len(positions)
}

func RegionPriceWithDiscount(farm Grid[rune], visited *Grid[bool], x int, y int) int {
	positions := Region(farm, visited, x, y)

	cornerPlots := NewGrid[string](farm.Width, farm.Height)
	cornerPlots.SetAll(" ")

	horizontal := []Direction{DirRight, DirLeft}
	vertical := []Direction{DirUp, DirDown}

	corners := 0
	region := farm.At(x, y)

	for _, current := range positions {
		cornerPlots.Set(current.X, current.Y, ".")

		// External corners
		for _, h := range horizontal {
			hx, hy := h.Apply(current.X, current.Y)

			if farm.IsCoordValid(hx, hy) && farm.At(hx, hy) == region {
				for _, v := range vertical {
					vx, vy := v.Apply(current.X, current.Y)
					if farm.IsCoordValid(vx, vy) && farm.At(vx, vy) == region && farm.At(hx, vy) != region {
						cornerPlots.Set(current.X, current.Y, "#")
						corners++
					}
				}
			} else {
				for _, v := range vertical {
					vx, vy := v.Apply(current.X, current.Y)
					if !farm.IsCoordValid(vx, vy) || farm.At(vx, vy) != region {
						cornerPlots.Set(current.X, current.Y, "#")
						corners++
					}
				}
			}
		}
	}

	return corners * len(positions)
}

func Price(farm Grid[rune]) int {
	price := 0

	visited := NewGrid[bool](farm.Width, farm.Height)

	for v := range visited.AllCells() {
		if v.Value {
			continue
		}
		price += RegionPrice(farm, visited, v.X, v.Y)
	}

	return price
}

func PriceWithDiscount(farm Grid[rune]) int {
	price := 0

	visited := NewGrid[bool](farm.Width, farm.Height)

	for v := range visited.AllCells() {
		if v.Value {
			continue
		}
		price += RegionPriceWithDiscount(farm, visited, v.X, v.Y)
	}

	return price
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	log.DefaultLogger.SetLevel(log.InfoLevel)

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})

	fmt.Println(input.Farm.Stringf(func(r rune) string { return string(r) }))

	callback(ctx, SolutionFound{Part: 1, Solution: Price(input.Farm)})
	callback(ctx, SolutionFound{Part: 2, Solution: PriceWithDiscount(input.Farm)})

}
