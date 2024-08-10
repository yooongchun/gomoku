package ai

import (
	"fmt"
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
	if !slice.Contain(shape.Name, XXXXX) {
		t.Errorf("Expected %s, got %v", XXXXX, shape.Name)
	}
	shape = GetShape(board, NewPoint(2, 8), DirectionVec[HORIZONTAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, XXXXX) {
		t.Errorf("Expected %s, got %v", XXXXX, shape.Name)
	}
	board[3][8] = CHESS_BLACK
	shape = GetShape(board, NewPoint(2, 8), DirectionVec[VERTICAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, OXXOOO) {
		t.Errorf("Expected %s, got %v", OXXOOO, shape.Name)
	}
	// 反转方向上是否能检测出来
	board[2][7] = CHESS_OBSTACLE
	shape = GetShape(board, NewPoint(2, 11), DirectionVec[HORIZONTAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, ZXXXXO) {
		t.Errorf("Expected %s, got %v", ZXXXXO, shape.Name)
	}
	board[2][7] = CHESS_BLACK
	board[2][11] = CHESS_OBSTACLE
	shape = GetShape(board, NewPoint(2, 10), DirectionVec[HORIZONTAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, ZXXXXO) {
		t.Errorf("Expected %s, got %v", ZXXXXO, shape.Name)
	}
	// 斜线方向上是否能检测出来
	board[1][6] = CHESS_BLACK
	board[2][8] = CHESS_EMPTY
	shape = GetShape(board, NewPoint(3, 8), DirectionVec[DIAGONAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, OXXXOO) {
		t.Errorf("Expected %s, got %v", OXXXOO, shape.Name)
	}
	board[0][5] = CHESS_OBSTACLE
	shape = GetShape(board, NewPoint(3, 8), DirectionVec[DIAGONAL], CHESS_BLACK)
	if !slice.Contain(shape.Name, ZXXXOO) {
		t.Errorf("Expected %s, got %v", ZXXXOO, shape.Name)
	}
	board[4][7] = CHESS_WHITE
	board[5][6] = CHESS_WHITE
	board[6][5] = CHESS_WHITE
	shape = GetShape(board, NewPoint(7, 4), DirectionVec[ANTI_DIAGONAL], CHESS_WHITE)
	if !slice.Contain(shape.Name, ZXXXOO) {
		t.Errorf("Expected %s, got %v", ZXXXOO, shape.Name)
	}

	fmt.Println(getBoardString(board, nil))
}
