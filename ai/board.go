package ai

import (
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/stream"
	"github.com/hashicorp/golang-lru/v2"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Board struct {
	size               int
	board              [][]TypeChess
	firstRole          TypeRole
	current            TypeChess
	history            []TypeHistory
	zobrist            *ZobristCache
	winnerCache        *lru.Cache[uint64, TypeChess]
	gameOverCache      *lru.Cache[uint64, bool]
	evaluateCache      *lru.Cache[uint64, TypeEvaluateCache]
	valuableMovesCache *lru.Cache[uint64, TypeValuableMoveCache]
	evaluateTime       time.Duration
	evaluate           *Evaluate
}

func NewBoard(size int, firstRole TypeRole) *Board {
	board := make([][]TypeChess, size)
	for i := range board {
		board[i] = make([]TypeChess, size)
	}
	cacheSize := 100 * 1024 * 1024 // 100MB
	winnerCache, _ := lru.New[uint64, TypeChess](cacheSize)
	gameOverCache, _ := lru.New[uint64, bool](cacheSize)
	evaCache, _ := lru.New[uint64, TypeEvaluateCache](cacheSize)
	valCache, _ := lru.New[uint64, TypeValuableMoveCache](cacheSize)

	return &Board{
		size:               size,
		board:              board,
		firstRole:          firstRole,
		current:            CHESS_BLACK,
		history:            make([]TypeHistory, 0),
		zobrist:            NewZobristCache(size),
		winnerCache:        winnerCache,
		gameOverCache:      gameOverCache,
		evaluateCache:      evaCache,
		valuableMovesCache: valCache,
		evaluateTime:       time.Duration(0),
		evaluate:           NewEvaluate(size),
	}
}
func Load(fileName string) *Board {
	// 按行读取内容
	textLines, err := fileutil.ReadFileByLine(fileName)
	if err != nil {
		_ = fmt.Errorf("error reading file:%v", err)
		return nil
	}
	// 过滤掉#开头的注释行
	textLines = stream.FromSlice(textLines).Filter(func(s string) bool { return !strings.HasPrefix(s, "#") }).ToSlice()
	// 加载meta信息
	var size, firstRole int
	_, err = fmt.Sscanf(textLines[0], "%d,%d", &size, &firstRole)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	board := NewBoard(size, TypeRole(firstRole))
	for i := 1; i < len(textLines); i++ {
		// 按行读取历史数据
		var x, y, chess int
		_, err = fmt.Sscanf(textLines[i], "%d,%d,%d", &x, &y, &chess)
		if err != nil {
			fmt.Printf("Error reading line %d: %s,%v\n", i, textLines[i], err)
		} else {
			board.Move(Point{x, y})
		}
	}
	return board
}

func (b *Board) toChess(role TypeRole) TypeChess {
	if role == b.firstRole {
		return CHESS_BLACK
	} else if role == b.secondRole() {
		return CHESS_WHITE
	}
	return CHESS_EMPTY
}

func (b *Board) toRole(chess TypeChess) TypeRole {
	if chess == CHESS_BLACK {
		return b.firstRole
	} else if chess == CHESS_WHITE {
		return b.secondRole()
	}
	return NOBODY
}

func (b *Board) secondRole() TypeRole {
	if b.firstRole == ROLE_HUMAN {
		return ROLE_AI
	}
	return ROLE_HUMAN
}

func (b *Board) WhoseTurn() TypeRole {
	return b.toRole(b.current)
}

func (b *Board) IsGameOver() bool {
	hash := b.zobrist.GetHash()

	val, ok := b.gameOverCache.Get(hash)
	if ok {
		return val
	}

	if b.GetWinner() != NOBODY {
		b.gameOverCache.Add(hash, true)
		return true
	}
	// 没有赢家但是还有空位，说明游戏还在进行中
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i][j] == CHESS_EMPTY {
				b.gameOverCache.Add(hash, false)
				return false
			}
		}
	}
	// 没有赢家并且没有空位，游戏结束
	b.gameOverCache.Add(hash, true)
	return true
}

func (b *Board) GetWinner() TypeRole {
	hash := b.zobrist.GetHash()

	chess, ok := b.winnerCache.Get(hash)
	if ok {
		if chess == CHESS_EMPTY {
			return 0
		}
		return b.toRole(chess)
	}
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i][j] == CHESS_EMPTY {
				continue
			}
			for _, v := range DirectionVec {
				count, nextX, nextY := 0, i, j
				for nextX >= 0 && nextX < b.size && nextY >= 0 && nextY < b.size && b.board[nextX][nextY] == b.board[i][j] {
					count++
					nextX = i + count*v.x
					nextY = j + count*v.y
				}
				if count >= 5 {
					b.winnerCache.Add(hash, b.board[i][j])
					return b.toRole(b.board[i][j])
				}
			}
		}
	}
	b.winnerCache.Add(hash, CHESS_EMPTY)
	return NOBODY
}

