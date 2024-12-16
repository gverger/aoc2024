package utils

type Braille [2][4]bool

func (b Braille) Rune() rune {
	lowEndian := [8]bool{b[0][0], b[0][1], b[0][2], b[1][0], b[1][1], b[1][2], b[0][3], b[1][3]}
	var v int
	for i, x := range lowEndian {
		if x {
			v += 1 << uint(i)
		}
	}
	return rune(v) + '\u2800'
}
