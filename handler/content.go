package handler

import (
	"makehugocontent/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	_ = filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, _ := os.ReadFile(path)
		fm := utils.ExtractFrontMatter(data)
		rows = append(rows, contentRow{
			Title:  fm["title"],
			Path:   strings.TrimPrefix(path, "content/"),
			Author: fm["author"],
			Date:   info.ModTime().Format("2006-01-02 15:04"),
		})
		return nil
	})

	utils.Render(w, "content_list.html", rows)
}

func EditPageHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	file := r.URL.Query().Get("file")
	full := filepath.Join("content", file)
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
	err := os.WriteFile(filepath.Join("content", file), []byte(content), 0644)
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
	os.Remove(filepath.Join("content", file))
	utils.RunHugo()
	http.Redirect(w, r, "/admin/content", http.StatusSeeOther)
}
