package ai

import (
	"math"
)

func Position2Coordinate(position, size int) (int, int) {
	return position / size, position % size
}

func Coordinate2Position(x, y, size int) int {
	return x*size + y
}

func IsLine(a, b, size int) bool {
	x1, y1 := Position2Coordinate(a, size)
	x2, y2 := Position2Coordinate(b, size)
	maxDistance := Config.InLineDistance
	return (x1 == x2 && math.Abs(float64(y1-y2)) < float64(maxDistance)) ||
		(y1 == y2 && math.Abs(float64(x1-x2)) < float64(maxDistance)) ||
		(math.Abs(float64(x1-x2)) == math.Abs(float64(y1-y2)) && math.Abs(float64(x1-x2)) < float64(maxDistance))
}

func IsAllInLine(p int, arr []int, size int) bool {
	for _, val := range arr {
		if !IsLine(p, val, size) {
			return false
		}
	}
	return true
}

func HasInLine(p int, arr []int, size int) bool {
	for _, val := range arr {
		if IsLine(p, val, size) {
			return true
		}
	}
	return false
}
