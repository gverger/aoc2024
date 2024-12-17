package day16

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day16"
	"github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

type App struct {
	app *cli.App
}

func NewApp(a *cli.App) *App {
	return &App{
		app: a,
	}
}

type Change struct {
	Event any
}

type Done struct{}

func (a *App) Run() {
	m := &model{
		changes: make(chan Change),
		done:    make(chan Done),
	}

	go day16.Run(context.Background(), m.callback)

	utils.Must(tea.NewProgram(m).Run())
}

type model struct {
	grid utils.Grid[day16.CellType]
	sol1 int
	sol2 int

	changes chan Change
	done    chan Done
	styles  Chars
}

type Chars struct {
}

func (m *model) callback(ctx context.Context, event any) {
	switch e := event.(type) {
	case day16.InputLoaded:
		log.Info().Interface("event", e).Msg("loaded")
		m.changes <- Change{Event: e}
	case day16.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
		m.changes <- Change{Event: e}
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
		switch e := msg.Event.(type) {
		case day16.InputLoaded:
			m.grid = *e.Input.Grid
		case day16.SolutionFound:
			if e.Part == 1 {
				m.sol1 = e.Solution
			} else {
				m.sol2 = e.Solution
			}
			m.grid = e.Grid
		}
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
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Solution 1: %v\n", m.sol1))
	sb.WriteString(fmt.Sprintf("Solution 2: %v\n", m.sol2))
	sb.WriteString(fmt.Sprintf("Grid:\n%v\n", m.grid.Stringf(func(ct day16.CellType) string {
		switch ct {
		case day16.Empty:
			return " "
		case day16.Wall:
			return "█"
		case day16.Footprints:
			return "X"
		}
		return "?"
	})))

	return sb.String()
}

var _ tea.Model = &model{}
