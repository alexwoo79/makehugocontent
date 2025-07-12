---  
title: "网站内容创作指南"  
date: 2025-07-12  
author: "Woo Alex"  
description: "如何使用 Go + Hugo 平台进行内容创作与管理"  
tags: ["hugo", "markdown", "go", "cms"]  
categories: ["指南"]  

---

# 🚀 网站内容创作指南（基于 Go + Hugo Web 平台）

## 🧭 项目概述

本系统是一个使用 Go 语言开发、结合 Hugo 静态站点生成器的内容创作平台，适合团队协作发布技术文档、博客文章或内部门户内容。系统支持多角色管理、Markdown 文本创作、部门权限控制和内容审核流程。

## 👤 用户角色与权限说明

| 角色    | 权限描述                                                                 |
|---------|--------------------------------------------------------------------------|
| `viewer` | 仅可浏览公共主页及归属部门文章                                          |
| `editor` | 可上传/编辑所属部门的文章，管理自己的 Markdown 内容                    |
| `admin`  | 拥有全部权限，包括用户管理、部门分配、内容管理、文章编辑与删除等操作    |

## 🛠️ 使用指南

### 1️⃣ 注册与登录

- 访问：`/admin/register`
- 默认注册角色为 `viewer`，如需创作权限请联系管理员升级为 `editor`
- 登录地址：`/admin/login`

### 2️⃣ Markdown 内容创作流程

#### Step 1: 进入上传页面

- 登录后，点击「上传文章」进入 `/admin/upload`
- 需拥有 `editor` 或 `admin` 权限

#### Step 2: 填写元信息

上传页面包含以下字段：

- **标题 Title**
- **作者 Author**
- **描述 Description**
- **标签 Tags**：用英文逗号分隔
- **部门分类**
- **内容 Content**：使用 Markdown 撰写，支持实时预览与图片拖拽

#### Step 3: 图片上传支持

- 拖入图片后会自动上传至服务器 `/uploads/`
- 编辑器自动插入链接：

```markdown
![示例图片](/uploads/2025/07/sample.jpg)
```

#### Step 4: 提交发布

- 系统保存 `.md` 文件至 Hugo 内容目录

路径示例：

```
content/posts/<部门>/<用户名>/<标题>.md
```

系统会自动触发 Hugo 构建，生成可预览页面。

### 3️⃣ 内容管理

- `/admin/content` 页面提供列表
- 可编辑 / 删除已上传内容（权限控制）

### 4️⃣ 用户管理（仅 admin）

- 页面：`/admin/users`
- 功能：修改角色、部门分配、用户删除

### 5️⃣ 首页浏览

- 所有角色均可访问主页 `/`
- `viewer` 无法访问 `/admin/*`

## 📁 内容结构

```plaintext
content/
  posts/
    pm/
      alice/项目流程优化.md
    cz/
      bob/流程梳理.md
    general/
      carol/员工手册.md
```

## 💡 推荐流程

```
注册账号 ➜ 联系管理员升级权限 ➜ 登录 ➜ 上传 Markdown ➜ 预览 ➜ 发布
```

## 🛠 技术栈

| 技术         | 说明                             |
|--------------|----------------------------------|
| Go           | 后端服务开发                     |
| SQLite       | 用户与部门信息存储               |
| Hugo         | 静态网站生成                     |
| Tailwind CSS | 前端样式系统                     |
| HTMX         | 实现局部无刷新与交互             |
| EasyMDE      | 支持图像拖拽的 Markdown 编辑器   |

## ✅ 未来计划

- 内容审批机制（草稿 / 发布）
- PDF / Word 导出
- Hugo 分类优化
- 多语言与主题支持

---

> 版本：v1.0  
> 作者：Woo Alex  
> 项目地址：[github.com/alexwoo79/makehugocontent](https://github.com/alexwoo79/makehugocontent)
