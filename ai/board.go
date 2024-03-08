package ai

import (
	"fmt"
	"strings"
	"time"
)

type Board struct {
	size               int
	board              [][]TypeRole
	firstRole          TypeRole
	role               TypeRole
	history            []TypeHistory
	zobrist            *ZobristCache
	winnerCache        *Cache
	gameOverCache      *Cache
	evaluateCache      *Cache
	valuableMovesCache *Cache
	evaluateTime       time.Duration
	evaluator          *Evaluate
}

func NewBoard(size int, firstRole TypeRole) *Board {
	board := make([][]TypeRole, size)
	for i := range board {
		board[i] = make([]TypeRole, size)
	}

	return &Board{
		size:               size,
		board:              board,
		firstRole:          firstRole,
		role:               firstRole,
		history:            make([]TypeHistory, 0),
		zobrist:            NewZobristCache(size),
		winnerCache:        NewCache(100 * 1024 * 1024), // 100MB cache
		gameOverCache:      NewCache(100 * 1024 * 1024), // 100MB cache
		evaluateCache:      NewCache(100 * 1024 * 1024), // 100MB cache
		valuableMovesCache: NewCache(100 * 1024 * 1024), // 100MB cache
		evaluateTime:       time.Duration(0),
		evaluator:          NewEvaluate(size),
	}
}
func (b *Board) isGameOver() bool {
	hash := b.zobrist.GetHash()

	val := b.gameOverCache.get(hash)
	if val != nil {
		return val.(bool)
	}

	if b.getWinner() != Chess.EMPTY {
		b.gameOverCache.put(hash, true)
		return true
	}
	// 没有赢家但是还有空位，说明游戏还在进行中
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i][j] == Chess.EMPTY {
				b.gameOverCache.put(hash, false)
				return false
			}
		}
	}
	// 没有赢家并且没有空位，游戏结束
	b.gameOverCache.put(hash, true)
	return true
}

func (b *Board) getWinner() TypeRole {
	hash := b.zobrist.GetHash()

	val := b.winnerCache.get(hash)
	if val != nil {
		return val.(TypeRole)
	}
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i][j] == Chess.EMPTY {
				continue
			}
			for _, vec := range DirectionVec {
				count := 0
				for i+vec.x*count >= 0 &&
					i+vec.x*count < b.size &&
					j+vec.y*count >= 0 &&
					j+vec.y*count < b.size &&
					b.board[i+vec.x*count][j+vec.y*count] == b.board[i][j] {
					count++
				}
				if count >= 5 {
					b.winnerCache.put(hash, b.board[i][j])
					return b.board[i][j]
				}
			}
		}
	}
	b.winnerCache.put(hash, Chess.EMPTY)
	return 0
}

func (b *Board) getValidMoves() []Point {
	moves := make([]Point, 0)
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i][j] == Chess.EMPTY {
				moves = append(moves, Point{i, j})
			}
		}
	}
	return moves
}

func (b *Board) put(i, j int, role TypeRole) bool {
	if role == Chess.EMPTY {
		role = b.role
	}
	if i < 0 || i >= b.size || j < 0 || j >= b.size {
		fmt.Println("Invalid move: out of boundary!", i, j)
		return false
	}
	if b.board[i][j] != Chess.EMPTY {
		fmt.Println("Invalid move: position not empty!", i, j)
		return false
	}
	b.board[i][j] = role
	b.history = append(b.history, TypeHistory{i, j, role})
	b.zobrist.TogglePiece(i, j, role)
	b.evaluator.move(i, j, role)
	b.role *= -1 // Switch role
	return true
}
func (b *Board) undo() bool {
	if len(b.history) == 0 {
		fmt.Println("No moves to undo!")
		return false
	}

	lastMove := b.history[len(b.history)-1]
	b.history = b.history[:len(b.history)-1]
	b.board[lastMove.x][lastMove.y] = Chess.EMPTY // Remove the piece from the board
	b.role = lastMove.role                        // Switch back to the previous player
	b.zobrist.TogglePiece(lastMove.x, lastMove.y, lastMove.role)
	b.evaluator.undo(lastMove.x, lastMove.y)
	return true
}

func (b *Board) getValuableMoves(role TypeRole, depth int, onlyThree, onlyFour bool) []Point {
	hash := b.zobrist.GetHash()

	moveCache := b.valuableMovesCache.get(hash)
	if moveCache != nil {
		prevMoveCache := moveCache.(TypeValuableMoveCache)
		if prevMoveCache.role == role && prevMoveCache.depth == depth && prevMoveCache.onlyThree == onlyThree && prevMoveCache.onlyFour == onlyFour {
			return prevMoveCache.moves
		}
	}

	moves := b.evaluator.getValuableMoves(role, depth, onlyThree, onlyFour)

	// Handle a special case, if the center point has not been played, then add the center point by default
	if !onlyThree && !onlyFour {
		center := b.size / 2
		if b.board[center][center] == Chess.EMPTY {
			moves = append(moves, Point{center, center})
		}
	}

	b.valuableMovesCache.put(hash, TypeValuableMoveCache{
		role:      role,
		moves:     moves,
		depth:     depth,
		onlyThree: onlyThree,
		onlyFour:  onlyFour,
	})

	return moves
}

func (b *Board) display(extraPoints []Point) string {
	extraPositions := make(map[int]bool, len(extraPoints))
	for _, point := range extraPoints {
		extraPositions[coordinate2Position(point.x, point.y, b.size)] = true
	}

	var result strings.Builder
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			position := coordinate2Position(i, j, b.size)
			if ok, exist := extraPositions[position]; ok && exist {
				result.WriteString("? ")
				continue
			}
			switch b.board[i][j] {
			case Chess.BLACK:
				result.WriteString("O ")
			case Chess.WHITE:
				result.WriteString("X ")
			default:
				result.WriteString("- ")
			}
		}
		result.WriteString("\n") // New line at the end of each row
	}
	return result.String()
}
func (b *Board) evaluate(role TypeRole) int {
	hash := b.zobrist.GetHash()

	evaCache := b.evaluateCache.get(hash)
	if evaCache != nil {
		prevCache := evaCache.(TypeEvaluateCache)
		if prevCache.role == role {
			return prevCache.score
		}
	}

	winner := b.getWinner()
	score := 0
	if winner != Chess.EMPTY {
		score = SCORE_FIVE * int(winner) * int(role)
	} else {
		score = b.evaluator.evaluate(role)
	}

	b.evaluateCache.put(hash, TypeEvaluateCache{
		role:  role,
		score: score,
	})

	return score
}
func (b *Board) reverse() *Board {
	newBoard := NewBoard(b.size, -b.firstRole)
	for _, move := range b.history {
		newBoard.put(move.x, move.y, -move.role)
	}
	return newBoard
}

func (b *Board) getHash() uint64 {
	return b.zobrist.GetHash()
}
