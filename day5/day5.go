package day5

import (
	"bufio"
	"context"
	"embed"
	"strconv"
	"strings"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type Ordering []int

type Input struct {
	Graph     Graph[int]
	Orderings []Ordering
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	g := NewGraph[int]()
	for scanner.Scan() && len(scanner.Text()) != 0 {
		numbers := strings.FieldsFunc(scanner.Text(), func(r rune) bool { return r == '|' })
		g.AddEdge(Must(strconv.Atoi(numbers[0])), Must(strconv.Atoi(numbers[1])))
	}

	orderings := make([]Ordering, 0)
	for scanner.Scan() {
		numbers := strings.FieldsFunc(scanner.Text(), func(r rune) bool { return r == ',' })
		ordering := make(Ordering, len(numbers))
		for i, v := range numbers {
			ordering[i] = Must(strconv.Atoi(v))
		}
		orderings = append(orderings, ordering)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	if len(orderings) == 0 {
		log.Fatal().Str("file", filename).Msg("No ordering")
	}

	return Input{
		Graph:     *g,
		Orderings: orderings,
	}
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

func isValid(ordering []int, g Graph[int]) bool {
	for i, first := range ordering[:len(ordering)-1] {
		second := ordering[i+1]
		if g.HasEdge(second, first) {
			return false
		}
	}
	return true
}

func reorder(ordering []int, g Graph[int]) []int {
	for i := 0; i < len(ordering)-1; i++ {
		if g.HasEdge(ordering[i+1], ordering[i]) {
			ordering[i], ordering[i+1] = ordering[i+1], ordering[i]
			if i >= 1 {
				i -= 2
			}
		}
	}
	return ordering
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})

	sum1 := 0
	for _, o := range input.Orderings {
		if isValid(o, input.Graph) {
			sum1 += o[len(o)/2]
		}
	}

	callback(ctx, SolutionFound{Part: 1, Solution: sum1})

	sum2 := 0
	for _, o := range input.Orderings {
		if isValid(o, input.Graph) {
			continue
		}

		o = reorder(o, input.Graph)
		sum2 += o[len(o)/2]
	}

	callback(ctx, SolutionFound{Part: 2, Solution: sum2})

}
