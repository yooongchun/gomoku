package ai

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
)

type Play struct {
	board *Board
}

func NewPlay(size int, firstRole TypeRole) *Play {
	board := NewBoard(size, firstRole)
	return &Play{
		board: board,
	}
}

func (p *Play) Play() {
	for {
		p.board.Display(nil)
		winner := p.board.GetWinner()
		if winner != NOBODY {
			var winnerName string
			if winner == ROLE_HUMAN {
				winnerName = "HUMAN"
			} else {
				winnerName = "AI"
			}
			logrus.Infoln("Game Over, Winner is ", winnerName)
			break
		}
		if p.board.IsGameOver() {
			logrus.Infoln("Game Over, No Winner")
			break
		}
		if p.board.WhoseTurn() == ROLE_HUMAN {
			point := p.getUserInput()
			p.board.Move(point)
		} else {
			// AI move
			point := p.getAiMove()
			p.board.Move(point)
		}
	}
}

func (p *Play) getUserInput() Point {
	var x, y int
	for {
		fmt.Print("Input your movement x, y: ")
		_, err := fmt.Scanf("%d,%d", &x, &y)
		if err != nil {
			logrus.Errorln("Invalid input", err)
			continue
		}
		if x < 0 || x >= p.board.size || y < 0 || y >= p.board.size {
			logrus.Errorln("Invalid input, out of range")
			continue
		}
		return Point{x, y}
	}
}

func (p *Play) getAiMove() Point {
	for {
		x := rand.Intn(p.board.size)
		y := rand.Intn(p.board.size)
		if p.board.board[x][y] == CHESS_EMPTY {
			return Point{x, y}
		}
	}
}
