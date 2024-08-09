package ai

import (
	"fmt"
	"time"
)

type PerformanceEnum struct {
	updateTime    time.Duration
	getPointsTime time.Duration
}

var performance = &PerformanceEnum{
	updateTime:    time.Duration(0),
	getPointsTime: time.Duration(0),
}

type Evaluate struct {
	size       int
	board      [][]TypeChess
	scores     map[TypeChess][][]int // [chess][x][y]->score
	history    []TypeHistory         // 记录历史 [position, role]
	scoreCache TypeScoreCache        // 缓存每个点位的分数，避免重复计算
	shapeCache TypeShapeCache        // 缓存每个形状对应的点位
}

func NewEvaluate(size int) *Evaluate {
	// 初始化board，带边框
	board := make([][]TypeChess, size+2)
	for i := 0; i < size+2; i++ {
		board[i] = make([]TypeChess, size+2)
		for j := 0; j < size+2; j++ {
			if i == 0 || j == 0 || i == size+1 || j == size+1 {
				board[i][j] = CHESS_OBSTACLE
			}
		}
	}
	// 黑白双方的得分
	scores := make(map[TypeChess][][]int, 2)
	for _, chess := range []TypeChess{CHESS_BLACK, CHESS_WHITE} {
		scores[chess] = make([][]int, size)
		for i := 0; i < size; i++ {
			scores[chess][i] = make([]int, size)
		}
	}
	// 缓存每个点位的形状，避免重复计算
	shapeCache := make(TypeShapeCache, 2)
	numDirection := len(DirectionEnum)
	for _, chess := range []TypeChess{CHESS_BLACK, CHESS_WHITE} {
		shapeCache[chess] = make(map[TypeDirection][][]TypeShapeField, numDirection)
		for _, direction := range DirectionEnum {
			shapeCache[chess][direction] = make([][]TypeShapeField, size)
			for i := 0; i < size; i++ {
				shapeCache[chess][direction][i] = make([]TypeShapeField, size)
				for j := 0; j < size; j++ {
					shapeCache[chess][direction][i][j] = ShapeEnum.None
				}
			}
		}
	}
	// 缓存score
	scoreCache := make(TypeScoreCache, 2)
	for _, chess := range []TypeChess{CHESS_BLACK, CHESS_WHITE} {
		scoreCache[chess] = make(map[TypeDirection][][]int, numDirection)
		for _, direction := range DirectionEnum {
			scoreCache[chess][direction] = make([][]int, size)
			for i := 0; i < size; i++ {
				scoreCache[chess][direction][i] = make([]int, size)
			}
		}
	}
	return &Evaluate{
		size:       size,
		board:      board,
		scores:     scores,
		history:    make([]TypeHistory, 0),
		scoreCache: scoreCache,
		shapeCache: shapeCache,
	}
}

func (e *Evaluate) opponent(chess TypeChess) TypeChess {
	if chess == CHESS_BLACK {
		return CHESS_WHITE
	}
	return CHESS_BLACK
}

func (e *Evaluate) Move(point Point, chess TypeChess) {
	// 清空记录
	x, y := point.x, point.y
	for _, d := range DirectionEnum {
		e.scoreCache[chess][d][x][y] = SCORE_NONE
		e.scoreCache[e.opponent(chess)][d][x][y] = SCORE_NONE
	}
	e.scores[CHESS_BLACK][x][y] = SCORE_NONE
	e.scores[CHESS_WHITE][x][y] = SCORE_NONE

	// 更新分数
	e.board[x+1][y+1] = chess
	e.updatePoint(x, y)
	e.history = append(e.history, TypeHistory{point, chess})
}

func (e *Evaluate) Undo(point Point) {
	x, y := point.x, point.y
	e.board[x+1][y+1] = CHESS_EMPTY // Adjust for the added wall
	e.updatePoint(x, y)
	e.history = e.history[:len(e.history)-1]
}

