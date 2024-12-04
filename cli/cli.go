package cli

import "github.com/phuslu/log"

type StartDay func(*App) Day

type AppConfig struct {
}

type Day interface {
	Run()
}

type App struct {
	Config AppConfig
	Day    Day

	daysRegistry map[int]Day
}

func NewApp(c AppConfig) *App {
	return &App{
		Config:       c,
		Day:          nil,
		daysRegistry: make(map[int]Day),
	}
}

func (a *App) RegisterDay(day int, dayApp Day) {
	log.Info().Int("day", day).Msg("Registering day")
	if _, ok := a.daysRegistry[day]; ok {
		log.Fatal().Msg("Already registered")
	}
	a.daysRegistry[day] = dayApp
}

func (a App) isRegistered(day int) bool {
	_, ok := a.daysRegistry[day]
	return ok
}

func (a App) Run(day int) {
	app, ok := a.daysRegistry[day]
	if !ok {
		log.Fatal().Int("day", day).Msg("No such day for cli")
	}

	app.Run()
}
