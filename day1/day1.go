package day1

import (
	"bufio"
	"embed"
	"image/color"
	"slices"
	"strconv"
	"strings"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/aoc2024/aoc"
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
		app:   a,
		quit:  make(chan bool),
		state: &State{},
	}
}

type App struct {
	app   *aoc.App
	state *State

	quit    chan bool
	running bool
}

type VerticalList struct {
	Items   []int
	Focused int

	MaxSize int
}

func (l VerticalList) DrawAt(x, y int32) {
	drawItem := func(item string, color color.RGBA) {
		rl.DrawText(item, x, y, 8, color)
		y += 16
	}
	drawList := func(list []int, color color.RGBA) {
		for _, item := range list {
			drawItem(printer.Sprintf("%d", item), color)
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
		drawItem(printer.Sprintf("%d", l.Items[l.Focused]), rl.DarkGreen)
		drawList(l.Items[l.Focused+1:l.MaxSize/2], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[len(l.Items)-l.MaxSize/2+1:], rl.Black)
	} else if l.Focused > len(l.Items)-l.MaxSize/2 {
		drawList(l.Items[:l.MaxSize/2], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[len(l.Items)-l.MaxSize/2+1:l.Focused], rl.Black)
		drawItem(printer.Sprintf("%d", l.Items[l.Focused]), rl.DarkGreen)
		drawList(l.Items[l.Focused+1:], rl.Black)
	} else {
		drawList(l.Items[:l.MaxSize/4], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[l.Focused-l.MaxSize/4+1:l.Focused], rl.Black)
		drawItem(printer.Sprintf("%d", l.Items[l.Focused]), rl.DarkGreen)
		drawList(l.Items[l.Focused+1:l.Focused+l.MaxSize/4], rl.Black)
		drawItem("...", rl.Black)
		drawList(l.Items[len(l.Items)-l.MaxSize/4+1:], rl.Black)
	}

}

type State struct {
	List1 VerticalList
	List2 VerticalList

	Distance    VerticalList
	DistanceSum int

	Counts        map[int]int
	SimilaritySum int
}

type Item struct {
	idx   int
	value int
}

func (a *App) Start() {
	if !a.running {
		a.running = true
		go a.Run()
	}
}

func (a *App) Run() {
	a.state = &State{}

	list1, list2 := ReadInput("sample.txt")
	slices.Sort(list1)
	slices.Sort(list2)

	a.state.List1 = VerticalList{
		Items:   list1,
		Focused: -1,
		MaxSize: 40,
	}
	a.state.List2 = VerticalList{
		Items:   list2,
		Focused: -1,
		MaxSize: 40,
	}
	a.state.Distance = VerticalList{
		Items:   make([]int, len(list1)),
		Focused: -1,
		MaxSize: 40,
	}
	a.state.Counts = make(map[int]int)
	for _, v := range list2 {
		a.state.Counts[v]++
	}

	for {
		select {
		case <-a.quit:
			log.Debug().Msg("Day 1 goroutine stopped")
			a.running = false
			return
		default:
		}

		time.Sleep(10 * time.Millisecond)
		if a.state.List1.Focused < len(a.state.List1.Items)-1 {
			a.state.List1.Focused++
			a.state.List2.Focused++
			value1 := a.state.List1.Items[a.state.List1.Focused]
			value2 := a.state.List2.Items[a.state.List2.Focused]
			a.state.Distance.Items[a.state.List1.Focused] = Abs(value1 - value2)
			a.state.Distance.Focused++
			a.state.DistanceSum += Abs(value1 - value2)
			a.state.SimilaritySum += value1 * a.state.Counts[value1]
		}

	}
}

func (a App) Title() string {
	return "Historian Hysteria"
}

func (a *App) Draw() {
	a.Start()
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

	area := rl.NewRectangle(10, 10, float32(rl.GetRenderWidth()-20), float32(rl.GetRenderHeight()-20))
	if gui.WindowBox(area, a.Title()) {
		a.Detach()
	}
	rl.BeginScissorMode(int32(area.X), int32(area.Y+24), area.ToInt32().Width, area.ToInt32().Height-24)
	// rl.PushMatrix()

	gui.Panel(rl.NewRectangle(80, 70, 60, 42*16), "List 1")
	a.state.List1.DrawAt(100, 100)
	gui.Panel(rl.NewRectangle(180, 70, 60, 42*16), "List 2")
	a.state.List2.DrawAt(200, 100)

	gui.Panel(rl.NewRectangle(280, 70, 60, 42*16), "Distance")
	a.state.Distance.DrawAt(300, 100)

	gui.Panel(rl.NewRectangle(500, 300, 200, 100), "Total Distance")
	rl.DrawText(printer.Sprintf("%d", a.state.DistanceSum), 520, 350, 32, rl.DarkGreen)

	gui.Panel(rl.NewRectangle(500, 500, 200, 100), "Total Similarity")
	rl.DrawText(printer.Sprintf("%d", a.state.SimilaritySum), 520, 550, 32, rl.DarkGreen)

	// rl.PopMatrix()
	rl.EndScissorMode()

	rl.EndDrawing()
}

func (a *App) Detach() {
	log.Info().Msg("Detaching")
	a.quit <- true
	a.app.Day = nil
}

func InputLineToPair(line string) (int, int) {
	values := strings.Fields(line)
	return Must(strconv.Atoi(values[0])), Must(strconv.Atoi(values[1]))
}

func ReadInput(filename string) ([]int, []int) {
	file, err := f.Open("input.txt")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	list1 := make([]int, 0)
	list2 := make([]int, 0)
	for scanner.Scan() {
		v1, v2 := InputLineToPair(scanner.Text())
		list1 = append(list1, v1)
		list2 = append(list2, v2)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err)
	}

	return list1, list2
}

func main() {
	list1, list2 := ReadInput("./sample.txt")

	slices.Sort(list1)
	slices.Sort(list2)

	log.Debug().Interface("list 1", list1).Msg("")
	log.Debug().Interface("list 2", list2).Msg("")

	distance := 0
	for i := range list1 {
		distance += Abs(list2[i] - list1[i])
	}

	log.Info().Int("distance", distance).Msg("Part 1")

	present := make(map[int]int)
	for _, v := range list2 {
		present[v]++
	}

	similarity := 0
	for _, v := range list1 {
		similarity += v * present[v]
	}

	log.Info().Int("similarity", similarity).Msg("Part 2")

}

func Minimum(list []int) int {
	if len(list) == 0 {
		log.Fatal().Msg("Empty list")
	}
	m := list[0]
	for _, value := range list[1:] {
		if value < m {
			m = value
		}
	}
	return m
}

func Must[T any](value T, err error) T {
	if err != nil {
		log.Fatal().Err(err)
	}
	return value
}

func Abs[T int](value T) T {
	if value < 0 {
		return -value
	}
	return value
}
