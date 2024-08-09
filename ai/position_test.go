package ai

import (
	"testing"
)

func TestPosition2Coordinate(t *testing.T) {
	x, y := Position2Coordinate(10, 5)
	if x != 2 || y != 0 {
		t.Errorf("Expected (2, 0) but got (%v, %v)", x, y)
	}
}

func TestCoordinate2Position(t *testing.T) {
	pos := Coordinate2Position(2, 0, 5)
	if pos != 10 {
		t.Errorf("Expected 10 but got %v", pos)
	}
}

func TestIsLine(t *testing.T) {
	if !IsLine(10, 15, 5) {
		t.Errorf("Expected true but got false")
	}
	if IsLine(10, 17, 5) {
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
	Config.InLineDistance = 5
	if !IsAllInLine(10, []int{15, 20}, 5) {
		t.Errorf("Expected true but got false")
	}
	if IsAllInLine(10, []int{15, 21}, 5) {
		t.Errorf("Expected false but got true")
	}
}
