package main

import (
	"os"
	"strconv"

	"github.com/gverger/aoc2024/aoc"
	"github.com/gverger/aoc2024/aoc/day4"
	"github.com/gverger/aoc2024/cli"
	day10cli "github.com/gverger/aoc2024/cli/day10"
	day11cli "github.com/gverger/aoc2024/cli/day11"
	day12cli "github.com/gverger/aoc2024/cli/day12"
	day13cli "github.com/gverger/aoc2024/cli/day13"
	day14cli "github.com/gverger/aoc2024/cli/day14"
	day15cli "github.com/gverger/aoc2024/cli/day15"
	day16cli "github.com/gverger/aoc2024/cli/day16"
	day17cli "github.com/gverger/aoc2024/cli/day17"
	day18cli "github.com/gverger/aoc2024/day18/cli"
	day20cli "github.com/gverger/aoc2024/day20/cli"
	day4cli "github.com/gverger/aoc2024/cli/day4"
	day5cli "github.com/gverger/aoc2024/cli/day5"
	day6cli "github.com/gverger/aoc2024/cli/day6"
	day7cli "github.com/gverger/aoc2024/cli/day7"
	day8cli "github.com/gverger/aoc2024/cli/day8"
	day9cli "github.com/gverger/aoc2024/cli/day9"
	"github.com/gverger/aoc2024/day1"
	"github.com/gverger/aoc2024/day2"
	"github.com/gverger/aoc2024/day3"
	"github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

func gui() {
	a := aoc.NewApp(
		aoc.AppConfig{
			WinWidth:  1600,
			WinHeight: 1000,
		},
	)

	a.RegisterDay(1, day1.NewApp(a))
	a.RegisterDay(2, day2.NewApp(a))
	a.RegisterDay(3, day3.NewApp(a))
	a.RegisterDay(4, day4.NewApp(a))

	a.Run()
}

func console(day string) {
	cli := cli.NewApp(cli.AppConfig{})
	cli.RegisterDay(4, day4cli.NewApp(cli))
	cli.RegisterDay(5, day5cli.NewApp(cli))
	cli.RegisterDay(6, day6cli.NewApp(cli))
	cli.RegisterDay(7, day7cli.NewApp(cli))
	cli.RegisterDay(8, day8cli.NewApp(cli))
	cli.RegisterDay(9, day9cli.NewApp(cli))
	cli.RegisterDay(10, day10cli.NewApp(cli))
	cli.RegisterDay(11, day11cli.NewApp(cli))
	cli.RegisterDay(12, day12cli.NewApp(cli))
	cli.RegisterDay(13, day13cli.NewApp(cli))
	cli.RegisterDay(14, day14cli.NewApp(cli))
	cli.RegisterDay(15, day15cli.NewApp(cli))
	cli.RegisterDay(16, day16cli.NewApp(cli))
	cli.RegisterDay(17, day17cli.NewApp(cli))
	cli.RegisterDay(18, day18cli.NewApp(cli))
	cli.RegisterDay(20, day20cli.NewApp(cli))

	d := utils.Must(strconv.Atoi(day))
	cli.Run(d)
}

func main() {
	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: false,
			},
		}
	}

	args := os.Args[1:]
	log.Debug().Interface("args", args).Msg("Running app")
	if len(args) > 1 && args[0] == "--cli" {
		console(args[1])
	} else {
		gui()
	}
}
