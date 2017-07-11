package common

import (
	"math"
	"strings"
)

func PadLeft(input string, padLength int, padString string) string {
	inputLength := len(input)
	padStringLength := len(padString)
	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	output := input + strings.Repeat(padString, int(repeat))
	output = output[:padLength]
	return output
}

func Map(strings []string, f func(string) string) []string {
	mapped := make([]string, len(strings))
	for i, str := range strings {
		mapped[i] = f(str)
	}
	return mapped
}

func All(strings []string, f func(string) bool) bool {
	for _, v := range strings {
		if !f(v) {
			return false
		}
	}
	return true
}
