package main

import (
	"fmt"
	"gomoku/ai"
)

func main() {
	z := ai.NewZobristCache(8)
	z.TogglePiece(1, 2, 0)
	x := z.GetHash()
	fmt.Println(z.GetHash(), x == max(5, 2, 3))
}
