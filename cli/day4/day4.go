package day4

import (
	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day4"
	"github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

type State struct {
	Input day4.Input

	Solution1 int
	Solution2 int

	Part1Done bool
	Part2Done bool

	IsDone bool
}

type App struct {
	app   *cli.App
	state *State

	running bool
}

func NewApp(a *cli.App) *App {
	return &App{
		app: a,

		state: &State{},
	}
}

func (a *App) Run() {
	events := make(chan utils.Event)

	go func() {
		for !a.state.IsDone {
			event := <-events
			switch e := event.(type) {
			case day4.InputLoaded:
				a.state.Input = e.Input
				log.Info().Msg("Input loaded")
			case day4.XMasFound:
				log.Info().Interface("event", e).Msg("New XMas")
			case day4.MasInXFound:
				log.Info().Interface("event", e).Msg("New Mas in X")
			case day4.SolutionFound:
				log.Info().Int("part", e.Part).Interface("solution", e.Solution).Msg("Solution found")
				switch e.Part {
				case 1:
					a.state.Solution1 = e.Solution
					a.state.Part1Done = true
				case 2:
					a.state.Solution2 = e.Solution
					a.state.Part2Done = true
				}
			default:
				log.Error().Interface("event", event).Msg("unrecognised event")
			}
			a.state.IsDone = a.state.Part1Done && a.state.Part2Done
		}
	}()

	day4.Run(events)
}
