package utils

import (
	// "embed"

	"net/http"
)

// go:embed "html/templates"
// var templateFS embed.FS

func Render(w http.ResponseWriter, templateName string, data interface{}) {
	// 构建模板路径
	// patterns := []string{
	// 	"../templates/base.html",
	// 	"../templates/partials/navbar.html",
	// 	filepath.Join("../templates", templateName),
	// }

	// // 使用嵌入式文件系统加载模板
	// tmpl, err := template.New("").ParseFS(templateFS, patterns...)
	// if err != nil {
	// 	log.Printf("模板加载失败: %v", err)
	// 	http.Error(w, "模板加载失败", http.StatusInternalServerError)
	// 	return
	// }

	// // 执行默认模板（假设每个页面模板都有一个 "content" 区块）
	// err = tmpl.ExecuteTemplate(w, "content", data)
	// if err != nil {
	// 	log.Printf("模板渲染失败: %v", err)
	// 	http.Error(w, "模板渲染失败", http.StatusInternalServerError)
	// }
}
