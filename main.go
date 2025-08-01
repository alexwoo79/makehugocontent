package main

import (
	"html/template"
	"log"
	"makehugocontent/database"
	"makehugocontent/internal/admin"
	"makehugocontent/router"
	"net/http"
	// "makehugocontent/internal/admin"
)

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func main() {
	// 1. 初始化数据库
	database.Init()
	database.InsertTestUsers()

	tmpl := template.New("").Funcs(template.FuncMap{
		"contains": contains,
	})
	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))

	// 2. 使用自定义 mux

	adminHandler := &admin.AdminHandler{UserDB: database.UserDB, Tmpl: tmpl}
	mux := http.NewServeMux()

	// 3. 注册静态资源
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("static/uploads"))))
	// 后台页面资源
	mux.Handle("/admin/css/", http.StripPrefix("/admin/css/", http.FileServer(http.Dir("static/css"))))
	mux.Handle("/admin/js/", http.StripPrefix("/admin/js/", http.FileServer(http.Dir("static/js"))))
	mux.Handle("/admin/webfonts/", http.StripPrefix("/admin/webfonts/", http.FileServer(http.Dir("static/webfonts"))))
	// 4. 注册 Hugo 页面
	mux.Handle("/", http.FileServer(http.Dir("hugo/public")))

	// 5. 注册后台业务路由
	adminHandler.RegisterRoutes(mux)

	// 6. 注册业务路由（将模板和数据库注入）
	router.RegisterRoutes(mux, tmpl, database.UserDB, database.DataDB)

	// 4. 启动服务器
	log.Println("Server started at :4000")
	http.ListenAndServe("0.0.0.0:4000", mux)
}
