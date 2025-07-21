makehugocontent/
├── main.go # 启动程序
├── router/
│ └── routes.go # 路由注册
├── handler/
│ ├── auth.go # 登录注册逻辑
│ ├── upload.go # 上传文件、处理表单
│ ├── content.go # 管理 content 文件夹（列出、编辑、删除）
│ ├── image.go # 图片上传 API
├── utils/
│ ├── render.go # 模板渲染封装
│ ├── markdown.go # Front Matter 提取、解析等
│ └── hugo.go # 调用 hugo 命令等
├── templates/
│ └── \*.html # 前端页面模板
├── static/
│ ├── css/style.css
│ └── js/...
├── database/
│ └── db.go # SQLite 初始化与操作
│ └── data.db
├── content/ # Hugo 内容目录
├── go.mod
