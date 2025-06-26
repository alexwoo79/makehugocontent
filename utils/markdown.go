package utils

import (
	"gopkg.in/yaml.v2"
	"strings"
	"log"
)

// 解析 front-matter，返回键值（支持 --- 和 +++ 分隔符）
func ExtractFrontMatter(data []byte) map[string]string {
	m := make(map[string]interface{}) // 使用 interface{} 以支持 YAML 结构
	lines := strings.Split(string(data), "\n")
	if len(lines) < 3 {
		return convertMapInterfaceToString(m)
	}

	// 检测使用的分隔符（--- 或 +++）
	var delimiter string
	switch lines[0] {
	case "---":
		delimiter = "---"
	case "+++":
		delimiter = "+++"
	default:
		return convertMapInterfaceToString(m)
	}

	// 提取分隔符之间的内容作为 YAML 数据
	var yamlData []byte
	for _, line := range lines[1:] {
		if line == delimiter {
			break
		}
		yamlData = append(yamlData, []byte(line+"\n")...)
	}

	// 使用 YAML 解析器解析数据
	if err := yaml.Unmarshal(yamlData, &m); err != nil {
		log.Printf("YAML 解析失败: %v", err)
		return convertMapInterfaceToString(m)
	}

	return convertMapInterfaceToString(m)
}

// 将 map[string]interface{} 转换为 map[string]string
func convertMapInterfaceToString(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if val, ok := v.(string); ok {
			result[k] = val
		} else {
			result[k] = ""
		}
	}
	return result
}