package ai

import (
	"github.com/duke-git/lancet/v2/slice"
	"testing"
)

func TestGetShape(t *testing.T) {
	board := make([][]TypeChess, 15)
	for i := range board {
		board[i] = make([]TypeChess, 15)
	}
	for j := 7; j < 12; j++ {
		board[2][j] = CHESS_BLACK
	}
	// 正常方向上是否能检测出来
	shape := GetShape(board, NewPoint(2, 7), DirectionVec[HORIZONTAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, IIIII) {
		t.Errorf("Expected %s, got %v", IIIII, shape.Name)
	}
	shape = GetShape(board, NewPoint(2, 8), DirectionVec[HORIZONTAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, IIIII) {
		t.Errorf("Expected %s, got %v", IIIII, shape.Name)
	}
	board[3][8] = CHESS_BLACK
	shape = GetShape(board, NewPoint(2, 8), DirectionVec[VERTICAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, OIIOOO) {
		t.Errorf("Expected %s, got %v", OIIOOO, shape.Name)
	}
	// 反转方向上是否能检测出来
	board[2][7] = CHESS_OBSTACLE
	shape = GetShape(board, NewPoint(2, 11), DirectionVec[HORIZONTAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, ZIIIIO) {
		t.Errorf("Expected %s, got %v", ZIIIIO, shape.Name)
	}
	board[2][7] = CHESS_BLACK
	board[2][11] = CHESS_OBSTACLE
	shape = GetShape(board, NewPoint(2, 10), DirectionVec[HORIZONTAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, ZIIIIO) {
		t.Errorf("Expected %s, got %v", ZIIIIO, shape.Name)
	}
}
