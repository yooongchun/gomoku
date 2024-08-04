package tests

import (
	"gomoku/ai"
	"testing"
)

func TestGetShapeFast(t *testing.T) {
	board := make([][]ai.TypeChess, 15)
	for i := range board {
		board[i] = make([]ai.TypeChess, 15)
	}

	// Set up a specific board situation
	board[7][7] = ai.CHESS_BLACK
	board[7][8] = ai.CHESS_BLACK
	board[7][9] = ai.CHESS_BLACK
	board[7][10] = ai.CHESS_BLACK
	board[7][11] = ai.CHESS_BLACK

	shape, _ := ai.GetShapeFast(board, 7, 7, 0, 1, ai.CHESS_BLACK)
	if shape != ai.Shapes.LiveFive {
		t.Errorf("Expected LiveFive, got %v", shape)
	}

	// Set up another board situation
	board[7][7] = ai.CHESS_BLACK
	board[7][8] = ai.CHESS_BLACK
	board[7][9] = ai.CHESS_EMPTY
	board[7][10] = ai.CHESS_BLACK
	board[7][11] = ai.CHESS_BLACK

	shape, _ = ai.GetShapeFast(board, 7, 7, 0, 1, ai.CHESS_BLACK)
	if shape != ai.Shapes.BlockFour {
		t.Errorf("Expected BlockFour, got %v", shape)
	}
}
