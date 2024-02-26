package tests

import (
	"gomoku/ai"
	"testing"
)

func TestGetShapeFast(t *testing.T) {
	board := make([][]ai.TypeRole, 15)
	for i := range board {
		board[i] = make([]ai.TypeRole, 15)
	}

	// Set up a specific board situation
	board[7][7] = ai.BLACK
	board[7][8] = ai.BLACK
	board[7][9] = ai.BLACK
	board[7][10] = ai.BLACK
	board[7][11] = ai.BLACK

	shape, _ := ai.GetShapeFast(board, 7, 7, 0, 1, ai.BLACK)
	if shape != ai.Shapes.LiveFive {
		t.Errorf("Expected LiveFive, got %v", shape)
	}

	// Set up another board situation
	board[7][7] = ai.BLACK
	board[7][8] = ai.BLACK
	board[7][9] = ai.EMPTY
	board[7][10] = ai.BLACK
	board[7][11] = ai.BLACK

	shape, _ = ai.GetShapeFast(board, 7, 7, 0, 1, ai.BLACK)
	if shape != ai.Shapes.BlockFour {
		t.Errorf("Expected BlockFour, got %v", shape)
	}
}
