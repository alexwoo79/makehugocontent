package handler

import (
	"fmt"
	"io"
	"makehugocontent/database"
	"makehugocontent/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	// 权限校验
	if !checkLogin(w, r) {
		return
	}

	if r.Method == http.MethodGet {
		utils.Render(w, "upload.html", nil)
		return
	}

	// POST 上传
	subfolder := r.FormValue("subfolder")
	filename := r.FormValue("filename")
	title := r.FormValue("title")
	author := r.FormValue("author")
	desc := r.FormValue("description")
	tags := r.FormValue("tags")
	category := r.FormValue("category")
	draft := r.FormValue("draft") == "true"

	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	front := fmt.Sprintf(`--- 
title : "%s"
date : "%s"
author : "%s"
description : "%s"
tags : [%s]
categories : ["%s"]
draft : %v
---

`, title, time.Now().Format(time.RFC3339), author, desc, tagArray(tags), category, draft)

	body := r.FormValue("body")
	content := front + body

	saveDir := filepath.Join("hugo/content/posts", subfolder)
	_ = os.MkdirAll(saveDir, os.ModePerm)
	full := filepath.Join(saveDir, filename)
	_ = os.WriteFile(full, []byte(content), 0644)

	// 记录上传
	cookie, _ := r.Cookie("session")
	database.DB.Exec(`INSERT INTO uploads(filename,filepath,user_id,author)
		VALUES(?,?,(SELECT id FROM users WHERE username=?),?)`,
		filename, filepath.Join(subfolder, filename), cookie.Value, author)

	// 触发构建
	utils.RunHugo()

	http.Redirect(w, r, "/admin/content", http.StatusSeeOther)
}

// 图片上传
func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) {
		return
	}
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "上传失败", 400)
		return
	}
	defer file.Close()

	_ = os.MkdirAll("static/uploads", os.ModePerm)
	dst := filepath.Join("static/uploads", time.Now().Format("20060102150405_")+header.Filename)
	out, _ := os.Create(dst)
	_, _ = io.Copy(out, file)
	out.Close()

	w.Write([]byte("/uploads/" + filepath.Base(dst)))
}

func tagArray(t string) string {
	out := []string{}
	for _, s := range strings.Split(t, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, fmt.Sprintf(`"%s"`, s))
		}
	}
	return strings.Join(out, ",")
}

// 简易登录检查
func checkLogin(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return false
	}
	return true
}

// 点击退出按钮 删除 session cookie
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// 清除 session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// 重定向到登录页
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
