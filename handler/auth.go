package handler

import (
	"html/template"
	"log"
	"net/http"
)

// 登录
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodGet {
	// 加载模板文件
	files := []string{
		"./html/base.tmpl",
		"./html/pages/login.tmpl",
		"./html/partials/navbar.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	// return
	// }

	// POST 请求：获取表单值
	if err := r.ParseForm(); err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("Received username: %s, password: %s", username, password)

	// 检查用户名和密码是否为空
	if username == "" || password == "" {
		log.Println("Username or password is empty")
		http.Error(w, "用户名或密码不能为空", http.StatusBadRequest)
		return
	}

	// 这里可以添加验证逻辑，比如校验用户名和密码是否正确
	// 示例重定向到管理页面
	http.Redirect(w, r, "/admin/upload", http.StatusSeeOther)
}

// 注册
// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		authTmpl.ExecuteTemplate(w, "register.html", nil)
// 		return
// 	}
// 	user := r.FormValue("username")
// 	pass := r.FormValue("password")
// 	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
//
// 	_, err := database.DB.Exec("INSERT INTO users(username,password,role)VALUES(?,?,?)", user, hash, "editor")
// 	if err != nil {
// 		authTmpl.ExecuteTemplate(w, "register.html", "注册失败："+err.Error())
// 		return
// 	}
// 	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
// }
