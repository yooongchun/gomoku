package main

import (
	"fmt"
	"gomoku/ai"
)

func main() {
	evaluate := ai.NewEvaluate(15)
	evaluate.Move(7, 7, ai.BLACK)
	evaluate.Move(7, 8, ai.WHITE)
	evaluate.Move(8, 8, ai.BLACK)
	evaluate.Move(8, 9, ai.WHITE)
	evaluate.Move(9, 9, ai.BLACK)
	evaluate.Move(9, 10, ai.WHITE)
	evaluate.Move(10, 10, ai.BLACK)
	evaluate.Display()
	fmt.Printf("O score=%d, X score=%d", evaluate.Evaluate(ai.BLACK), evaluate.Evaluate(ai.WHITE))

}
