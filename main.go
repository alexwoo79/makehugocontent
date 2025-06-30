package main

import (
	"html/template"
	"log"
	"makehugocontent/database"
	"makehugocontent/router"
	"net/http"
)

func main() {
	// 1. 初始化数据库
	database.Init()

	// 2. 使用自定义 mux
	mux := http.NewServeMux()

	tmpl := template.Must(template.P)

	// 3. 注册静态资源
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("hugo/public/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("hugo/public/js"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("hugo/public/images"))))

	// 后台页面资源
	mux.Handle("/admin/css/", http.StripPrefix("/admin/css/", http.FileServer(http.Dir("static/css"))))
	mux.Handle("/admin/js/", http.StripPrefix("/admin/js/", http.FileServer(http.Dir("static/js"))))

	// 4. 注册 Hugo 页面
	mux.Handle("/", http.FileServer(http.Dir("hugo/public")))

	// 5. 注册后台业务路由
	router.RegisterRoutes(mux)

	// 6. 启动服务
	log.Println("Serving Hugo site at http://localhost:8080/")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
