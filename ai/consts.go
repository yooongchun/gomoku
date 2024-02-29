package ai

const (
	SCORE_FIVE        = 10000000
	SCORE_FOUR        = 1000000
	SCORE_FOUR_FOUR   = 1000000
	SCORE_FOUR_THREE  = 1000000
	SCORE_THREE_THREE = 1000000
	SCORE_BLOCK_FOUR  = 1500
	SCORE_THREE       = 1000
	SCORE_BLOCK_THREE = 150
	SCORE_TWO_TWO     = 200
	SCORE_TWO         = 100
	SCORE_BLOCK_TWO   = 15
	SCORE_ONE         = 10
	SCORE_BLOCK_ONE   = 1
	SCORE_NONE        = 0
)

const (
	BLACK    = 1
	WHITE    = -1
	EMPTY    = 0
	OBSTACLE = 2
)

const (
	HORIZONTAL = iota
	VERTICAL
	DIAGONAL
	ANTI_DIAGONAL
)
