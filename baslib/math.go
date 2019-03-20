package baslib

import (
	"math"
)

func Pow(a, b float64) float64 {
	return math.Pow(a, b)
}

func Sin(a float64) float64 {
	return math.Sin(a)
}

func Cos(a float64) float64 {
	return math.Cos(a)
}

func Tan(a float64) float64 {
	return math.Tan(a)
}

func Atn(x float64) float64 {
	return math.Atan(x)
}

func Sqr(v float64) float64 {
	return math.Sqrt(v)
}

func Log(v float64) float64 {
	return math.Log(v)
}
