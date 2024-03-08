package ai

import (
	"math/rand"
)

type ZobristCache struct {
	size         int
	zobristTable [][]map[TypeRole]uint64
	hash         uint64
}

func NewZobristCache(size int) *ZobristCache {
	z := &ZobristCache{
		size: size,
		hash: 0,
	}
	z.zobristTable = initializeZobristTable(size)
	return z
}

func initializeZobristTable(size int) [][]map[TypeRole]uint64 {
	table := make([][]map[TypeRole]uint64, size)
	for i := 0; i < size; i++ {
		table[i] = make([]map[TypeRole]uint64, size)
		for j := 0; j < size; j++ {
			table[i][j] = map[TypeRole]uint64{
				Chess.BLACK: randomBit64(), // black
				Chess.WHITE: randomBit64(), // white
			}
		}
	}
	return table
}

func randomBit64() uint64 {
	return rand.Uint64()
}

func (z *ZobristCache) TogglePiece(x, y int, role TypeRole) {
	z.hash ^= z.zobristTable[x][y][role]
}

func (z *ZobristCache) GetHash() uint64 {
	return z.hash
}
