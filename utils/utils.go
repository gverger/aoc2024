package utils

import (
	"fmt"

	"github.com/phuslu/log"
)

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

func MustSucceed(err error) {
	if err != nil {
		log.Fatal().Err(err)
	}
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

func Assert(condition bool, msg string, args ...any) {
	if condition {
		return
	}

	log.Fatal().Msg(fmt.Sprint("Assertion failed: ", fmt.Sprintf(msg, args...)))
}

func AssertNoErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatal().Err(err)
}

func MapTo[T any, U any](list []T, mapper func(T) U) []U {
	mappedValues := make([]U, len(list))
	for i, v := range list {
		mappedValues[i] = mapper(v)
	}
	return mappedValues
}

func Filter[T any](list []T, keepIt func(T) bool) []T {
	filteredValues := make([]T, 0, len(list))
	for _, v := range list {
		if keepIt(v) {
			filteredValues = append(filteredValues, v)
		}
	}
	return filteredValues
}

func Mod(a, b int) int {
    return (a % b + b) % b
}
