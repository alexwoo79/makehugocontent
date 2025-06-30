package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type AdminHandler struct {
	Tmpl   *template.Template
	UserDB *sql.DB
}

type User struct {
	Username string
	Role     string
}

func (a *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/users", a.usersHandler)
	mux.HandleFunc("/admin/update-role", a.updateRoleHandler)
}

func (a *AdminHandler) usersHandler(w http.ResponseWriter, r *http.Request) {
	_, role := getCurrentUser(r)
	if role != "admin" {
		http.Error(w, "管理员专用", http.StatusForbidden)
		return
	}

	rows, err := a.UserDB.Query("SELECT username, role FROM users ORDER BY username")
	if err != nil {
		http.Error(w, "数据库查询失败", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.Username, &u.Role)
		users = append(users, u)
	}

	err = a.Tmpl.ExecuteTemplate(w, "users.html", users)
	if err != nil {
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
}

func (a *AdminHandler) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	// 限制为 POST 方法
	if r.Method != http.MethodPost {
		http.Error(w, "只允许 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	// 确保是 HTMX 请求（防止浏览器或非预期客户端直接 POST）
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "非法请求：仅支持 HTMX", http.StatusBadRequest)
		return
	}

	// 权限校验
	_, role := getCurrentUser(r)
	if role != "admin" {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}

	// 获取并验证表单参数
	if err := r.ParseForm(); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	newRole := r.FormValue("role")
	if username == "" || newRole == "" {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 执行数据库更新
	_, err := a.UserDB.Exec("UPDATE users SET role=? WHERE username=?", newRole, username)
	if err != nil {
		http.Error(w, "数据库更新失败", http.StatusInternalServerError)
		return
	}

	// 返回更新后的表单 HTML，HTMX 自动替换页面
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<form hx-post="/admin/update-role" hx-target="this" hx-swap="outerHTML">
			<input type="hidden" name="username" value="%s">
			<select name="role" class="border rounded p-1">
				<option value="editor" %s>editor</option>
				<option value="viewer" %s>viewer</option>
				<option value="admin" %s>admin</option>
			</select>
			<button class="ml-2 px-2 py-1 bg-green-600 text-white rounded">更新</button>
		</form>`,
		username,
		boolToSelected(newRole == "editor"),
		boolToSelected(newRole == "viewer"),
		boolToSelected(newRole == "admin"),
	)
}

func boolToSelected(ok bool) string {
	if ok {
		return "selected"
	}
	return ""
}

func getCurrentUser(r *http.Request) (string, string) {
	c, err := r.Cookie("session")
	if err != nil {
		return "", ""
	}
	parts := strings.Split(c.Value, "|")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
