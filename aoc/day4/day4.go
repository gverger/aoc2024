package day4

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/aoc2024/aoc"
	"github.com/gverger/aoc2024/day4"
	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

var printer = message.NewPrinter(language.French)

func NewApp(a *aoc.App) *App {
	return &App{
		app:    a,
		events: make(chan Event, 10),
		state:  &State{},
	}
}

type App struct {
	app   *aoc.App
	state *State

	events chan Event
	cancel func()
}

// Init implements aoc.Day.
func (a *App) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	go a.Listen(ctx)
	go day4.Run(ctx, a.events)
}

type State struct {
	Input day4.Input

	Solution1 int
	Solution2 int

	Part1Done bool
	Part2Done bool

	IsDone bool
}

func (a *App) Listen(ctx context.Context) {
	for !a.state.IsDone {
		event := <-a.events
		switch e := event.(type) {
		case day4.InputLoaded:
			a.state.Input = e.Input
			log.Info().Msg("Input loaded")
		case day4.XMasFound:
			log.Info().Interface("event", e).Msg("New XMas")
		case day4.MasInXFound:
			log.Info().Interface("event", e).Msg("New Mas in X")
		case day4.SolutionFound:
			log.Info().Int("part", e.Part).Interface("solution", e.Solution).Msg("Solution found")
			switch e.Part {
			case 1:
				a.state.Solution1 = e.Solution
				a.state.Part1Done = true
			case 2:
				a.state.Solution2 = e.Solution
				a.state.Part2Done = true
			}
		default:
			log.Error().Interface("event", event).Msg("unrecognised event")
		}
		a.state.IsDone = a.state.Part1Done && a.state.Part2Done
		time.Sleep(1 * time.Second)
	}
}

func (a App) Title() string {
	return "CrissCross"
}

func (a *App) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

	area := rl.NewRectangle(10, 10, float32(rl.GetRenderWidth()-20), float32(rl.GetRenderHeight()-20))
	if gui.WindowBox(area, a.Title()) {
		a.Detach()
	}
	rl.BeginScissorMode(int32(area.X), int32(area.Y+24), area.ToInt32().Width, area.ToInt32().Height-24)

	// Draw here

	rl.EndScissorMode()

	rl.EndDrawing()
}

func (a *App) Detach() {
	log.Info().Msg("Detaching")
	a.cancel()
	a.app.Day = nil
}

type Input struct {
	Grid Grid[rune]
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	return Input{}
}
