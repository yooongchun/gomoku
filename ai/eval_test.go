package ai

import (
	"fmt"
	"testing"
)

func TestUpdateSinglePoint(t *testing.T) {
	e := NewEvaluate(15)
	p := NewPoint(0, 1)
	e.board[p.x+1][p.y+1] = CHESS_BLACK
	score := e.updateSinglePoint(p.x, p.y, CHESS_BLACK)
	if score != SCORE_NONE {
		t.Errorf("score %d should be 0", score)
	}
	// 眠2
	p = NewPoint(1, 2)
	e.board[p.x+1][p.y+1] = CHESS_BLACK
	score = e.updateSinglePoint(p.x, p.y, CHESS_BLACK)
	if score != SCORE_BLOCK_TWO {
		t.Errorf("score %d should be 0", score)
	}
	// 眠3
	p = NewPoint(2, 3)
	e.board[p.x+1][p.y+1] = CHESS_BLACK
	score = e.updateSinglePoint(p.x, p.y, CHESS_BLACK)
	if score != ShapeEnum.BlockThree.Score {
		t.Errorf("score %d should be %d", score, ShapeEnum.BlockThree.Score)
	}
	// 冲4活3
	p = NewPoint(3, 4)
	e.board[p.x+1][p.y+1] = CHESS_BLACK
	p = NewPoint(3, 3)
	e.board[p.x+1][p.y+1] = CHESS_BLACK
	p = NewPoint(4, 3)
	e.board[p.x+1][p.y+1] = CHESS_BLACK
	score = e.updateSinglePoint(2, 3, CHESS_BLACK)
	if score < ShapeEnum.FourThree.Score {
		t.Errorf("score %d should be %d", score, ShapeEnum.FourThree.Score)
	}
	// 双活3
	e.board[4][3] = CHESS_BLACK
	e.board[1][2] = CHESS_EMPTY
	score = e.updateSinglePoint(3, 3, CHESS_BLACK)
	if score < ShapeEnum.DoubleThree.Score {
		t.Errorf("score %d should be %d", score, ShapeEnum.DoubleThree.Score)
	}
	// 双冲4
	e.board[2][3] = CHESS_EMPTY
	e.board[1][4] = CHESS_WHITE
	e.board[2][4] = CHESS_BLACK
	e.board[4][1] = CHESS_BLACK
	score = e.updateSinglePoint(3, 3, CHESS_BLACK)
	if score < ShapeEnum.DoubleFour.Score {
		t.Errorf("score %d should be %d", score, ShapeEnum.DoubleFour.Score)
	}
	// 连5
	e.board[4][2] = CHESS_BLACK
	score = e.updateSinglePoint(3, 3, CHESS_BLACK)
	if score < ShapeEnum.Five.Score {
		t.Errorf("score %d should be %d", score, ShapeEnum.Five.Score)
	}
	fmt.Println(getBoardString(e.board, nil, nil))

}
