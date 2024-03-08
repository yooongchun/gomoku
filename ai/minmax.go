package ai

import (
	"fmt"
	"math"
)

var (
	cacheHits = struct {
		search int
		total  int
		hit    int
	}{}
	onlyThreeThreshold = 6
	cache              = NewCache(100 * 1024 * 1024)
)

type CacheEntry struct {
	depth     int
	value     int
	move      *Point
	role      TypeRole
	path      []Point
	onlyThree bool
	onlyFour  bool
}

func copySlice[T any](s []T) []T {
	copySlice := make([]T, len(s))
	copy(copySlice, s)
	return copySlice
}

func helper(board *Board, role TypeRole, path []Point, depth, cDepth, alpha, beta int, onlyThree, onlyFour bool) (int, *Point, []Point) {
	// depth 表示总深度，cDepth表示当前搜索深度
	cacheHits.search++
	if cDepth >= depth || board.isGameOver() {
		return board.evaluate(role), nil, copySlice(path)
	}
	hash := board.getHash()
	prevCache := cache.get(hash).(*CacheEntry)
	if prevCache != nil && prevCache.role == role {
		if (math.Abs(float64(prevCache.value)) >= SCORE_FIVE || prevCache.depth >= depth-cDepth) && prevCache.onlyThree == onlyThree && prevCache.onlyFour == onlyFour {
			cacheHits.hit++
			return prevCache.value, prevCache.move, append(copySlice(path), prevCache.path...)
		}
	}
	value := -SCORE_MAX
	var move *Point
	bestPath := copySlice(path)
	bestDepth := 0
	points := board.getValuableMoves(role, cDepth, onlyThree || cDepth > onlyThreeThreshold, onlyFour)
	if cDepth == 0 {
		fmt.Println("points:", points)
	}
	if len(points) == 0 {
		return board.evaluate(role), nil, copySlice(path)
	}
	for d := cDepth + 1; d <= depth; d++ {
		if d%2 != 0 {
			continue
		}
		breakAll := false
		for _, point := range points {
			board.put(point.x, point.y, role)
			newPath := append(copySlice(path), point)
			currentValue, _, currentPath := helper(board, -role, newPath, d, cDepth+1, -beta, -alpha, onlyThree, onlyFour)
			currentValue = -currentValue
			board.undo()
			if currentValue >= SCORE_FIVE || d == depth {
				if (currentValue > value) || (currentValue <= -SCORE_FIVE && value <= -SCORE_FIVE && len(currentPath) > bestDepth) {
					value = currentValue
					move = &point
					bestPath = currentPath
					bestDepth = len(currentPath)
				}
			}
			alpha = int(math.Max(float64(alpha), float64(value)))
			if alpha >= SCORE_FIVE {
				breakAll = true
				break
			}
			if alpha >= beta {
				break
			}
		}
		if breakAll {
			break
		}
	}
	if (cDepth < onlyThreeThreshold || onlyThree || onlyFour) && (prevCache == nil || prevCache.depth < depth-cDepth) {
		cacheHits.total++
		cache.put(hash, &CacheEntry{
			depth:     depth - cDepth,
			value:     value,
			move:      move,
			role:      role,
			path:      bestPath[cDepth:],
			onlyThree: onlyThree,
			onlyFour:  onlyFour,
		})
	}
	return value, move, bestPath
}
func _minmax(board *Board, role TypeRole, path []Point, depth, cDepth, alpha, beta int) (int, *Point, []Point) {
	return helper(board, role, path, depth, cDepth, alpha, beta, false, false)
}
func vct(board *Board, role TypeRole, path []Point, depth, cDepth, alpha, beta int) (int, *Point, []Point) {
	return helper(board, role, path, depth, cDepth, alpha, beta, true, false)
}

func vcf(board *Board, role TypeRole, path []Point, depth, cDepth, alpha, beta int) (int, *Point, []Point) {
	return helper(board, role, path, depth, cDepth, alpha, beta, false, true)
}

func minmax(board *Board, role TypeRole, depth int, enableVCT bool) (int, *Point, []Point) {
	if enableVCT {
		vctDepth := depth + 8
		// 先看自己有没有杀棋
		value, move, bestPath := vct(board, role, []Point{}, vctDepth, 0, -SCORE_MAX, SCORE_MAX)
		if value >= SCORE_FIVE {
			return value, move, bestPath
		}
		value, move, bestPath = _minmax(board, role, []Point{}, depth, 0, -SCORE_MAX, SCORE_MAX)
		// 假设对方有杀棋，先按自己的思路走，走完之后看对方是不是还有杀棋
		// 如果对方没有了，那么就说明走的是对的
		// 如果对方还是有，那么要对比对方的杀棋路径和自己没有走棋时的长短
		// 如果走了棋之后路径变长了，说明走的是对的
		// 如果走了棋之后，对方杀棋路径长度没变，甚至更短，说明走错了，此时就优先封堵对方
		board.put(move.x, move.y, role)
		value2, move2, bestPath2 := vct(board.reverse(), role, []Point{}, vctDepth, 0, -SCORE_MAX, SCORE_MAX)
		board.undo()
		if value < SCORE_FIVE && value2 == SCORE_FIVE && len(bestPath2) > len(bestPath) {
			_, _, bestPath3 := vct(board.reverse(), role, []Point{}, vctDepth, 0, -SCORE_MAX, SCORE_MAX)
			if len(bestPath2) <= len(bestPath3) {
				return value, move2, bestPath2 // value2 是被挡住的，所以这里还是用value
			}
		}
		return value, move, bestPath
	} else {
		return _minmax(board, role, []Point{}, depth, 0, -SCORE_MAX, SCORE_MAX)
	}
}
