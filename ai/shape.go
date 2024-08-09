package ai

import (
	"github.com/duke-git/lancet/v2/strutil"
	"reflect"
	"strconv"
	"strings"
)

func NewPoint(x, y int) Point {
	return Point{x, y}
}

// countShape function
func countShape(board [][]TypeChess, x, y, offsetX, offsetY int, role TypeChess) (int, int, int, int, int, int) {
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
		nx, ny := x+i*offsetX, y+i*offsetY
		if nx < 0 || ny < 0 || nx >= len(board) || ny >= len(board[0]) {
			break
		}
		currentRole := board[nx][ny]
		if currentRole == CHESS_OBSTACLE || currentRole == opponent {
			break
		}
		totalLength++
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
		} else if currentRole == CHESS_EMPTY {
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

func getShapeFields() []TypeShapeField {
	allFields := make([]TypeShapeField, 0)
	v := reflect.ValueOf(*ShapeEnum)
	for i := 0; i < v.NumField(); i++ {
		allFields = append(allFields, v.Field(i).Interface().(TypeShapeField))
	}
	return allFields
}

func toggleChess(chess TypeChess) TypeChess {
	if chess == CHESS_BLACK {
		return CHESS_WHITE
	}
	return CHESS_BLACK
}

func transBoardWithObstacle(board [][]TypeChess, chess TypeChess) [][]TypeChess {
	newBoard := make([][]TypeChess, len(board))
	for i := 0; i < len(board); i++ {
		newBoard[i] = make([]TypeChess, len(board[i]))
		for j := 0; j < len(board[i]); j++ {
			newBoard[i][j] = board[i][j]
			if board[i][j] == toggleChess(chess) {
				newBoard[i][j] = CHESS_OBSTACLE
			}
		}
	}
	return newBoard
}

// GetShape 字符串匹配方式实现形状检测
func GetShape(board [][]TypeChess, pos Point, dirVec Point, chess TypeChess) (shape TypeShapeField) {
	// 将chess之外的棋子转换为OBSTACLE
	board = transBoardWithObstacle(board, chess)
	// 取出指定方向上的字符串，前后各延伸4个位置
	line := []string{4: strconv.Itoa(int(chess)), 8: ""}
	for i := 1; i < 5; i++ {
		x, y := pos.x+i*dirVec.x, pos.y+i*dirVec.y
		if isPosValid(board, x, y) {
			line[4+i] = strconv.Itoa(int(board[x][y]))
		}
		x, y = pos.x-i*dirVec.x, pos.y-i*dirVec.y
		if isPosValid(board, x, y) {
			line[4-i] = strconv.Itoa(int(board[x][y]))
		}
	}
	// 目标字符串
	str := strings.Join(line, "")
	// 遍历所有形状，查找匹配的形状
	shape = ShapeEnum.None
	for _, field := range getShapeFields() {
		if field.Name == nil {
			continue
		}
		for _, name := range field.Name {
			if strings.Contains(str, name) {
				shape = field
				return
			}
			// 反转字符串再匹配一次
			if strings.Contains(strutil.Reverse(str), name) {
				shape = field
				return
			}
		}
	}
	return
}

func isPosValid(board [][]TypeChess, x, y int) bool {
	return x >= 0 && y >= 0 && x < len(board) && y < len(board[0])
}

// GetShapeFast 使用遍历位置的方式实现的形状检测，速度较快，大约是字符串速度的2倍 但理解起来会稍微复杂一些
func GetShapeFast(board [][]TypeChess, x, y, offsetX, offsetY int, role TypeChess) (TypeShapeField, int) {
	// 有一点点优化效果：跳过为空的节点（左右两边空位为2）
	if board[x+offsetX][y+offsetY] == 0 &&
		board[x-offsetX][y-offsetY] == 0 &&
		board[x+2*offsetX][y+2*offsetY] == 0 &&
		board[x-2*offsetX][y-2*offsetY] == 0 {
		return ShapeEnum.None, 1
	}

	selfCount := 1
	totalLength := 1
	shape := ShapeEnum.None

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
		return ShapeEnum.Five, selfCount
	}
	if noEmptySelfCount == 4 {
		if (rightEmpty >= 1 || rOneEmptySelfCount > rNoEmptySelfCount) && (leftEmpty >= 1 || lOneEmptySelfCount > lNoEmptySelfCount) {
			return ShapeEnum.Four, selfCount
		} else if !(rightEmpty == 0 && leftEmpty == 0) {
			return ShapeEnum.BlockFour, selfCount
		}
	}
	if OneEmptySelfCount == 4 {
		return ShapeEnum.BlockFour, selfCount
	}
	if noEmptySelfCount == 3 {
		if (rightEmpty >= 2 && leftEmpty >= 1) || (rightEmpty >= 1 && leftEmpty >= 2) {
			return ShapeEnum.Three, selfCount
		} else {
			return ShapeEnum.BlockThree, selfCount
		}
	}
	if OneEmptySelfCount == 3 {
		if rightEmpty >= 1 && leftEmpty >= 1 {
			return ShapeEnum.Three, selfCount
		} else {
			return ShapeEnum.BlockThree, selfCount
		}
	}
	if (noEmptySelfCount == 2 || OneEmptySelfCount == 2) && totalLength > 5 {
		shape = ShapeEnum.Two
	}

	return shape, selfCount
}

// IsFive function
func IsFive(shape TypeShapeField) bool {
	return shape.Code == ShapeEnum.Five.Code
}

// IsFour function
func IsFour(shape TypeShapeField) bool {
	return shape.Code == ShapeEnum.Four.Code || shape.Code == ShapeEnum.BlockFour.Code
}

// GetAllShapesOfPoint function
func GetAllShapesOfPoint(shapeCache TypeShapeCache, x, y int, chess TypeChess) []TypeShapeField {
	chessList := []TypeChess{chess}
	if chess == CHESS_EMPTY {
		chessList = []TypeChess{CHESS_BLACK, CHESS_WHITE}
	}
	var result []TypeShapeField
	for _, c := range chessList {
		for _, d := range DirectionEnum {
			shape := shapeCache[c][d][x][y]
			if shape.Code != ShapeEnum.None.Code {
				result = append(result, shape)
			}
		}
	}
	return result
}
