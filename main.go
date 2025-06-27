package main

import (
	"log"
	"makehugocontent/database"
	"makehugocontent/router"
	"net/http"
)

func main() {
	// 1. 初始化数据库
	database.Init()

	// 2. 设置hanlers
	mux := http.NewServeMux()

	// 2. 提供静态资源
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("hugo/static/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("hugo/static/js"))))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("hugo/static/uploads"))))

	// 3. Hugo 生成的静态站点
	mux.Handle("/", http.FileServer(http.Dir("hugo/public")))

	// 4. 注册业务路由
	router.RegisterRoutes(mux)

	log.Println("Server running at http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
	// log.Printf("Rendering template: %s", "content_list.html")
}
