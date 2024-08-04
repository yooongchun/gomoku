package ai

import (
	"math/rand"
)

type ZobristCache struct {
	size         int
	zobristTable [][]map[TypeChess]uint64
	hash         uint64
}

func NewZobristCache(size int) *ZobristCache {
	z := &ZobristCache{
		size: size,
		hash: 0,
	}
	z.zobristTable = z.initializeZobristTable(size)
	return z
}

func (z *ZobristCache) initializeZobristTable(size int) [][]map[TypeChess]uint64 {
	table := make([][]map[TypeChess]uint64, size)
	for i := 0; i < size; i++ {
		table[i] = make([]map[TypeChess]uint64, size)
		for j := 0; j < size; j++ {
			table[i][j] = map[TypeChess]uint64{
				CHESS_BLACK: rand.Uint64(), // black
				CHESS_WHITE: rand.Uint64(), // white
			}
		}
	}
	return table
}

func (z *ZobristCache) TogglePiece(x, y int, chess TypeChess) {
	z.hash ^= z.zobristTable[x][y][chess]
}

func (z *ZobristCache) GetHash() uint64 {
	return z.hash
}
