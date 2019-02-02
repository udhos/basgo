package baslib

func Sgn(v float64) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	}
	return 0
}
