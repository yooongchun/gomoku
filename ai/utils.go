package ai

import (
	"runtime"
	"strconv"
	"strings"
)

func getBoardString(board [][]TypeChess, lastOp *Point, extraPoints ...Point) string {
	size := len(board)
	extraPositions := make(map[int]bool, len(extraPoints))
	for _, point := range extraPoints {
		extraPositions[Coordinate2Position(point.x, point.y, size)] = true
	}
	var result strings.Builder
	prefix := "  "
	result.WriteString(prefix + " ")
	for i := 0; i < size; i++ {
		result.WriteString(strings.ToUpper(strconv.FormatInt(int64(i), 32)) + prefix)
	}
	result.WriteString("\n")

	for i := 0; i < size; i++ {
		result.WriteString(strings.ToUpper(strconv.FormatInt(int64(i), 32)) + prefix)
		for j := 0; j < size; j++ {
			position := Coordinate2Position(i, j, size)
			if ok, exist := extraPositions[position]; ok && exist {
				result.WriteString("?" + prefix)
				continue
			}

			op := ""
			switch board[i][j] {
			case CHESS_BLACK:
				op = "X"
			case CHESS_WHITE:
				op = "O"
			case CHESS_OBSTACLE:
				op = "#"
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
	} else if chess == CHESS_WHITE {
		return CHESS_BLACK
	}
	return chess
}

// chess2str 将棋子转换为字符串 跟const中定义的数字有关，不可随意修改
func chess2str(chess, target TypeChess) string {
	if chess == CHESS_EMPTY {
		return "0"
	} else if chess == target {
		return "1"
	}
	return "2"
}

func ifPresent[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func ifPresentOrNil[T any](condition bool, trueValue T) *T {
	if condition {
		return &trueValue
	}
	return nil
}

func getFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
