package internal

import (
	"fmt"
	"strings"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"github.com/xuri/excelize/v2"
)

// 解析 Excel 文件
func ConfigParser_Excel(fullPath string) (map[string]string, *zservice.Error) {
	// 解析
	valueMap := map[string]string{}
	if e := func() *zservice.Error {

		// 更新配置
		excel, e := excelize.OpenFile(fullPath)
		if e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_OpenFileErr)
		}
		defer func() {
			if e := excel.Close(); e != nil {
				zservice.LogError(e)
			}
		}()
		if excel.SheetCount == 0 {
			return zservice.NewError("excel sheet count is 0").SetCode(zglobal.Code_Zauth_config_ExcelNoContent)
		}

		sheets := excel.GetSheetList()
		rows, e := excel.Rows(sheets[0])
		if e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_Zauth_config_ParserFail)
		}

		// 开始循环 excel
		rowsIndex := 0 // 当前行
		// exportType := []string{}  // 导出类型 当前需要导出全部配置，禁用导出类型处理
		valueTypes := []string{}  // 值类型
		valueFields := []string{} // 值字段名
		for rows.Next() {
			rowsIndex++
			cols, e := rows.Columns()

			if e != nil {
				return zservice.NewError(e).SetCode(zglobal.Code_Zauth_config_ParserFail)
			}
			if len(cols) == 0 {
				continue
			}

			// 第一行是字段说明，不处理
			// 第二行是字段描述
			// 第三行是自动导出方案
			// 第四行是字段类型
			// 第五行是字段名
			// 第六行开始为数据
			if rowsIndex <= 2 { // 前三行不处理
				continue
			}

			if rowsIndex == 3 {
				// exportType = cols
				continue
			}

			if rowsIndex == 4 {
				valueTypes = cols
				continue
			}
			if rowsIndex == 5 {
				valueFields = cols
				continue
			}

			data := map[string]any{}
			count := len(cols)
			for i := 0; i < count; i++ {
				// 导出类型判定
				// 默认没有填写导出，c 表示客户端专用，不导出
				// if len(exportType) > i && strings.ToLower(exportType[i]) == "c" {
				// 	// 注释忽略导出类型
				// 	// continue
				// }

				// 字段获取
				if len(valueFields) <= i { // 字段不够参数长度, 不继续导出
					break
				}
				field := valueFields[i]
				if field == "" { // 字段为空，不导出
					continue
				}

				// 类型获取
				vtype := "string"
				if len(valueTypes) > i {
					vtype = valueTypes[i]
				}

				value := cols[i]

				if vtype == "number" {
					if strings.Contains(value, ".") {
						vtype = "float"
					} else {
						vtype = "int"
					}
				}
				// 去掉空格
				vtype = strings.ToLower(strings.ReplaceAll(vtype, " ", ""))

				// 类型转换
				switch vtype {
				case "int", "int32", "int64":
					data[field] = zservice.StringToInt64(value)
				case "float", "float32", "float64":
					data[field] = zservice.StringToFloat64(value)
				case "int[]", "int32[]", "int64[]":
					arr := strings.Split(value, ",")
					list := []int64{}
					for i := 0; i < len(arr); i++ {
						list = append(list, zservice.StringToInt64(arr[i]))
					}
					data[field] = list
				case "float[]", "float32[]", "float64[]":
					arr := strings.Split(value, ",")
					list := []float64{}
					for i := 0; i < len(arr); i++ {
						list = append(list, zservice.StringToFloat64(arr[i]))
					}
					data[field] = list

				case "bool":
					data[field] = zservice.StringToBoolean(value)
				case "string[]":
					data[field] = zservice.StringSplit(value, ",")
				default: // 默认 string 类型
					data[field] = value
				}
				continue
			}
			if len(data) > 0 {
				valueMap[fmt.Sprint(data["id"])] = string(zservice.JsonMustMarshal(data))
			}

		}

		return nil
	}(); e != nil {
		return nil, e
	}

	return valueMap, nil
}
