package cli

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gverger/aoc2024/cli"
	"github.com/gverger/aoc2024/day20"
	"github.com/gverger/aoc2024/utils"
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
	commonStyle := lipgloss.NewStyle().Padding(0).Width(1)
	m := &model{
		changes: make(chan Change),
		done:    make(chan Done),
		styles: Chars{
			Fall:       commonStyle.Foreground(lipgloss.Color("#888888")).Render("█"),
			Empty:      commonStyle.Foreground(lipgloss.Color("#444444")).Render("."),
			Footprints: commonStyle.Foreground(lipgloss.Color("#009944")).Render("X"),
			Block:      commonStyle.Foreground(lipgloss.Color("#990044")).Render("█"),
		},
	}

	go day20.Run(context.Background(), m.callback)

	utils.Must(tea.NewProgram(m).Run())
}

type model struct {
	grid utils.Grid[day20.CellType]
	sol1 int
	sol2 int

	changes chan Change
	done    chan Done
	styles  Chars
}

type Chars struct {
	Fall       string
	Empty      string
	Footprints string
	Block      string
}

func (m *model) callback(ctx context.Context, event any) {
	switch e := event.(type) {
	case day20.InputLoaded:
		m.changes <- Change{Event: e}
	case day20.SolutionFound:
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
		case day20.InputLoaded:
			m.grid = *e.Input.Grid
		case day20.SolutionFound:
			if e.Part == 1 {
				m.sol1 = e.Solution
			} else {
				m.sol2 = e.Solution
			}
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

	sb.WriteString(m.grid.Stringf(func(ct day20.CellType) string {
		switch ct {

		case day20.Empty:
			return " "
		case day20.Wall:
			return "#"
		case day20.Footprints:
			return "X"
		default:
			panic(fmt.Sprintf("unexpected day20.CellType: %#v", ct))
		}
	}))
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("Solution 1: %v\n", m.sol1))
	sb.WriteString(fmt.Sprintf("Solution 2: %v\n", m.sol2))

	return sb.String()
}

var _ tea.Model = &model{}
