package day11

import (
	"bufio"
	"context"
	"embed"
	"math"
	"strconv"
	"strings"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed *.txt
var f embed.FS

type Input struct {
	Numbers []int
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
	Assert(len(lines) == 1, "one line input")

	numbers := MapTo(strings.Fields(lines[0]), func(c string) int { return Must(strconv.Atoi(c)) })

	return Input{Numbers: numbers}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

type ComputedResult struct {
	stone int
	depth int
}

func digits(n int) int {
	return int(math.Floor(math.Log10(float64(n)) + 1))
}

func countStones(stone int, depth int, cache map[ComputedResult]int) int {
	log.Debug().Int("stone", stone).Int("depth", depth).Msg("counting")
	input := ComputedResult{stone: stone, depth: depth}
	if result, ok := cache[input]; ok {
		return result
	}
	if depth == 0 {
		return 1
	}

	if stone == 0 {
		result := countStones(1, depth-1, cache)
		cache[input] = result
		return result
	}

	nbDigits := digits(stone)
	if nbDigits%2 != 0 {
		result := countStones(stone*2024, depth-1, cache)
		cache[input] = result
		return result
	}

	divider := int(math.Pow10(nbDigits / 2))
	log.Debug().Int("div", divider).Int("first", stone/divider).Int("second", stone%divider).Msg("divide")
	result1 := countStones(stone%divider, depth-1, cache)
	result2 := countStones(stone/divider, depth-1, cache)
	cache[input] = result1 + result2

	return result1 + result2
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	log.DefaultLogger.SetLevel(log.InfoLevel)

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})

	cache := make(map[ComputedResult]int)

	sum := 0
	for _, n := range input.Numbers {
		sum += countStones(n, 25, cache)
	}
	callback(ctx, SolutionFound{Part: 1, Solution: sum})

	sum = 0
	for _, n := range input.Numbers {
		sum += countStones(n, 75, cache)
	}
	callback(ctx, SolutionFound{Part: 2, Solution: sum})
}
