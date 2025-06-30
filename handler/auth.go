package handler

import (
	"database/sql"
	"html/template"
	"makehugocontent/database"
	"net/http"

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

	http.Redirect(w, r, "/admin/content", http.StatusSeeOther)
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

	_, err := database.UserDB.Exec("INSERT INTO users(username,password,role)VALUES(?,?,?)", user, hash, "editor")
	if err != nil {
		a.Tmpl.ExecuteTemplate(w, "register.html", "注册失败："+err.Error())
		return
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
