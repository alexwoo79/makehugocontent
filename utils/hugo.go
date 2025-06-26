package utils

import (
	"log"
	"os/exec"
)

func RunHugo() {
	go func() {
		cmd := exec.Command("hugo", "--minify")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Hugo 构建失败:", err, string(out))
		} else {
			log.Println("Hugo 构建成功")
		}
	}()
}
