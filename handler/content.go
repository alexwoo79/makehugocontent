package handler

import (
	"log"
	"makehugocontent/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 定义内容路径为常量
const (
	CONTENT_PATH = "hugo/content/posts"
)

type contentRow struct {
	Title  string
	Path   string
	Author string
	Date   string
}

func ContentListHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}

	var rows []contentRow
	_ = filepath.Walk(CONTENT_PATH, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, _ := os.ReadFile(path)
		fm := utils.ExtractFrontMatter(data)
		rows = append(rows, contentRow{
			Title:  fm["title"],
			Path:   strings.TrimPrefix(path, CONTENT_PATH+"/"),
			Author: fm["author"],
			Date:   info.ModTime().Format("2006-01-02 15:04"),
		})
		return nil
	})

	// 确保使用正确的模板路径
	templatePath := "content_list.html"
	log.Printf("Rendering template from: %s", templatePath)
	utils.Render(w, templatePath, rows)
}

func EditPageHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	file := r.URL.Query().Get("file")
	full := filepath.Join(CONTENT_PATH, file)
	data, err := os.ReadFile(full)
	if err != nil {
		http.Error(w, "读取失败", 500)
		return
	}
	utils.Render(w, "edit.html", map[string]string{
		"FilePath": file,
		"Content":  string(data),
	})
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	r.ParseForm()
	file := r.FormValue("filepath")
	content := r.FormValue("content")
	err := os.WriteFile(filepath.Join(CONTENT_PATH, file), []byte(content), 0644)
	if err != nil {
		http.Error(w, "保存失败", 500)
		return
	}
	utils.RunHugo()
	http.Redirect(w, r, "/admin/content", http.StatusSeeOther)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	file := r.URL.Query().Get("file")
	os.Remove(filepath.Join(CONTENT_PATH, file))
	utils.RunHugo()
	http.Redirect(w, r, "/admin/content", http.StatusSeeOther)
}
