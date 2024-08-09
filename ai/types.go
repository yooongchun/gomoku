package ai

type TypeChess int
type TypeRole int
type TypeDirection int
type TypeScoreCache map[TypeChess]map[TypeDirection][][]int            // [chess][direction][x][y] -> score
type TypeShapeCache map[TypeChess]map[TypeDirection][][]TypeShapeField // [chess][direction][x][y] -> shape

type Point struct {
	x int
	y int
}

type TypeHistory struct {
	point Point
	chess TypeChess
}

type TypeEvaluateCache struct {
	chess TypeChess
	score int
}

type TypeValuableMoveCache struct {
	role      TypeChess
	moves     []Point
	depth     int
	onlyThree bool
	onlyFour  bool
}

type TypeShapeField struct {
	Code  int
	Score int
	Name  []string
}

// TypeShapeAll 可取的形状
type TypeShapeAll struct {
	Five        TypeShapeField //11111
	Four        TypeShapeField //011110
	DoubleFour  TypeShapeField
	FourThree   TypeShapeField
	DoubleThree TypeShapeField
	RushFour    TypeShapeField //10111|11011|11101|211110|211101|211011|210111|011112|101112|110112|111012
	Three       TypeShapeField //011100|011010|010110|001110
	BlockThree  TypeShapeField //211100|211010|210110|001112|010112|011012
	DoubleTwo   TypeShapeField
	Two         TypeShapeField //001100|011000|000110|010100|001010
	BlockTwo    TypeShapeField
	One         TypeShapeField
	BlockOne    TypeShapeField
	None        TypeShapeField
}
