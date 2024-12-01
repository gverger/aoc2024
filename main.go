package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/phuslu/log"
)

var days = []string{
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

type AppConfig struct {
	WinWidth  int
	WinHeight int
}

type App struct {
	Config AppConfig
}

func runDay(day int) {
	log.Info().Str("Day", days[day]).Msg("Running day...")
}

func (a App) drawMainPanel() {

	const (
		buttonsPerRow = 5
		buttonWidth   = 100
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

	for day := 0; day < 25; day++ {
		x := day / buttonsPerRow
		y := day % buttonsPerRow
		if gui.Button(rl.NewRectangle(panelX+float32(padding+x*buttonWidth+x*gutter), panelY+float32(padding+y*buttonHeight+y*gutter), buttonWidth, buttonHeight), days[day]) {
			runDay(day)
		}
	}

}

func (a App) Init() {
	rl.InitWindow(int32(a.Config.WinWidth), int32(a.Config.WinHeight), "Advent of Code - 2024")
	rl.SetTargetFPS(60)
}

func (a App) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

	a.drawMainPanel()

	rl.EndDrawing()
}

func main() {
	a := App{
		Config: AppConfig{
			WinWidth:  1600,
			WinHeight: 1000,
		},
	}

	a.Init()

	for !rl.WindowShouldClose() {
		// Draw
		//----------------------------------------------------------------------------------
		a.Draw()
	}
	//
	rl.CloseWindow()

}
