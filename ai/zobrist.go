package ai

import (
	"math/rand"
)

type ZobristCache struct {
	size         int
	zobristTable [][][]uint64
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

func (z *ZobristCache) initializeZobristTable(size int) [][][]uint64 {
	table := make([][][]uint64, size)
	for i := 0; i < size; i++ {
		table[i] = make([][]uint64, size)
		for j := 0; j < size; j++ {
			table[i][j] = []uint64{
				z.randomBitString(), // black
				z.randomBitString(), // white
			}
		}
	}
	return table
}

func (z *ZobristCache) randomBitString() uint64 {
	return rand.Uint64()
}

func (z *ZobristCache) TogglePiece(x, y, role int) {
	z.hash ^= z.zobristTable[x][y][role]
}

func (z *ZobristCache) GetHash() uint64 {
	return z.hash
}
