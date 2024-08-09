package ai

import (
	"fmt"
	"strings"
)

func getBoardString(board [][]TypeChess, lastOp *Point, extraPoints []Point) string {
	size := len(board)
	extraPositions := make(map[int]bool, len(extraPoints))
	for _, point := range extraPoints {
		extraPositions[Coordinate2Position(point.x, point.y, size)] = true
	}
	var result strings.Builder
	prefix := "  "
	result.WriteString(prefix + " ")
	for i := 0; i < size; i++ {
		result.WriteString(fmt.Sprintf("%X", i) + prefix)
	}
	result.WriteString("\n")

	for i := 0; i < size; i++ {
		result.WriteString(fmt.Sprintf("%X", i) + prefix)
		for j := 0; j < size; j++ {
			position := Coordinate2Position(i, j, size)
			if ok, exist := extraPositions[position]; ok && exist {
				result.WriteString("?" + prefix)
				continue
			}

			op := ""
			switch board[i][j] {
			case CHESS_BLACK:
				op = "O"
			case CHESS_WHITE:
				op = "X"
			default:
				op = "-"
			}
			if j == 0 && lastOp != nil && lastOp.y == 0 && lastOp.x == i {
				op = "[" + op + "]"
			} else if lastOp != nil && (j+1 == lastOp.y) && i == lastOp.x {
				op = op + " ["
			} else if lastOp != nil && j == lastOp.y && i == lastOp.x {
				op = op + "] "
			} else {
				op = op + "  "
			}

			result.WriteString(op)
		}
		result.WriteString("\n") // New line at the end of each row
	}
	return result.String()
}

func togglePiece(chess TypeChess) TypeChess {
	if chess == CHESS_BLACK {
		return CHESS_WHITE
	}
	return CHESS_BLACK
}
