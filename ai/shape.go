package ai

import (
	"regexp"
)

var patterns = map[string]*regexp.Regexp{
	"five":       regexp.MustCompile("11111"),
	"blockFive":  regexp.MustCompile("211111|111112"),
	"four":       regexp.MustCompile("011110"),
	"blockFour":  regexp.MustCompile("10111|11011|11101|211110|211101|211011|210111|011112|101112|110112|111012"),
	"three":      regexp.MustCompile("011100|011010|010110|001110"),
	"blockThree": regexp.MustCompile("211100|211010|210110|001112|010112|011012"),
	"two":        regexp.MustCompile("001100|011000|000110|010100|001010"),
}

var shapes = map[string]int{
	"FIVE":        5,
	"BLOCK_FIVE":  50,
	"FOUR":        4,
	"FOUR_FOUR":   44, // 双冲四
	"FOUR_THREE":  43, // 冲四活三
	"THREE_THREE": 33, // 双三
	"BLOCK_FOUR":  40,
	"THREE":       3,
	"BLOCK_THREE": 30,
	"TWO_TWO":     22, // 双活二
	"TWO":         2,
	"NONE":        0,
}

var performance = map[string]int{
	"five":       0,
	"blockFive":  0,
	"four":       0,
	"blockFour":  0,
	"three":      0,
	"blockThree": 0,
	"two":        0,
	"none":       0,
	"total":      0,
}

// getShape function
func getShape(board [][]int, x, y, offsetX, offsetY, role int) (int, int, int, int) {
	opponent := -role
	emptyCount := 0
	selfCount := 1
	opponentCount := 0
	shape := shapes["NONE"]

	// 跳过为空的节点
	if board[x+offsetX+1][y+offsetY+1] == 0 &&
		board[x-offsetX+1][y-offsetY+1] == 0 &&
		board[x+2*offsetX+1][y+2*offsetY+1] == 0 &&
		board[x-2*offsetX+1][y-2*offsetY+1] == 0 {
		return shapes["NONE"], selfCount, opponentCount, emptyCount
	}

	// two 类型占比超过一半，做一下优化
	// 活二是不需要判断特别严谨的
	for i := -3; i <= 3; i++ {
		if i == 0 {
			continue
		}
		nx, ny := x+i*offsetX+1, y+i*offsetY+1
		if nx < 0 || ny < 0 || nx >= len(board) || ny >= len(board[0]) {
			continue
		}
		currentRole := board[nx][ny]
		if currentRole == 2 {
			opponentCount++
		} else if currentRole == role {
			selfCount++
		} else if currentRole == 0 {
			emptyCount++
		}
	}
	if selfCount == 2 {
		if opponentCount == 0 {
			return shapes["TWO"], selfCount, opponentCount, emptyCount
		} else {
			return shapes["NONE"], selfCount, opponentCount, emptyCount
		}
	}

	emptyCount = 0
	selfCount = 1
	opponentCount = 0
	resultString := "1"

	for i := 1; i <= 5; i++ {
		nx, ny := x+i*offsetX+1, y+i*offsetY+1
		if nx < 0 || ny < 0 || nx >= len(board) || ny >= len(board[0]) {
			break
		}
		currentRole := board[nx][ny]
		if currentRole == 2 {
			resultString += "2"
		} else if currentRole == 0 {
			resultString += "0"
		} else {
			if currentRole == role {
				resultString += "1"
			} else {
				resultString += "2"
			}
		}
		if currentRole == 2 || currentRole == opponent {
			opponentCount++
			break
		}
		if currentRole == 0 {
			emptyCount++
		}
		if currentRole == role {
			selfCount++
		}
	}

	for i := 1; i <= 5; i++ {
		nx, ny := x-i*offsetX+1, y-i*offsetY+1
		if nx < 0 || ny < 0 || nx >= len(board) || ny >= len(board[0]) {
			break
		}
		currentRole := board[nx][ny]
		if currentRole == 2 {
			resultString = "2" + resultString
		} else if currentRole == 0 {
			resultString = "0" + resultString
		} else {
			if currentRole == role {
				resultString = "1" + resultString
			} else {
				resultString = "2" + resultString
			}
		}
		if currentRole == 2 || currentRole == opponent {
			opponentCount++
			break
		}
		if currentRole == 0 {
			emptyCount++
		}
		if currentRole == role {
			selfCount++
		}
	}

	if patterns["five"].MatchString(resultString) {
		shape = shapes["FIVE"]
		performance["five"]++
		performance["total"]++
	} else if patterns["four"].MatchString(resultString) {
		shape = shapes["FOUR"]
		performance["four"]++
		performance["total"]++
	} else if patterns["blockFour"].MatchString(resultString) {
		shape = shapes["BLOCK_FOUR"]
		performance["blockFour"]++
		performance["total"]++
	} else if patterns["three"].MatchString(resultString) {
		shape = shapes["THREE"]
		performance["three"]++
		performance["total"]++
	} else if patterns["blockThree"].MatchString(resultString) {
		shape = shapes["BLOCK_THREE"]
		performance["blockThree"]++
		performance["total"]++
	} else if patterns["two"].MatchString(resultString) {
		shape = shapes["TWO"]
		performance["two"]++
		performance["total"]++
	}

	if selfCount <= 1 || len(resultString) < 5 {
		return shape, selfCount, opponentCount, emptyCount
	}

	return shape, selfCount, opponentCount, emptyCount
}

