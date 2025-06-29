package router

import (
	"makehugocontent/handler"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	// Auth
	mux.HandleFunc("/admin/login", handler.LoginHandler)
	mux.HandleFunc("/admin/register", handler.RegisterHandler)

	// Upload
	mux.HandleFunc("/admin/upload", handler.UploadPageHandler)
	mux.HandleFunc("/api/upload-image", handler.ImageUploadHandler)
	mux.HandleFunc("/admin/logout", handler.LogoutHandler)

	// Content 管理
	mux.HandleFunc("/admin/content", handler.ContentListHandler)
	mux.HandleFunc("/admin/edit", handler.EditPageHandler)
	mux.HandleFunc("/admin/update", handler.UpdateHandler)
	mux.HandleFunc("/admin/delete", handler.DeleteHandler)
	mux.HandleFunc("/admin/export", handler.ExportCSVHandler)
}
