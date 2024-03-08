package ai

import (
	"math"
)

func position2Coordinate(position, size int) (int, int) {
	return position / size, position % size
}

func coordinate2Position(x, y, size int) int {
	return x*size + y
}

func isLine(a, b, size int) bool {
	x1, y1 := position2Coordinate(a, size)
	x2, y2 := position2Coordinate(b, size)
	maxDistance := Config.InLineDistance
	return (x1 == x2 && math.Abs(float64(y1-y2)) < float64(maxDistance)) ||
		(y1 == y2 && math.Abs(float64(x1-x2)) < float64(maxDistance)) ||
		(math.Abs(float64(x1-x2)) == math.Abs(float64(y1-y2)) && math.Abs(float64(x1-x2)) < float64(maxDistance))
}

func isAllInLine(p int, arr []int, size int) bool {
	for _, val := range arr {
		if !isLine(p, val, size) {
			return false
		}
	}
	return true
}

func hasInLine(p int, points []Point, size int) bool {
	for _, pt := range points {
		if isLine(p, coordinate2Position(pt.x, pt.y, size), size) {
			return true
		}
	}
	return false
}
