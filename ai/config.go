package ai

// Config 一些全局配置放在这里，其中有一些配置是用来控制一些不稳定的功能是否开启的，比如缓存，只搜索一条线上的点位等。
type Config struct {
	EnableCache    bool // 是否开启缓存
	OnlyInLine     bool // 是否只搜索一条线上的点位，一种优化方式。
	InlineCount    int  // 最近多少个点位能算作
	InLineDistance int  // 判断点位是否在一条线上的最大距离
}

// NewConfig 创建一个新的配置
func NewConfig() *Config {
	return &Config{
		EnableCache:    true,
		OnlyInLine:     false,
		InlineCount:    4,
		InLineDistance: 5,
	}
}

var config = NewConfig()

var directions = &Directions{
	Horizontal:   Point{X: 0, Y: 1},  // 水平 -
	Vertical:     Point{X: 1, Y: 0},  // 垂直 |
	Diagonal:     Point{X: 1, Y: 1},  // 斜线 /
	AntiDiagonal: Point{X: 1, Y: -1}, // 反斜线 \
}

var scores = &Score{
	Five:       10000000,
	BlockFive:  10000000,
	Four:       100000,
	FourFour:   100000,
	FourThree:  100000,
	ThreeThree: 50000,
	BlockFour:  1500,
	Three:      1000,
	BlockThree: 150,
	TwoTwo:     200,
	Two:        100,
	BlockTwo:   15,
	One:        10,
	BlockOne:   1,
	None:       0,
}

type Role struct {
	Black int
	White int
}

var role = &Role{
	Black: 1,
	White: -1,
}
