package ai

// 形状转换分数，注意这里的分数是当前位置还没有落子的分数
func getRealShapeScore(shape TypeShape) int {
	switch shape {
	case Shapes.LiveFive:
		return ScoreFive
	case Shapes.LiveFour:
		return ScoreLiveFour
	case Shapes.FourFour:
		return ScoreFourFour
	case Shapes.FourThree:
		return ScoreFourThree
	case Shapes.ThreeThree:
		return ScoreThreeThree
	case Shapes.BlockFour:
		return ScoreBlockFour
	case Shapes.LiveThree:
		return ScoreLiveThree
	case Shapes.BlockThree:
		return ScoreBlockThree
	case Shapes.TwoTwo:
		return ScoreTwoTwo
	case Shapes.LiveTwo:
		return ScoreLiveTwo
	case Shapes.BlockTwo:
		return ScoreBlockTwo
	case Shapes.LiveOne:
		return ScoreOne
	case Shapes.BlockOne:
		return ScoreBlockOne
	default:
		return ScoreNone
	}
}

func direction2index(vec Vector) int {
	if vec.x == 0 { // |
		return VERTICAL
	}
	if vec.y == 0 { // -
		return HORIZONTAL
	}
	if vec.x == vec.y { // \
		return DIAGONAL
	}
	return ANTI_DIAGONAL // /
}

type PerformanceEnum struct {
	updateTime    int
	getPointsTime int
}

var performance = &PerformanceEnum{
	updateTime:    0,
	getPointsTime: 0,
}

type Evaluate struct {
	size        int
	board       [][]TypeRole
	blackScores [][]int
	whiteScores [][]int
	history     []TypeHistory  // 记录历史 [position, role]
	shapeCache  TypeShapeCache // 缓存每个点位的分数，避免重复计算
	pointCache  TypePointCache // 缓存每个形状对应的点位
}

func NewEvaluate(size int) *Evaluate {
	board := make([][]TypeRole, size+2)
	for i := 0; i < size+2; i++ {
		board[i] = make([]TypeRole, size+2)
		for j := 0; j < size+2; j++ {
			if i == 0 || j == 0 || i == size+1 || j == size+1 {
				board[i][j] = Chess.OBSTACLE
			}
		}
	}
	blackScores := make([][]int, size)
	whiteScores := make([][]int, size)
	for i := 0; i < size; i++ {
		blackScores[i] = make([]int, size)
		whiteScores[i] = make([]int, size)
	}
	// 缓存每个点位的分数，避免重复计算
	shapeCache := make(TypeShapeCache)
	for _, r := range []TypeRole{Chess.BLACK, Chess.WHITE} {
		shapeCache[r] = make(map[TypeDirection]map[int]map[int]TypeShape)
		for _, direction := range []TypeDirection{HORIZONTAL, VERTICAL, DIAGONAL, ANTI_DIAGONAL} {
			shapeCache[r][direction] = make(map[int]map[int]TypeShape)
			for i := 0; i < size; i++ {
				shapeCache[r][direction][i] = make(map[int]TypeShape)
				for j := 0; j < size; j++ {
					shapeCache[r][direction][i][j] = Shapes.None
				}
			}
		}
	}
	// 缓存每个形状对应的点位
	pointCache := make(TypePointCache)
	for _, r := range Roles {
		pointCache[r] = make(map[TypeShape]map[int]bool)
		for _, shape := range ShapeFields {
			pointCache[r][shape] = make(map[int]bool)
		}
	}
	return &Evaluate{
		size:        size,
		board:       board,
		blackScores: blackScores,
		whiteScores: whiteScores,
		history:     make([]TypeHistory, 0),
	}
}

func (e *Evaluate) move(x, y int, role TypeRole) {
	// 清空记录
	for _, d := range Directions {
		e.shapeCache[role][d][x][y] = Shapes.None
		e.shapeCache[-role][d][x][y] = Shapes.None
	}
	e.blackScores[x][y] = 0
	e.whiteScores[x][y] = 0

	// 更新分数
	e.board[x+1][y+1] = role
	e.updatePoint(x, y)
	e.history = append(e.history, TypeHistory{Coordinate2Position(x, y, e.size), role})
}
func (e *Evaluate) undo(x, y int) {
	e.board[x+1][y+1] = Chess.EMPTY // Adjust for the added wall
	e.updatePoint(x, y)
	e.history = e.history[:len(e.history)-1]
}

/*
	getPointsInLine(role) {
	  for (let r of [role, -role]) {
	    for (let point of last2Points) {
	      const [x, y] = position2Coordinate(point, this.size);
	      for (let [ox, oy] of allDirections) {
	        for (let sign of [1, -1]) { // -1 for negative direction, 1 for positive direction
	          for (let step = 1; step <= config.inLineDistance; step++) {
	            const [nx, ny] = [x + sign * step * ox, y + sign * step * oy]; // +1 to adjust for wall
	            const position = coordinate2Position(nx, ny, this.size);

              // 检测是否到达边界
              if (nx < 0 || nx >= this.size || ny < 0 || ny >= this.size) {
                break;
              }
              if (this.board[nx + 1][ny + 1] !== 0) {
                continue;
              }
              if (processed[position] === r) continue;
              processed[position] = r;
              for (let direction of [0, 1, 2, 3]) {
                const shape = this.shapeCache[r][direction][nx][ny];
                // 到达边界停止，但是注意到达对方棋子不能停止
                if (shape) {
                  pointsInLine[shape].add(coordinate2Position(nx, ny, this.size));
                  hasPointsInLine = true;
                }
              }
            }
          }
        }
      }
    }
    if (hasPointsInLine) {
      return pointsInLine;
    }
    return false;
  }

*/

// 只返回和最后几步在一条直线上的点位。
// 这么做有一点问题：
// 1. 因为己方可能会由于防守暂时离开原来的线，这样就会导致己方被中断，只能增加最后几步的长度，比如不是取最后一步，而是最后3步
// 2. 如果不是取最后1步，取的步数太多了，反而还不如直接返回所有点位。
func (e *Evaluate) getPointInLine() {
	pointsInLine := make(map[TypeShape]map[int]bool)
	hasPointInLine := false
	for _, shape := range ShapeFields {
		pointsInLine[shape] = make(map[int]bool)
	}
	last2Points := e.history[len(e.history)-Config.InlineCount:]
	processed := make(map[int]TypeRole) // 已经处理过的点位
	// 在last2Points中查找是否有点位在一条线上
	for _, r := range Roles {
		for _, point := range last2Points {
			x, y := Position2Coordinate(point.position, e.size)
			for _, direction := range Directions {
				for _, sign := range []int{1, -1} {
					for step := 1; step <= Config.InLineDistance; step++ {
						nx, ny := x+sign*step*direction.x, y+sign*step*direction.y
						position := Coordinate2Position(nx, ny, e.size)
						// 检测是否到达边界
						if nx < 0 || nx >= e.size || ny < 0 || ny >= e.size {
							break
						}
						if e.board[nx+1][ny+1] != Chess.EMPTY {
							continue
						}
						if processed[position] == r {
							continue
						}
						processed[position] = r
						for _, direction := range Directions {
							shape := e.shapeCache[r][direction2index(direction)][nx][ny]
							if shape != Shapes.None {
								pointsInLine[shape][position] = true
								hasPointInLine = true
							}
						}
					}
				}
			}
		}

	}

}
func (e *Evaluate) updatePoint(x, y int) {
	// 更新当前点位的分数
	// 更新当前点位的分数

}
