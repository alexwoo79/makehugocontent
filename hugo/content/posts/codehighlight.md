---
title: "code highlight and style"
date: "2025-06-24T23:20:54+08:00"
author: "Alex"
description: "code高亮和样式"
tags: []
categories: [""]
draft: false
---

- 代码框显示效果如下：

```go {linenos=inline hl_lines=[3,"6-8"]}
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        fmt.Println("Value of i:", i)
    }
}
```
- 代码块默认具有highlight的样式，可以调用。详情可以参考一下文档。
[代码语法高亮chroma设置的文档链接](https://xyproto.github.io/splash/docs/)
[github库](https://github.com/alecthomas/chroma/tree/master)

以下是使用方法：
```go {linenos=false hl_lines=[3,"6-8"] style=emacs} 来调整`
... some code here...

`{`中的参数含义`}`
`linenos=false` -- 没有行号
`style=emacs`  -- 样式

- 没有行号的显示方式

```go {linenos=false hl_lines=[3,"6-8"] style=emacs}
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        fmt.Println("Value of i:", i)
    }
}
```

- 短代码实现macos样式的代码块
下图是一个使用md中使用shortcode的例题：
![alt text](/posts/image-6.png)


短代码模版要定义到shortcodes目录下 layouts/shortcodes/terminal.html
```html {linenos=inline style=emacs}
{{ $lang := .Get "lang" | default "text" }}
{{ $title := .Get "title" | default "Terminal" }}

<div class="terminal">
  <div class="terminal-header">
    <span class="terminal-dot red"></span>
    <span class="terminal-dot yellow"></span>
    <span class="terminal-dot green"></span>
    <span class="terminal-title">{{ $lang }}</span>
  </div>
  <div class="terminal-body">
    {{ highlight .Inner $lang "linenos=inline" }}
  </div>
</div>
```
下面是展示的效果，可以看到是macos样式。（通过定义css来实现。保存在self_style.css文件中,`layouts/static/css/self_style.css`）

{{< terminal lang="go" >}}
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        fmt.Println("Value of i:", i)
    }
}
{{< /terminal >}}