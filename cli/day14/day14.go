package day14

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day14"
	"github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

type State struct {
	Input day14.Input
}

type App struct {
	app *cli.App
}

func NewApp(a *cli.App) *App {
	return &App{
		app: a,
	}
}

type Change struct {
	Turn int
	Grid utils.Grid[int]
}

type Done struct{}

func (a *App) Run() {
	m := &model{
		changes: make(chan Change),
		done:    make(chan Done),
	}

	go day14.Run(context.Background(), m.callback)

	utils.Must(tea.NewProgram(m).Run())
}

type model struct {
	turn int
	grid utils.Grid[int]

	changes chan Change
	done    chan Done
}

func (m *model) callback(ctx context.Context, event any) {
	switch e := event.(type) {
	case day14.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
		g := utils.NewGrid[int](uint(e.Width), uint(e.Height))
		for _, r := range e.Input.Robots {
			x := utils.Mod(r.Position.X, e.Width)
			y := utils.Mod(r.Position.Y, e.Height)
			g.Set(x, y, g.At(x, y)+1)
		}
		m.changes <- Change{Turn: 0, Grid: *g}
	case day14.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
	case day14.StateUpdated:
		g := utils.NewGrid[int](uint(e.Width), uint(e.Height))
		for _, p := range e.Positions {
			x := utils.Mod(p.X, e.Width)
			y := utils.Mod(p.Y, e.Height)
			g.Set(x, y, g.At(x, y)+1)
		}
		m.changes <- Change{Turn: e.Turn, Grid: *g}

		for r := range g.AllCells() {
			if r.Value > 1 {
				return
			}
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func waitForChange(change chan Change) tea.Cmd {
	return func() tea.Msg {
		return <-change
	}
}

func waitForDone(done chan Done) tea.Cmd {
	return func() tea.Msg {
		return <-done
	}
}

// Init implements tea.Model.
func (m *model) Init() tea.Cmd {
	return tea.Batch(waitForChange(m.changes), waitForDone(m.done))
}

// Update implements tea.Model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Change:
		m.grid = msg.Grid
		m.turn = msg.Turn
		return m, waitForChange(m.changes)
	case Done:
		return m, tea.Quit
	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}

// View implements tea.Model.
func (m model) View() string {
	return fmt.Sprintf("Turn: %d\n%s\n", m.turn, m.grid.StringDots(func(i int) bool {
		return i > 0
	}))
}

var _ tea.Model = &model{}
