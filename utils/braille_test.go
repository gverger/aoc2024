package utils_test

import (
	"testing"

	"github.com/gverger/aoc2024/utils"
	"github.com/matryer/is"
)

func TestBraille(t *testing.T) {
	is := is.New(t)

	b := utils.Braille{{true, false, true, true}, {false, true, true, false}}

	is.Equal(string(b.Rune()), "⡵")
}
