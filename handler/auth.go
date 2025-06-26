package handler

import (
	"html/template"
	"makehugocontent/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var authTmpl = template.Must(template.ParseFiles(
	"templates/login.html",
	"templates/register.html"),
)

// 登录
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		authTmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	user := r.FormValue("username")
	pass := r.FormValue("password")

	var hash string
	err := database.DB.QueryRow("SELECT password FROM users WHERE username=?", user).Scan(&hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) != nil {
		authTmpl.ExecuteTemplate(w, "login.html", "用户名或密码错误")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: user, Path: "/"})
	http.Redirect(w, r, "/admin/upload", http.StatusSeeOther)
}

// 注册
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		authTmpl.ExecuteTemplate(w, "register.html", nil)
		return
	}
	user := r.FormValue("username")
	pass := r.FormValue("password")
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	_, err := database.DB.Exec("INSERT INTO users(username,password,role)VALUES(?,?,?)", user, hash, "editor")
	if err != nil {
		authTmpl.ExecuteTemplate(w, "register.html", "注册失败："+err.Error())
		return
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
