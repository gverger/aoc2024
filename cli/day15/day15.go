package day15

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	Grid utils.Grid[day15.CellType]
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
		// log.Info().Interface("event", e).Msg("loaded")
		a.changes <- Change{Grid: *e.Input.Grid}
	case day15.SolutionFound:
		log.Info().Interface("event", e).Msg("solution")
		a.done <- Done{}
	case day15.Moved:
		a.changes <- Change{Grid: e.Grid}
		// time.Sleep(1 * time.Second)
	}
}

func (a *App) Run() {
	commonStyle := lipgloss.NewStyle().Padding(0).Width(1)
	p := tea.NewProgram(&model{
		app: a,
		styles: styles{
			player:              commonStyle.Foreground(lipgloss.Color("5")).Render("@"),
			box:                 commonStyle.Foreground(lipgloss.Color("#008833")).Render("O"),
			boxLeft:             commonStyle.Foreground(lipgloss.Color("#008833")).Render("["),
			boxRight:            commonStyle.Foreground(lipgloss.Color("#008833")).Render("]"),
			wall:                commonStyle.Foreground(lipgloss.Color("#888888")).Render("#"),
			empty:               commonStyle.Foreground(lipgloss.Color("#444444")).Render("."),
			highlightedBoxLeft:  commonStyle.Background(lipgloss.Color("#008833")).Render("["),
			highlightedBoxRight: commonStyle.Background(lipgloss.Color("#008833")).Render("]"),
		},
	})

	go day15.Run(context.Background(), a.callback)

	utils.Must(p.Run())
}

type model struct {
	grid  utils.Grid[day15.CellType]
	app   *App
	moves int

	styles styles
}
type styles struct {
	player   string
	box      string
	boxLeft  string
	boxRight string
	wall     string
	empty    string

	highlightedBoxLeft  string
	highlightedBoxRight string
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
	return fmt.Sprintf("Moves: %d\n%v\n", m.moves,
		m.grid.Stringf(func(s day15.CellType) string {
			if s&day15.Highlighted != 0 {
				s = s - day15.Highlighted
				switch s {
				case day15.Box | day15.Left:
					return m.styles.highlightedBoxLeft
				case day15.Box | day15.Right:
					return m.styles.highlightedBoxRight
				}
			}
			switch s {
			case day15.Empty:
				return m.styles.empty
			case day15.Wall:
				return m.styles.wall
			case day15.Box:
				return m.styles.box
			case day15.Box | day15.Left:
				return m.styles.boxLeft
			case day15.Box | day15.Right:
				return m.styles.boxRight
			case day15.Player:
				return m.styles.player
			default:
				return "?"
			}
		}))
}
