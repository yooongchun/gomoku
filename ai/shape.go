package ai

import "reflect"

type TypeChess int
type TypeRole int
type TypeShape int
type TypeDirection int
type TypeShapeCache map[TypeChess]map[TypeDirection]map[int]map[int]TypeShape
type TypePointCache map[TypeChess]map[TypeShape]map[int]bool

type Point struct {
	x int
	y int
}

type TypeHistory struct {
	point Point
	chess TypeChess
}

type TypeEvaluateCache struct {
	chess TypeChess
	score int
}

type TypeValuableMoveCache struct {
	role      TypeChess
	moves     []Point
	depth     int
	onlyThree bool
	onlyFour  bool
}

// ShapeEnum 可取的形状
type ShapeEnum struct {
	FIVE        TypeShape //11111
	FOUR        TypeShape //011110
	FOUR_FOUR   TypeShape
	FOUR_THREE  TypeShape
	THREE_THREE TypeShape
	BLOCK_FOUR  TypeShape //10111|11011|11101|211110|211101|211011|210111|011112|101112|110112|111012
	THREE       TypeShape //011100|011010|010110|001110
	BLOCK_THREE TypeShape //211100|211010|210110|001112|010112|011012
	TWO_TWO     TypeShape
	TWO         TypeShape //001100|011000|000110|010100|001010
	BLOCK_TWO   TypeShape
	ONE         TypeShape
	BLOCK_ONE   TypeShape
	NONE        TypeShape
}

var DirectionVec = []Point{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
var DirectionEnum = []TypeDirection{HORIZONTAL, VERTICAL, DIAGONAL, ANTI_DIAGONAL}
var Shapes = &ShapeEnum{
	FIVE:        5,
	FOUR:        4,
	FOUR_FOUR:   44,
	FOUR_THREE:  43,
	THREE_THREE: 33,
	BLOCK_FOUR:  40,
	THREE:       3,
	BLOCK_THREE: 30,
	TWO_TWO:     22,
	TWO:         2,
	BLOCK_TWO:   20,
	ONE:         1,
	BLOCK_ONE:   10,
	NONE:        0,
}
var ShapeFields []TypeShape

func init() {
	// 初始化 ShapeFields
	ShapeFields = make([]TypeShape, 0)
	v := reflect.ValueOf(*Shapes)
	for i := 0; i < v.NumField(); i++ {
		ShapeFields = append(ShapeFields, v.Field(i).Interface().(TypeShape))
	}
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

// GetShapeFast 使用遍历位置的方式实现的形状检测，速度较快，大约是字符串速度的2倍 但理解起来会稍微复杂一些
func GetShapeFast(board [][]TypeChess, x, y, offsetX, offsetY int, role TypeChess) (TypeShape, int) {
	// 有一点点优化效果：跳过为空的节点（左右两边空位为2）
	if board[x+offsetX][y+offsetY] == 0 &&
		board[x-offsetX][y-offsetY] == 0 &&
		board[x+2*offsetX][y+2*offsetY] == 0 &&
		board[x-2*offsetX][y-2*offsetY] == 0 {
		return Shapes.NONE, 1
	}

	selfCount := 1
	totalLength := 1
	shape := Shapes.NONE

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
		return Shapes.FIVE, selfCount
	}
	if noEmptySelfCount == 4 {
		if (rightEmpty >= 1 || rOneEmptySelfCount > rNoEmptySelfCount) && (leftEmpty >= 1 || lOneEmptySelfCount > lNoEmptySelfCount) {
			return Shapes.FOUR, selfCount
		} else if !(rightEmpty == 0 && leftEmpty == 0) {
			return Shapes.BLOCK_FOUR, selfCount
		}
	}
	if OneEmptySelfCount == 4 {
		return Shapes.BLOCK_FOUR, selfCount
	}
	if noEmptySelfCount == 3 {
		if (rightEmpty >= 2 && leftEmpty >= 1) || (rightEmpty >= 1 && leftEmpty >= 2) {
			return Shapes.THREE, selfCount
		} else {
			return Shapes.BLOCK_THREE, selfCount
		}
	}
	if OneEmptySelfCount == 3 {
		if rightEmpty >= 1 && leftEmpty >= 1 {
			return Shapes.THREE, selfCount
		} else {
			return Shapes.BLOCK_THREE, selfCount
		}
	}
	if (noEmptySelfCount == 2 || OneEmptySelfCount == 2) && totalLength > 5 {
		shape = Shapes.TWO
	}

	return shape, selfCount
}

// IsFive function
func IsFive(shape TypeShape) bool {
	return shape == Shapes.FIVE
}

// IsFour function
func IsFour(shape TypeShape) bool {
	return shape == Shapes.FOUR || shape == Shapes.BLOCK_FOUR
}

// GetAllShapesOfPoint function
func GetAllShapesOfPoint(shapeCache TypeShapeCache, x, y int, chess TypeChess) []TypeShape {
	chesses := []TypeChess{chess}
	if chess == CHESS_EMPTY {
		chesses = []TypeChess{CHESS_BLACK, CHESS_WHITE}
	}
	var result []TypeShape
	for _, r := range chesses {
		for _, d := range []TypeDirection{HORIZONTAL, VERTICAL, DIAGONAL, ANTI_DIAGONAL} {
			shape := shapeCache[r][d][x][y]
			if shape != Shapes.NONE {
				result = append(result, shape)
			}
		}
	}
	return result
}
