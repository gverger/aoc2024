package day5

import (
	"context"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day5"
	"github.com/phuslu/log"
)

type State struct {
	Input day5.Input
}

type App struct {
	app *cli.App
}

func NewApp(a *cli.App) *App {
	return &App{
		app: a,
	}
}

func callback(ctx context.Context, event any) {
	log.Info().Interface("event", event).Msg("callback")
}

func (a *App) Run() {
	day5.Run(context.Background(), callback)
}
