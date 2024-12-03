package templates

import (
	"bufio"
	"embed"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/aoc2024/aoc"
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
}

func (a *App) Start() {
	if !a.running {
		a.running = true
		go a.Run()
	}
}

func (a *App) Run() {
	// init

	for {
		select {
		case <-a.quit:
			log.Debug().Msg("Day 1 goroutine stopped")
			a.running = false
			return
		default:
		}

		// steps
	}
}

func (a App) Title() string {
	return "Title of the day"
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

	// Draw here

	rl.EndScissorMode()

	rl.EndDrawing()
}

func (a *App) Detach() {
	log.Info().Msg("Detaching")
	a.quit <- true
	a.app.Day = nil
}

type Input struct {
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	return Input{}
}
