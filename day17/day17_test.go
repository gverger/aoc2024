package day17_test

import (
	"testing"

	"github.com/gverger/aoc2024/day17"
	"github.com/matryer/is"
)

func TestBst(t *testing.T) {
	is := is.New(t)

	c := day17.Computer{C: 9}

	c.Run(2, 6)

	is.Equal(c.B, 1)
}

func TestOut(t *testing.T) {
	is := is.New(t)

	c := day17.Computer{A: 10}

	c.Run(5, 0, 5, 1, 5, 4)

	is.Equal(c.Out, "0,1,2")
}

func TestJnz(t *testing.T) {
	is := is.New(t)
	c := day17.Computer{A: 2024}

	c.Run(0, 1, 5, 4, 3, 0)

	is.Equal(c.Out, "4,2,5,6,7,7,7,7,3,1,0")
	is.Equal(c.A, 0)
}

func TestBxl(t *testing.T) {
	is := is.New(t)
	c := day17.Computer{B: 29}

	c.Run(1, 7)

	is.Equal(c.B, 26)
}

func TestBxc(t *testing.T) {
	is := is.New(t)
	c := day17.Computer{B: 2024, C: 43690}

	c.Run(4, 0)

	is.Equal(c.B, 44354)
}
