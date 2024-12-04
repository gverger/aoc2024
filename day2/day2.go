package day2

import (
	"bufio"
	"embed"
	"strconv"
	"strings"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/aoc2024/aoc"
	"github.com/phuslu/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Report struct {
	Levels []int
}

type Input struct {
	Reports []Report
}

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

var printer = message.NewPrinter(language.French)

func NewApp(a *aoc.App) *App {
	return &App{
		app:   a,
		quit:  make(chan bool),
		state: &State{},
	}
}

type App struct {
	app   *aoc.App
	state *State

	quit    chan bool
	running bool
}

type State struct {
	NbOfSafeReportsP1 int
	NbOfSafeReportsP2 int
}

func (a *App) Start() {
	if !a.running {
		a.running = true
		go a.Run()
	}
}

func (a *App) Run() {
	a.state = &State{}

	input := ReadInput("input.txt")
	a.state.NbOfSafeReportsP1 = len(Filter(input.Reports, func(r Report) bool { return r.IsSafeP1() }))
	a.state.NbOfSafeReportsP2 = len(Filter(input.Reports, func(r Report) bool { return r.IsSafeP2() }))

	for {
		select {
		case <-a.quit:
			log.Debug().Msg("Day 2 goroutine stopped")
			a.running = false
			return
		default:
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (a App) Title() string {
	return "Red-Nosed Reports"
}

func (a *App) Draw() {
	a.Start()
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

	area := rl.NewRectangle(10, 10, float32(rl.GetRenderWidth()-20), float32(rl.GetRenderHeight()-20))
	if gui.WindowBox(area, a.Title()) {
		a.Detach()
	}
	rl.BeginScissorMode(int32(area.X), int32(area.Y+24), area.ToInt32().Width, area.ToInt32().Height-24)

	gui.Panel(rl.NewRectangle(500, 300, 200, 100), "Safe Reports Part 1")
	rl.DrawText(printer.Sprintf("%d", a.state.NbOfSafeReportsP1), 520, 350, 32, rl.DarkGreen)

	gui.Panel(rl.NewRectangle(500, 500, 200, 100), "Safe Reports Part 2")
	rl.DrawText(printer.Sprintf("%d", a.state.NbOfSafeReportsP2), 520, 550, 32, rl.DarkGreen)

	rl.EndScissorMode()

	rl.EndDrawing()
}

func (a *App) Detach() {
	log.Info().Msg("Detaching")
	a.quit <- true
	a.app.Day = nil
}

func InputLineToReport(line string) Report {
	values := strings.Fields(line)
	return Report{
		Levels: MapTo(values, func(value string) int {
			return Must(strconv.Atoi(value))
		}),
	}
}

func (s Report) IsSafeP2() bool {
	if s.IsSafeP1() {
		return true
	}

	levels := make([]int, len(s.Levels))
	copy(levels, s.Levels)

	for i := range levels {
		copy(levels[i:], levels[i+1:])
		safer := Report{Levels: levels[:len(levels)-1]}
		if safer.IsSafeP1() {
			return true
		}
		copy(levels, s.Levels)
	}
	return false
}

func (s Report) IsSafeP1() bool {
	if len(s.Levels) < 2 {
		return true
	}
	increasing := true
	if s.Levels[0] > s.Levels[1] {
		increasing = false
	}
	for i := 1; i < len(s.Levels); i++ {
		up := s.Levels[i]-s.Levels[i-1] > 0
		if up && !increasing || !up && increasing {
			return false
		}
		diff := Abs(s.Levels[i] - s.Levels[i-1])
		if diff < 1 || diff > 3 {
			return false
		}
	}
	return true
}

func ReadInput(filename string) Input {
	file, err := f.Open(filename)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer file.Close()

	input := Input{
		Reports: make([]Report, 0),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input.Reports = append(input.Reports, InputLineToReport(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	return input
}

func Abs[T int](value T) T {
	if value > 0 {
		return value
	}
	return -value
}

func Must[T any](value T, err error) T {
	if err != nil {
		log.Fatal().Err(err)
	}
	return value
}

func MapTo[T any, U any](list []T, mapper func(T) U) []U {
	mappedValues := make([]U, len(list))
	for i, v := range list {
		mappedValues[i] = mapper(v)
	}
	return mappedValues
}

func Filter[T any](list []T, keepIt func(T) bool) []T {
	filteredValues := make([]T, 0, len(list))
	for _, v := range list {
		if keepIt(v) {
			filteredValues = append(filteredValues, v)
		}
	}
	return filteredValues
}
