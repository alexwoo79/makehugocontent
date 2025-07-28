package utils

import (
	"log"
	"os/exec"
)

func RunHugo() {
	go func() {
		log.Println("开始执行 Hugo 构建...")

		// 使用系统环境中的 hugo 命令，避免依赖本地相对路径
		cmd := exec.Command("hugo", "--minify")
		cmd.Dir = "hugo" // 设置工作目录为 hugo 文件夹

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("构建失败: %v\n输出: %s", err, out)
		} else {
			log.Println("构建完成")
		}
	}()
}