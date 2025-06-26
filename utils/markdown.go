package utils

import "strings"

// 解析 +++ front-matter，返回键值
func ExtractFrontMatter(data []byte) map[string]string {
	m := map[string]string{}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 3 || lines[0] != "+++" {
		return m
	}
	for _, line := range lines[1:] {
		if line == "+++" {
			break
		}
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
			m[key] = val
		}
	}
	return m
}
