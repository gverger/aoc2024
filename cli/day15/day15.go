package day15

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day15"
	"github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

type State struct {
	Input day15.Input
}

type App struct {
	app *cli.App

	changes chan Change
	done    chan Done
}

type Change struct {
	Grid utils.Grid[string]
}

type Done struct {
}

func NewApp(a *cli.App) *App {
	return &App{
		app:     a,
		changes: make(chan Change),
		done:    make(chan Done),
	}
}

func (a *App) callback(ctx context.Context, event any) {
	switch e := event.(type) {
	case day15.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
	case day15.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
		a.done <- Done{}
	case day15.Moved:
		a.changes <- Change{Grid: *utils.MapGrid(e.Grid, func(ct day15.CellType) string {
			switch ct {
			case day15.Empty:
				return "."
			case day15.Wall:
				return "#"
			case day15.Box:
				return "O"
			case day15.Player:
				return "@"
			default:
				return "?"
			}
		})}
	}
}

func (a *App) Run() {
	p := tea.NewProgram(&model{
		app: a,
	})

	go day15.Run(context.Background(), a.callback)

	utils.Must(p.Run())
}

type model struct {
	grid  utils.Grid[string]
	app   *App
	moves int
}

func waitForChange(sub chan Change) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func waitForDone(sub chan Done) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(waitForChange(m.app.changes), waitForDone(m.app.done))
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Change:
		m.grid = msg.Grid
		m.moves++
		return m, waitForChange(m.app.changes)
	case Done:
		return m, tea.Quit
	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Moves: %d\n%v", m.moves, m.grid)
}
