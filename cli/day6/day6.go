package day6

import (
	"context"
	"fmt"
	"strings"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day6"
	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

type State struct {
	Input day6.Input

	currentGrid  *Grid[day6.CellType]
	currentGuard *day6.Guard
}

type App struct {
	app *cli.App

	state *State
}

func NewApp(a *cli.App) *App {
	return &App{
		app:   a,
		state: &State{},
	}
}

func (a *App) callback(ctx context.Context, event any) {
	switch e := event.(type) {
	case day6.InputLoaded:
		a.state.Input = e.Input
		a.state.currentGrid = a.state.Input.Grid.Clone()
		a.state.currentGuard = &day6.Guard{
			X:   e.Input.Guard.X,
			Y:   e.Input.Guard.Y,
			Dir: e.Input.Guard.Dir,
		}
		log.Info().Interface("event", e).Msg("callback")
	case day6.GuardMoved:
		a.state.currentGuard = &day6.Guard{X: e.X, Y: e.Y, Dir: e.Dir}
		x, y := e.OldX, e.OldY
		for x != e.X && y != e.Y {
			a.state.currentGrid.Set(x, y, day6.FootPrintsCell)
		}
	case day6.GuardTurned:
		a.state.currentGuard = &day6.Guard{
			X:   a.state.currentGuard.X,
			Y:   a.state.currentGuard.Y,
			Dir: e.Dir,
		}
	case day6.SolutionFound:
		log.Info().Interface("event", e).Msg("SOLUTION")
	case day6.VisitedGrid:
		// display(&e.Grid, &e.Guard)
	}
	if a.state.currentGrid != nil {
		// display(a.state.currentGrid, a.state.currentGuard)
		// time.Sleep(100 * time.Millisecond)
	}
}

func display(g *Grid[day6.CellType], guard *day6.Guard) {
	fmt.Println(strings.Repeat("-", int(g.Width)+2))
	for y := 0; y < int(g.Height); y++ {
		fmt.Print("|")
		for x := 0; x < int(g.Width); x++ {
			if guard.X == x && guard.Y == y {
				switch guard.Dir {
				case DirUp:
					fmt.Print("^")
				case DirDown:
					fmt.Print("v")
				case DirRight:
					fmt.Print(">")
				case DirLeft:
					fmt.Print("<")
				}
			} else {
				switch g.At(x, y) {
				case day6.EmptyCell:
					fmt.Print(" ")
				case day6.ObstacleCell:
					fmt.Print("#")
				case day6.GuardCell:
					fmt.Print("o")
				case day6.FootPrintsCell:
					fmt.Print("X")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println(strings.Repeat("-", int(g.Width)+2))
}

func (a *App) Run() {
	day6.Run(context.Background(), a.callback)
}