// 只返回和最后几步在一条直线上的点位。
// 这么做有一点问题：
// 1. 因为己方可能会由于防守暂时离开原来的线，这样就会导致己方被中断，只能增加最后几步的长度，比如不是取最后一步，而是最后3步
// 2. 如果不是取最后1步，取的步数太多了，反而还不如直接返回所有点位。
func (e *Evaluate) getPointInLine() map[int]map[int]bool {
	pointsInLine := make(map[int]map[int]bool)
	hasPointInLine := false
	for _, shape := range getShapeFields() {
		pointsInLine[shape.Code] = make(map[int]bool)
	}
	last2History := e.history[len(e.history)-Config.InlineCount:]
	processed := make(map[int]TypeChess) // 已经处理过的点位
	// 在last2Points中查找是否有点位在一条线上
	for _, chess := range []TypeChess{CHESS_BLACK, CHESS_WHITE} {
		for _, his := range last2History {
			for _, dire := range DirectionEnum {
				vec := DirectionVec[dire]
				for _, sign := range []int{1, -1} {
					for step := 1; step <= Config.InLineDistance; step++ {
						nx, ny := his.point.x+sign*step*vec.x, his.point.y+sign*step*vec.y
						position := Coordinate2Position(nx, ny, e.size)
						// 检测是否到达边界
						if nx < 0 || nx >= e.size || ny < 0 || ny >= e.size {
							break
						}
						if e.board[nx+1][ny+1] != CHESS_EMPTY {
							continue
						}
						if processed[position] == chess {
							continue
						}
						processed[position] = chess
						for _, d := range DirectionEnum {
							shape := e.shapeCache[chess][d][nx][ny]
							if shape.Code != ShapeEnum.None.Code {
								pointsInLine[shape.Code][position] = true
								hasPointInLine = true
							}
						}
					}
				}
			}
		}
	}
	if hasPointInLine {
		return pointsInLine
	}
	return nil
}

func (e *Evaluate) getPoints(role TypeChess, depth int, vct, vcf bool) map[int]map[int]bool {
	var first TypeChess
	if depth%2 == 0 {
		first = role
	} else {
		first = role * -1
	}
	start := time.Now()
	if Config.OnlyInLine && len(e.history) >= Config.InlineCount {
		pointsInLine := e.getPointInLine()
		if pointsInLine != nil {
			performance.getPointsTime += time.Since(start)
			return pointsInLine
		}
	}

	points := make(map[int]map[int]bool)
	for _, shape := range getShapeFields() {
		points[shape.Code] = make(map[int]bool)
	}
	lastPoints := e.history[len(e.history)-4:]
	for _, r := range []TypeChess{CHESS_BLACK, CHESS_WHITE} {
		for i := 0; i < e.size; i++ {
			for j := 0; j < e.size; j++ {
				fourCount, blockFourCount, threeCount := 0, 0, 0
				for _, direction := range DirectionEnum {
					if e.board[i+1][j+1] != CHESS_EMPTY {
						continue
					}
					shape := e.shapeCache[r][direction][i][j]
					if shape.Code != ShapeEnum.None.Code {
						continue
					}
					if vcf {
						if r == first && !IsFour(shape) && !IsFive(shape) {
							continue
						}
						if r == -first && IsFive(shape) {
							continue
						}
					}
					point := i*e.size + j
					if vct {
						if depth%2 == 0 {
							if depth == 0 && r != first {
								continue
							}
							if shape.Code != ShapeEnum.Three.Code && !IsFour(shape) && !IsFive(shape) {
								continue
							}
							if shape.Code == ShapeEnum.Three.Code && r != first {
								continue
							}
							if depth == 0 && r != first {
								continue
							}
							if depth > 0 {
								if shape.Code == ShapeEnum.Three.Code && len(GetAllShapesOfPoint(e.shapeCache, i, j, r)) == 1 {
									continue
								}
								if shape.Code == ShapeEnum.BlockFour.Code && len(GetAllShapesOfPoint(e.shapeCache, i, j, r)) == 1 {
									continue
								}
							}
						} else {
							if shape.Code != ShapeEnum.Three.Code && !IsFour(shape) && !IsFive(shape) {
								continue
							}
							if shape.Code == ShapeEnum.Three.Code && r == -first {
								continue
							}
							if depth > 1 {
								if shape.Code == ShapeEnum.BlockFour.Code && len(GetAllShapesOfPoint(e.shapeCache, i, j, CHESS_EMPTY)) == 1 {
									continue
								}
								if shape.Code == ShapeEnum.BlockFour.Code && !HasInLine(point, lastPoints, e.size) {
									continue
								}
							}
						}
					}
					if vcf {
						if !IsFour(shape) && !IsFive(shape) {
							continue
						}
					}
					if depth > 2 && (shape.Code == ShapeEnum.Two.Code || shape.Code == ShapeEnum.DoubleTwo.Code || shape.Code == ShapeEnum.DoubleThree.Code) && !HasInLine(point, lastPoints, e.size) {
						continue
					}
					points[shape.Code][point] = true
					if shape.Code == ShapeEnum.Four.Code {
						fourCount++
					} else if shape.Code == ShapeEnum.BlockFour.Code {
						blockFourCount++
					} else if shape.Code == ShapeEnum.Three.Code {
						threeCount++
					}
					var unionShape TypeShapeField
					if fourCount >= 2 {
						unionShape = ShapeEnum.DoubleFour
					} else if blockFourCount > 0 && threeCount > 0 {
						unionShape = ShapeEnum.FourThree
					} else if threeCount >= 2 {
						unionShape = ShapeEnum.DoubleThree
					}
					if unionShape.Code != ShapeEnum.None.Code {
						points[unionShape.Code][point] = true
					}
				}
			}
		}
	}
	performance.getPointsTime += time.Since(start)
	return points
}

