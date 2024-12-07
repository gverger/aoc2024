package day7

import (
	"context"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day7"
	"github.com/phuslu/log"
)

type State struct {
	Input day7.Input
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
	case day7.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
	case day7.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")

	}
}

func (a *App) Run() {
	day7.Run(context.Background(), a.callback)
}
