// Home 路由
package handler

import (
	"makehugocontent/internal/admin"
	"net/http"
)

func (a *AuthHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// 获取当前用户的 session
	username, role := admin.GetCurrentUser(r)
	if username == "" {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	// 根据角色设置提示信息
	var message string

	// 登录后的角色信息处理
	switch role {
	case "viewer":
		message = "你可以浏览主页和上传新内容，但无法进行文章编辑。"
	case "editor":
		message = "你可以编辑文章和上传新内容。"
	case "admin":
		message = "你可以管理用户、编辑文章和上传新内容。"
	default:
		message = "未知角色，无法访问此页面。"
	}

	// 渲染 home.html 模板，并传递用户信息和提示信息
	a.Tmpl.ExecuteTemplate(w, "home.html", map[string]any{
		"Username": username,
		"Role":     role,
		"Message":  message,
	})
}
