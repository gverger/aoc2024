package day13

import (
	"bufio"
	"context"
	"embed"
	"regexp"
	"strconv"

	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type Axes struct {
	X int
	Y int
}
type Machine struct {
	A     Axes
	B     Axes
	Price Axes
}

type Input struct {
	Machines []Machine
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	numbers := regexp.MustCompile(`\d+`)

	machines := make([]Machine, 0)

	for scanner.Scan() {
		a := numbers.FindAllString(scanner.Text(), 2)
		Assert(scanner.Scan(), "cannot scan after '%v'", a)

		b := numbers.FindAllString(scanner.Text(), 2)
		Assert(scanner.Scan(), "cannot scan after '%v'", b)

		price := numbers.FindAllString(scanner.Text(), 2)

		machines = append(machines, Machine{
			A:     Axes{X: Must(strconv.Atoi(a[0])), Y: Must(strconv.Atoi(a[1]))},
			B:     Axes{X: Must(strconv.Atoi(b[0])), Y: Must(strconv.Atoi(b[1]))},
			Price: Axes{X: Must(strconv.Atoi(price[0])), Y: Must(strconv.Atoi(price[1]))},
		})

		if !scanner.Scan() {
			break
		}
	}
	AssertNoErr(scanner.Err(), "reading input file")

	return Input{Machines: machines}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int
}

func IsParallel(m Machine) bool {
	return m.A.X*m.B.Y-m.B.X*m.A.Y == 0
}

func IsSameLine(m Machine) bool {
	return m.Price.X == m.Price.Y && IsParallel(m)
}

func IntegerIntersection(m Machine) (int, int, bool) {
	// m.A.X . a + m.B.X . b = m.P.X
	// m.A.Y . a + m.B.Y . b = m.P.Y

	denom := m.A.X*m.B.Y - m.B.X*m.A.Y
	a := m.Price.X*m.B.Y - m.Price.Y*m.B.X
	if a%denom != 0 {
		return 0, 0, false
	}
	b := m.Price.Y*m.A.X - m.Price.X*m.A.Y
	if b%denom != 0 {
		return 0, 0, false
	}

	return a / denom, b / denom, true
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	log.DefaultLogger.SetLevel(log.InfoLevel)

	input := ReadInput("sample.txt")
	callback(ctx, InputLoaded{Input: input})

	sum := 0
	for _, m := range input.Machines {
		if IsParallel(m) {
			log.Info().Interface("machine", m).Msg("parallel")
		}
		if IsSameLine(m) {
			log.Info().Interface("machine", m).Msg("same line")
		}
		if a, b, ok := IntegerIntersection(m); ok {
			sum += 3*a + b
		}
	}

	callback(ctx, SolutionFound{Part: 1, Solution: sum})

	added := 10000000000000
	sum = 0
	for _, m := range input.Machines {
		if IsParallel(m) {
			log.Info().Interface("machine", m).Msg("parallel")
			continue
		}
		if IsSameLine(m) {
			log.Error().Interface("machine", m).Msg("same line: not implemented")
			continue
		}
		m.Price = Axes{
			X: m.Price.X + added,
			Y: m.Price.Y + added,
		}
		if a, b, ok := IntegerIntersection(m); ok {
			sum += 3*a + b
		}
	}
	callback(ctx, SolutionFound{Part: 2, Solution: sum})

}
