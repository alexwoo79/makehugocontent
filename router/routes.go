package router

import (
	"database/sql"
	"html/template"
	"makehugocontent/handler"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, tmpl *template.Template, userDB, dataDB *sql.DB) {
	// 1. 构造 handler 实例
	// Auth
	authHandler := &handler.AuthHandler{
		Tmpl:   tmpl,
		UserDB: userDB,
	}
	uploadHandler := &handler.UploadHandler{
		Tmpl:   tmpl,
		DataDB: dataDB,
	}
	contentHandler := &handler.ContentHandler{
		Tmpl:   tmpl,
		DataDB: dataDB,
	}
	mux.HandleFunc("/admin/login", authHandler.LoginHandler)
	mux.HandleFunc("/admin/register", authHandler.RegisterHandler)

	// Upload
	mux.HandleFunc("/admin/upload", uploadHandler.UploadPageHandler)
	mux.HandleFunc("/api/upload-image", uploadHandler.ImageUploadHandler)
	mux.HandleFunc("/admin/logout", uploadHandler.LogoutHandler)

	// Content 管理
	mux.HandleFunc("/admin/content", contentHandler.ContentListHandler)
	mux.HandleFunc("/admin/edit", contentHandler.EditPageHandler)
	mux.HandleFunc("/admin/update", contentHandler.UpdateHandler)
	mux.HandleFunc("/admin/delete", contentHandler.DeleteHandler)
	mux.HandleFunc("/admin/export", handler.ExportCSVHandler)
}
