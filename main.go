package main

import (
	"github.com/gverger/aoc2024/aoc"
	"github.com/gverger/aoc2024/day1"
)

func main() {
	a := aoc.NewApp(
		aoc.AppConfig{
			WinWidth:  1600,
			WinHeight: 1000,
		},
	)

	a.RegisterDay(1, day1.NewApp(a))
	a.RegisterDay(2, day1.NewApp(a))
	a.RegisterDay(3, day1.NewApp(a))
	a.RegisterDay(4, day1.NewApp(a))
	a.RegisterDay(5, day1.NewApp(a))
	a.RegisterDay(6, day1.NewApp(a))
	a.RegisterDay(7, day1.NewApp(a))
	a.RegisterDay(8, day1.NewApp(a))
	a.RegisterDay(9, day1.NewApp(a))
	a.RegisterDay(10, day1.NewApp(a))
	a.RegisterDay(11, day1.NewApp(a))
	a.RegisterDay(12, day1.NewApp(a))

	a.Run()
}
