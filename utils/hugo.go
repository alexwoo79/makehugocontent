package utils

import (
	"log"
	"os"
	"os/exec"
)

func RunHugo() {
	go func() {
		log.Println("开始执行 Hugo 构建...")
		// 切换到 hugo 目录
		err := os.Chdir("hugo")
		if err != nil {
			log.Printf("切换目录失败: %v\n", err)
			return
		}
		// 执行构建命令
		cmd := exec.Command("hugo", "--minify")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("构建失败: %v\n%s", err, out)
		} else {
			log.Println("构建完成")
		}
	}()
}
