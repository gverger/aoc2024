package day8

import (
	"context"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day8"
	"github.com/phuslu/log"
)

type State struct {
	Input day8.Input
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
	case day8.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
	case day8.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
	}
}

func (a *App) Run() {
	day8.Run(context.Background(), a.callback)
}
