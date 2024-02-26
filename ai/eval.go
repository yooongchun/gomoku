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
	board       [][]int
	blackScores [][]int
	whiteScores [][]int
	history     [][]int        // 记录历史 [position, role]
	shapeCache  TypeShapeCache // 缓存每个点位的分数，避免重复计算
}

func NewEvaluate(size int) *Evaluate {
	board := make([][]int, size+2)
	for i := 0; i < size+2; i++ {
		board[i] = make([]int, size+2)
		for j := 0; j < size+2; j++ {
			if i == 0 || j == 0 || i == size+1 || j == size+1 {
				board[i][j] = 2
			}
		}
	}
	blackScores := make([][]int, size)
	whiteScores := make([][]int, size)
	for i := 0; i < size; i++ {
		blackScores[i] = make([]int, size)
		whiteScores[i] = make([]int, size)
	}
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
	/*
	    // 缓存每个形状对应的点位
	    // 结构： pointsCache[role][shape] = Set(direction1, direction2);
	    this.pointsCache = {}
	    for (let role of [1, -1]) {
	      this.pointsCache[role] = {};
	      for (let key of Object.keys(shapes)) {
	        const shape = shapes[key];
	        this.pointsCache[role][shape] = new Set();
	      }
	    }
	  }
	*/
	return &Evaluate{
		size:        size,
		board:       board,
		blackScores: blackScores,
		whiteScores: whiteScores,
		history:     make([][]int, 0),
	}
}

/*
	move(x, y, role) {
	  // 清空记录
	  for (const d of [0, 1, 2, 3]) {
	    this.shapeCache[role][d][x][y] = 0;
	    this.shapeCache[-role][d][x][y] = 0;
	  }
	  this.blackScores[x][y] = 0;
	  this.whiteScores[x][y] = 0;

	  // 更新分数
	  this.board[x + 1][y + 1] = role; // Adjust for the added wall
	  this.updatePoint(x, y);
	  this.history.push([coordinate2Position(x, y, this.size), role]);
	}
*/
func (e *Evaluate) move(x, y, role int) {
	// 清空记录
	for _, d := range []int{0, 1, 2, 3} {
		e.shapeCache[role][d][x][y] = 0
		e.shapeCache[-role][d][x][y] = 0
	}
	e.blackScores[x][y] = 0
	e.whiteScores[x][y] = 0

	// 更新分数
	e.board[x+1][y+1] = role
	e.updatePoint(x, y)
	e.history = append(e.history, []int{Coordinate2Position(x, y, e.size), role})
}

func (e *Evaluate) updatePoint(x, y int) {
	// 更新当前点位的分数
	// 更新当前点位的分数

}
