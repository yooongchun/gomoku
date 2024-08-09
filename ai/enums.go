package ai

var ShapeEnum = &TypeShapeAll{
	Five:        TypeShapeField{5, SCORE_FIVE, []string{XXXXX}},
	Four:        TypeShapeField{4, SCORE_FOUR, []string{OXXXXO}},
	RushFour:    TypeShapeField{40, SCORE_BLOCK_FOUR, []string{XXOXX, XXXOX, ZXXXXO}},
	Three:       TypeShapeField{3, SCORE_THREE, []string{OXXXOO, OXXOXO}},
	BlockThree:  TypeShapeField{30, SCORE_BLOCK_THREE, []string{OXOXX, XOOXX, OXXOX, XOXOX, XXXOO, ZOXXXO, ZXOXXO, ZXXOXO, ZXXXOO}},
	Two:         TypeShapeField{2, SCORE_TWO, []string{OOXXOO, OXXOOO, OXOXOO, OXOOXO}},
	BlockTwo:    TypeShapeField{20, SCORE_BLOCK_TWO, []string{XXOOO, XOOXO, OXOXO, XOOOX, XOXOO, OXXOO, ZOOXXO, ZOXOXO, ZOXXOO, ZXOOXO, ZXOXOO, ZXXOOO}},
	DoubleFour:  TypeShapeField{44, SCORE_FOUR_FOUR, nil},
	FourThree:   TypeShapeField{43, SCORE_FOUR_THREE, nil},
	DoubleThree: TypeShapeField{33, SCORE_THREE_THREE, nil},
	DoubleTwo:   TypeShapeField{22, SCORE_TWO_TWO, nil},
	None:        TypeShapeField{0, SCORE_NONE, nil},
}

var DirectionVec = map[TypeDirection]Point{HORIZONTAL: {0, 1}, VERTICAL: {1, 0}, DIAGONAL: {1, 1}, ANTI_DIAGONAL: {1, -1}}
var DirectionEnum = []TypeDirection{HORIZONTAL, VERTICAL, DIAGONAL, ANTI_DIAGONAL}
