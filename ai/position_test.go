package ai

import (
	"testing"
)

func TestPosition2Coordinate(t *testing.T) {
	x, y := position2Coordinate(10, 5)
	if x != 2 || y != 0 {
		t.Errorf("Expected (2, 0) but got (%v, %v)", x, y)
	}
}

func TestCoordinate2Position(t *testing.T) {
	pos := coordinate2Position(2, 0, 5)
	if pos != 10 {
		t.Errorf("Expected 10 but got %v", pos)
	}
}

func TestIsLine(t *testing.T) {
	if !isLine(10, 15, 5) {
		t.Errorf("Expected true but got false")
	}
	if isLine(10, 17, 5) {
		t.Errorf("Expected false but got true")
	}
}

/*
0  1  2  3  4
5  6  7  8  9
10 11 12 13 14
15 16 17 18 19
20 21 22 23 24
*/

func TestIsAllInLine(t *testing.T) {
	config.inLineDistance = 5
	if !isAllInLine(10, []int{15, 20}, 5) {
		t.Errorf("Expected true but got false")
	}
	if isAllInLine(10, []int{15, 21}, 5) {
		t.Errorf("Expected false but got true")
	}
}

func TestHasInLine(t *testing.T) {
	config.inLineDistance = 5
	if !hasInLine(10, []int{15, 21}, 5) {
		t.Errorf("Expected true but got false")
	}
	if hasInLine(10, []int{24, 23}, 5) {
		t.Errorf("Expected false but got true")
	}
}
