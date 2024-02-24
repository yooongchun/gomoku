package ai

// 形状转换分数，注意这里的分数是当前位置还没有落子的分数
func getRealShapeScore(shape int) int {
	switch shape {
	case shapes["FIVE"]:
		return scores.Four
	case shapes["BLOCK_FIVE"]:
		return scores.BlockFour
	case shapes["FOUR"]:
		return scores.Three
	case shapes["FOUR_FOUR"]:
		return scores.Three
	case shapes["FOUR_THREE"]:
		return scores.Three
	case shapes["THREE_THREE"]:
		return scores.ThreeThree / 10
	case shapes["BLOCK_FOUR"]:
		return scores.BlockThree
	case shapes["THREE"]:
		return scores.Two
	case shapes["BLOCK_THREE"]:
		return scores.BlockTwo
	case shapes["TWO_TWO"]:
		return scores.TwoTwo / 10
	case shapes["TWO"]:
		return scores.One
	default:
		return scores.None
	}
}

func direction2index(ox, oy int) int {
	if ox == 0 { // |
		return 0
	}
	if oy == 0 { // -
		return 1
	}
	if ox == oy { // \
		return 2
	}
	// /
	return 3
}

type Performance struct {
	updateTime    int
	getPointsTime int
}

var performanceData = &Performance{
	updateTime:    0,
	getPointsTime: 0,
}

type Evaluate struct {
	size        int
	board       [][]int
	blackScores [][]int
	whiteScores [][]int
	history     [][]int                             // 记录历史 [position, role]
	shapeCache  map[int]map[int]map[int]map[int]int // 缓存每个点位的分数，避免重复计算
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
	shapeCache := make(map[int]map[int]map[int]map[int]int)
	for _, r := range []int{role.Black, role.White} {
		shapeCache[r] = make(map[int]map[int]map[int]int)
		for _, direction := range []int{0, 1, 2, 3} {
			shapeCache[r][direction] = make(map[int]map[int]int)
			for i := 0; i < size; i++ {
				shapeCache[r][direction][i] = make(map[int]int)
				for j := 0; j < size; j++ {
					shapeCache[r][direction][i][j] = shapes["NONE"]
				}
			}
		}

	}
	/*
	  initPoints() {
	    // 缓存每个点位的分数，避免重复计算
	    // 结构： [role][direction][x][y] = shape
	    this.shapeCache = {};
	    for (let role of [1, -1]) {
	      this.shapeCache[role] = {};
	      for (let direction of [0, 1, 2, 3]) {
	        this.shapeCache[role][direction] = Array.from({ length: this.size }).map(() => Array.from({ length: this.size }).fill(shapes.NONE));
	      }
	    }
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
	e.history = append(e.history, []int{coordinate2Position(x, y, e.size), role})
}
