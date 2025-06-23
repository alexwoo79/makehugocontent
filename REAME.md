# Hugo Admin Web Panel (Go + SQLite + EasyMDE)

本项目是一个基于 [Go](https://golang.org) 与 [Hugo](https://gohugo.io) 的 Markdown 网站管理系统，提供后台登录、文章管理、图片上传、草稿控制、多语言支持等功能，适用于 Hugo 博客或文档项目的定制化后台。

---

## 功能特性

- 用户注册与登录（基于 SQLite）  
- 多角色权限控制（admin / editor）  
- Markdown 编辑器（[EasyMDE](https://github.com/Ionaru/easy-markdown-editor)）  
- 上传文章并自动生成 Front Matter  
- 图片上传 + 自动插图  
- 支持草稿/分类/标签/多语言  
- Hugo 自动构建并预览

---

## 目录结构
```plaintext
makehugocontent/
├── main.go
├── users.db # 程序自动创建，勿提交
├── static/
│ └── uploads/ # 图片上传目录（初次需创建）
├── content/
│ ├── en/
│ └── zh/
├── public/ # Hugo 静态文件输出
├── templates/
│ ├── login.html
│ ├── register.html
│ └── upload.html
├── config.toml
├── .gitignore
└── README.md
```
---

## 快速开始

### 依赖环境

- Go ≥ 1.18  
- Hugo（安装请参考 https://gohugo.io/getting-started/install/）

### 启动服务

```bash
git clone https://github.com/alexwoo79/makehugocontent.git
cd makehugocontent
go mod init makehugocontent
go get github.com/mattn/go-sqlite3 golang.org/x/crypto/bcrypt
go run main.go
```

访问 http://localhost:8080/

使用说明
注册和登录
访问 /admin/register 注册新用户（默认为 editor 角色）

访问 /admin/login 登录后台管理

上传文章
进入上传页面 /admin/upload

填写标题、作者、描述、标签、分类、语言和草稿状态

使用 EasyMDE 编辑 Markdown 正文

支持图片上传自动插入

提交后会自动生成 Markdown 文件，并调用 Hugo 构建生成静态页面

多语言支持
已支持英文（en）和中文（zh）内容目录和 Hugo 多语言配置。
你可以根据需求扩展语言或内容。

许可证
MIT © 2025 Alex Woo


