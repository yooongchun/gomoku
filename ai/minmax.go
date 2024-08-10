package ai

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
)

type hits struct {
	search int
	total  int
	hit    int
}
type TypeCacheField struct {
	depth     int // 剩余搜索深度
	value     int
	move      Point
	chess     TypeChess
	path      []Point // 剩余搜索路径
	onlyThree bool
	onlyFour  bool
}

type MinMax struct {
	onlyThree          bool
	onlyFour           bool
	onlyThreeThreshold int
	cacheHits          hits
	cache              *lru.Cache[uint64, TypeCacheField]
}

func NewMinMax(onlyThree, onlyFour bool) *MinMax {
	cache, _ := lru.New[uint64, TypeCacheField](100_0000)
	return &MinMax{
		cache:              cache,
		onlyThree:          onlyThree,
		onlyFour:           onlyFour,
		onlyThreeThreshold: 6,
	}
}

func (mm *MinMax) search(board *Board, chess TypeChess, depth, cDepth, alpha, beta int, path []Point) (fValue int, fMove *Point, fPath []Point) {
	// depth 表示总深度 cDepth 表示当前搜索深度
	mm.cacheHits.search++
	// 搜索到底或者游戏结束
	if cDepth >= depth || board.IsGameOver() {
		fValue = board.Evaluate(chess)
		fMove = nil
		fPath = path
		return
	}
	hash := board.zobrist.GetHash()
	prev, ok := mm.cache.Get(hash)
	if ok && prev.chess == chess {
		// 不能连5的，minmax和vct vcf缓存不能通用
		if (abs(prev.value) >= SCORE_FIVE || prev.depth >= depth-cDepth) && prev.onlyThree == mm.onlyThree && prev.onlyFour == mm.onlyFour {
			mm.cacheHits.hit++
			fValue = prev.value
			fMove = &prev.move
			fPath = append(path, prev.path...)
			return
		}
	}
	value := -SCORE_MAX
	move := Point{}
	var bestPath []Point
	bestDepth := 0
	points := board.GetValuableMoves(chess, cDepth, mm.onlyThree || cDepth > mm.onlyThreeThreshold, mm.onlyFour)
	if cDepth == 0 {
		fmt.Println("points:", points)
	}
	if len(points) == 0 {
		fValue = board.Evaluate(chess)
		fMove = nil
		fPath = path
		return
	}
breakLoop:
	for d := cDepth + 1; d <= depth; d++ {
		// 迭代加深的过程中只找自己方能赢的解，因此只需要搜索偶数层即可
		if d%2 != 0 {
			continue
		}
		for _, point := range points {
			board.Move(point)
			newPath := append(path, point)
			currValue, _, currPath := mm.search(board, togglePiece(chess), d, cDepth+1, -beta, -alpha, newPath)
			currValue = -currValue
			board.Undo()
			// 迭代加深过程中，除了能赢的棋，其他都不要
			// 原因是除了必胜的，其他评估不准。比如必输的棋，由于走的步数偏少，也会变成没有输，比如5步之后输了，但是1步肯定不会输，这时候1步的分数是不准确的，显然不能选择
			if currValue >= SCORE_FIVE || d == depth {
				// 必输的棋，也要挣扎一下，选择最长的路径
				if currValue > value || (currValue <= -SCORE_FIVE && value <= -SCORE_FIVE && len(currPath) > bestDepth) {
					value = currValue
					move = point
					bestDepth = len(currPath)
				}
			}
			alpha = max(alpha, value)
			if alpha >= SCORE_FIVE {
				// 自己赢了也结束，但是对方赢了还是要继续搜索
				break breakLoop
			}
			// 剪枝
			if alpha >= beta {
				break
			}
		}
	}
	// 缓存
	if (cDepth < mm.onlyThreeThreshold || mm.onlyThree || mm.onlyFour) && (!ok || prev.depth < depth-cDepth) {
		mm.cacheHits.total++
		mm.cache.Add(hash, TypeCacheField{
			depth:     depth - cDepth,
			value:     value,
			move:      move,
			chess:     chess,
			path:      bestPath[cDepth:],
			onlyThree: mm.onlyThree,
			onlyFour:  mm.onlyFour,
		})
	}
	fValue = value
	fMove = &move
	fPath = bestPath
	return
}

func abs(val int) int {
	if val > 0 {
		return val
	}
	return -val
}

var minmax = NewMinMax(false, false)
var vct = NewMinMax(true, false)
var vcf = NewMinMax(false, true)

func MinMaxSearch(board *Board, chess TypeChess, depth int, enableVCT bool) (score int, move *Point, path []Point) {
	if enableVCT {
		vctDepth := 8 + depth
		// 先看自己有没有杀棋
		score, move, path = vct.search(board, chess, vctDepth, 0, -SCORE_MAX, SCORE_MAX, []Point{})
		if score >= SCORE_FIVE {
			return
		}
		score, move, path = minmax.search(board, chess, depth, 0, -SCORE_MAX, SCORE_MAX, []Point{})
		// 假设对方有杀棋，先按自己的思路走，走完之后看对方是不是还有杀棋
		// 如果对方没有了，那么就说明走的是对的
		// 如果对方还是有，那么要对比对方的杀棋路径和自己没有走棋时的长短
		// 如果走了棋之后路径变长了，说明走的是对的
		// 如果走了棋之后，对方杀棋路径长度没变，甚至更短，说明走错了，此时就优先封堵对方
		board.Move(*move)
		value2, move2, bestPath2 := vct.search(board.reverse(), chess, vctDepth, 0, -SCORE_MAX, SCORE_MAX, []Point{})
		board.Undo()
		if score < SCORE_FIVE && value2 == SCORE_FIVE && len(bestPath2) > len(path) {
			_, _, bestPath3 := vct.search(board.reverse(), chess, vctDepth, 0, -SCORE_MAX, SCORE_MAX, []Point{})
			if len(bestPath2) <= len(bestPath3) {
				// value2 是被挡住的，所以这里还是用value
				move = move2
				path = bestPath2
				return
			}
		}
		return
	} else {
		return minmax.search(board, chess, depth, 0, -SCORE_MAX, SCORE_MAX, []Point{})
	}
}
