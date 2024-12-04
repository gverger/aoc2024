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
		app:     a,
		events:  make(chan Event, 100),
		state:   &State{},
		actions: make(chan Action, 100),
	}
}

type Tile struct {
	value string
	color rl.Color
}

type App struct {
	app   *aoc.App
	state *State

	events chan Event
	cancel func()

	actions        chan Action
	currentActions []Action

	cells  *Grid[*Tile]
	offset float32
}

func (a *App) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.state = &State{}
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

type Action interface {
	Tick()
	IsDone() bool
}

type Point struct {
	X int
	Y int
}

type highlight struct {
	Cells    []*Tile
	Color    rl.Color
	Duration int

	currentTime int
}

func Highlight(tiles []*Tile, color rl.Color, duration time.Duration) *highlight {
	return &highlight{
		Cells:       tiles,
		Color:       color,
		Duration:    int(duration * 60 / time.Second),
		currentTime: 0,
	}
}

func (h *highlight) IsDone() bool {
	return h.currentTime > h.Duration
}

func (h *highlight) Tick() {
	color := h.Color
	if h.currentTime >= h.Duration {
		color = rl.Black
	}

	for _, c := range h.Cells {
		c.color = color
	}

	h.currentTime++
}

var _ Action = &highlight{}

func (a *App) Listen(ctx context.Context) {
	for !a.state.IsDone {
		event := <-a.events
		switch e := event.(type) {
		case day4.InputLoaded:
			a.state.Input = e.Input
			g := e.Input.Grid
			a.cells = NewGrid[*Tile](g.Width, g.Height)
			for y := 0; y < int(g.Height); y++ {
				for x := 0; x < int(g.Width); x++ {
					a.cells.Set(x, y, &Tile{
						value: g.At(x, y),
						color: rl.Black,
					})
				}
			}
		case day4.XMasFound:
			g := a.state.Input.Grid
			a.offset = float32(e.Y) / float32(g.Height)

			points := make([]*Tile, 4)
			x, y := e.X, e.Y
			for i := 0; i < 4; i++ {
				points[i] = a.cells.At(x, y)
				x, y = e.Dir.Apply(x, y)
			}
			a.actions <- Highlight(points, rl.Green, 200*time.Millisecond)
			a.state.Solution1++
			time.Sleep(1 * time.Millisecond)
		case day4.MasInXFound:
			g := a.state.Input.Grid
			a.offset = float32(e.Y) / float32(g.Height)

			x, y := e.X, e.Y
			points := []*Tile{
				a.cells.At(x, y),
				a.cells.At(x-1, y-1),
				a.cells.At(x-1, y+1),
				a.cells.At(x+1, y+1),
				a.cells.At(x+1, y-1),
			}
			a.actions <- Highlight(points, rl.Blue, 200*time.Millisecond)
			a.state.Solution2++
			time.Sleep(1 * time.Millisecond)
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
			a.offset = 0
			time.Sleep(1 * time.Second)
		default:
			log.Error().Interface("event", event).Msg("unrecognised event")
		}
		a.state.IsDone = a.state.Part1Done && a.state.Part2Done
	}
}

func (a App) Title() string {
	return "Ceres Search"
}

func (a *App) Draw() {
	found := true
	for found {
		select {
		case action := <-a.actions:
			a.currentActions = append(a.currentActions, action)
		default:
			found = false
		}
	}

	for _, action := range a.currentActions {
		if !action.IsDone() {
			action.Tick()
		}
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

	area := rl.NewRectangle(10, 10, float32(rl.GetRenderWidth()-20), float32(rl.GetRenderHeight()-20))
	if gui.WindowBox(area, a.Title()) {
		a.Detach()
	}
	rl.BeginScissorMode(int32(area.X), int32(area.Y+24), area.ToInt32().Width, area.ToInt32().Height-24)

	g := a.state.Input.Grid
	totalHeight := 12 * g.Height
	distance := max(0, float32(totalHeight) - area.Height + 24)

	if a.cells != nil {
		minx := area.X + 100
		miny := area.Y + 24 - a.offset*distance
		for j := 0; j < int(g.Height); j++ {
			y := miny + float32(12*j)
			for i := 0; i < int(g.Width); i++ {
				x := minx + float32(7*i)
				cell := a.cells.At(i, j)
				rl.DrawTextEx(a.app.Font, cell.value, rl.NewVector2(x, y), 16, 0, cell.color)
			}
		}
	}

	xSolutionPanel := area.X + area.Width - 350
	gui.Panel(rl.NewRectangle(xSolutionPanel, 300, 200, 100), "Part 1: XMAS")
	part1Col := rl.Black
	if a.state.Part1Done {
		part1Col = rl.DarkGreen
	}
	rl.DrawText(printer.Sprintf("%d", a.state.Solution1), int32(xSolutionPanel)+20, 350, 32, part1Col)

	gui.Panel(rl.NewRectangle(xSolutionPanel, 500, 200, 100), "Part 2: MAX in X")
	part2Col := rl.Black
	if a.state.Part2Done {
		part2Col = rl.DarkBlue
	}
	rl.DrawText(printer.Sprintf("%d", a.state.Solution2), int32(xSolutionPanel)+20, 550, 32, part2Col)

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
