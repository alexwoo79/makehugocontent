package handler

import (
	"database/sql"
	"html/template"
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

type ContentHandler struct {
	Tmpl   *template.Template
	DataDB *sql.DB
}

type contentRow struct {
	Title  string
	Path   string
	Author string
	Date   string
}

func (c *ContentHandler) ContentListHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	if !checkRole(r, []string{"admin", "editor"}) {
		http.Error(w, "管理员专用", http.StatusSeeOther)
		return
	}

	rows, err := utils.ScanDir(CONTENT_PATH)
	if err != nil {
		http.Error(w, "内容读取失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = utils.SyncToDB(rows, database.DataDB)

	// 默认按时间降序
	utils.SortRows(rows, utils.ByDate, true)

	var list []contentRow
	for _, r := range rows {
		// ✅ 1. 使用 FM.Date 而不是 ParsedDate
		date := r.FM.ParsedDate.Format("2006-01-02 15:04")
		if r.FM.ParsedDate.IsZero() {
			fi, _ := os.Stat(filepath.Join(CONTENT_PATH, r.FileName))
			date = fi.ModTime().Format("2006-01-02 15:04")
		}

		// ✅ 2. r.FileName 已是相对文件名，直接用即可
		list = append(list, contentRow{
			Title:  r.FM.Title,
			Path:   r.FileName, // 例如 "post1.md"
			Author: r.FM.Author,
			Date:   date,
		})
	}
	// ✅ 使用统一模板直接渲染 list 数据
	err = c.Tmpl.ExecuteTemplate(w, "content_list.html", list)
	if err != nil {
		log.Println("模板渲染失败:", err)
		http.Error(w, "页面渲染失败", http.StatusInternalServerError)
	}
}

func (c *ContentHandler) EditPageHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	file := r.URL.Query().Get("file")
	full := filepath.Join(CONTENT_PATH, file)

	// 检查文件是否存在
	if _, err := os.Stat(full); os.IsNotExist(err) {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}

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

func (c *ContentHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	if !checkRole(r, []string{"admin", "editor"}) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	file := r.FormValue("filepath")
	newContent := r.FormValue("content")
	newContent = strings.ReplaceAll(newContent, "\r\n", "\n")
	newContent = strings.ReplaceAll(newContent, "\r", "\n")

	fullPath := filepath.Join(CONTENT_PATH, file)

	// 将编辑器中的全部内容写入文件
	err := os.WriteFile(fullPath, []byte(newContent), 0644)
	if err != nil {
		http.Error(w, "保存失败", http.StatusInternalServerError)
		return
	}

	// 读取保存后的文件以检查 Front Matter 分隔符
	savedData, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "检查保存文件失败", http.StatusInternalServerError)
		return
	}

	savedLines := strings.Split(string(savedData), "\n")
	if len(savedLines) < 3 {
		http.Error(w, "保存后文件格式错误，缺少 Front Matter", http.StatusInternalServerError)
		return
	}

	// 检查第一个分隔符，忽略空行和空白
	savedStart := -1
	var savedDelimiter string
	for i, line := range savedLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" || trimmed == "+++" {
			savedDelimiter = trimmed
			savedStart = i
			break
		}
		if trimmed != "" {
			http.Error(w, "保存后文件开头缺少有效 Front Matter 分隔符", http.StatusInternalServerError)
			return
		}
	}

	if savedStart == -1 || savedDelimiter == "" {
		http.Error(w, "保存后未找到有效的 Front Matter 分隔符", http.StatusInternalServerError)
		return
	}

	// 检查第二个分隔符
	savedEnd := -1
	for i := savedStart + 1; i < len(savedLines); i++ {
		if strings.TrimSpace(savedLines[i]) == savedDelimiter {
			savedEnd = i
			break
		}
	}

	if savedEnd == -1 || savedEnd <= savedStart {
		http.Error(w, "保存后 Front Matter 结构错误：找不到结束分隔符", http.StatusInternalServerError)
		return
	}

	// 执行 Hugo 构建
	utils.RunHugo()

	// 重定向到内容列表页面
	http.Redirect(w, r, "/admin/content", http.StatusSeeOther)
}

func (c *ContentHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	file := r.URL.Query().Get("file")
	os.Remove(filepath.Join(CONTENT_PATH, file))
	utils.RunHugo()
	http.Redirect(w, r, "/admin/content", http.StatusSeeOther)
}
