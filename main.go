package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	tmpl        = template.Must(template.ParseGlob("templates/*.html"))
	db          *sql.DB
	hugoContent = "content/posts"
	imageFolder = "static/uploads"
)

func main() {
	initDB()

	http.HandleFunc("/admin/login", loginHandler)
	http.HandleFunc("/admin/register", registerHandler)
	http.HandleFunc("/admin/upload", authMiddleware(uploadHandler))
	http.HandleFunc("/api/upload-image", authMiddleware(imageUploadHandler))

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("static/uploads"))))

	log.Println("服务启动: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	createSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);
	`
	_, err = db.Exec(createSQL)
	if err != nil {
		log.Fatal("初始化数据库失败: ", err)
	}
}

// authMiddleware 简单session检查，cookie里session值为用户名
func authMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		// 这里可以添加更复杂权限校验，简化只判断登录
		fn(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	user := r.FormValue("username")
	pass := r.FormValue("password")

	var hash string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", user).Scan(&hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) != nil {
		tmpl.ExecuteTemplate(w, "login.html", "登录失败")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: user,
		Path:  "/",
	})
	http.Redirect(w, r, "/admin/upload", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "register.html", nil)
		return
	}
	user := r.FormValue("username")
	pass := r.FormValue("password")
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	_, err := db.Exec("INSERT INTO users(username, password, role) VALUES (?, ?, ?)", user, hash, "editor")
	if err != nil {
		tmpl.ExecuteTemplate(w, "register.html", "注册失败："+err.Error())
		return
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "upload.html", nil)
		return
	}

	// 读取表单字段
	subfolder := r.FormValue("subfolder")
	filename := r.FormValue("filename")
	category := r.FormValue("category")
	draft := r.FormValue("draft") == "true"
	title := r.FormValue("title")
	author := r.FormValue("author")
	description := r.FormValue("description")
	tags := r.FormValue("tags")
	body := r.FormValue("body")

	front := fmt.Sprintf(`---
title: "%s"
date: "%s"
author: "%s"
description: "%s"
tags: [%s]
categories: ["%s"]
draft: %v
---

`, title, time.Now().Format(time.RFC3339), author, description, parseTags(tags), category, draft)

	full := front + body

	// 确保文件名后缀
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	savePath := filepath.Join(hugoContent)
	if subfolder != "" {
		savePath = filepath.Join(savePath, filepath.Clean(subfolder))
	}
	err := os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		http.Error(w, "创建目录失败: "+err.Error(), 500)
		return
	}

	fullPath := filepath.Join(savePath, filename)
	err = os.WriteFile(fullPath, []byte(full), 0644)
	if err != nil {
		http.Error(w, "保存文件失败: "+err.Error(), 500)
		return
	}

	// 异步构建 hugo
	go buildHugo()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func imageUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB限制
	if err != nil {
		http.Error(w, "解析表单失败", 400)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "获取文件失败", 400)
		return
	}
	defer file.Close()

	err = os.MkdirAll(imageFolder, os.ModePerm)
	if err != nil {
		http.Error(w, "创建目录失败", 500)
		return
	}

	filename := time.Now().Format("20060102150405_") + filepath.Base(header.Filename)
	dstPath := filepath.Join(imageFolder, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "保存文件失败", 500)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "保存文件失败", 500)
		return
	}

	// 返回图片相对路径给前端
	url := "/uploads/" + filename
	w.Write([]byte(url))
}

func parseTags(input string) string {
	var tags []string
	for _, t := range strings.Split(input, ",") {
		trim := strings.TrimSpace(t)
		if trim != "" {
			tags = append(tags, fmt.Sprintf(`"%s"`, trim))
		}
	}
	return strings.Join(tags, ", ")
}

func buildHugo() {
	log.Println("开始执行 Hugo 构建...")
	cmd := exec.Command("hugo", "--minify")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("构建失败: %v\n%s", err, out)
	} else {
		log.Println("构建完成")
	}
}
