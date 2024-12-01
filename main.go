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

	a.Run()
}
