package tests

import (
	"gomoku/ai"
	"testing"
)

func TestZobristCache(t *testing.T) {
	z := ai.NewZobristCache(8)

	if z.GetHash() != 0 {
		t.Errorf("Expected 0 but got %v", z.GetHash())
	}

	z.TogglePiece(1, 2, ai.CHESS_BLACK)
	hashAfterToggle := z.GetHash()
	if hashAfterToggle == 0 {
		t.Errorf("Expected non-zero hash but got 0")
	}

	z.TogglePiece(1, 2, ai.CHESS_BLACK)
	if z.GetHash() != 0 {
		t.Errorf("Expected 0 but got %v", z.GetHash())
	}

	z.TogglePiece(1, 2, ai.CHESS_BLACK)
	if z.GetHash() != hashAfterToggle {
		t.Errorf("Expected %v but got %v", hashAfterToggle, z.GetHash())
	}
}
