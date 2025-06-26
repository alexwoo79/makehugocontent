package router

import (
	"makehugocontent/handler"
	"net/http"
)

func RegisterRoutes() {
	// Auth
	http.HandleFunc("/admin/login", handler.LoginHandler)
	http.HandleFunc("/admin/register", handler.RegisterHandler)

	// Upload
	http.HandleFunc("/admin/upload", handler.UploadPageHandler)
	http.HandleFunc("/api/upload-image", handler.ImageUploadHandler)
	http.HandleFunc("/admin/logout", handler.LogoutHandler)

	// Content 管理
	http.HandleFunc("/admin/content", handler.ContentListHandler)
	http.HandleFunc("/admin/edit", handler.EditPageHandler)
	http.HandleFunc("/admin/update", handler.UpdateHandler)
	http.HandleFunc("/admin/delete", handler.DeleteHandler)
}
