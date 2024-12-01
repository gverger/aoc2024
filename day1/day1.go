package day1

import (
	"bufio"
	"os"
	"slices"
	"strconv"
	"strings"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/aoc2024/aoc"
	"github.com/phuslu/log"
)

func NewApp(a *aoc.App) App {
	return App{
		app: a,
	}
}

var _ aoc.Day = App{}

type App struct {
	app *aoc.App
}

// Title implements aoc.Day.
func (a App) Title() string {
	return "Historian Hysteria"
}

// Draw implements aoc.Day.
func (a App) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

	rl.NewRectangle(10, 10, 1000, 400)

	rl.EndDrawing()
}

func InputLineToPair(line string) (int, int) {
	values := strings.Fields(line)
	return Must(strconv.Atoi(values[0])), Must(strconv.Atoi(values[1]))
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	list1 := make([]int, 0)
	list2 := make([]int, 0)
	for scanner.Scan() {
		v1, v2 := InputLineToPair(scanner.Text())
		list1 = append(list1, v1)
		list2 = append(list2, v2)
	}
	slices.Sort(list1)
	slices.Sort(list2)

	log.Debug().Interface("list 1", list1).Msg("")
	log.Debug().Interface("list 2", list2).Msg("")

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	distance := 0
	for i := range list1 {
		distance += Abs(list2[i] - list1[i])
	}

	log.Info().Int("distance", distance).Msg("Part 1")

	present := make(map[int]int)
	for _, v := range list2 {
		present[v]++
	}

	similarity := 0
	for _, v := range list1 {
		similarity += v * present[v]
	}

	log.Info().Int("similarity", similarity).Msg("Part 2")

}

func Minimum(list []int) int {
	if len(list) == 0 {
		log.Fatal().Msg("Empty list")
	}
	m := list[0]
	for _, value := range list[1:] {
		if value < m {
			m = value
		}
	}
	return m
}

func Must[T any](value T, err error) T {
	if err != nil {
		log.Fatal().Err(err)
	}
	return value
}

func Abs[T int](value T) T {
	if value < 0 {
		return -value
	}
	return value
}
