package utils

import "github.com/phuslu/log"

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
