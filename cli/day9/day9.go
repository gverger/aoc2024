package day9

import (
	"context"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day9"
	"github.com/phuslu/log"
)

type State struct {
	Input day9.Input
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
	case day9.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
	case day9.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
	}
}

func (a *App) Run() {
	day9.Run(context.Background(), a.callback)
}