// 当一个位置发生变时候，要更新这个位置的四个方向上得分，更新规则是：
// 1. 如果这个位置是空的，那么就重新计算这个位置的得分
// 2. 如果碰到了边界或者对方的棋子，那么就停止计算
// 3. 如果超过2个空位，那么就停止计算
// 4. 要更新自己的和对方的得分
func (e *Evaluate) updatePoint(x, y int) {
	// 更新当前点位的分数
	start := time.Now()
	e.updateSinglePoint(x, y, CHESS_BLACK)
	e.updateSinglePoint(x, y, CHESS_WHITE)

	for _, dir := range DirectionEnum {
		vec := DirectionVec[dir]
		for _, sign := range []int{1, -1} { // -1 for negative direction, 1 for positive direction
			for step := 1; step <= 5; step++ {
				reachEdge := false
				for _, chess := range []TypeChess{CHESS_BLACK, CHESS_WHITE} {
					nx, ny := x+sign*step*vec.x+1, y+sign*step*vec.y+1 // +1 to adjust for wall
					if e.board[nx][ny] == CHESS_OBSTACLE {             // 到达边界停止
						reachEdge = true
						break
					} else if e.board[nx][ny] == toggleChess(chess) { // 达到对方棋子，则转换角色
						continue
					} else if e.board[nx][ny] == CHESS_EMPTY {
						//[sign * ox, sign * oy]
						// 这里不能跳过，可能会在悔棋时漏掉一些待更新的点位
						e.updateSinglePoint(nx-1, ny-1, chess, dir) // -1 to adjust back from wall
					}
				}
				if reachEdge {
					break
				}
			}
		}
	}
	performance.updateTime += time.Since(start)
}

/*
计算单个点的得分
计算原理是：
遍历四个方向，生成四个方向上的字符串，用patters来匹配字符串, 匹配到的话，就将对应的得分加到scores上
四个方向的字符串生成规则是：向两边都延伸5个位置，如果遇到边界或者对方的棋子，就停止延伸
在更新周围棋子时，只有一个方向需要更新，因此可以传入direction参数，只更新一个方向
*/
func (e *Evaluate) updateSinglePoint(x, y int, chess TypeChess, direction ...TypeDirection) int {
	// 遍历的方向：有则传入，没有则遍历所有方向
	var directions []TypeDirection
	if direction != nil {
		directions = []TypeDirection{direction[0]}
	} else {
		directions = DirectionEnum
	}

	// 只取当前角色的缓存
	shapeCache := e.shapeCache[chess]
	// Clear cache
	for _, dir := range directions {
		shapeCache[dir][x][y] = ShapeEnum.None
	}

	score := 0
	blockFourCount := 0
	threeCount := 0
	twoCount := 0

	// Calculate existing scores
	for _, dir := range DirectionEnum {
		shape := shapeCache[dir][x][y]
		score += shape.Score
		switch shape.Code {
		case ShapeEnum.BlockFour.Code:
			blockFourCount++
		case ShapeEnum.Three.Code:
			threeCount++
		case ShapeEnum.Two.Code:
			twoCount++
		}
	}

	for _, dir := range directions {
		vec := DirectionVec[dir]
		shape := GetShape(e.board, Point{x + 1, y + 1}, vec, chess)
		fmt.Printf("direction=%v, shape=%v", dir, shape)
		if shape.Code == ShapeEnum.None.Code {
			continue
		}
		// Cache only single shapes, complex shapes like double Three are not cached
		shapeCache[dir][x][y] = shape
		switch shape.Code {
		case ShapeEnum.BlockFour.Code:
			blockFourCount++
		case ShapeEnum.Three.Code:
			threeCount++
		case ShapeEnum.Two.Code:
			twoCount++
		}
		if blockFourCount >= 2 {
			shape = ShapeEnum.DoubleFour
		} else if blockFourCount > 0 && threeCount > 0 {
			shape = ShapeEnum.FourThree
		} else if threeCount >= 2 {
			shape = ShapeEnum.DoubleThree
		} else if twoCount >= 2 {
			shape = ShapeEnum.DoubleTwo
		}
		score += shape.Score
	}
	e.scores[chess][x][y] = score
	return score
}

