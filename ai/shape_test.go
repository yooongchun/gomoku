package ai

import (
	"testing"
)

func TestGetShapeFast(t *testing.T) {
	board := make([][]TypeRole, 15)
	for i := range board {
		board[i] = make([]TypeRole, 15)
	}

	// Set up a specific board situation
	board[7][7] = BLACK
	board[7][8] = BLACK
	board[7][9] = BLACK
	board[7][10] = BLACK
	board[7][11] = BLACK

	shape, _ := GetShapeFast(board, 7, 7, 0, 1, BLACK)
	if shape != Shapes.FIVE {
		t.Errorf("Expected LiveFive, got %v", shape)
	}

	// Set up another board situation
	board[7][7] = BLACK
	board[7][8] = BLACK
	board[7][9] = EMPTY
	board[7][10] = BLACK
	board[7][11] = BLACK

	shape, _ = GetShapeFast(board, 7, 7, 0, 1, BLACK)
	if shape != Shapes.BLOCK_FOUR {
		t.Errorf("Expected BlockFour, got %v", shape)
	}
}