func (b *Board) GetValidMoves() []Point {
	moves := make([]Point, 0)
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i][j] == CHESS_EMPTY {
				moves = append(moves, Point{i, j})
			}
		}
	}
	return moves
}

func (b *Board) togglePiece() {
	b.current = togglePiece(b.current)
}
func (b *Board) Move(point Point) {
	if point.x < 0 || point.x >= b.size || point.y < 0 || point.y >= b.size {
		logrus.Errorf("Invalid move: out of board! %v", point)
		return
	}
	if b.board[point.x][point.y] != CHESS_EMPTY {
		logrus.Errorf("Invalid move: position not empty! %v", point)
		return
	}
	b.board[point.x][point.y] = b.current
	b.history = append(b.history, TypeHistory{point, b.current})
	b.zobrist.TogglePiece(point.x, point.y, b.current)
	b.evaluate.Move(point, b.current)
	b.togglePiece()
}
func (b *Board) Undo() {
	if len(b.history) == 0 {
		logrus.Warning("No moves to undo!")
		return
	}
	lastMove := b.history[len(b.history)-1]
	point := lastMove.point
	b.history = b.history[:len(b.history)-1]
	b.board[point.x][point.y] = CHESS_EMPTY // Remove the piece from the board
	b.zobrist.TogglePiece(point.x, point.y, lastMove.chess)
	b.evaluate.Undo(point)
	b.togglePiece()
}

func (b *Board) GetValuableMoves(chess TypeChess, depth int, onlyThree, onlyFour bool) []Point {
	hash := b.zobrist.GetHash()

	moveCache, ok := b.valuableMovesCache.Get(hash)
	if ok {
		prevMoveCache := moveCache
		if prevMoveCache.chess == chess && prevMoveCache.depth == depth && prevMoveCache.onlyThree == onlyThree && prevMoveCache.onlyFour == onlyFour {
			return prevMoveCache.moves
		}
	}

	moves := b.evaluate.getMoves(chess, depth, onlyThree, onlyFour)

	// Handle a special case, if the center point has not been played, then add the center point by default
	if !onlyThree && !onlyFour {
		center := b.size / 2
		if b.board[center][center] == CHESS_EMPTY {
			moves = append(moves, Point{center, center})
		}
	}

	b.valuableMovesCache.Add(hash, TypeValuableMoveCache{
		chess:     chess,
		moves:     moves,
		depth:     depth,
		onlyThree: onlyThree,
		onlyFour:  onlyFour,
	})

	return moves
}

func (b *Board) Display(extraPoints ...Point) {
	var lastOp *Point
	if len(b.history) > 0 {
		lastOp = &b.history[len(b.history)-1].point
	}
	// 显示棋盘
	fmt.Print(getBoardString(b.board, lastOp, extraPoints...))
	// 显示当前各方的得分
	bs, ws := b.evaluate.getScore()
	fmt.Printf("Score:\t[BLACK] : [WHITE] = %d : %d\n", bs, ws)
}

func (b *Board) Save(name ...string) {
	historyDir := "../history"
	fileName := fmt.Sprintf("%s/play-%d.txt", historyDir, time.Now().Unix())
	if len(name) > 0 {
		fileName = fmt.Sprintf("%s/%s", historyDir, name[0])
	}
	// Get the directory from the file path
	dir := filepath.Dir(fileName)

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory if it does not exist
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 第一行保存meta信息，包括size,firstRole
	_, _ = file.WriteString(fmt.Sprintf("#meta info:size,firstRole\n%d,%d\n#history info:x,y,chess\n", b.size, b.firstRole))
	// Write the board state to the file
	for _, op := range b.history {
		_, _ = file.WriteString(fmt.Sprintf("%d,%d,%d\n", op.point.x, op.point.y, op.chess))
	}
	fmt.Printf("Game save to %s\n", fileName)
}

func (b *Board) Evaluate(chess TypeChess) int {
	hash := b.zobrist.GetHash()

	evaCache, ok := b.evaluateCache.Get(hash)
	if ok {
		prevCache := evaCache
		if prevCache.chess == chess {
			return prevCache.score
		}
	}
	score := b.evaluate.Evaluate(chess)
	b.evaluateCache.Add(hash, TypeEvaluateCache{
		chess: chess,
		score: score,
	})

	return score
}

func (b *Board) reverse() *Board {
	newBoard := NewBoard(b.size, b.secondRole())
	for _, op := range b.history {
		newBoard.Move(op.point)
	}
	return newBoard
}
