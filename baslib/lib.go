package baslib

import (
	"fmt"
	"time"
)

func Sgn(v float64) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	}
	return 0
}

func Date() string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%02d-%02d-%04d", m, d, y)
}

func Time() string {
	h, m, s := time.Now().Clock()
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func Timer() float64 {
	now := time.Now()
	y, m, d := now.Date()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	elapsed := now.Sub(midnight)
	return elapsed.Seconds()
}