func (e *Evaluate) Evaluate(chess TypeChess) int {
	blackScore, whiteScore := 0, 0
	for i := 0; i < e.size; i++ {
		for j := 0; j < e.size; j++ {
			blackScore += e.scores[CHESS_BLACK][i][j]
			whiteScore += e.scores[CHESS_WHITE][i][j]
		}
	}
	if chess == CHESS_BLACK {
		return blackScore - whiteScore
	}
	return whiteScore - blackScore
}

/**
 * 获取有价值的点位
 * @param {*} role 当前角色
 * @param {*} onlyThrees 只返回 活三、冲四、活四
 * @param {*} maxCount 最多返回多少个点位，这个参数只会裁剪活三以下的点位
 * @returns
 */
func (e *Evaluate) getValuableMoves(role TypeChess, depth int, onlyThree, onlyFour bool) []Point {
	moves := e.getMoves(role, depth, onlyThree, onlyFour)
	movePacks := make([]Point, 0)
	for m := range moves {
		movePacks = append(movePacks, Point{m / e.size, m % e.size})
	}
	return movePacks
}
func (e *Evaluate) getMoves(role TypeChess, depth int, onlyThree bool, onlyFour bool) map[int]bool {
	points := e.getPoints(role, depth, onlyThree, onlyFour)
	fives, ok := points[ShapeEnum.Five.Code]
	if !ok {
		fives = make(map[int]bool)
	}
	if len(fives) > 0 {
		return fives
	}
	fours, ok := points[ShapeEnum.Four.Code]
	if !ok {
		fours = make(map[int]bool)
	}
	blockFours, ok := points[ShapeEnum.BlockFour.Code]
	if !ok {
		blockFours = make(map[int]bool)
	}
	if onlyFour || len(fours) > 0 {
		return mergeMaps(fours, blockFours)
	}
	fourFours, ok := points[ShapeEnum.DoubleFour.Code]
	if !ok {
		fourFours = make(map[int]bool)
	}
	if len(fourFours) > 0 {
		return mergeMaps(fourFours, blockFours)
	}
	threes, ok := points[ShapeEnum.Three.Code]
	if !ok {
		threes = make(map[int]bool)
	}
	fourThrees, ok := points[ShapeEnum.FourThree.Code]
	if !ok {
		fourThrees = make(map[int]bool)
	}
	if len(fourThrees) > 0 {
		return mergeMaps(fourThrees, blockFours, threes)
	}
	threeThrees, ok := points[ShapeEnum.DoubleThree.Code]
	if !ok {
		threeThrees = make(map[int]bool)
	}
	if len(threeThrees) > 0 {
		return mergeMaps(threeThrees, blockFours, threes)
	}
	if onlyThree {
		return mergeMaps(blockFours, threes)
	}
	blockThrees, ok := points[ShapeEnum.DoubleThree.Code]
	if !ok {
		blockThrees = make(map[int]bool)
	}
	twoTwos, ok := points[ShapeEnum.DoubleTwo.Code]
	if !ok {
		twoTwos = make(map[int]bool)
	}
	twos, ok := points[ShapeEnum.Two.Code]
	if !ok {
		twos = make(map[int]bool)
	}
	return mergeMaps(blockFours, threes, blockThrees, twoTwos, twos)
}

func mergeMaps(maps ...map[int]bool) map[int]bool {
	result := make(map[int]bool)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
