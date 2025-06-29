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

	// 2. 提供静态资源
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("static/data"))))
	http.Handle("/svg/", http.StripPrefix("/svg/", http.FileServer(http.Dir("static/svg"))))
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("static/lib"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("hugo/static/uploads"))))

	// 3. Hugo 生成的静态站点
	http.Handle("/", http.FileServer(http.Dir("hugo/public")))

	// 4. 注册业务路由
	router.RegisterRoutes()

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Printf("Rendering template: %s", "content_list.html")
}
