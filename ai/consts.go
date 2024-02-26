package ai

const (
	ScoreFive       = 10000000
	ScoreLiveFour   = 1000000
	ScoreFourFour   = 1000000
	ScoreFourThree  = 1000000
	ScoreThreeThree = 1000000
	ScoreBlockFour  = 1500
	ScoreLiveThree  = 1000
	ScoreBlockThree = 150
	ScoreTwoTwo     = 200
	ScoreLiveTwo    = 100
	ScoreBlockTwo   = 15
	ScoreOne        = 10
	ScoreBlockOne   = 1
	ScoreNone       = 0
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
