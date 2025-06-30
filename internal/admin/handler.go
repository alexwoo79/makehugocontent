package admin

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type AdminHandler struct {
	DB   *sql.DB
	Tmpl *template.Template
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
	user, role := getCurrentUser(r)
	if role != "admin" {
		http.Error(w, "管理员专用", http.StatusForbidden)
		return
	}

	rows, err := a.DB.Query("SELECT username, role FROM users ORDER BY username")
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
	a.Tmpl.ExecuteTemplate(w, "users.html", users)
}

func (a *AdminHandler) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	user, role := getCurrentUser(r)
	if role != "admin" {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}
	r.ParseForm()
	targetUser := r.FormValue("username")
	newRole := r.FormValue("role")

	_, err := a.DB.Exec("UPDATE users SET role=? WHERE username=?", newRole, targetUser)
	if err != nil {
		http.Error(w, "数据库更新失败", http.StatusInternalServerError)
		return
	}

	// 立即返回新的表单片段 (HTMX)
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
		targetUser,
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
