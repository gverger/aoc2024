package day7

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"math"
	"strconv"
	"strings"

	. "github.com/gverger/aoc2024/utils"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type Equation struct {
	Result int
	Terms  []int
}
type Input struct {
	Equations []Equation
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	input := Input{
		Equations: make([]Equation, 0),
	}

	for scanner.Scan() {
		elements := strings.SplitN(scanner.Text(), ":", 2)
		value := Must(strconv.Atoi(elements[0]))
		input.Equations = append(input.Equations, Equation{
			Result: value,
			Terms:  MapTo(strings.Split(strings.TrimSpace(elements[1]), " "), func(s string) int { return Must(strconv.Atoi(s)) }),
		})

	}
	AssertNoErr(scanner.Err())

	return input
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

func solve1(e Equation) (string, bool) {
	if len(e.Terms) == 1 {
		return strconv.Itoa(e.Result), e.Result == e.Terms[0]
	}

	last := e.Terms[len(e.Terms)-1]
	if e.Result%last == 0 {
		if line, ok := solve1(Equation{Result: e.Result / last, Terms: e.Terms[:len(e.Terms)-1]}); ok {
			return line + " * " + strconv.Itoa(last), true
		}
	}

	if e.Result-last >= 0 {
		if line, ok := solve1(Equation{Result: e.Result - last, Terms: e.Terms[:len(e.Terms)-1]}); ok {
			return line + " + " + strconv.Itoa(last), true
		}
	}

	return "", false
}

func solve2(e Equation) (string, bool) {
	if len(e.Terms) == 1 {
		return strconv.Itoa(e.Result), e.Result == e.Terms[0]
	}

	last := e.Terms[len(e.Terms)-1]
	if e.Result%last == 0 {
		if line, ok := solve2(Equation{Result: e.Result / last, Terms: e.Terms[:len(e.Terms)-1]}); ok {
			return line + " * " + strconv.Itoa(last), true
		}
	}

	if e.Result-last >= 0 {
		if line, ok := solve2(Equation{Result: e.Result - last, Terms: e.Terms[:len(e.Terms)-1]}); ok {
			return line + " + " + strconv.Itoa(last), true
		}
	}

	terms := make([]int, len(e.Terms)-1)
	copy(terms, e.Terms)

	if finishesWith(e.Result, last) {
		result := e.Result
		for i := 0; i < digits(last); i++ {
			result /= 10
		}
		if line, ok := solve2(Equation{Result: result, Terms: e.Terms[:len(e.Terms)-1]}); ok {
			return line + " || " + strconv.Itoa(last), true
		}
	}

	return "", false
}

func digits(n int) int {
	return int(math.Floor(math.Log10(float64(n)) + 1))
}

func finishesWith(number, end int) bool {
	for end != 0 {
		uNumber := number % 10
		uEnd := end % 10
		if uNumber != uEnd {
			return false
		}
		number /= 10
		end /= 10
	}
	return true
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})

	sum1 := 0
	for _, e := range input.Equations {
		if line, ok := solve1(e); ok {
			fmt.Println(e.Result, "=", line)
			sum1 += e.Result
		}
	}

	callback(ctx, SolutionFound{Part: 1, Solution: sum1})

	sum2 := 0
	for _, e := range input.Equations {
		if line, ok := solve2(e); ok {
			fmt.Println(e.Result, "=", line)
			sum2 += e.Result
		}
	}

	callback(ctx, SolutionFound{Part: 2, Solution: sum2})
}
