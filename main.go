package main

import "gomoku/ai"

func main() {
	//ui.Show()
	ai.NewPlay(15, ai.ROLE_HUMAN).Play()
}
