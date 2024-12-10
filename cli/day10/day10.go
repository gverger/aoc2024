package day10

import (
	"context"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day10"
	"github.com/phuslu/log"
)

type State struct {
	Input day10.Input
}

type App struct {
	app *cli.App
}

func NewApp(a *cli.App) *App {
	return &App{
		app: a,
	}
}

func (a *App) callback(ctx context.Context, event any) {
	switch e := event.(type) {
	case day10.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
	case day10.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
	}
}

func (a *App) Run() {
	day10.Run(context.Background(), a.callback)
}
