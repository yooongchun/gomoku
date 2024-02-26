package ai

import "reflect"

type TypeRole int
type TypeShape int
type TypeDirection int
type TypeShapeCache map[TypeRole]map[TypeDirection]map[int]map[int]TypeShape
type TypePointCache map[TypeRole]map[TypeShape]map[int]bool

type TypeHistory struct {
	position int
	role     TypeRole
}

// ShapeEnum 可取的形状
type ShapeEnum struct {
	LiveFive   TypeShape //11111
	LiveFour   TypeShape //011110
	FourFour   TypeShape
	FourThree  TypeShape
	ThreeThree TypeShape
	BlockFour  TypeShape //10111|11011|11101|211110|211101|211011|210111|011112|101112|110112|111012
	LiveThree  TypeShape //011100|011010|010110|001110
	BlockThree TypeShape //211100|211010|210110|001112|010112|011012
	TwoTwo     TypeShape
	LiveTwo    TypeShape //001100|011000|000110|010100|001010
	BlockTwo   TypeShape
	LiveOne    TypeShape
	BlockOne   TypeShape
	None       TypeShape
}

type Point struct {
	x int
	y int
}

type Vector Point

type ChessEnum struct {
	WHITE    TypeRole
	BLACK    TypeRole
	EMPTY    TypeRole
	OBSTACLE TypeRole
}

var Chess = &ChessEnum{
	BLACK:    BLACK,
	WHITE:    WHITE,
	EMPTY:    EMPTY,
	OBSTACLE: OBSTACLE,
}
var Roles = []TypeRole{Chess.BLACK, Chess.WHITE}
var Directions = []TypeDirection{HORIZONTAL, VERTICAL, DIAGONAL, ANTI_DIAGONAL}
var Shapes = &ShapeEnum{
	LiveFive:   5,
	LiveFour:   4,
	FourFour:   44,
	FourThree:  43,
	ThreeThree: 33,
	BlockFour:  40,
	LiveThree:  3,
	BlockThree: 30,
	TwoTwo:     22,
	LiveTwo:    2,
	BlockTwo:   20,
	LiveOne:    1,
	BlockOne:   10,
	None:       0,
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
func countShape(board [][]TypeRole, x, y, offsetX, offsetY int, role TypeRole) (int, int, int, int, int, int) {
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
		if currentRole == Chess.OBSTACLE || currentRole == opponent {
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
		} else if currentRole == Chess.EMPTY {
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
func GetShapeFast(board [][]TypeRole, x, y, offsetX, offsetY int, role TypeRole) (TypeShape, int) {
	// 有一点点优化效果：跳过为空的节点（左右两边空位为2）
	if board[x+offsetX][y+offsetY] == 0 &&
		board[x-offsetX][y-offsetY] == 0 &&
		board[x+2*offsetX][y+2*offsetY] == 0 &&
		board[x-2*offsetX][y-2*offsetY] == 0 {
		return Shapes.None, 1
	}

	selfCount := 1
	totalLength := 1
	shape := Shapes.None

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
		return Shapes.LiveFive, selfCount
	}
	if noEmptySelfCount == 4 {
		if (rightEmpty >= 1 || rOneEmptySelfCount > rNoEmptySelfCount) && (leftEmpty >= 1 || lOneEmptySelfCount > lNoEmptySelfCount) {
			return Shapes.LiveFour, selfCount
		} else if !(rightEmpty == 0 && leftEmpty == 0) {
			return Shapes.BlockFour, selfCount
		}
	}
	if OneEmptySelfCount == 4 {
		return Shapes.BlockFour, selfCount
	}
	if noEmptySelfCount == 3 {
		if (rightEmpty >= 2 && leftEmpty >= 1) || (rightEmpty >= 1 && leftEmpty >= 2) {
			return Shapes.LiveThree, selfCount
		} else {
			return Shapes.BlockThree, selfCount
		}
	}
	if OneEmptySelfCount == 3 {
		if rightEmpty >= 1 && leftEmpty >= 1 {
			return Shapes.LiveThree, selfCount
		} else {
			return Shapes.BlockThree, selfCount
		}
	}
	if (noEmptySelfCount == 2 || OneEmptySelfCount == 2) && totalLength > 5 {
		shape = Shapes.LiveTwo
	}

	return shape, selfCount
}

// IsFive function
func IsFive(shape TypeShape) bool {
	return shape == Shapes.LiveFive
}

// IsFour function
func IsFour(shape TypeShape) bool {
	return shape == Shapes.LiveFour || shape == Shapes.BlockFour
}

// GetAllShapesOfPoint function
func GetAllShapesOfPoint(shapeCache TypeShapeCache, x, y int, role TypeRole) []TypeShape {
	roles := []TypeRole{role}
	if role == Chess.EMPTY {
		roles = Roles
	}
	var result []TypeShape
	for _, r := range roles {
		for _, d := range []TypeDirection{HORIZONTAL, VERTICAL, DIAGONAL, ANTI_DIAGONAL} {
			shape := shapeCache[r][d][x][y]
			if shape != Shapes.None {
				result = append(result, shape)
			}
		}
	}
	return result
}
