package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// AdminHandler 处理 users.db 数据库中的用户管理相关操作
type AdminHandler struct {
	UserDB *sql.DB
	Tmpl   *template.Template
}

/* ---------- 数据模型 ---------- */

type User struct {
	Username    string
	Role        string
	Departments []string
}

type Department struct {
	ID   int
	Name string
}

/* ---------- 路由注册 ---------- */

func (a *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/users", a.usersHandler)
	mux.HandleFunc("/admin/update-role", a.updateRoleHandler)
}

/* ---------- 页面：用户列表 ---------- */

func (a *AdminHandler) usersHandler(w http.ResponseWriter, r *http.Request) {
	_, role := currentUser(r)
	if role != "admin" {
		http.Error(w, "管理员专用", http.StatusForbidden)
		return
	}

	// 1. 查询所有用户
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
		u.Departments, _ = a.getDepartmentsForUser(u.Username)
		users = append(users, u)
	}

	// 2. 查询所有部门（渲染多选框）
	allDepts, _ := a.getAllDepartments()

	data := map[string]any{
		"Users":       users,
		"Departments": allDepts,
	}
	a.Tmpl.ExecuteTemplate(w, "users.html", data)
}

/* ---------- 接口：修改角色 / 部门 ---------- */

func (a *AdminHandler) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	_, role := currentUser(r)
	if role != "admin" {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}
	r.ParseForm()

	username := r.FormValue("username")
	newRole := r.FormValue("role")
	selectedDepts := r.Form["departments"] // 多选框 names="departments"

	// 1) 更新角色
	_, err := a.UserDB.Exec("UPDATE users SET role=? WHERE username=?", newRole, username)
	if err != nil {
		http.Error(w, "角色更新失败", http.StatusInternalServerError)
		return
	}

	// 2) 更新部门映射
	_, _ = a.UserDB.Exec(`DELETE FROM user_departments 
					  WHERE user_id=(SELECT id FROM users WHERE username=?)`, username)
	for _, deptName := range selectedDepts {
		_, _ = a.UserDB.Exec(`
			INSERT OR IGNORE INTO user_departments(user_id, dept_id)
			SELECT u.id, d.id 
			FROM users u, departments d
			WHERE u.username=? AND d.name=?`,
			username, deptName)
	}

	// 3) 返回整行片段 (HTMX)：重渲染该用户行
	deptsStr := strings.Join(selectedDepts, ", ")
	fmt.Fprintf(w, `
	<td class="border px-4 py-2">%s</td>
	<td class="border px-4 py-2">%s</td>
	<td class="border px-4 py-2">%s</td>`,
		username, newRole, deptsStr)
}

/* ---------- 数据库辅助 ---------- */

func (a *AdminHandler) getDepartmentsForUser(username string) ([]string, error) {
	rows, err := a.UserDB.Query(`
		SELECT d.name
		  FROM departments d
		  JOIN user_departments ud ON d.id = ud.dept_id
		  JOIN users u ON u.id = ud.user_id
		 WHERE u.username = ?`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		list = append(list, name)
	}
	return list, nil
}

func (a *AdminHandler) getAllDepartments() ([]Department, error) {
	rows, err := a.UserDB.Query("SELECT id, name FROM departments ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Department
	for rows.Next() {
		var d Department
		rows.Scan(&d.ID, &d.Name)
		list = append(list, d)
	}
	return list, nil
}

/* ---------- 工具函数 ---------- */

// currentUser 从 session cookie 中解析 "username|role"
func currentUser(r *http.Request) (string, string) {
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
