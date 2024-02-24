package ai

// Config 一些全局配置放在这里，其中有一些配置是用来控制一些不稳定的功能是否开启的，比如缓存，只搜索一条线上的点位等。
type Config struct {
	enableCache    bool // 是否开启缓存
	onlyInLine     bool // 是否只搜索一条线上的点位，一种优化方式。
	inlineCount    int  // 最近多少个点位能算作
	inLineDistance int  // 判断点位是否在一条线上的最大距离
}

// NewConfig 创建一个新的配置
func NewConfig() *Config {
	return &Config{
		enableCache:    true,
		onlyInLine:     false,
		inlineCount:    4,
		inLineDistance: 5,
	}
}

var config = NewConfig()
