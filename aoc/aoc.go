package aoc

import (
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/phuslu/log"
)

var days = []string{
	"",
	"Sun 01",
	"Mon 02",
	"Tue 03",
	"Wed 04",
	"Thu 05",
	"Fri 06",
	"Sat 07",
	"Sun 08",
	"Mon 09",
	"Tue 10",
	"Wed 11",
	"Thu 12",
	"Fri 13",
	"Sat 14",
	"Sun 15",
	"Mon 16",
	"Tue 17",
	"Wed 18",
	"Thu 19",
	"Fri 20",
	"Sat 21",
	"Sun 22",
	"Mon 23",
	"Tue 24",
	"Wed 25",
}

type StartDay func(*App) Day

type AppConfig struct {
	WinWidth  int
	WinHeight int
}

type Day interface {
	Draw()
	Title() string
}

type App struct {
	Config AppConfig
	Day    Day

	daysRegistry map[int]Day
	Font         rl.Font
}

func NewApp(c AppConfig) *App {
	return &App{
		Config:       c,
		Day:          nil,
		daysRegistry: make(map[int]Day),
	}
}

func (a *App) RegisterDay(day int, dayApp Day) {
	log.Info().Int("day", day).Msg("Registering day")
	if _, ok := a.daysRegistry[day]; ok {
		log.Fatal().Msg("Already registered")
	}
	a.daysRegistry[day] = dayApp
}

func (a App) isRegistered(day int) bool {
	_, ok := a.daysRegistry[day]
	return ok
}

func (a *App) switchToDay(day int) {
	if a.Day != nil {
		log.Warn().Msg("Already switched to day...")
		return
	}
	log.Info().Str("Day", days[day]).Msg("Running day...")
	a.Day = a.daysRegistry[day]
}

func (a *App) drawMainPanel() {

	const (
		buttonsPerRow = 5
		buttonWidth   = 150
		buttonHeight  = 70
		gutter        = 30
		padding       = 100
	)

	buttonsPerCol := (24 / (buttonsPerRow)) + 1

	panelWidth := buttonsPerRow*buttonWidth + (buttonsPerRow-1)*gutter + padding*2
	panelHeight := (buttonsPerCol)*buttonHeight + (buttonsPerCol-1)*gutter + padding*2

	panelX := float32(a.Config.WinWidth-panelWidth) / 2
	panelY := float32(a.Config.WinHeight-panelHeight) / 2

	gui.Panel(rl.NewRectangle(panelX, panelY, float32(panelWidth), float32(panelHeight)), "Pick your day")

	for day := 1; day <= 25; day++ {
		if !a.isRegistered(day) {
			continue
		}
		x := (day - 1) % buttonsPerRow
		y := (day - 1) / buttonsPerRow
		title := fmt.Sprintf("%s\n%s", days[day], a.daysRegistry[day].Title())
		if gui.Button(rl.NewRectangle(
			panelX+float32(padding+x*buttonWidth+x*gutter),
			panelY+float32(padding+y*buttonHeight+y*gutter),
			buttonWidth,
			buttonHeight,
		), title) {
			a.switchToDay(day)
		}
	}

}

func (a *App) Init() {
	rl.SetConfigFlags(rl.TextureFilterLinear)
	rl.InitWindow(int32(a.Config.WinWidth), int32(a.Config.WinHeight), "Advent of Code - 2024")
	rl.SetTargetFPS(60)
	// a.Font = rl.LoadFontEx("resources/fonts/mecha.png", 16, nil)
}

func (a *App) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

	a.drawMainPanel()

	rl.EndDrawing()
}

func (a *App) Run() {
	a.Init()

	for !rl.WindowShouldClose() {
		// Draw
		//----------------------------------------------------------------------------------
		if a.Day == nil {
			a.Draw()
		} else {
			a.Day.Draw()
		}
	}
	//
	rl.CloseWindow()
}
