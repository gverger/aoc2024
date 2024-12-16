package utils

func ClearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}
