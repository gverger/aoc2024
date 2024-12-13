package day12

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day12"
	"github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

type State struct {
	Input day12.Input
}

type App struct {
	app *cli.App

	nodes []Node
}

func NewApp(a *cli.App) *App {
	return &App{
		app:   a,
		nodes: make([]Node, 0),
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
	case day12.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
	case day12.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
	}
}

func (a *App) Run() {
	day12.Run(context.Background(), a.callback)

	a.nodes = append(a.nodes, Node{Id: "0", ParentIds: make([]string, 0), Info: "Root", ShortInfo: "Root"})

	g := Graph{Nodes: a.nodes}

	file, _ := os.OpenFile("/tmp/graph.json", os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)
	utils.MustSucceed(encoder.Encode(g))
}
