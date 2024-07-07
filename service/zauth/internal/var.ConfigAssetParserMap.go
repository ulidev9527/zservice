package internal

// 配置文件解析器映射表
var ConfigAssetParserMap = map[uint32]IAssetConfigParserFN{
	1: AssetConfigParser_Excel,
}