// countShape function
func countShape(board [][]int, x, y, offsetX, offsetY, role int) (int, int, int, int, int, int) {
	opponent := -role

	innerEmptyCount := 0 // 棋子中间的内部空位
	tempEmptyCount := 0
	selfCount := 0
	totalLength := 0

	sideEmptyCount := 0 // 边上的空位
	noEmptySelfCount := 0
	OneEmptySelfCount := 0

	// right
	for i := 1; i <= 5; i++ {
		nx, ny := x+i*offsetX+1, y+i*offsetY+1
		if nx < 0 || ny < 0 || nx >= len(board) || ny >= len(board[0]) {
			break
		}
		currentRole := board[nx][ny]
		if currentRole == 2 || currentRole == opponent {
			break
		}
		if currentRole == role {
			selfCount++
			sideEmptyCount = 0
			if tempEmptyCount > 0 {
				innerEmptyCount += tempEmptyCount
				tempEmptyCount = 0
			}
			if innerEmptyCount == 0 {
				noEmptySelfCount++
				OneEmptySelfCount++
			} else if innerEmptyCount == 1 {
				OneEmptySelfCount++
			}
		}
		totalLength++
		if currentRole == 0 {
			tempEmptyCount++
			sideEmptyCount++
		}
		if sideEmptyCount >= 2 {
			break
		}
	}
	if innerEmptyCount == 0 {
		OneEmptySelfCount = 0
	}
	return selfCount, totalLength, noEmptySelfCount, OneEmptySelfCount, innerEmptyCount, sideEmptyCount
}

// getShapeFast function
func getShapeFast(board [][]int, x, y, offsetX, offsetY, role int) (int, int) {
	if board[x+offsetX+1][y+offsetY+1] == 0 &&
		board[x-offsetX+1][y-offsetY+1] == 0 &&
		board[x+2*offsetX+1][y+2*offsetY+1] == 0 &&
		board[x-2*offsetX+1][y-2*offsetY+1] == 0 {
		return shapes["NONE"], 1
	}

	selfCount := 1
	totalLength := 1
	shape := shapes["NONE"]

	leftEmpty := 0
	rightEmpty := 0
	noEmptySelfCount := 1
	OneEmptySelfCount := 1

	lSelfCount, lTotalLength, lNoEmptySelfCount, lOneEmptySelfCount, _, lSideEmptyCount := countShape(board, x, y, -offsetX, -offsetY, role)
	rSelfCount, rTotalLength, rNoEmptySelfCount, rOneEmptySelfCount, _, rSideEmptyCount := countShape(board, x, y, offsetX, offsetY, role)

	selfCount = lSelfCount + rSelfCount + 1
	totalLength = lTotalLength + rTotalLength + 1
	noEmptySelfCount = lNoEmptySelfCount + rNoEmptySelfCount + 1
	OneEmptySelfCount = max(lOneEmptySelfCount+rNoEmptySelfCount, lNoEmptySelfCount+rOneEmptySelfCount) + 1
	rightEmpty = rSideEmptyCount
	leftEmpty = lSideEmptyCount

	if totalLength < 5 {
		return shape, selfCount
	}
	if noEmptySelfCount >= 5 {
		if rightEmpty > 0 && leftEmpty > 0 {
			return shapes["FIVE"], selfCount
		} else {
			return shapes["BLOCK_FIVE"], selfCount
		}
	}
	if noEmptySelfCount == 4 {
		if (rightEmpty >= 1 || rOneEmptySelfCount > rNoEmptySelfCount) && (leftEmpty >= 1 || lOneEmptySelfCount > lNoEmptySelfCount) {
			return shapes["FOUR"], selfCount
		} else if !(rightEmpty == 0 && leftEmpty == 0) {
			return shapes["BLOCK_FOUR"], selfCount
		}
	}
	if OneEmptySelfCount == 4 {
		return shapes["BLOCK_FOUR"], selfCount
	}
	if noEmptySelfCount == 3 {
		if (rightEmpty >= 2 && leftEmpty >= 1) || (rightEmpty >= 1 && leftEmpty >= 2) {
			return shapes["THREE"], selfCount
		} else {
			return shapes["BLOCK_THREE"], selfCount
		}
	}
	if OneEmptySelfCount == 3 {
		if rightEmpty >= 1 && leftEmpty >= 1 {
			return shapes["THREE"], selfCount
		} else {
			return shapes["BLOCK_THREE"], selfCount
		}
	}
	if (noEmptySelfCount == 2 || OneEmptySelfCount == 2) && totalLength > 5 {
		shape = shapes["TWO"]
	}

	return shape, selfCount
}

// isFive function
func isFive(shape int) bool {
	return shape == shapes["FIVE"] || shape == shapes["BLOCK_FIVE"]
}

// isFour function
func isFour(shape int) bool {
	return shape == shapes["FOUR"] || shape == shapes["BLOCK_FOUR"]
}

// getAllShapesOfPoint function
func getAllShapesOfPoint(shapeCache map[int]map[int]map[int]map[int]int, x, y, role int) []int {
	roles := []int{role}
	if role == 0 {
		roles = []int{1, -1}
	}
	var result []int
	for _, r := range roles {
		for d := 0; d <= 3; d++ {
			shape := shapeCache[r][d][x][y]
			if shape > 0 {
				result = append(result, shape)
			}
		}
	}
	return result
}
