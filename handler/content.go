package handler

import (
	"log"
	"makehugocontent/database"
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

	// Step 1: 获取全部内容行（含 FrontMatter）
	rows, err := utils.ScanDir(CONTENT_PATH)
	if err != nil {
		http.Error(w, "内容读取失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = utils.SyncToDB(rows, database.DB) // ⚠️ 引入你的 database 包
	// Step 2: 排序（按时间ByDate,ByTitle,ByAuthor）
	utils.SortRows(rows, utils.ByDate, true)

	// Step 3: 转换成模板需要的数据结构（你原本使用 contentRow）
	var list []contentRow
	for _, r := range rows {
		date := r.FM.ParsedDate.Format("2006-01-02 15:04")
		if r.FM.ParsedDate.IsZero() {
			fi, _ := os.Stat(filepath.Join(CONTENT_PATH, r.FileName))
			date = fi.ModTime().Format("2006-01-02 15:04")
		}
		list = append(list, contentRow{
			Title:  r.FM.Title,
			Path:   strings.TrimPrefix(r.FileName, CONTENT_PATH+"/"),
			Author: r.FM.Author,
			Date:   date,
		})
	}

	// 渲染模板
	templatePath := "content_list.html"
	log.Printf("Rendering template from: %s", templatePath)
	utils.Render(w, templatePath, list)
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
	newContent := r.FormValue("content")

	// 读取原始文件内容以提取 Front Matter
	fullPath := filepath.Join(CONTENT_PATH, file)
	originalData, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "读取原始文件失败", 500)
		return
	}

	// 检测 Front Matter 分隔符
	originalLines := strings.Split(string(originalData), "\n")
	var delimiter string
	if len(originalLines) > 0 {
		switch originalLines[0] {
		case "---":
			delimiter = "---"
		case "+++":
			delimiter = "+++"
		}
	}

	var frontMatter string
	if delimiter != "" && len(originalLines) >= 3 {
		// 提取 Front Matter 内容
		for i := 1; i < len(originalLines); i++ {
			if originalLines[i] == delimiter {
				frontMatter = strings.Join(originalLines[1:i], "\n")
				break
			}
		}
	}

	// 重组内容，保留 Front Matter 和分隔符
	var updatedContent string
	if delimiter != "" && frontMatter != "" {
		// 保留原有 Front Matter 并更新正文
		updatedContent = delimiter + "\n" + frontMatter + "\n" + delimiter + "\n" + newContent
	} else {
		// 如果没有 Front Matter，则直接写入新内容
		updatedContent = newContent
	}

	// 写入更新后的内容
	err = os.WriteFile(fullPath, []byte(updatedContent), 0644)
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
