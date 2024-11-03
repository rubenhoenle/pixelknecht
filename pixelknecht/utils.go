package main

import (
	"math"
)

func CompareFloat(a, b float64) bool {
	epsilon := 1e-10
	return math.Abs(a-b) < epsilon
}
