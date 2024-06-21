package internal

// 配置文件解析器映射表
var ConfigFileParserMap = map[uint32]FileParserFN{
	1: ConfigParser_Excel,
}
