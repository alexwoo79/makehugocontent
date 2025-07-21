package handler

import (
	"database/sql"
	"html/template"
	"makehugocontent/database"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Tmpl   *template.Template
	UserDB *sql.DB
}

// 登录
func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.Tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	user := r.FormValue("username")
	pass := r.FormValue("password")

	var hash, role string
	err := database.UserDB.QueryRow("SELECT password, role FROM users WHERE username=?", user).Scan(&hash, &role)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) != nil {
		a.Tmpl.ExecuteTemplate(w, "login.html", "用户名或密码错误")
		return
	}

	// ✅ 正确设置包含角色信息的 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: user + "|" + role,
		Path:  "/",
	})

	// 登录成功后，重定向到 /home 页面
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// 注册
func (a *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.Tmpl.ExecuteTemplate(w, "register.html", nil)
		return
	}
	user := r.FormValue("username")
	pass := r.FormValue("password")
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	_, err := database.UserDB.Exec("INSERT INTO users(username,password,role)VALUES(?,?,?)", user, hash, "viewer")
	if err != nil {
		a.Tmpl.ExecuteTemplate(w, "register.html", "注册失败："+err.Error())
		return
	}

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// 验证用户权限

func checkRole(r *http.Request, allowedRoles []string) bool {
	// 获取当前用户的 session
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	// 解析 session 内容
	parts := strings.Split(c.Value, "|")
	if len(parts) != 2 {
		return false
	}

	role := parts[1] // 获取 role

	// 判断当前用户角色是否在允许的角色列表中
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}
