package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// AdminHandler 管理用户页面
type AdminHandler struct {
	UserDB *sql.DB
	Tmpl   *template.Template
}

// 用户结构体
type User struct {
	Username     string
	Role         string
	Departments  []string
	SelectedDept map[string]bool
}

// 部门结构体
type Department struct {
	ID   int
	Name string
}

// 注册路由
func (a *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/users", a.usersHandler)
	mux.HandleFunc("/admin/update-role", a.updateRoleHandler)
}

// 用户列表页面
func (a *AdminHandler) usersHandler(w http.ResponseWriter, r *http.Request) {
	_, role := currentUser(r)
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
		if err := rows.Scan(&u.Username, &u.Role); err != nil {
			continue
		}
		depts, _ := a.getDepartmentsForUser(u.Username)
		u.Departments = depts
		u.SelectedDept = make(map[string]bool)
		for _, d := range depts {
			u.SelectedDept[d] = true
		}
		users = append(users, u)
	}

	allDepts, _ := a.getAllDepartments()

	data := map[string]any{
		"Users":       users,
		"Departments": allDepts,
	}
	err = a.Tmpl.ExecuteTemplate(w, "users.html", data)
	if err != nil {
		http.Error(w, "模板渲染失败: "+err.Error(), http.StatusInternalServerError)
	}
}

// 更新用户角色和部门
func (a *AdminHandler) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "非法请求", http.StatusMethodNotAllowed)
		return
	}

	_, role := currentUser(r)
	if role != "admin" {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}
	r.ParseForm()

	username := r.FormValue("username")
	newRole := r.FormValue("role")
	selectedDepts := r.Form["departments"]

	_, err := a.UserDB.Exec("UPDATE users SET role=? WHERE username=?", newRole, username)
	if err != nil {
		http.Error(w, "角色更新失败", http.StatusInternalServerError)
		return
	}

	_, _ = a.UserDB.Exec(`
		DELETE FROM user_departments 
		WHERE user_id=(SELECT id FROM users WHERE username=?)`, username)

	for _, deptName := range selectedDepts {
		_, _ = a.UserDB.Exec(`
			INSERT OR IGNORE INTO user_departments(user_id, dept_id)
			SELECT u.id, d.id 
			FROM users u, departments d
			WHERE u.username=? AND d.name=?`,
			username, deptName)
	}
	depts, _ := a.getDepartmentsForUser(username)
	// deptsStr := strings.Join(depts, ", ")
	// Compute selected attributes for role dropdown
	adminSelected := ""
	userSelected := ""
	if newRole == "admin" {
		adminSelected = "selected"
	}
	if newRole == "user" {
		userSelected = "selected"
	}

	// Return the complete <tr> HTML
	fmt.Fprintf(w, `
<tr hx-target="this" hx-swap="outerHTML">
  <form hx-post="/admin/update-role" class="flex items-center gap-2">
    <td class="border px-4 py-2">%s</td>
    <td class="border px-4 py-2">
      <input type="hidden" name="username" value="%s">
      <select name="role" class="border rounded p-1">
        <option value="admin" %s>Admin</option>
        <option value="user" %s>User</option>
      </select>
    </td>
    <td class="border px-4 py-2">
      <select name="departments[]" multiple size="3" class="border rounded p-1">
        %s
      </select>
    </td>
    <td class="border px-4 py-2">
      <button class="px-3 py-1 bg-blue-600 text-white rounded">保存</button>
    </td>
  </form>
</tr>`,
		username, username,
		adminSelected, userSelected,
		a.generateDeptOptions(depts))
}

// 查询某用户所属部门
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
		if err := rows.Scan(&name); err == nil {
			list = append(list, name)
		}
	}
	return list, nil
}

// 获取所有部门
func (a *AdminHandler) getAllDepartments() ([]Department, error) {
	rows, err := a.UserDB.Query("SELECT id, name FROM departments ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Department
	for rows.Next() {
		var d Department
		if err := rows.Scan(&d.ID, &d.Name); err == nil {
			list = append(list, d)
		}
	}
	return list, nil
}

// 从 cookie 中获取当前用户信息
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

// 生成部门选项 HTML
func (a *AdminHandler) generateDeptOptions(selectedDepts []string) string {
	depts, _ := a.getAllDepartments()
	var options strings.Builder
	for _, dept := range depts {
		selected := ""
		for _, sel := range selectedDepts {
			if sel == dept.Name {
				selected = "selected"
				break
			}
		}
		options.WriteString(fmt.Sprintf(`<option value="%s" %s>%s</option>`, dept.Name, selected, dept.Name))
	}
	return options.String()
}
