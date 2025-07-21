package utils

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// 定义模板路径为项目根目录下的 templates 文件夹
const TEMPLATE_DIR = "templates"

func Render(w http.ResponseWriter, templateName string, data interface{}) {
	// 获取当前工作目录用于调试
	pwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", pwd)

	// 构造完整的模板路径并检查其是否存在
	templatePath := filepath.Join(TEMPLATE_DIR, templateName)
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("模板文件不存在: %s", templatePath)
		http.Error(w, "模板文件不存在", http.StatusInternalServerError)
		return
	}

	absTemplatePath, _ := filepath.Abs(templatePath)
	log.Printf("正在加载模板: %s", absTemplatePath)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("模板加载失败: %v", err)
		http.Error(w, "模板加载失败", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}