package ai

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Play struct {
	board *Board
	cnt   int
}

func NewPlay(size int, firstRole TypeRole) *Play {
	board := NewBoard(size, firstRole)
	return &Play{
		board: board,
		cnt:   0,
	}
}

func (p *Play) Play() {
	for {
		// Clear the previous board display
		for i := 0; i < p.cnt; i++ {
			fmt.Printf("\033[%dA\r\033[K", 1)
		}
		p.board.Display()
		p.cnt = p.board.size + 2
		winner := p.board.GetWinner()
		if winner != NOBODY {
			var winnerName string
			if winner == ROLE_HUMAN {
				winnerName = "HUMAN"
			} else {
				winnerName = "AI"
			}
			logrus.Infoln("Game Over, Winner is ", winnerName)
			p.board.Save()
			break
		}
		if p.board.IsGameOver() {
			logrus.Infoln("Game Over, No Winner")
			p.board.Save()
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
	re, err := regexp.Compile(`^\s*(\d+)\s*,\s*(\d+)\s*$`)
	if err != nil {
		panic(err)
	}
	fmt.Print("Input your movement x,y: ")
	p.cnt += 1
	for {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			text := input.Text()
			if text == "exit" {
				// 退出游戏
				fmt.Println("Game Over, Goodbye!!")
				os.Exit(0)
			} else if strings.HasPrefix(text, "save") {
				// 保存游戏
				p.board.Save()
				os.Exit(0)
			} else if !re.MatchString(text) {
				// 输入不合法
				fmt.Printf("\033[%dA\r\033[K", 1)
				fmt.Print("Invalid input, please input again: ")
			} else {
				// 判断输入的坐标是否合法
				parts := strings.Split(text, ",")
				x, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
				y, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
				if isPosValid(p.board.board, x, y) && p.board.board[x][y] == CHESS_EMPTY {
					return Point{x, y}
				}
				fmt.Printf("\033[%dA\r\033[K", 1)
				fmt.Print("The position has been occupied or not valid, please input again: ")
			}
		}
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
