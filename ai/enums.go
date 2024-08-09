package ai

var ShapeEnum = &TypeShapeAll{
	Five:        TypeShapeField{5, SCORE_FIVE, []string{IIIII}},
	Four:        TypeShapeField{4, SCORE_FOUR, []string{OIIIIO}},
	DoubleFour:  TypeShapeField{44, SCORE_FOUR_FOUR, nil},
	FourThree:   TypeShapeField{43, SCORE_FOUR_THREE, nil},
	DoubleThree: TypeShapeField{33, SCORE_THREE_THREE, nil},
	BlockFour:   TypeShapeField{40, SCORE_BLOCK_FOUR, []string{IIOII, IIIOI, ZIIIIO, ZIIIOI, ZIIOII, ZIOIII}},
	Three:       TypeShapeField{3, SCORE_THREE, []string{OIIIOO, OIIOIO}},
	BlockThree:  TypeShapeField{30, SCORE_BLOCK_THREE, []string{ZIIIOO, ZIIOIO, ZIOIIO}},
	DoubleTwo:   TypeShapeField{22, SCORE_TWO_TWO, nil},
	Two:         TypeShapeField{2, SCORE_TWO, []string{OOIIOO, OIIOOO, OOOIIO, OIOIOO, OOIOIO}},
	BlockTwo:    TypeShapeField{20, SCORE_BLOCK_TWO, []string{OOIIIZ, OIOIIZ, OIIOIZ}},
	None:        TypeShapeField{0, SCORE_NONE, nil},
}

var DirectionVec = map[TypeDirection]Point{HORIZONTAL: {0, 1}, VERTICAL: {1, 0}, DIAGONAL: {1, 1}, ANTI_DIAGONAL: {1, -1}}
var DirectionEnum = []TypeDirection{HORIZONTAL, VERTICAL, DIAGONAL, ANTI_DIAGONAL}
