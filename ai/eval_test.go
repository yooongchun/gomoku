package ai

import (
	"fmt"
	"testing"
)

func TestUpdateSinglePoint(t *testing.T) {
	e := NewEvaluate(15)
	p0 := NewPoint(0, 1)
	e.board[p0.x+1][p0.y+1] = CHESS_BLACK
	score := e.updateSinglePoint(p0.x, p0.y, CHESS_BLACK)
	fmt.Println("score=", score)
	p1 := NewPoint(1, 2)
	e.board[p1.x+1][p1.y+1] = CHESS_BLACK
	score = e.updateSinglePoint(p1.x, p1.y, CHESS_BLACK)
	fmt.Println("score=", score)
	fmt.Println(getBoardString(e.board, nil, nil))
}
