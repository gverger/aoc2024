package day3

import (
	"bufio"
	"embed"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/raylib-go/easings"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/aoc2024/aoc"
	. "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

var printer = message.NewPrinter(language.French)

func NewApp(a *aoc.App) *App {
	return &App{
		app:     a,
		quit:    make(chan bool, 1),
		state:   &State{},
		actions: make(chan Action),
	}
}

type App struct {
	app   *aoc.App
	state *State

	quit    chan bool
	running bool

	actions       chan Action
	currentAction Action
}

type Mult struct {
	A int
	B int
}

func (m Mult) Value() int {
	return m.A * m.B
}

func (m Mult) String() string {
	return fmt.Sprintf("%dx%d", m.A, m.B)
}

type State struct {
	Input   Input
	LinePos *Pos

	Mults   VerticalList[Mult]
	MultSum int
}

func (a *App) Start() {
	a.state.Input = ReadInput("input.txt")
	if !a.running {
		a.running = true
		go a.Run()
	}
}

func (a *App) Run() {
	// init
	a.state.LinePos = &Pos{X: 400, Y: 200}
	a.state.Mults = VerticalList[Mult]{
		Items:   make([]Mult, 0),
		Focused: 0,
		MaxSize: 40,
	}

	idxEnd := 0
	re := regexp.MustCompile(`mul\(\d\d?\d?,\d\d?\d?\)`)
	numbers := regexp.MustCompile(`\d+`)
	speed := 2000.0
	line := a.state.Input.Line

	for {
		select {
		case <-a.quit:
			log.Debug().Msg("Day 1 goroutine stopped")
			a.running = false
			return
		default:
		}

		// steps
		loc := re.FindStringIndex(line[idxEnd:])
		log.Info().Msg("Reading")

		if idxEnd < len(line) && loc != nil {
			prefix := line[:idxEnd+loc[0]]
			text := line[idxEnd+loc[0] : idxEnd+loc[1]]
			n := numbers.FindAllString(text, 2)

			do := strings.LastIndex(prefix, "do()")
			dont := strings.LastIndex(prefix, "don't()")

			idxEnd += loc[1]
			if do >= dont {
				prefixSize := int(rl.MeasureTextEx(a.app.Font, prefix, 32, 0).X)

				nextX := 400 - int(prefixSize)
				a.actions <- Move(a.state.LinePos, Pos{X: nextX, Y: a.state.LinePos.Y}, time.Duration(speed)*time.Millisecond)
				speed *= 0.9
				a.actions <- &ActionFunc{
					Run: func(actionFunc *ActionFunc) {
						mult := Mult{A: Must(strconv.Atoi(n[0])), B: Must(strconv.Atoi(n[1]))}
						a.state.Mults.Items = append(a.state.Mults.Items, mult)
						a.state.MultSum += mult.A * mult.B
						actionFunc.done = true
					},
				}
			}
		} else {
			log.Info().Int("sum", a.state.MultSum).Msg("response")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (a App) Title() string {
	return "Mull It Over"
}

type Pos struct {
	X int
	Y int
}

type ActionFunc struct {
	done bool
	Run  func(*ActionFunc)
}

func (a *ActionFunc) IsDone() bool {
	return a.done
}

// Tick implements Action.
func (a *ActionFunc) Tick() {
	a.Run(a)
}

type wait struct {
	Duration    int
	currentTime int
}

func (w *wait) IsDone() bool {
	return w.currentTime >= w.Duration
}

// Tick implements Action.
func (w *wait) Tick() {
	w.currentTime++
}

func Wait(duration time.Duration) *wait {
	return &wait{
		Duration:    int(duration * 60 / time.Second), // FPS = 60
		currentTime: 0,
	}
}

var _ Action = &wait{}

type move struct {
	From     Pos
	To       Pos
	Duration int

	element     *Pos
	currentTime int
}

func Move(el *Pos, to Pos, duration time.Duration) *move {
	return &move{
		To:          to,
		Duration:    int(duration * 60 / time.Second), // FPS = 60
		element:     el,
		currentTime: 0,
	}
}

func (m *move) Tick() {
	if m.currentTime == 0 {
		m.From = Pos{
			X: m.element.X,
			Y: m.element.Y,
		}
	}
	if m.currentTime >= m.Duration {
		m.element.X = m.To.X
		m.element.Y = m.To.Y
	} else {
		m.element.X = int(easings.SineInOut(float32(m.currentTime), float32(m.From.X), float32(m.To.X-m.From.X), float32(m.Duration)))
		m.element.Y = int(easings.SineInOut(float32(m.currentTime), float32(m.From.Y), float32(m.To.Y-m.From.Y), float32(m.Duration)))

		m.currentTime++
	}
}

func (m move) IsDone() bool {
	return m.currentTime >= m.Duration
}

var _ Action = &move{}

type Action interface {
	Tick()
	IsDone() bool
}

type VerticalList[T fmt.Stringer] struct {
	Items   []T
	Focused int

	MaxSize int
}

func (l VerticalList[T]) DrawAt(x, y int32) {
	drawItem := func(item string, color color.RGBA) {
		rl.DrawText(item, x, y, 8, color)
		y += 16
	}
	drawList := func(list []T, color color.RGBA) {
		for _, item := range list {
			drawItem(item.String(), color)
		}
	}

	if len(l.Items) <= l.MaxSize {
		drawList(l.Items, rl.Black)
		return
	}

	if l.Focused < 0 {
		drawList(l.Items[:l.MaxSize/2], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[len(l.Items)-l.MaxSize/2+1:], rl.Black)
	} else if l.Focused < l.MaxSize/2 {
		drawList(l.Items[:l.Focused], rl.Black)
		drawItem(l.Items[l.Focused].String(), rl.DarkGreen)
		drawList(l.Items[l.Focused+1:l.MaxSize/2], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[len(l.Items)-l.MaxSize/2+1:], rl.Black)
	} else if l.Focused > len(l.Items)-l.MaxSize/2 {
		drawList(l.Items[:l.MaxSize/2], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[len(l.Items)-l.MaxSize/2+1:l.Focused], rl.Black)
		drawItem(l.Items[l.Focused].String(), rl.DarkGreen)
		drawList(l.Items[l.Focused+1:], rl.Black)
	} else {
		drawList(l.Items[:l.MaxSize/4], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[l.Focused-l.MaxSize/4+1:l.Focused], rl.Black)
		drawItem(l.Items[l.Focused].String(), rl.DarkGreen)
		drawList(l.Items[l.Focused+1:l.Focused+l.MaxSize/4], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[len(l.Items)-l.MaxSize/4+1:], rl.Black)
	}

}
func (a *App) Draw() {
	a.Start()

	if a.currentAction == nil || a.currentAction.IsDone() {
		a.currentAction = nil
		select {
		case a.currentAction = <-a.actions:
		default:
		}
	}

	if a.currentAction != nil {
		a.currentAction.Tick()
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))
	txtDim := rl.MeasureTextEx(a.app.Font, "mul(888,888)", 32, 4)

	area := rl.NewRectangle(10, 10, float32(rl.GetRenderWidth()-20), float32(rl.GetRenderHeight()-20))
	if gui.WindowBox(area, a.Title()) {
		a.Detach()
	}
	rl.BeginScissorMode(int32(area.X), int32(area.Y+24), area.ToInt32().Width, area.ToInt32().Height-24)

	// Draw here

	gui.Panel(rl.NewRectangle(400, 176, txtDim.X, txtDim.Y+24), fmt.Sprintf("Scanner"))
	rl.BeginScissorMode(400, 200, int32(txtDim.X), int32(txtDim.Y))
	line := a.state.Input.Line
	rl.DrawTextEx(a.app.Font, line, rl.NewVector2(float32(a.state.LinePos.X), float32(a.state.LinePos.Y)), 32, 0, rl.Black)
	rl.EndScissorMode()

	gui.Panel(rl.NewRectangle(400, 300, 60, 42*16), "Operations")
	a.state.Mults.DrawAt(410, 330)

	gui.Panel(rl.NewRectangle(500, 300, 200, 100), "Total Sum")
	rl.DrawText(printer.Sprintf("%d", a.state.MultSum), 520, 350, 32, rl.DarkGreen)

	rl.EndScissorMode()

	rl.EndDrawing()
}

func (a *App) Detach() {
	log.Info().Msg("Detaching")
	a.quit <- true
	select {
	case <-a.actions:
	default:
	}
	a.app.Day = nil
}

type Input struct {
	Line string
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	return Input{
		Line: strings.Join(lines, " "),
	}
}
