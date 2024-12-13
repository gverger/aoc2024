package day13

import (
	"context"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day13"
	"github.com/phuslu/log"
)

type State struct {
	Input day13.Input
}

type App struct {
	app *cli.App
}

func NewApp(a *cli.App) *App {
	return &App{
		app:   a,
	}
}

type Node struct {
	Id        string   `json:"id"`
	ParentIds []string `json:"parentIds"`
	Info      string   `json:"info"`
	ShortInfo string   `json:"shortInfo"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
}

func (a *App) callback(ctx context.Context, event any) {
	switch e := event.(type) {
	case day13.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
	case day13.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
	}
}

func (a *App) Run() {
	day13.Run(context.Background(), a.callback)
}
