package ai

import (
	"fmt"
	"strings"
	"time"
)

type Board struct {
	Size               int
	Board              [][]TypeRole
	FirstRole          TypeRole
	Role               TypeRole
	History            []TypeHistory
	Zobrist            *ZobristCache
	WinnerCache        *Cache
	GameOverCache      *Cache
	EvaluateCache      *Cache
	ValuableMovesCache *Cache
	EvaluateTime       time.Duration
	Evaluator          *Evaluate
}

func NewBoard(size int, firstRole TypeRole) *Board {
	board := make([][]TypeRole, size)
	for i := range board {
		board[i] = make([]TypeRole, size)
	}

	return &Board{
		Size:               size,
		Board:              board,
		FirstRole:          firstRole,
		Role:               firstRole,
		History:            make([]TypeHistory, 0),
		Zobrist:            NewZobristCache(size),
		WinnerCache:        NewCache(100 * 1024 * 1024), // 100MB cache
		GameOverCache:      NewCache(100 * 1024 * 1024), // 100MB cache
		EvaluateCache:      NewCache(100 * 1024 * 1024), // 100MB cache
		ValuableMovesCache: NewCache(100 * 1024 * 1024), // 100MB cache
		EvaluateTime:       time.Duration(0),
		Evaluator:          NewEvaluate(size),
	}
}
func (b *Board) IsGameOver() bool {
	hash := b.Zobrist.GetHash()

	val := b.GameOverCache.get(hash)
	if val != nil {
		return val.(bool)
	}

	if b.GetWinner() != Chess.EMPTY {
		b.GameOverCache.put(hash, true)
		return true
	}
	// 没有赢家但是还有空位，说明游戏还在进行中
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Board[i][j] == Chess.EMPTY {
				b.GameOverCache.put(hash, false)
				return false
			}
		}
	}
	// 没有赢家并且没有空位，游戏结束
	b.GameOverCache.put(hash, true)
	return true
}

func (b *Board) GetWinner() TypeRole {
	hash := b.Zobrist.GetHash()

	val := b.WinnerCache.get(hash)
	if val != nil {
		return val.(TypeRole)
	}
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Board[i][j] == Chess.EMPTY {
				continue
			}
			for _, vec := range DirectionVec {
				count := 0
				for i+vec.x*count >= 0 &&
					i+vec.x*count < b.Size &&
					j+vec.y*count >= 0 &&
					j+vec.y*count < b.Size &&
					b.Board[i+vec.x*count][j+vec.y*count] == b.Board[i][j] {
					count++
				}
				if count >= 5 {
					b.WinnerCache.put(hash, b.Board[i][j])
					return b.Board[i][j]
				}
			}
		}
	}
	b.WinnerCache.put(hash, Chess.EMPTY)
	return 0
}

func (b *Board) GetValidMoves() []Point {
	moves := make([]Point, 0)
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Board[i][j] == Chess.EMPTY {
				moves = append(moves, Point{i, j})
			}
		}
	}
	return moves
}

func (b *Board) Put(i, j int, role TypeRole) bool {
	if role == Chess.EMPTY {
		role = b.Role
	}
	if i < 0 || i >= b.Size || j < 0 || j >= b.Size {
		fmt.Println("Invalid move: out of boundary!", i, j)
		return false
	}
	if b.Board[i][j] != Chess.EMPTY {
		fmt.Println("Invalid move: position not empty!", i, j)
		return false
	}
	b.Board[i][j] = role
	b.History = append(b.History, TypeHistory{i, j, role})
	b.Zobrist.TogglePiece(i, j, role)
	b.Evaluator.move(i, j, role)
	b.Role *= -1 // Switch role
	return true
}
func (b *Board) Undo() bool {
	if len(b.History) == 0 {
		fmt.Println("No moves to undo!")
		return false
	}

	lastMove := b.History[len(b.History)-1]
	b.History = b.History[:len(b.History)-1]
	b.Board[lastMove.x][lastMove.y] = Chess.EMPTY // Remove the piece from the board
	b.Role = lastMove.role                        // Switch back to the previous player
	b.Zobrist.TogglePiece(lastMove.x, lastMove.y, lastMove.role)
	b.Evaluator.undo(lastMove.x, lastMove.y)
	return true
}

func (b *Board) GetValuableMoves(role TypeRole, depth int, onlyThree, onlyFour bool) []Point {
	hash := b.Zobrist.GetHash()

	moveCache := b.ValuableMovesCache.get(hash)
	if moveCache != nil {
		prevMoveCache := moveCache.(TypeValuableMoveCache)
		if prevMoveCache.role == role && prevMoveCache.depth == depth && prevMoveCache.onlyThree == onlyThree && prevMoveCache.onlyFour == onlyFour {
			return prevMoveCache.moves
		}
	}

	moves := b.Evaluator.getValuableMoves(role, depth, onlyThree, onlyFour)

	// Handle a special case, if the center point has not been played, then add the center point by default
	if !onlyThree && !onlyFour {
		center := b.Size / 2
		if b.Board[center][center] == Chess.EMPTY {
			moves = append(moves, Point{center, center})
		}
	}

	b.ValuableMovesCache.put(hash, TypeValuableMoveCache{
		role:      role,
		moves:     moves,
		depth:     depth,
		onlyThree: onlyThree,
		onlyFour:  onlyFour,
	})

	return moves
}

func (b *Board) Display(extraPoints []Point) string {
	extraPositions := make(map[int]bool, len(extraPoints))
	for _, point := range extraPoints {
		extraPositions[coordinate2Position(point.x, point.y, b.Size)] = true
	}

	var result strings.Builder
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			position := coordinate2Position(i, j, b.Size)
			if ok, exist := extraPositions[position]; ok && exist {
				result.WriteString("? ")
				continue
			}
			switch b.Board[i][j] {
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
func (b *Board) Evaluate(role TypeRole) int {
	hash := b.Zobrist.GetHash()

	evaCache := b.EvaluateCache.get(hash)
	if evaCache != nil {
		prevCache := evaCache.(TypeEvaluateCache)
		if prevCache.role == role {
			return prevCache.score
		}
	}

	winner := b.GetWinner()
	score := 0
	if winner != Chess.EMPTY {
		score = SCORE_FIVE * int(winner) * int(role)
	} else {
		score = b.Evaluator.evaluate(role)
	}

	b.EvaluateCache.put(hash, TypeEvaluateCache{
		role:  role,
		score: score,
	})

	return score
}
func (b *Board) Reverse() *Board {
	newBoard := NewBoard(b.Size, -b.FirstRole)
	for _, move := range b.History {
		newBoard.Put(move.x, move.y, -move.role)
	}
	return newBoard
}
