package ai

// ConfigEnum 一些全局配置放在这里，其中有一些配置是用来控制一些不稳定的功能是否开启的，比如缓存，只搜索一条线上的点位等。
type ConfigEnum struct {
	EnableCache    bool // 是否开启缓存
	OnlyInLine     bool // 是否只搜索一条线上的点位，一种优化方式。
	InlineCount    int  // 最近多少个点位能算作
	InLineDistance int  // 判断点位是否在一条线上的最大距离
}

// NewConfig 创建一个新的配置
func NewConfig() *ConfigEnum {
	return &ConfigEnum{
		EnableCache:    true,
		OnlyInLine:     false,
		InlineCount:    4,
		InLineDistance: 5,
	}
}

var Config = NewConfig()
